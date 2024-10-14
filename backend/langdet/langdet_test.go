package langdet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Detect(t *testing.T) {
	t.Parallel()
	assert.Equal(t, LangJapanese, Detect(`
    ※無言申請は原則お断りさせていただいております。 \
    -加入条件- 1.勝率50%以上。 \
    2.Tier10艦艇1隻以上所有。 \
    3.ランダム戦闘数1000戦以上。 \
    4.Discordの導入必須。聞き専可(コミュニケーションが取れることが必須)。 \
    5.日本語で意思疎通可能。 \
    6.ゲーム利用規約(EULA)を遵守。 \
    7.クランの平穏を乱さないこと。 \
    -除隊対象- 一定期間ログインが無い クランの名誉を害する行為
    `))
	assert.Equal(t, LangJapanese, Detect(`
    『-K2-』神風‐sではTyphoonリーグを目指すクランとなります。 \
    クラン加入については下記のDiscordの招待URLを通して面接申請をお願いします。 \
    もしDiscordをお持ちでない方は、ゲーム内チャットにてmyouko02もしくは、gaku0083、Orca_0313までご連絡してください。
    `))
}
