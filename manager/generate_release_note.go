package manager

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/google/go-github/v29/github"
	"github.com/juju/errors"
	"github.com/ngaut/log"
	"github.com/you06/releaser/pkg/git"
	"github.com/you06/releaser/pkg/parser"
	"github.com/you06/releaser/pkg/types"
	"github.com/you06/releaser/pkg/utils"
)

func (m *Manager) runGenerateReleaseNote() error {
	if err := m.initRepo(); err != nil {
		return errors.Trace(err)
	}
	var errs []error

	for _, repo := range m.Repos {
		errs = append(errs, m.generateReleaseNoteRepo(repo))
	}

	for _, err := range errs {
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (m *Manager) generateReleaseNoteRepo(repo types.Repo) error {
	var errs []error
	if m.Opt.Version == "all" {
		milestones, err := m.PullCollector.ListAllOpenedMilestones(repo)
		if err != nil {
			return errors.Trace(err)
		}
		for _, milestone := range milestones {
			errs = append(errs, m.generateReleaseNoteRepoMilestone(repo, milestone))
		}
	} else {
		milestone, err := m.PullCollector.GetVersionMilestone(repo, m.Opt.Version)
		if err != nil {
			return errors.Trace(err)
		}
		errs = append(errs, m.generateReleaseNoteRepoMilestone(repo, milestone))
	}
	for _, err := range errs {
		if err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (m *Manager) generateReleaseNoteRepoMilestone(repo types.Repo, milestone *github.Milestone) error {
	gitClient := git.New(m.Config, &git.Config{
		Github: m.Github,
		User:   m.User,
		Base:   m.RelaseNoteRepo,
		Head:   types.Repo{Owner: m.User.GetLogin(), Repo: m.RelaseNoteRepo.Repo},
		Dir:    fmt.Sprintf("%s-%s", repo.Repo, milestone.GetTitle()),
	})
	// get release notes in PR
	_, pulls, err := m.PullCollector.ListAllMilestoneIssues(repo, milestone)
	if err != nil {
		return errors.Trace(err)
	}
	releaseNotes, err := m.NoteCollector.ListReleaseNote(repo, milestone.GetTitle())
	if err != nil {
		fmt.Printf("get release notes error %+v\n", err)
	}
	var defaultLangReleaseNote *parser.ReleaseNoteLang
	for _, releaseNote := range releaseNotes {
		if releaseNote.Lang == m.Config.PullLanguage {
			defaultLangReleaseNote = &releaseNote
		}
	}
	if defaultLangReleaseNote == nil {
		dir := strings.ReplaceAll(m.Config.ReleaseNotePath, "{repo}", repo.Repo)
		defaultLangReleaseNote = &parser.ReleaseNoteLang{
			Lang:    m.Config.PullLanguage,
			Path:    path.Join(dir, fmt.Sprintf("%s.md", milestone.GetTitle())),
			Version: milestone.GetTitle(),
		}
	}

	for _, pull := range pulls {
		note, has := hasReleaseNote(pull.GetBody())
		if has {
			inRepo := false
			for _, releaseNote := range defaultLangReleaseNote.Notes {
				if releaseNote.PullNumber == pull.GetNumber() {
					inRepo = true
					releaseNote.Note = note
				}
			}
			if !inRepo {
				defaultLangReleaseNote.Notes = append(defaultLangReleaseNote.Notes, parser.ReleaseNote{
					Repo:       repo,
					PullNumber: pull.GetNumber(),
					Note:       note,
				})
			}
		}
	}

	if err := gitClient.Clone(); err != nil {
		return errors.Trace(err)
	}

	defer func() {
		if err := gitClient.Clear(); err != nil {
			log.Error(err)
		}
	}()

	branch := fmt.Sprintf("%s-%s-%s", repo.Owner, repo.Repo, milestone.GetTitle())
	if err := gitClient.Checkout(branch); err != nil {
		if err := gitClient.CheckoutNew(branch); err != nil {
			return errors.Trace(err)
		}
	}

	log.Info(defaultLangReleaseNote.Path)
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

	title := fmt.Sprintf("update %s release notes", milestone.GetTitle())
	if _, err := gitClient.CreatePull(title, branch); err != nil {
		return errors.Trace(err)
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
