package service

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	filename   = "example.png"
	base64Data = "abc123"
)

type mockScreenshotRepo struct {
	mock.Mock
}

func (m *mockScreenshotRepo) Save(path string, base64Data string) error {
	args := m.Called(path, base64Data)
	return args.Error(0)
}

func TestScreenshot_SaveForAuto_Success(t *testing.T) {
	t.Parallel()

	// 期待されるメソッド呼び出しと戻り値の設定
	screenshotPath := filepath.Join("screenshot", filename)
	mockRepo := &mockScreenshotRepo{}
	mockRepo.On("Save", screenshotPath, base64Data).Return(nil)

	// Screenshot インスタンスの作成
	s := NewScreenshot(mockRepo, func(_ context.Context, _ runtime.SaveDialogOptions) (string, error) {
		return screenshotPath, nil
	})

	// テスト実行
	err := s.SaveForAuto(filename, base64Data)

	// 結果の検証
	assert.NoError(t, err)
}

func TestScreenshot_SaveWithDialog_Success(t *testing.T) {
	t.Parallel()

	// 期待されるメソッド呼び出しと戻り値の設定
	screenshotPath := filepath.Join("directory", filename)
	mockRepo := &mockScreenshotRepo{}
	mockRepo.On("Save", screenshotPath, base64Data).Return(nil)

	// Screenshot インスタンスの作成
	s := NewScreenshot(mockRepo, func(_ context.Context, _ runtime.SaveDialogOptions) (string, error) {
		return screenshotPath, nil
	})

	// テスト実行
	err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	assert.NoError(t, err)
}

func TestScreenshot_SaveWithDialog_Failure_Error(t *testing.T) {
	t.Parallel()

	// Screenshot インスタンスの作成
	mockRepo := &mockScreenshotRepo{}
	s := NewScreenshot(mockRepo, func(_ context.Context, _ runtime.SaveDialogOptions) (string, error) {
		//nolint:goerr113
		return "", fmt.Errorf("file dialog error")
	})

	// テスト実行
	err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	assert.EqualError(t, err, "S400 SaveDialog file dialog error")
	mockRepo.AssertNotCalled(t, "Save")
}

func TestScreenshot_SaveWithDialog_Failure_Cancel(t *testing.T) {
	t.Parallel()

	// Screenshot インスタンスの作成
	mockRepo := &mockScreenshotRepo{}
	s := NewScreenshot(mockRepo, func(_ context.Context, _ runtime.SaveDialogOptions) (string, error) {
		return "", nil
	})

	// テスト実行
	err := s.SaveWithDialog(context.Background(), filename, base64Data)

	// 結果の検証
	assert.EqualError(t, err, "S401 Canceled")
	mockRepo.AssertNotCalled(t, "Save")
}
