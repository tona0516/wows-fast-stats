package persistence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_SaveAndGetUserSetting(t *testing.T) {
	tmpDir := t.TempDir()
	api := NewAPI(func() (string, error) {
		return tmpDir, nil
	})

	us := UserSetting{}
	us.Version = 1
	us.Required.InstallPath = "/some/path"
	us.Optional.FontSize = "large"
	us.Optional.SendReport = true

	// 保存テスト
	err := api.SaveUserSetting(us)
	assert.NoError(t, err)

	// 読み込みテスト
	loaded, err := api.GetUserSetting()
	assert.NoError(t, err)
	assert.Equal(t, us.Version, loaded.Version)
	assert.Equal(t, us.Required.InstallPath, loaded.Required.InstallPath)
	assert.Equal(t, us.Optional.FontSize, loaded.Optional.FontSize)
	assert.Equal(t, us.Optional.SendReport, loaded.Optional.SendReport)

	// ファイルの中身も検証
	path := filepath.Join(tmpDir, "wows-fast-stats", "user_setting.json")
	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	var fileContent UserSetting
	err = json.Unmarshal(data, &fileContent)
	assert.NoError(t, err)
	assert.Equal(t, us.Required.InstallPath, fileContent.Required.InstallPath)
}

func TestAPI_SaveAndGetAlertPlayers(t *testing.T) {
	tmpDir := t.TempDir()
	api := NewAPI(func() (string, error) {
		return tmpDir, nil
	})

	ap := AlertPlayers{
		{
			AccountID:   123,
			AccountName: "PlayerOne",
			Icon:        "bi-check-circle-fill",
			Message:     "Watch out!",
		},
		{
			AccountID:   456,
			AccountName: "PlayerTwo",
			Icon:        "bi-exclamation-triangle-fill",
			Message:     "Dangerous player",
		},
	}

	// 保存テスト
	err := api.SaveAlertPlayers(ap)
	assert.NoError(t, err)

	// 読み込みテスト
	loaded, err := api.GetAlertPlayers()
	assert.NoError(t, err)
	assert.Len(t, loaded, 2)
	assert.Equal(t, ap[0].AccountID, loaded[0].AccountID)
	assert.Equal(t, ap[0].AccountName, loaded[0].AccountName)
	assert.Equal(t, ap[1].Icon, loaded[1].Icon)
	assert.Equal(t, ap[1].Message, loaded[1].Message)

	// ファイルの中身も検証
	path := filepath.Join(tmpDir, "wows-fast-stats", "alert_player.json")
	data, err := os.ReadFile(path)
	assert.NoError(t, err)
	var fileContent AlertPlayers
	err = json.Unmarshal(data, &fileContent)
	assert.NoError(t, err)
	assert.Equal(t, ap[0].AccountName, fileContent[0].AccountName)
}
