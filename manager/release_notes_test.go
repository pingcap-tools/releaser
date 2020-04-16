package manager

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type hasReleaseNoteRes struct {
	note string
	has  bool
}

func testHasReleaseNote(raw string) hasReleaseNoteRes {
	note, has := hasReleaseNote(raw)
	return hasReleaseNoteRes{note, has}
}

func TestHasReleaseNote(t *testing.T) {
	assert.Equal(t, testHasReleaseNote("Bug Fix \nRelease note\n- NA"), hasReleaseNoteRes{"", false}, "case 1")
	assert.Equal(t, testHasReleaseNote("Bug Fix \nRelease note\n- N/A"), hasReleaseNoteRes{"", false}, "case 1")
	assert.Equal(t, testHasReleaseNote("Bug Fix \nRelease note\n- (N/A)"), hasReleaseNoteRes{"", false}, "case 1")
}
