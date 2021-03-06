package manager

import (
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/google/go-github/v30/github"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/releaser/pkg/git"
	"github.com/you06/releaser/pkg/parser"
	"github.com/you06/releaser/pkg/types"
	"github.com/you06/releaser/pkg/utils"
)

var (
	numPattern = regexp.MustCompile(`^.*?(\d+)$`)
)

func (m *Manager) runGenerateReleaseNote() error {
	if err := m.initRepo(); err != nil {
		return errors.Trace(err)
	}
	var errs []error

	for _, project := range m.Products {
		errs = append(errs, m.generateReleaseNoteProduct(project))
	}

	for _, err := range errs {
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (m *Manager) generateReleaseNoteProduct(product types.Product) error {
	// Do not process empty product
	if len(product.Repos) == 0 {
		return nil
	}

	var (
		errs       []error
		milestones []*github.Milestone
	)

	// get target milestones
	if m.Opt.Version == "all" {
		// get milestones by first product
		var err error
		milestones, err = m.PullCollector.ListAllOpenedMilestones(product.Repos[0])
		if err != nil {
			return errors.Trace(err)
		}
	} else {
		milestone, err := m.PullCollector.GetVersionMilestone(product.Repos[0], m.Opt.Version)
		if err != nil {
			return errors.Trace(err)
		}
		milestones = append(milestones, milestone)
	}

	for _, milestone := range milestones {
		errs = append(errs, m.generateReleaseNoteProductMilestone(product, milestone))
	}

	for _, err := range errs {
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (m *Manager) generateReleaseNoteProductMilestone(product types.Product, milestone *github.Milestone) error {
	releaseNotes, err := m.NoteCollector.ListReleaseNote(product, milestone.GetTitle())
	if err != nil {
		return errors.Errorf("get release notes error %+v\n", err)
	}
	var defaultLangReleaseNote *parser.ReleaseNoteLang
	for _, releaseNote := range releaseNotes {
		if releaseNote.Lang == m.Config.PullLanguage {
			defaultLangReleaseNote = &releaseNote
			// break because the meomory address will be replaced by next item
			// the pointer will indicate to the wrong item
			break
		}
	}

	if defaultLangReleaseNote == nil {
		dir := strings.ReplaceAll(m.Config.ReleaseNotePath, "{product}", product.Name)
		defaultLangReleaseNote = &parser.ReleaseNoteLang{
			Name:               product.Name,
			Lang:               m.Config.PullLanguage,
			Path:               path.Join(dir, fmt.Sprintf("%s.md", strings.TrimLeft(m.Opt.Version, "v"))),
			Version:            milestone.GetTitle(),
			ReleaseNoteClasses: make(map[string][]parser.RepoReleaseNotes),
			Structure:          product.Structure,
		}
	}

	for _, repo := range product.Repos {
		rename, ok := product.Renames[repo]
		if !ok {
			rename = repo
		}
		if err := m.makeReleaseNoteRepoMilestone(product, repo, rename, defaultLangReleaseNote); err != nil {
			return errors.Trace(err)
		}
	}

	gitClient := git.New(m.Config, &git.Config{
		Github: m.Github,
		User:   m.User,
		Base:   m.RelaseNoteRepo,
		Head:   types.Repo{Owner: m.User.GetLogin(), Repo: m.RelaseNoteRepo.Repo},
		Dir:    fmt.Sprintf("%s-%s", product.Name, milestone.GetTitle()),
	})
	if err := gitClient.Clone(); err != nil {
		return errors.Trace(err)
	}
	defer func() {
		if err := gitClient.Clear(); err != nil {
			log.Error(err)
		}
	}()

	branch := fmt.Sprintf("%s-%s", "update", strings.TrimLeft(m.Opt.Version, "v"))
	if err := gitClient.Checkout(branch); err != nil {
		if err := gitClient.CheckoutNew(branch); err != nil {
			return errors.Trace(err)
		}
	}

	if err := gitClient.WriteFileContent(defaultLangReleaseNote.Path,
		defaultLangReleaseNote.String()); err != nil {
		return errors.Trace(err)
	}

	commitMessage := fmt.Sprintf("update %s release notes at %s", milestone.GetTitle(), now())
	if err := gitClient.Commit(commitMessage); err != nil {
		return errors.Trace(err)
	}

	if err := gitClient.Push(branch); err != nil {
		return errors.Trace(err)
	}

	title := fmt.Sprintf("update %s %s release notes", product.Name, milestone.GetTitle())
	if _, err := gitClient.CreatePull(title, branch); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (m *Manager) makeReleaseNoteRepoMilestone(product types.Product, repo, rename types.Repo, releaseNote *parser.ReleaseNoteLang) error {
	if releaseNote == nil {
		return errors.New("releaseNote cannot be nil")
	}

	milestone, err := m.PullCollector.GetVersionMilestone(repo, m.Opt.Version)
	if err != nil {
		fmt.Printf("Find milestone in %s failed", repo)
		return nil
	}

	// get release notes in PR
	_, pulls, err := m.PullCollector.ListAllMilestoneContents(repo, milestone)
	if err != nil {
		return errors.Trace(err)
	}

	ref := version2ref(m.Opt.Version)

	for _, pull := range pulls {
		if pull.GetBase().GetRef() != ref {
			continue
		}
		if !pull.GetMerged() {
			continue
		}
		note, has := hasReleaseNote(pull.GetBody())
		if has {
			releaseNoteType := getReleaseNoteType(pull, product)
			if _, ok := releaseNote.ReleaseNoteClasses[releaseNoteType]; !ok {
				releaseNote.ReleaseNoteClasses[releaseNoteType] = make([]parser.RepoReleaseNotes, 1)
			}

			var repoReleaseNote *parser.RepoReleaseNotes
			for i := range releaseNote.ReleaseNoteClasses[releaseNoteType] {
				if releaseNote.ReleaseNoteClasses[releaseNoteType][i].Repo.Repo == repo.Repo {
					repoReleaseNote = &releaseNote.ReleaseNoteClasses[releaseNoteType][i]
					break
				}
			}
			if repoReleaseNote == nil {
				releaseNote.ReleaseNoteClasses[releaseNoteType] = append(releaseNote.ReleaseNoteClasses[releaseNoteType], parser.RepoReleaseNotes{
					Repo:   repo,
					Rename: rename,
				})
				repoReleaseNote = &releaseNote.ReleaseNoteClasses[releaseNoteType][len(releaseNote.ReleaseNoteClasses[releaseNoteType])-1]
			}

			inRepo := false
			for _, releaseNote := range repoReleaseNote.Notes {
				if releaseNote.PullNumber == pull.GetNumber() {
					inRepo = true
					// update note
					releaseNote.Note = note
				}
			}
			if !inRepo {
				repoReleaseNote.Notes = append(repoReleaseNote.Notes, parser.ReleaseNote{
					Repo:       repo,
					PullNumber: pull.GetNumber(),
					Note:       note,
				})
			}
		}
	}

	return nil
}

func (m *Manager) initRepo() error {
	_, err := m.getRepo(types.Repo{Owner: m.User.GetLogin(), Repo: m.RelaseNoteRepo.Repo})
	if err != nil {
		if strings.Contains(err.Error(), "Not Found") {
			if err := m.forkRepo(m.RelaseNoteRepo); err != nil {
				return errors.Trace(err)
			}
		}
	}
	return nil
}

func (m *Manager) getRepo(repo types.Repo) (*github.Repository, error) {
	ctx, _ := utils.NewTimeoutContext()
	r, _, err := m.Github.Repositories.Get(ctx, repo.Owner, repo.Repo)
	return r, errors.Trace(err)
}

func (m *Manager) forkRepo(repo types.Repo) error {
	ctx, _ := utils.NewTimeoutContext()
	_, _, err := m.Github.Repositories.CreateFork(ctx, repo.Owner, repo.Repo, &github.RepositoryCreateForkOptions{})
	return errors.Trace(err)
}

func now() string {
	return time.Now().Format("2006-01-02T15:04:05")
}

func getReleaseNoteType(pull *github.PullRequest, product types.Product) string {
	for _, label := range pull.Labels {
		if class, ok := product.Label2Type[label.GetName()]; ok {
			return class
		}
	}
	return "Others"
}

func version2ref(version string) string {
	p1 := strings.Split(version, ".")[0]
	num := numPattern.FindStringSubmatch(p1)
	if len(num) == 2 {
		return fmt.Sprintf("release-%s.0", num[1])
	}
	return ""
}
