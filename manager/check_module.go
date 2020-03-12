package manager

import (
	"github.com/juju/errors"
	"github.com/ngaut/log"
)

func (m *Manager) runCheckModule() error {
	ds, err := m.DependencyCollector.GetDependencies(m.Repos[0], m.Opt.Version)
	if err != nil {
		return errors.Trace(err)
	}
	for _, di := range ds {
		log.Infof("%s-%s\n", di.Repo, di.Name)
		for _, d := range di.Dependencies {
			log.Info(d.Name, d.Version)
		}
	}
	return nil
}
