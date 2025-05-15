package persistence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAPI_SaveAndGetUserSetting(t *testing.T) {
	t.Parallel()

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
