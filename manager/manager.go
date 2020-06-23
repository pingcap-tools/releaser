package manager

import (
	"regexp"
	"strings"

	"github.com/google/go-github/v30/github"
	"github.com/juju/errors"
	"github.com/nlopes/slack"
	"github.com/you06/releaser/config"
	"github.com/you06/releaser/pkg/dependency"
	"github.com/you06/releaser/pkg/note"
	"github.com/you06/releaser/pkg/pull"
	"github.com/you06/releaser/pkg/types"
)

var (
	structurePattern = regexp.MustCompile(`(.*)?:\s?(.*)`)
)

// Manager struct
type Manager struct {
	Config   *config.Config
	Opt      *Option
	User     *github.User
	Repos    []types.Repo
	Products []types.Product

	RelaseNoteRepo      types.Repo
	Github              *github.Client
	Slack               *slack.Client
	NoteCollector       *note.Collector
	PullCollector       *pull.Collector
	DependencyCollector *dependency.Dependency
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
	products, err := parseProducts(cfg.Products)
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
	user, err := getGithubUser(githubClient)
	if err != nil {
		return nil, errors.Trace(err)
	}

	m := Manager{
		Config:         cfg,
		Opt:            opt,
		Repos:          repos,
		Products:       products,
		RelaseNoteRepo: relaseNoteRepo,
		Github:         githubClient,
		User:           user,
		Slack:          initSlackClient(cfg.SlackToken),
		NoteCollector:  note.New(githubClient, cfg, relaseNoteRepo),
		PullCollector:  pull.New(githubClient, cfg),
		DependencyCollector: dependency.New(&dependency.Config{
			Config: cfg,
			Github: githubClient,
			User:   user,
		}),
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
	case types.SubCmdGenerateReleaseNote:
		return errors.Trace(m.runGenerateReleaseNote())
	case types.SubCmdCheckModule:
		return errors.Trace(m.runCheckModule())
	default:
		return errors.New("invalid sub command")
	}
}

func parseProducts(products []config.Product) ([]types.Product, error) {
	var p []types.Product

	for _, product := range products {
		repos, err := parseRepos(product.Repos)
		if err != nil {
			return nil, errors.Trace(err)
		}
		renames, err := parseRenames(product.Rename)
		if err != nil {
			return nil, errors.Trace(err)
		}
		structure, err := parseStructure(product.Structure)
		if err != nil {
			return nil, errors.Trace(err)
		}
		p = append(p, types.Product{
			Name:      product.Name,
			Repos:     repos,
			Renames:   renames,
			Structure: structure,
		})
	}

	return p, nil
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

func parseRenames(renameSources map[string]string) (map[types.Repo]types.Repo, error) {
	renameRepo := make(map[types.Repo]types.Repo)
	for k, v := range renameSources {
		kRepo, err := parseRepo(k)
		if err != nil {
			return nil, errors.Trace(err)
		}
		vRepo, err := parseRepo(v)
		if err != nil {
			return nil, errors.Trace(err)
		}
		renameRepo[kRepo] = vRepo
	}
	return renameRepo, nil
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

func parseStructure(structure []string) ([]types.ProductItem, error) {
	var ps []types.ProductItem
	for _, item := range structure {
		structureMatch := structurePattern.FindStringSubmatch(item)
		if len(structureMatch) != 3 {
			repo, err := parseRepo(item)
			if err != nil {
				return nil, errors.Trace(err)
			}
			ps = append(ps, types.ProductItem{
				Repo: repo,
			})
			continue
		}
		title, repoRaws := structureMatch[1], structureMatch[2]
		var children []types.ProductItem
		for _, repoRaw := range strings.Split(repoRaws, ",") {
			repoRaw = strings.Trim(repoRaw, " ")
			repo, err := parseRepo(repoRaw)
			if err != nil {
				return nil, errors.Trace(err)
			}
			children = append(children, types.ProductItem{
				Repo: repo,
			})
		}
		ps = append(ps, types.ProductItem{
			Title:    title,
			Children: children,
		})
	}
	return ps, nil
}
