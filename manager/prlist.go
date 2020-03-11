package manager

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/olekukonko/tablewriter"
)

func (m *Manager) runRRList() error {
	tableString := strings.Builder{}
	table := tablewriter.NewWriter(&tableString)
	table.SetHeader([]string{"Repo", "PR", "Author", "Title"})
	for _, repo := range m.Repos {
		pulls, err := m.PullCollector.ListPRList(repo, m.Opt.Version)
		if err != nil {
			if strings.Contains(err.Error(), "milestone not found") {
				fmt.Printf("can not find milestone %s in %s\n", m.Opt.Version, repo.String())
				continue
			}
			return errors.Trace(err)
		}
		for _, pull := range pulls {
			var (
				repo    = repo.String()
				pullStr = fmt.Sprintf("%d", pull.GetNumber())
				author  = pull.GetUser().GetLogin()
				title   = pull.GetTitle()
			)
			table.Append([]string{repo, pullStr, author, title})
		}
	}
	table.Render()
	fmt.Println(tableString.String())
	return nil
}
