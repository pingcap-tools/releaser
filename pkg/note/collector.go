package note

import (
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/juju/errors"
	"github.com/you06/releaser/config"
	"github.com/you06/releaser/pkg/parser"
	"github.com/you06/releaser/pkg/types"
	"github.com/you06/releaser/pkg/utils"
)

var langs = []string{"cn", "en", "jp"}
var (
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
func (c *Collector) ListReleaseNote(repo types.Repo, version string) ([]parser.ReleaseNoteLang, error) {
	var (
		filePath = strings.ReplaceAll(c.Config.ReleaseNotePath, "{repo}", repo.Repo)
		notes    []parser.ReleaseNoteLang
	)
	version = strings.ToLower(strings.Trim(version, "v"))
	contents, err := c.ListContents(repo, filePath, version)
	if err != nil {
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
		releaseNotes, err := c.ParseContent(repo, fullPath)
		if err != nil {
			return notes, errors.Trace(err)
		}
		notes = append(notes, parser.ReleaseNoteLang{
			Lang:  lang,
			Path:  fullPath,
			Notes: releaseNotes,
		})
	}
	return notes, nil
}

// ListContents list contents in a path
func (c *Collector) ListContents(repo types.Repo, filePath, version string) ([]*github.RepositoryContent, error) {
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
func (c *Collector) ParseContent(repo types.Repo, fullPath string) ([]parser.ReleaseNote, error) {
	var notes []parser.ReleaseNote
	content, err := c.GetFileContent(repo, fullPath)
	if err != nil {
		return notes, errors.Trace(err)
	}
	for _, line := range strings.Split(strings.ReplaceAll(content, "\r", ""), "\n") {
		match := releaseNoteLine.FindStringSubmatch(strings.Trim(line, " "))
		if len(match) != 5 {
			continue
		}
		note, owner, repo, numberStr := match[1], match[2], match[3], match[4]
		number, err := strconv.Atoi(numberStr)
		if err != nil {
			return notes, errors.Trace(err)
		}
		notes = append(notes, parser.ReleaseNote{
			Repo: types.Repo{
				Owner: owner,
				Repo:  repo,
			},
			PullNumber: number,
			Note:       note,
		})
	}
	return notes, nil
}

// GetFileContent gets content of file and decode it to string
func (c *Collector) GetFileContent(repo types.Repo, p string) (string, error) {
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
