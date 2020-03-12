package dependency

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	case1 := `module github.com/pingcap/tidb

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/blacktear23/go-proxyprotocol v0.0.0-20180807104634-af7a81e8dd0d
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
)`

	p := parseGoMod(case1)
	assert.Equal(t, p.Name, "github.com/pingcap/tidb", "module name parse")
	assert.Equal(t, len(p.Dependencies), 3, "dependencies count")
	assert.Equal(t, p.Dependencies[0].Name, "github.com/BurntSushi/toml", "dependency name")
	assert.Equal(t, p.Dependencies[1].Name, "github.com/blacktear23/go-proxyprotocol", "dependency name")
	assert.Equal(t, p.Dependencies[2].Name, "github.com/codahale/hdrhistogram", "dependency name")
	assert.Equal(t, p.Dependencies[0].Version, "v0.3.1", "dependency version")
	assert.Equal(t, p.Dependencies[1].Version, "v0.0.0-20180807104634-af7a81e8dd0d", "dependency version")
	assert.Equal(t, p.Dependencies[2].Version, "v0.0.0-20161010025455-3a0bb77429bd", "dependency version")
}
