package usecase

import (
	"context"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/mocks"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	filename   = "example.png"
	base64Data = "abc123"
)

func TestScreenshot_SaveForAuto_正常系(t *testing.T) {
	t.Parallel()

	// 期待されるメソッド呼び出しと戻り値の設定
	screenshotPath := filepath.Join("screenshot", filename)
	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("SaveScreenshot", screenshotPath, base64Data).Return(nil)

	// Screenshot インスタンスの作成
	s := NewScreenshot(mockLocalFile)
	s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
		return screenshotPath, nil
	}

	// テスト実行
	err := s.SaveForAuto(filename, base64Data)

	// 結果の検証
	require.NoError(t, err)
}

func TestScreenshot_SaveWithDialog_正常系(t *testing.T) {
	t.Parallel()

	// 期待されるメソッド呼び出しと戻り値の設定
	screenshotPath := filepath.Join("directory", filename)
	mockLocalFile := &mocks.LocalFileInterface{}
	mockLocalFile.On("SaveScreenshot", screenshotPath, base64Data).Return(nil)

	// Screenshot インスタンスの作成
	s := NewScreenshot(mockLocalFile)
	s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
		return screenshotPath, nil
	}

	// テスト実行
	saved, err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	require.True(t, saved)
	require.NoError(t, err)
}

func TestScreenshot_SaveWithDialog_異常系(t *testing.T) {
	t.Parallel()

	// Screenshot インスタンスの作成
	mockLocalFile := &mocks.LocalFileInterface{}
	s := NewScreenshot(mockLocalFile)
	s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
		return "", failure.New(apperr.WailsError)
	}

	// テスト実行
	saved, err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	assert.False(t, saved)
	code, ok := failure.CodeOf(err)
	assert.True(t, ok)
	assert.Equal(t, apperr.WailsError, code)
	mockLocalFile.AssertNotCalled(t, "SaveScreenshot")
}

func TestScreenshot_SaveWithDialog_異常系_キャンセル(t *testing.T) {
	t.Parallel()

	// Screenshot インスタンスの作成
	mockLocalFile := &mocks.LocalFileInterface{}
	s := NewScreenshot(mockLocalFile)
	s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
		return "", nil
	}

	// テスト実行
	saved, err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	assert.False(t, saved)
	require.NoError(t, err)
	mockLocalFile.AssertNotCalled(t, "SaveScreenshot")
}
