package manager

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/you06/releaser/pkg/types"
)

type packageVersion struct {
	version string
	repo    types.Repo
}

func (m *Manager) runCheckModule() error {
	var (
		packages   []*types.Package
		versionMap = make(map[string]packageVersion)
	)

	for _, repo := range m.Repos {
		batch, err := m.DependencyCollector.GetDependencies(repo, m.Opt.Version)
		if err != nil {
			return errors.Trace(err)
		}
		packages = append(packages, batch...)
	}

	for _, p := range packages {
		fmt.Printf("%s %s: %s\n", p.Repo, p.Type, p.URL)
	}
	fmt.Println("-----------------------")

	for _, p := range packages {
		for _, dependency := range p.Dependencies {
			existVersion, ok := versionMap[dependency.Name]
			if !ok {
				versionMap[dependency.Name] = packageVersion{
					version: dependency.Version,
					repo:    p.Repo,
				}
			} else {
				if existVersion.version != dependency.Version {
					fmt.Printf("%s | %s: %s, %s: %s\n", dependency.Name,
						existVersion.repo, existVersion.version, p.Repo, dependency.Version)
				}
			}
		}
	}

	return nil
}
