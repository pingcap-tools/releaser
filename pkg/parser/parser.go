package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/juju/errors"
	"github.com/you06/releaser/pkg/types"
)

var (
	releaseNoteLine = regexp.MustCompile(`^- ?(.*?)([a-zA-Z0-9-]+)\/([a-zA-Z0-9-]+)#([0-9]+).*$`)
)

// ReleaseNoteLang collects all release notes of a language
type ReleaseNoteLang struct {
	Lang      string
	Path      string
	RepoNotes []RepoReleaseNotes
	Version   string
}

// RepoReleaseNotes defines release notes in a repo
type RepoReleaseNotes struct {
	Repo  types.Repo
	Notes []ReleaseNote
}

// ReleaseNote is single release note
type ReleaseNote struct {
	Repo       types.Repo
	PullNumber int
	Note       string
}

// ParseContent parse content
func ParseContent(content string) ([]ReleaseNote, error) {
	var notes []ReleaseNote
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
		notes = append(notes, ReleaseNote{
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

// String ...
func (r ReleaseNote) String() string {
	return fmt.Sprintf("%s %s#%d", r.Note, r.Repo.String(), r.PullNumber)
}

// String ...
func (r RepoReleaseNotes) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s", r.Repo.Repo)
	b.WriteString("\n\n")
	for _, note := range r.Notes {
		fmt.Fprintf(&b, "- %s\n", note.String())
	}
	return b.String()
}

// String ...
func (r ReleaseNoteLang) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s", r.Version)
	b.WriteString("\n\n")
	for _, repoNotes := range r.RepoNotes {
		fmt.Fprintf(&b, "%s\n", repoNotes.String())
	}
	return b.String()
}
