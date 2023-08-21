package service

import (
	"context"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
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
	mockLocalFile := &mockLocalFile{}
	mockLocalFile.On("SaveScreenshot", screenshotPath, base64Data).Return(nil)

	// Screenshot インスタンスの作成
	s := NewScreenshot(mockLocalFile)
	s.SaveFileDialogFunc = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
		return screenshotPath, nil
	}

	// テスト実行
	err := s.SaveForAuto(filename, base64Data)

	// 結果の検証
	assert.NoError(t, err)
}

func TestScreenshot_SaveWithDialog_正常系(t *testing.T) {
	t.Parallel()

	// 期待されるメソッド呼び出しと戻り値の設定
	screenshotPath := filepath.Join("directory", filename)
	mockLocalFile := &mockLocalFile{}
	mockLocalFile.On("SaveScreenshot", screenshotPath, base64Data).Return(nil)

	// Screenshot インスタンスの作成
	s := NewScreenshot(mockLocalFile)
	s.SaveFileDialogFunc = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
		return screenshotPath, nil
	}

	// テスト実行
	err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	assert.NoError(t, err)
}

func TestScreenshot_SaveWithDialog_異常系(t *testing.T) {
	t.Parallel()

	// Screenshot インスタンスの作成
	mockLocalFile := &mockLocalFile{}
	s := NewScreenshot(mockLocalFile)
	s.SaveFileDialogFunc = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
		return "", failure.New(apperr.WailsError)
	}

	// テスト実行
	err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	code, ok := failure.CodeOf(err)
	assert.True(t, ok)
	assert.Equal(t, code, apperr.WailsError)
	mockLocalFile.AssertNotCalled(t, "SaveScreenshot")
}

func TestScreenshot_SaveWithDialog_異常系_キャンセル(t *testing.T) {
	t.Parallel()

	// Screenshot インスタンスの作成
	mockLocalFile := &mockLocalFile{}
	s := NewScreenshot(mockLocalFile)
	s.SaveFileDialogFunc = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
		return "", nil
	}

	// テスト実行
	err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	code, ok := failure.CodeOf(err)
	assert.True(t, ok)
	assert.Equal(t, code, apperr.UserCanceled)
	mockLocalFile.AssertNotCalled(t, "SaveScreenshot")
}
