package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasReleaseNote(t *testing.T) {
	noReleaseNoteSample1 := `Bug Fix
Release note
- NA`

	note1, has1 := hasReleaseNote(noReleaseNoteSample1)
	assert.Equal(t, note1, "", "case 1")
	assert.Equal(t, has1, false, "case 1")
}
