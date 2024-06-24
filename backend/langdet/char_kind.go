package langdet

import "unicode"

type charKind string

const (
	charKindKana   charKind = "kana"
	charKindHangul charKind = "hangul"
	charKindKanji  charKind = "kanji"
	charKindOther  charKind = "other"
)

func newCharKind(r rune) charKind {
	// ひらがな/カタカナ
	if unicode.In(r, unicode.Hiragana, unicode.Katakana) {
		return charKindKana
	}

	// ハングル
	if unicode.Is(unicode.Hangul, r) {
		return charKindHangul
	}

	// 漢字
	if unicode.Is(unicode.Han, r) {
		return charKindKanji
	}

	return charKindOther
}
