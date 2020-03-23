package note

import (
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v30/github"
	"github.com/juju/errors"
	"github.com/you06/releaser/config"
	"github.com/you06/releaser/pkg/parser"
	"github.com/you06/releaser/pkg/types"
	"github.com/you06/releaser/pkg/utils"
)

var langs = []string{"cn", "en", "jp"}
var (
	commentPattern  = regexp.MustCompile(`<!--[^>]*-->`)
	repoPattern     = regexp.MustCompile(`^## (.*)$`)
	releaseNoteLine = regexp.MustCompile(`^- ?(.*?)([a-zA-Z0-9-]+)\/([a-zA-Z0-9-]+)#([0-9]+).*$`)
)

// Collector for collect release notes
type Collector struct {
	github         *github.Client
	Config         *config.Config
	relaseNoteRepo types.Repo
}

// New creates Collector instance
func New(g *github.Client, config *config.Config, relaseNoteRepo types.Repo) *Collector {
	return &Collector{g, config, relaseNoteRepo}
}

// ListReleaseNote lists release notes
func (c *Collector) ListReleaseNote(product types.Product, version string) ([]parser.ReleaseNoteLang, error) {
	var (
		filePath = strings.ReplaceAll(c.Config.ReleaseNotePath, "{product}", product.Name)
		notes    []parser.ReleaseNoteLang
	)
	version = strings.ToLower(strings.Trim(version, "v"))
	contents, err := c.ListContents(filePath, version)
	if err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return notes, nil
		}
		return notes, errors.Trace(err)
	}
	for _, content := range contents {
		var (
			name     = content.GetName()
			fullPath = path.Join(filePath, name)
		)
		if content.GetType() != "file" {
			continue
		}
		lang, match := matchLang(name, version)
		if !match {
			continue
		}
		releaseNotes, err := c.ParseContent(fullPath)
		if err != nil {
			return notes, errors.Trace(err)
		}
		notes = append(notes, parser.ReleaseNoteLang{
			Lang:      lang,
			Path:      fullPath,
			Version:   version,
			RepoNotes: releaseNotes,
		})
	}
	return notes, nil
}

// ListContents list contents in a path
func (c *Collector) ListContents(filePath, version string) ([]*github.RepositoryContent, error) {
	// FIXME: should use full name of the repo
	ctx, _ := utils.NewTimeoutContext()
	// TODO: what will happen if there are more than 100 files?
	_, contents, _, err := c.github.Repositories.GetContents(ctx,
		c.relaseNoteRepo.Owner, c.relaseNoteRepo.Repo, filePath, &github.RepositoryContentGetOptions{})
	if err != nil {
		return nil, errors.Trace(err)
	}
	return contents, nil
}

// ParseContent parses content and get all release notes
func (c *Collector) ParseContent(fullPath string) ([]parser.RepoReleaseNotes, error) {
	var (
		repos      []parser.RepoReleaseNotes
		reposNotes *parser.RepoReleaseNotes
	)
	content, err := c.GetFileContent(fullPath)
	if err != nil {
		return repos, errors.Trace(err)
	}

	content = commentPattern.ReplaceAllString(content, "")
	content = strings.ReplaceAll(content, "\r", "")

	for _, line := range strings.Split(content, "\n") {
		line = strings.Trim(line, " ")

		repoMatch := repoPattern.FindStringSubmatch(line)
		if len(repoMatch) == 2 {
			// push old repo into repos
			if reposNotes.Repo.Repo != "" {
				repos = append(repos, *reposNotes)
			}
			// compose new repo
			reposNotes = &parser.RepoReleaseNotes{
				Repo: types.Repo{Repo: repoMatch[1]},
			}
			continue
		}

		match := releaseNoteLine.FindStringSubmatch(line)
		if len(match) != 5 {
			continue
		}
		note, owner, repo, numberStr := match[1], match[2], match[3], match[4]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return repos, errors.Trace(err)
		}
		reposNotes.Notes = append(reposNotes.Notes, parser.ReleaseNote{
			Repo: types.Repo{
				Owner: owner,
				Repo:  repo,
			},
			PullNumber: number,
			Note:       note,
		})
	}
	return repos, nil
}

// GetFileContent gets content of file and decode it to string
func (c *Collector) GetFileContent(p string) (string, error) {
	ctx, _ := utils.NewTimeoutContext()
	content, _, _, err := c.github.Repositories.GetContents(ctx,
		c.relaseNoteRepo.Owner, c.relaseNoteRepo.Repo, p, &github.RepositoryContentGetOptions{})
	if err != nil {
		return "", errors.Trace(err)
	}
	if content.GetType() != "file" {
		return "", errors.New("content is not a file")
	}
	decoded, err := content.GetContent()
	if err != nil {
		return "", errors.Trace(err)
	}
	return decoded, nil
}

func matchLang(name, version string) (string, bool) {
	name = strings.ToLower(name)
	if !strings.Contains(name, version) {
		return "", false
	}
	for _, lang := range langs {
		if strings.Contains(name, lang) {
			return lang, true
		}
	}
	return "en", true
}
