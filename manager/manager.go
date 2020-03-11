package manager

import (
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/juju/errors"
	"github.com/you06/releaser/config"
	"github.com/you06/releaser/pkg/note"
	"github.com/you06/releaser/pkg/pull"
	"github.com/you06/releaser/pkg/types"
)

// Manager struct
type Manager struct {
	Config         *config.Config
	Opt            *Option
	Repos          []types.Repo
	RelaseNoteRepo types.Repo
	Github         *github.Client
	NoteCollector  *note.Collector
	PullCollector  *pull.Collector
}

// Option for usage
type Option struct {
	Version string
}

// New create releaser manager
func New(cfg *config.Config, opt *Option) (*Manager, error) {
	repos, err := parseRepos(cfg.Repos)
	if err != nil {
		return nil, errors.Trace(err)
	}
	relaseNoteRepo, err := parseRepo(cfg.ReleaseNoteRepo)
	if err != nil {
		return nil, errors.Trace(err)
	}
	githubClient, err := initGithubClient(cfg.GithubToken)
	if err != nil {
		return nil, errors.Trace(err)
	}

	m := Manager{
		Config:         cfg,
		Opt:            opt,
		Repos:          repos,
		RelaseNoteRepo: relaseNoteRepo,
		Github:         githubClient,
		NoteCollector:  note.New(githubClient, cfg, relaseNoteRepo),
		PullCollector:  pull.New(githubClient, cfg),
	}
	if _, err := m.GetReleaseNoteRepos(); err != nil {
		return nil, errors.Trace(err)
	}

	return &m, nil
}

// Run start sub commands
func (m *Manager) Run(subCommand string) error {
	switch subCommand {
	case types.SubCmdPRList:
		return errors.Trace(m.runRRList())
	case types.SubCmdReleaseNotes:
		return errors.Trace(m.runReleaseNotes())
	default:
		return errors.New("invalid sub command")
	}
}

func parseRepos(repoStrs []string) ([]types.Repo, error) {
	var (
		repos []types.Repo
	)

	for _, repoStr := range repoStrs {
		repo, err := parseRepo(repoStr)
		if err != nil {
			return repos, errors.Trace(err)
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

func parseRepo(repo string) (types.Repo, error) {
	var (
		p = strings.Split(repo, "/")
		r types.Repo
	)
	if len(p) != 2 {
		return r, errors.Errorf("repo %s not valid", repo)
	}

	r.Owner, r.Repo = p[0], p[1]

	return r, nil
}
