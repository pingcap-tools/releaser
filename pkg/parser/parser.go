package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/you06/releaser/pkg/types"
)

const (
	FOUR_SPACE = "    "
	OTHER_TYPE = "Others"
	HEADER     = `---
title: {name} {version} Release Notes
category: Releases
aliases: ['/docs/dev/releases/{version}/']
---

# {name} {version} Release Notes

Release date: {date}

TiDB version: {version}`
	DATE_FORMAT = "January 02, 2006"
)

var (
	releaseNoteLine = regexp.MustCompile(`^- ?(.*?)([a-zA-Z0-9-]+)\/([a-zA-Z0-9-]+)#([0-9]+).*$`)
)

// ReleaseNoteLang collects all release notes of a language
type ReleaseNoteLang struct {
	Name               string
	Lang               string
	Path               string
	ReleaseNoteClasses map[string][]RepoReleaseNotes
	Structure          []types.ProductItem
	Version            string
}

// RepoReleaseNotes defines release notes in a repo
type RepoReleaseNotes struct {
	Repo   types.Repo
	Rename types.Repo
	Notes  []ReleaseNote
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

// RemoveAllNote remote all notes
func (r *ReleaseNoteLang) RemoveAllNote() {
	for k := range r.ReleaseNoteClasses {
		for i := range r.ReleaseNoteClasses[k] {
			r.ReleaseNoteClasses[k][i].Notes = nil
		}
	}
}

// String ...
func (r ReleaseNote) String() string {
	return fmt.Sprintf("%s [#%d](https://github.com/%s/pull/%d)", Ucfirst(r.Note), r.PullNumber, r.Repo.String(), r.PullNumber)
}

// String ...
func (r RepoReleaseNotes) String() string {
	var b strings.Builder
	fmt.Fprintf(&b, "## %s", r.Rename.Repo)
	if len(r.Notes) > 0 {
		b.WriteString("\n\n")
	} else {
		b.WriteString("\n")
	}
	for _, note := range r.Notes {
		fmt.Fprintf(&b, "- %s\n", note.String())
	}
	return b.String()
}

func (r ReleaseNoteLang) formatHeader() string {
	header := strings.ReplaceAll(HEADER, "{name}", r.Name)
	header = strings.ReplaceAll(header, "{version}", r.Version)
	header = strings.ReplaceAll(header, "{date}", time.Now().Format(DATE_FORMAT))
	return header
}

// String ...
func (r ReleaseNoteLang) String() string {
	var b strings.Builder

	b.WriteString(r.formatHeader())
	b.WriteString("\n\n")

	var haveReleaseNote func(item types.ProductItem, repos []RepoReleaseNotes) bool
	haveReleaseNote = func(item types.ProductItem, repos []RepoReleaseNotes) bool {
		if item.Title == "" {
			var repo RepoReleaseNotes
			for _, r := range repos {
				if r.Repo == item.Repo {
					repo = r
					break
				}
			}
			return len(repo.Notes) != 0
		}

		for _, projectItem := range item.Children {
			if haveReleaseNote(projectItem, repos) {
				return true
			}
		}
		return false
	}

	var writeProjectItems func(b *strings.Builder, depth int, structure []types.ProductItem, repos []RepoReleaseNotes)
	writeProjectItems = func(b *strings.Builder, depth int, structure []types.ProductItem, repos []RepoReleaseNotes) {
		for _, projectItem := range structure {
			// skip repos don't have release notes
			if !haveReleaseNote(projectItem, repos) {
				continue
			}

			// write header
			if depth == 0 {
				b.WriteString("+ ")
			} else if depth == 1 {
				b.WriteString(FOUR_SPACE)
				b.WriteString("- ")
			} else {
				for i := 0; i < depth; i++ {
					b.WriteString(FOUR_SPACE)
				}
				b.WriteString("* ")
			}

			// recursive go deeper
			if projectItem.Title != "" {
				fmt.Fprintf(b, "%s\n\n", projectItem.Title)
				writeProjectItems(b, depth+1, projectItem.Children, repos)
				continue
			}

			// find repo
			var repo RepoReleaseNotes
			for _, r := range repos {
				if r.Repo == projectItem.Repo {
					repo = r
					break
				}
			}

			// list repos
			if repo.Rename.Repo != "" {
				b.WriteString(repo.Rename.Repo)
			} else {
				b.WriteString(repo.Repo.Repo)
			}
			b.WriteString("\n\n")

			for _, note := range repo.Notes {
				indent := depth + 1
				if indent == 1 {
					b.WriteString(FOUR_SPACE)
					b.WriteString("- ")
				} else {
					for i := 0; i < indent; i++ {
						b.WriteString(FOUR_SPACE)
					}
					b.WriteString("* ")
				}
				b.WriteString(note.String())
				b.WriteString("\n")
			}
			b.WriteString("\n")
		}
	}

	for tp, repos := range r.ReleaseNoteClasses {
		if tp == OTHER_TYPE {
			continue
		}
		fmt.Fprintf(&b, "## %s\n\n", tp)
		writeProjectItems(&b, 0, r.Structure, repos)
	}
	if repos, ok := r.ReleaseNoteClasses[OTHER_TYPE]; ok {
		fmt.Fprintf(&b, "## %s\n\n", OTHER_TYPE)
		writeProjectItems(&b, 0, r.Structure, repos)
	}

	// remove last "\n" character, a little hack
	res := b.String()
	return res[:len(res)-1]
}

func Ucfirst(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}
