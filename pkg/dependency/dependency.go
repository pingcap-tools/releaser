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

const (
	cargo = "Cargo.toml"
	gomod = "go.mod"
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

	ref := version
	contents, err := d.ListContents(repo, version)
	if err != nil {
		if strings.Contains(err.Error(), "No commit found") {
			ref, err = d.GetVersionRef(repo, version)
			if err != nil {
				return packages, errors.Trace(err)
			}
			contents, err = d.ListContents(repo, ref)
			if err != nil {
				return packages, errors.Trace(err)
			}
		} else {
			return packages, errors.Trace(err)
		}
	}

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
			content, err := d.GetContent(repo, ref, filename)
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
			p.URL = content.GetHTMLURL()
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

	var (
		refs    []*github.Reference
		batch   []*github.Reference
		page    = 0
		perpage = 100
	)

	for page == 0 || len(batch) == perpage {
		page++
		batch, _, err := d.Github.Git.ListRefs(ctx, repo.Owner, repo.Repo, &github.ReferenceListOptions{
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perpage,
			},
		})
		if err != nil {
			return "", errors.Trace(err)
		}
		refs = append(refs, batch...)
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
		repo.Owner, repo.Repo, "", &github.RepositoryContentGetOptions{
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
	case cargo:
		c := types.NewCargo()
		if err := c.Parse(fileContent); err != nil {
			return nil, errors.Trace(err)
		}
		p = c.ToPackage()
		p.Type = cargo
	case gomod:
		p = parseGoMod(fileContent)
		p.Type = gomod
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
