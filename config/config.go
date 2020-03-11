package config

import (
	"io"
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/juju/errors"
	"github.com/ngaut/log"
)

// Config is cherry picker config struct
type Config struct {
	GithubToken     string   `toml:"github-token"`
	Repos           []string `toml:"repos"`
	ReleaseNoteRepo string   `toml:"release-note-repo"`
	ReleaseNotePath string   `toml:"release-note-path"`
	PullLanguage    string   `toml:"pull-language"`
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

// Print Config
func (c *Config) Print(writer ...io.Writer) {
	if len(writer) == 0 {
		log.Infof("%+v\n", *c)
	}
}
