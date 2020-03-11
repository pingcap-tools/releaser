package manager

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/juju/errors"
	"github.com/olekukonko/tablewriter"
)

const (
	iconHasNote = "√"
	iconNoNote  = "x"
)

var (
	releaseNoteStart     = regexp.MustCompile(`^.*release ?note.*$`)
	releaseNoteListMatch = regexp.MustCompile(`^- ?(.*)$`)
	releaseNoteNAMatch   = regexp.MustCompile(`^\s*(na\.?|no need\.?|none\.?)\s*$`)
)

func (m *Manager) runReleaseNotes() error {
	for _, repo := range m.Repos {
		pulls, err := m.PullCollector.ListPRList(repo, m.Opt.Version)
		if err != nil {
			if strings.Contains(err.Error(), "milestone not found") {
				fmt.Printf("can not find milestone %s in %s\n", m.Opt.Version, repo.String())
				continue
			}
			return errors.Trace(err)
		}
		releaseNotes, err := m.NoteCollector.ListReleaseNote(repo, m.Opt.Version)
		if err != nil {
			fmt.Printf("get release notes error %+v\n", err)
		}
		var (
			langs       []string
			tableString = strings.Builder{}
			table       = tablewriter.NewWriter(&tableString)
		)
		for _, releaseNote := range releaseNotes {
			langs = append(langs, releaseNote.Lang)
		}
		table.SetHeader(append([]string{"Repo", "PR", "Author", "Title", "Release Note"}, langs...))
		for _, pull := range pulls {
			var (
				repo                  = repo.String()
				pullStr               = fmt.Sprintf("%d", pull.GetNumber())
				author                = pull.GetUser().GetLogin()
				title                 = pull.GetTitle()
				langStatus            []string
				_, pullHasReleaseNote = hasReleaseNote(pull.GetBody())
			)
			if pullHasReleaseNote {
				langStatus = append(langStatus, iconHasNote)
			} else {
				langStatus = append(langStatus, iconNoNote)
			}
			for _, releaseNote := range releaseNotes {
				hasNote := false
				for _, note := range releaseNote.Notes {
					if note.PullNumber == pull.GetNumber() {
						hasNote = true
					}
				}
				if hasNote {
					langStatus = append(langStatus, iconHasNote)
				} else {
					langStatus = append(langStatus, iconNoNote)
				}
			}
			table.Append(append([]string{repo, pullStr, author, title}, langStatus...))
		}
		table.Render()
		fmt.Println(tableString.String())
	}

	return nil
}

func hasReleaseNote(body string) (string, bool) {
	var (
		findReleaseNoteStart = false
		releaseNoteLines     []string
	)
	for _, line := range strings.Split(strings.ReplaceAll(body, "\r", ""), "\n") {
		line = strings.Trim(line, " ")
		if findReleaseNoteStart {
			releaseNoteLines = append(releaseNoteLines, strings.Trim(line, " "))
		}
		if releaseNoteStart.MatchString(strings.ToLower(line)) {
			findReleaseNoteStart = true
		}
	}
	for _, line := range releaseNoteLines {
		listMatch := releaseNoteListMatch.FindStringSubmatch(line)
		if len(listMatch) == 2 {
			line = listMatch[1]
		}
		if releaseNoteNAMatch.MatchString(strings.ToLower(line)) {
			return "", false
		}
		if line != "" {
			return line, true
		}
	}
	return "", false
}
