package config

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
)

// Config is cherry picker config struct
type Config struct {
	GithubToken string   `toml:"github-token"`
	Repos       []string `toml:"repos"`
}

// New inits config by default
func New() *Config {
	return &Config{}
}

// Read from file
func (c *Config) Read(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Trace(err)
	}
	_, err = toml.Decode(string(file), c)
	return errors.Trace(err)
}
