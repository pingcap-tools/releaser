package dependency

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/releaser/config"
	"github.com/you06/releaser/pkg/git"
	"github.com/you06/releaser/pkg/types"
	"github.com/you06/releaser/pkg/utils"
)

var dependencyFiles = []string{"go.mod", "Cargo.toml"}

var (
	gomodRemoveComment          = regexp.MustCompile(`(.*)\/{2}.*$`)
	gomodModuleMatch            = regexp.MustCompile(`module\s(.*?)\s*$`)
	gomodRequireStart           = regexp.MustCompile(`require\s+\(\s*$`)
	gomodPartEnd                = regexp.MustCompile(`\)\s*$`)
	gomodDependencyMatchComment = regexp.MustCompile(`\s*(.*?)\s(.*?)\s*$`)
)

// Dependency struct
type Dependency struct {
	Config *config.Config
	Github *github.Client
	User   *github.User
}

// Config struct
type Config struct {
	Config *config.Config
	Github *github.Client
	User   *github.User
}

// New creates Dependency instance
func New(cfg *Config) *Dependency {
	return &Dependency{
		Config: cfg.Config,
		Github: cfg.Github,
		User:   cfg.User,
	}
}

// GetDependencies get all dependencies in a version
func (d *Dependency) GetDependencies(repo types.Repo, version string) ([]*types.Package, error) {
	var packages []*types.Package

	sha, err := d.GetVersionSHA(repo, version)
	if err != nil {
		// TODO: specify the error information
		sha, err = d.GetVersionRef(repo, version)
	}
	if err != nil {
		return packages, errors.Trace(err)
	}

	contents, err := d.ListContents(repo, sha)
	if err != nil {
		return packages, errors.Trace(err)
	}
	sha = strings.Split(sha, "\n")[0]

	for _, c := range contents {
		if c.GetType() != "file" {
			continue
		}
		var (
			filename = c.GetName()
			match    = false
		)
		for _, f := range dependencyFiles {
			if f == filename {
				match = true
			}
		}
		if match {
			content, err := d.GetContent(repo, sha, filename)
			if err != nil {
				return packages, errors.Trace(err)
			}
			contentStr, err := content.GetContent()
			if err != nil {
				return packages, errors.Trace(err)
			}
			p, err := parsePackage(repo, filename, contentStr)
			if err != nil {
				return packages, errors.Trace(err)
			}
			packages = append(packages, p)
		}
	}

	return packages, nil
}

// GetVersionSHA get SHA by a version
// TODO: it's too slow, find a better way to get sha
func (d *Dependency) GetVersionSHA(repo types.Repo, version string) (string, error) {
	gitClient := git.New(d.Config, &git.Config{
		User: d.User,
		Base: repo,
		Dir:  fmt.Sprintf("%s-%s", repo.Repo, version),
	})
	if err := gitClient.Clone(); err != nil {
		return "", errors.Trace(err)
	}
	defer func() {
		if err := gitClient.Clear(); err != nil {
			log.Error(err)
		}
	}()
	sha, err := gitClient.CheckTagSHA(version)
	return sha, errors.Trace(err)
}

// GetVersionRef gets ref by a version
func (d *Dependency) GetVersionRef(repo types.Repo, version string) (string, error) {
	version = strings.TrimLeft(version, "v")
	ctx, _ := utils.NewTimeoutContext()
	refs, _, err := d.Github.Git.GetRefs(ctx, repo.Owner, repo.Repo, "")
	if err != nil {
		return "", errors.Trace(err)
	}

	if sha, match := matchRef(refs, version); match {
		return sha, nil
	}
	versionSlice := strings.Split(version, ".")
	version = strings.Join(versionSlice[:len(versionSlice)-1], ".")
	if sha, match := matchRef(refs, version); match {
		return sha, nil
	}
	return "", errors.New("ref not found")
}

func matchRef(refs []*github.Reference, version string) (string, bool) {
	for _, ref := range refs {
		if strings.Contains(ref.GetRef(), version) {
			return ref.GetObject().GetSHA(), true
		}
	}
	return "", false
}

// ListContents list contents in a ref
func (d *Dependency) ListContents(repo types.Repo, sha string) ([]*github.RepositoryContent, error) {
	ctx, _ := utils.NewTimeoutContext()
	// TODO: what will happen if there are more than 100 files?
	_, contents, _, err := d.Github.Repositories.GetContents(ctx,
		repo.Owner, repo.Repo, "/", &github.RepositoryContentGetOptions{
			Ref: sha,
		})
	if err != nil {
		return nil, errors.Trace(err)
	}
	return contents, nil
}

// GetContent get specific content in a ref
func (d *Dependency) GetContent(repo types.Repo, sha, filename string) (*github.RepositoryContent, error) {
	ctx, _ := utils.NewTimeoutContext()
	content, _, _, err := d.Github.Repositories.GetContents(ctx,
		repo.Owner, repo.Repo, filename, &github.RepositoryContentGetOptions{
			Ref: sha,
		})
	if err != nil {
		return nil, errors.Trace(err)
	}
	return content, nil
}

func parsePackage(repo types.Repo, filename, fileContent string) (*types.Package, error) {
	var p *types.Package
	switch filename {
	case "Cargo.toml":
		c := types.NewCargo()
		if err := c.Parse(fileContent); err != nil {
			return nil, errors.Trace(err)
		}
		p = c.ToPackage()
	case "go.mod":
		p = parseGoMod(fileContent)
	}

	if p != nil {
		p.Repo = repo
		return p, nil
	}
	// unreachable
	return nil, errors.New("unreachable code")
}

// TODO: consider replace pattern
func parseGoMod(content string) *types.Package {
	var (
		p         types.Package
		inRequire = false
	)

	for _, line := range strings.Split(strings.ReplaceAll(content, "\r", ""), "\n") {
		line = strings.Trim(line, " ")
		if line == "" {
			continue
		}
		removeCommentMatch := gomodRemoveComment.FindStringSubmatch(line)
		if len(removeCommentMatch) == 2 {
			line = removeCommentMatch[1]
			line = strings.Trim(line, " ")
			if line == "" {
				continue
			}
		}

		moduleMatch := gomodModuleMatch.FindStringSubmatch(line)
		if len(moduleMatch) == 2 {
			p.Name = moduleMatch[1]
		}

		if inRequire {
			dependencyMatch := gomodDependencyMatchComment.FindStringSubmatch(line)
			if len(dependencyMatch) == 3 {
				p.Dependencies = append(p.Dependencies, types.Dependency{
					Name:    dependencyMatch[1],
					Version: dependencyMatch[2],
				})
			}
		}

		if gomodRequireStart.MatchString(line) {
			inRequire = true
		}
		if gomodPartEnd.MatchString(line) {
			inRequire = false
		}
	}
	return &p
}
