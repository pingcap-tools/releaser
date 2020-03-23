package pull

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v30/github"
	"github.com/juju/errors"
	"github.com/you06/releaser/config"
	"github.com/you06/releaser/pkg/types"
	"github.com/you06/releaser/pkg/utils"
)

// Collector for collect pulls
type Collector struct {
	github *github.Client
	Config *config.Config
}

// New creates Collector instance
func New(g *github.Client, config *config.Config) *Collector {
	return &Collector{g, config}
}

// ListPRList lists PR list in a version
func (c *Collector) ListPRList(repo types.Repo, version string) ([]*github.PullRequest, error) {
	milestone, err := c.GetVersionMilestone(repo, version)
	if err != nil {
		return []*github.PullRequest{}, errors.Trace(err)
	}
	_, pulls, err := c.ListAllMilestoneIssues(repo, milestone)
	if err != nil {
		return []*github.PullRequest{}, errors.Trace(err)
	}
	return pulls, nil
}

// GetVersionMilestone gets milestone
func (c *Collector) GetVersionMilestone(repo types.Repo, version string) (*github.Milestone, error) {
	milestones, err := c.ListAllMilestones(repo)
	if err != nil {
		return nil, errors.Trace(err)
	}
	for _, milestone := range milestones {
		if strings.Contains(strings.ToLower(milestone.GetTitle()), strings.ToLower(version)) {
			return milestone, nil
		}
	}
	return nil, errors.New("milestone not found")
}

// ListAllMilestones lists milestones list in a version
func (c *Collector) ListAllMilestones(repo types.Repo) ([]*github.Milestone, error) {
	var (
		page    = 0
		perpage = 100
		all     []*github.Milestone
		batch   []*github.Milestone
		err     error
	)
	for page == 0 || len(batch) == perpage {
		page++
		ctx, _ := utils.NewTimeoutContext()
		batch, _, err = c.github.Issues.ListMilestones(ctx, repo.Owner, repo.Repo, &github.MilestoneListOptions{
			State: "all",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perpage,
			},
		})
		if err != nil {
			return []*github.Milestone{}, errors.Trace(err)
		}
		all = append(all, batch...)
	}
	return all, nil
}

// ListAllOpenedMilestones lists milestones in opened state
func (c *Collector) ListAllOpenedMilestones(repo types.Repo) ([]*github.Milestone, error) {
	var (
		page    = 0
		perpage = 100
		all     []*github.Milestone
		batch   []*github.Milestone
		err     error
	)
	for page == 0 || len(batch) == perpage {
		page++
		ctx, _ := utils.NewTimeoutContext()
		batch, _, err = c.github.Issues.ListMilestones(ctx, repo.Owner, repo.Repo, &github.MilestoneListOptions{
			State: "open",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perpage,
			},
		})
		if err != nil {
			return []*github.Milestone{}, errors.Trace(err)
		}
		all = append(all, batch...)
	}
	return all, nil
}

// ListAllMilestoneIssues lists issues and pull requests in a milestone
func (c *Collector) ListAllMilestoneIssues(repo types.Repo, milestone *github.Milestone) ([]*github.Issue, []*github.PullRequest, error) {
	var (
		issues []*github.Issue
		pulls  []*github.PullRequest
	)
	if milestone.GetID() == 0 {
		return issues, pulls, errors.New("milestone with id 0")
	}

	all, err := c.ListAllIssuesFrom(repo, milestone)
	if err != nil {
		return issues, pulls, errors.Trace(err)
	}

	for _, item := range all {
		if item.IsPullRequest() {
			// ctx, _ := utils.NewTimeoutContext()
			// pull, _, err := c.github.PullRequests.Get(ctx, repo.Owner, repo.Repo, item.GetNumber())
			// if err != nil {
			// 	return issues, pulls, errors.Trace(err)
			// }
			pulls = append(pulls, issue2pull(item))
		} else {
			issues = append(issues, item)
		}
	}

	return issues, pulls, err
}

// ListAllIssuesFrom lists issues from
func (c *Collector) ListAllIssuesFrom(repo types.Repo, milestone *github.Milestone) ([]*github.Issue, error) {
	var (
		page    = 0
		perpage = 100
		all     []*github.Issue
		batch   []*github.Issue
	)

	for page == 0 || len(batch) == perpage {
		page++
		ctx, _ := utils.NewTimeoutContext()
		batch, _, err := c.github.Issues.ListByRepo(ctx, repo.Owner, repo.Repo, &github.IssueListByRepoOptions{
			Milestone: fmt.Sprintf("%d", milestone.GetNumber()),
			State:     "all",
			ListOptions: github.ListOptions{
				Page:    page,
				PerPage: perpage,
			},
		})
		if err != nil {
			return all, errors.Trace(err)
		}
		all = append(all, batch...)
	}

	return all, nil
}

// issue2pull transfer issue to pull with some common fields
// for those fields exist in pulls only
// you should get request by another API call
func issue2pull(issue *github.Issue) *github.PullRequest {
	if !issue.IsPullRequest() {
		return nil
	}

	var labels []*github.Label
	for _, label := range issue.Labels {
		labels = append(labels, label)
	}

	pull := github.PullRequest{
		ID:          issue.ID,
		Number:      issue.Number,
		State:       issue.State,
		Locked:      issue.Locked,
		Title:       issue.Title,
		Body:        issue.Body,
		CreatedAt:   issue.CreatedAt,
		UpdatedAt:   issue.UpdatedAt,
		ClosedAt:    issue.ClosedAt,
		Labels:      labels,
		User:        issue.User,
		URL:         issue.URL,
		HTMLURL:     issue.HTMLURL,
		IssueURL:    issue.URL,
		CommentsURL: issue.CommentsURL,
		Milestone:   issue.Milestone,
		NodeID:      issue.NodeID,
	}
	return &pull
}
