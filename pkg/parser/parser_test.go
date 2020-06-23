package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigLetter(t *testing.T) {
	assert.Equal(t, Ucfirst(""), "")
	assert.Equal(t, Ucfirst("``"), "``")
	assert.Equal(t, Ucfirst("()"), "()")
	assert.Equal(t, Ucfirst("[()]"), "[()]")
	assert.Equal(t, Ucfirst("{()}"), "{()}")
	assert.Equal(t, Ucfirst("🐸"), "🐸")
	assert.Equal(t, Ucfirst("å"), "å")
	assert.Equal(t, Ucfirst("虵"), "虵")
	assert.Equal(t, Ucfirst("hA"), "HA")
}
