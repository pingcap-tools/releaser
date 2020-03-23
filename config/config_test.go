package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg := New()
	assert.Equal(t, cfg.Read("../config.example.toml"), nil, "read config")
	assert.Equal(t, cfg.Products, []Product{
		{
			Name:  "tidb",
			Repos: []string{"pingcap/tidb", "tikv/tikv", "pingcap/pd"},
		},
	}, "read config")
}
