package service

import (
	"context"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	filename   = "example.png"
	base64Data = "abc123"
)

var errShowDialog = errors.New("file dialog error")

func TestScreenshot_SaveForAuto_正常系(t *testing.T) {
	t.Parallel()

	// 期待されるメソッド呼び出しと戻り値の設定
	screenshotPath := filepath.Join("screenshot", filename)
	mockLocalFile := &mockLocalFile{}
	mockLocalFile.On("SaveScreenshot", screenshotPath, base64Data).Return(nil)

	// Screenshot インスタンスの作成
	s := NewScreenshot(mockLocalFile, func(_ context.Context, _ runtime.SaveDialogOptions) (string, error) {
		return screenshotPath, nil
	})

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
	s := NewScreenshot(mockLocalFile, func(_ context.Context, _ runtime.SaveDialogOptions) (string, error) {
		return screenshotPath, nil
	})

	// テスト実行
	err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	assert.NoError(t, err)
}

func TestScreenshot_SaveWithDialog_異常系(t *testing.T) {
	t.Parallel()

	// Screenshot インスタンスの作成
	mockLocalFile := &mockLocalFile{}
	s := NewScreenshot(mockLocalFile, func(_ context.Context, _ runtime.SaveDialogOptions) (string, error) {
		return "", errShowDialog
	})

	// テスト実行
	err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	assert.EqualError(t, err, apperr.New(apperr.ShowDialog, errShowDialog).Error())
	mockLocalFile.AssertNotCalled(t, "SaveScreenshot")
}

func TestScreenshot_SaveWithDialog_異常系_キャンセル(t *testing.T) {
	t.Parallel()

	// Screenshot インスタンスの作成
	mockLocalFile := &mockLocalFile{}
	s := NewScreenshot(mockLocalFile, func(_ context.Context, _ runtime.SaveDialogOptions) (string, error) {
		return "", nil
	})

	// テスト実行
	err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	assert.EqualError(t, err, apperr.UserCanceled.String())
	mockLocalFile.AssertNotCalled(t, "SaveScreenshot")
}
