package langdet

import (
	"regexp"
)

func Detect(text string) Lang {
	charKinds := charKinds(text)
	kana := charKinds[charKindKana]
	hangul := charKinds[charKindHangul]
	kanji := charKinds[charKindKanji]

	if hangul > (kana + kanji) {
		return LangKorea
	}

	kanaRate := float64(kana) / (float64(kana) + float64(kanji))
	if kanaRate > 0.2 {
		return LangJapanese
	}

	ewc := englishWordCount(text)
	if kanji > ewc {
		return LangChinese
	}

	return LangOther
}

func charKinds(text string) map[charKind]uint {
	dist := make(map[charKind]uint)
	for _, r := range text {
		ck := newCharKind(r)
		dist[ck] += 1
	}

	return dist
}

func englishWordCount(text string) uint {
	wordPattern := `[a-zA-Z]+`
	re := regexp.MustCompile(wordPattern)
	words := re.FindAllString(text, -1)

	return uint(len(words))
}
