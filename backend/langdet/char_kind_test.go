package langdet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()
	assert.Equal(t, charKindKana, newCharKind([]rune("あ")[0]))
	assert.Equal(t, charKindHangul, newCharKind([]rune("가")[0]))
	assert.Equal(t, charKindKanji, newCharKind([]rune("你")[0]))
	assert.Equal(t, charKindOther, newCharKind([]rune("a")[0]))
}
