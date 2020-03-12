package manager

import (
	"context"

	"github.com/google/go-github/v29/github"
	"github.com/juju/errors"
	"github.com/you06/releaser/pkg/utils"
	"golang.org/x/oauth2"
)

// GetReleaseNoteRepos gets repos info
func (m *Manager) GetReleaseNoteRepos() ([]*github.Repository, error) {
	var githubRepos []*github.Repository

	for _, repo := range m.Repos {
		ctx, _ := utils.NewTimeoutContext()
		githubRepo, _, err := m.Github.Repositories.Get(ctx, repo.Owner, repo.Repo)
		if err != nil {
			return githubRepos, errors.Trace(err)
		}
		githubRepos = append(githubRepos, githubRepo)
	}

	return githubRepos, nil
}

func initGithubClient(token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(context.Background(), ts)

	return github.NewClient(tc), nil
}

func (m *Manager) getGithubUser() (*github.User, error) {
	ctx, _ := utils.NewTimeoutContext()
	user, _, err := m.Github.Users.Get(ctx, "")
	return user, errors.Trace(err)
}
