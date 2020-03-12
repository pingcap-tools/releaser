package types

import (
	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
)

// Cargo ...
type Cargo struct {
	Package      CargoPackage      `toml:"package"`
	Dependencies CargoDependencies `toml:"dependencies"`
}

// CargoPackage ...
type CargoPackage struct {
	Name        string   `toml:"name"`
	Version     string   `toml:"version"`
	Authors     []string `toml:"authors"`
	Description string   `toml:"description"`
	License     string   `toml:"license"`
	Keywords    []string `toml:"keywords"`
	Homepage    string   `toml:"homepage"`
	Repository  string   `toml:"repository"`
	Readme      string   `toml:"readme"`
	Edition     string   `toml:"edition"`
	Publish     bool     `toml:"publish"`
}

// CargoDependencies ...
type CargoDependencies map[string]interface{}

// NewCargo init empty Cargo
func NewCargo() *Cargo {
	return &Cargo{}
}

// Parse toml
func (c *Cargo) Parse(t string) error {
	_, err := toml.Decode(string(t), c)
	return errors.Trace(err)
}

// ToPackage transfer Cargo to common Package
func (c *Cargo) ToPackage() *Package {
	var dependencies []Dependency

	for name, dependency := range c.Dependencies {
		d := Dependency{
			Name: name,
		}
		switch val := dependency.(type) {
		case string:
			d.Version = val
		}
		if d.Version != "" {
			dependencies = append(dependencies, d)
		}
	}

	return &Package{
		Name:         c.Package.Name,
		Dependencies: dependencies,
	}
}
