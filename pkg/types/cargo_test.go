package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	case1 := `[package]
name = "tikv"
version = "4.1.0-alpha"
authors = ["The TiKV Authors"]
description = "A distributed transactional key-value database powered by Rust and Raft"
license = "Apache-2.0"
keywords = ["KV", "distributed-systems", "raft"]
homepage = "https://tikv.org"
repository = "https://github.com/tikv/tikv/"
readme = "README.md"
edition = "2018"
publish = false

[dependencies]
async-stream = "0.2"
batch-system = { path = "components/batch-system", default-features = false }`

	cargo := NewCargo()
	err := cargo.Parse(case1)
	assert.Empty(t, err, "err must be nil")
	assert.Equal(t, cargo.Dependencies["async-stream"].(string), "0.2", "assert string type")
	assert.Equal(t, cargo.Dependencies["batch-system"].(map[string]interface{})["path"].(string), "components/batch-system", "assert map type, string field")
	assert.Equal(t, cargo.Dependencies["batch-system"].(map[string]interface{})["default-features"].(bool), false, "assert map type, false field")

	pkg := cargo.ToPackage()
	assert.Equal(t, len(pkg.Dependencies), 1, "pkg count")
	assert.Equal(t, pkg.Dependencies[0].Name, "async-stream", "pkg name")
	assert.Equal(t, pkg.Dependencies[0].Version, "0.2", "pkg version")
}
