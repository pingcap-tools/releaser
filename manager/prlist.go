package manager

import (
	"fmt"
	"strings"

	"github.com/juju/errors"
	"github.com/olekukonko/tablewriter"
	"github.com/you06/releaser/pkg/types"
)

func (m *Manager) runRRList() error {
	var (
		tableString      = strings.Builder{}
		table            = tablewriter.NewWriter(&tableString)
		noMilestoneRepos types.Repos
	)
	table.SetHeader([]string{"Repo", "PR", "Author", "Title"})
	for _, repo := range m.Repos {
		pulls, err := m.PullCollector.ListPRList(repo, m.Opt.Version)
		if err != nil {
			if strings.Contains(err.Error(), "milestone not found") {
				noMilestoneRepos = append(noMilestoneRepos, repo)
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

	if len(noMilestoneRepos) > 0 {
		fmt.Printf("No milestone repos: %s\n", noMilestoneRepos)
	}
	return nil
}
