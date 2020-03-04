package manager

import (
	"github.com/google/go-github/v29/github"
	"github.com/juju/errors"
	"github.com/you06/releaser/pkg/utils"
	"golang.org/x/oauth2"
)

// GetReleaseNoteRepo gets repo info
func (m *Manager) GetReleaseNoteRepo() (*github.Repository, error) {
	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()
	repo, _, err := m.Github.Repositories.Get(ctx, m.ReleaseNoteRepo.Owner, m.ReleaseNoteRepo.Repo)
	return repo, errors.Trace(err)
}

func initGithubClient(token string) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	ctx, cancel := utils.NewTimeoutContext()
	defer cancel()

	tc := oauth2.NewClient(ctx, ts)
	if err := ctx.Err(); err != nil {
		return nil, errors.Trace(err)
	}

	return github.NewClient(tc), nil
}
