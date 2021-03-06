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
	assert.Equal(t, testHasReleaseNote("Bug Fix \nRelease note\n- N/A"), hasReleaseNoteRes{"", false}, "case 2")
	assert.Equal(t, testHasReleaseNote("Bug Fix \nRelease note\n- (N/A)"), hasReleaseNoteRes{"", false}, "case 3")
	assert.Equal(t, testHasReleaseNote("### Release note <!-- bugfixes or new feature need a release note -->\n\n* No release note"), hasReleaseNoteRes{"", false}, "case 4")
	assert.Equal(t, testHasReleaseNote("### Release note <!-- bugfixes or new feature need a release note -->\n- No release note."), hasReleaseNoteRes{"", false}, "case 5")
	assert.Equal(t, testHasReleaseNote("### Release note <!-- bugfixes or new feature need a release note -->\n- `No release note`"), hasReleaseNoteRes{"", false}, "case 6")
	assert.Equal(t, testHasReleaseNote("### Release note <!-- bugfixes or new feature need a release note -->\n- `No release note`."), hasReleaseNoteRes{"", false}, "case 6")
}
