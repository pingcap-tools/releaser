package manager

import (
	"regexp"
	"strings"
)

const (
	iconHasNote = "√"
	iconNoNote  = "x"
)

var shouldExistLangs = []string{"cn", "en"}

var (
	commentPattern       = regexp.MustCompile(`<!--[^>]*-->`)
	releaseNoteStart     = regexp.MustCompile(`^.*release ?note.*$`)
	releaseNoteListMatch = regexp.MustCompile(`^- ?(.*)$`)
	releaseNoteNAMatch   = regexp.MustCompile(`^\s*(na\.?|no need\.?|none\.?|no\.?|no. it's trivial.)\s*$`)
	titlePattern         = regexp.MustCompile(`^\#{1,3}\ .*$`)
)

func (m *Manager) runReleaseNotes() error {
	return nil
}

// func (m *Manager) runReleaseNotes() error {
// 	var errs []error
// 	for _, product := range m.Products {
// 		pulls, err := m.PullCollector.ListPRList(repo, m.Opt.Version)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "milestone not found") {
// 				fmt.Printf("can not find milestone %s in %s\n", m.Opt.Version, repo.String())
// 				continue
// 			}
// 			return errors.Trace(err)
// 		}
// 		releaseNotes, err := m.NoteCollector.ListReleaseNote(product, m.Opt.Version)
// 		if err != nil {
// 			fmt.Printf("get release notes error %+v\n", err)
// 		}
// 		var (
// 			langs            []string
// 			tableString      strings.Builder
// 			table            = tablewriter.NewWriter(&tableString)
// 			slackTableString strings.Builder
// 			slackTable       = tablewriter.NewWriter(&slackTableString)
// 			rows             [][]string
// 		)
// 		for _, releaseNote := range releaseNotes {
// 			langs = append(langs, releaseNote.Lang)
// 		}
// 		table.SetHeader(append([]string{"Repo", "PR", "Author", "Title", "Release Note"}, langs...))
// 		slackTable.SetHeader(append([]string{"Repo", "PR", "Author", "Title", "Release Note"}, langs...))
// 		for _, pull := range pulls {
// 			var (
// 				repo                  = repo.String()
// 				pullStr               = fmt.Sprintf("%d", pull.GetNumber())
// 				author                = pull.GetUser().GetLogin()
// 				title                 = pull.GetTitle()
// 				langStatus            []string
// 				_, pullHasReleaseNote = hasReleaseNote(pull.GetBody())
// 			)
// 			if pullHasReleaseNote {
// 				langStatus = append(langStatus, iconHasNote)
// 			} else {
// 				langStatus = append(langStatus, iconNoNote)
// 			}
// 			for _, releaseNote := range releaseNotes {
// 				hasNote := false
// 				for _, note := range releaseNote.Notes {
// 					if note.PullNumber == pull.GetNumber() {
// 						hasNote = true
// 					}
// 				}
// 				if hasNote {
// 					langStatus = append(langStatus, iconHasNote)
// 				} else {
// 					langStatus = append(langStatus, iconNoNote)
// 				}
// 			}
// 			rows = append(rows, (append([]string{repo, pullStr, author, title}, langStatus...)))
// 		}

// 		for _, row := range rows {
// 			table.Append(row)
// 			allOk := true
// 			for i := 0; i < len(langs); i++ {
// 				if row[len(row)-i-1] == iconNoNote {
// 					allOk = false
// 				}
// 			}
// 			if !allOk && row[len(row)-len(langs)-1] == iconHasNote {
// 				slackTable.Append(row)
// 			}
// 		}

// 		table.Render()
// 		slackTable.Render()

// 		if len(langs) != 2 {
// 			slackTableString.Reset()
// 			var missingLangs []string
// 			for _, lang := range shouldExistLangs {
// 				has := false
// 				for _, existLang := range langs {
// 					if lang == existLang {
// 						has = true
// 					}
// 				}
// 				if !has {
// 					missingLangs = append(missingLangs, lang)
// 				}
// 			}
// 			existLangs := strings.Join(langs, ", ")
// 			if len(langs) == 0 {
// 				existLangs = "None"
// 			}
// 			fmt.Fprintf(&slackTableString, "Found language: %s. Missing: %s.", existLangs, strings.Join(missingLangs, ", "))
// 		}

// 		fmt.Println(tableString.String())
// 		errs = append(errs, m.SendMessage(fmt.Sprintf("```%s Version:%s\n%s```", repo.String(), m.Opt.Version, slackTableString.String())))
// 	}

// 	for _, err := range errs {
// 		if err != nil {
// 			return errors.Trace(err)
// 		}
// 	}
// 	return nil
// }

func hasReleaseNote(body string) (string, bool) {
	var (
		findReleaseNoteStart = false
		releaseNoteLines     []string
	)

	body = commentPattern.ReplaceAllString(body, "")
	body = strings.ReplaceAll(body, "\r", "")

	for _, line := range strings.Split(body, "\n") {
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
		if titlePattern.MatchString(strings.ToLower(line)) {
			return "", false
		}
		if line != "" {
			return removeHeader(line), true
		}
	}
	return "", false
}

func removeHeader(line string) string {
	origin := line
	line = strings.TrimLeft(line, " ")
	line = strings.TrimLeft(line, "*")
	line = strings.TrimLeft(line, "-")
	if line == origin {
		return line
	}
	return removeHeader(line)
}
