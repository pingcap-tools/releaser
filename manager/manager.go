package manager

import (
	"strings"

	"github.com/google/go-github/v29/github"
	"github.com/juju/errors"
	"github.com/you06/releaser/pkg/note"
	"github.com/you06/releaser/pkg/pull"
	"github.com/you06/releaser/pkg/types"
)

// Config struct
type Config struct {
	SubCommand      string
	Version         string
	GithubToken     string
	ReleaseNoteRepo string
}

// Manager struct
type Manager struct {
	Config          *Config
	ReleaseNoteRepo types.Repo
	Github          *github.Client
	NoteCollector   *note.Collector
	PullCollector   *pull.Collector
}

// New create releaser manager
func New(cfg *Config) (*Manager, error) {
	repo, err := parseRepo(cfg.ReleaseNoteRepo)
	if err != nil {
		return nil, errors.Trace(err)
	}
	githubClient, err := initGithubClient(cfg.GithubToken)
	if err != nil {
		return nil, errors.Trace(err)
	}

	m := Manager{
		Config:          cfg,
		ReleaseNoteRepo: repo,
		Github:          githubClient,
		NoteCollector:   note.New(githubClient),
		PullCollector:   pull.New(githubClient),
	}
	if _, err := m.GetReleaseNoteRepo(); err != nil {
		return nil, errors.Trace(err)
	}

	return &m, nil
}

// Run start sub commands
func (m *Manager) Run() error {
	switch m.Config.SubCommand {
	case types.SubCmdPRList:
		return errors.Trace(m.runRRList())
	default:
		return errors.New("invalid sub command")
	}
}

func parseRepo(repo string) (types.Repo, error) {
	var (
		p = strings.Split(repo, "/")
		r types.Repo
	)
	if len(p) != 2 {
		return r, errors.Errorf("repo %s not valid", repo)
	}

	r.Owner = p[0]
	r.Repo = p[1]

	return r, nil
}
