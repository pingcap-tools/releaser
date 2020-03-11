package types

import (
	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
)

// Cargo ...
type Cargo struct {
	Package      `toml:"package"`
	Dependencies `toml:"dependencies"`
}

// Package ...
type Package struct {
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

// Dependencies ...
type Dependencies map[string]interface{}

// NewCargo init empty Cargo
func NewCargo() *Cargo {
	return &Cargo{}
}

// Parse toml
func (c *Cargo) Parse(t string) error {
	_, err := toml.Decode(string(t), c)
	return errors.Trace(err)
}
