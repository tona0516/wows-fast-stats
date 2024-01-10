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

func TestScreenshot_SaveForAuto(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		screenshotPath := filepath.Join("screenshot", filename)
		mockLocalFile := &mocks.LocalFileInterface{}
		mockLocalFile.On("SaveScreenshot", screenshotPath, base64Data).Return(nil)

		s := NewScreenshot(mockLocalFile, nil)
		s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
			return screenshotPath, nil
		}

		// テスト
		err := s.SaveForAuto(filename, base64Data)

		// アサーション
		require.NoError(t, err)
		mockLocalFile.AssertExpectations(t)
	})
}

func TestScreenshot_SaveWithDialog(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()
		// 準備
		screenshotPath := filepath.Join("directory", filename)
		mockLocalFile := &mocks.LocalFileInterface{}
		mockLocalFile.On("SaveScreenshot", screenshotPath, base64Data).Return(nil)

		s := NewScreenshot(mockLocalFile, nil)
		s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
			return screenshotPath, nil
		}

		// テスト
		saved, err := s.SaveWithDialog(context.Background(), filename, base64Data)

		// アサーション
		require.True(t, saved)
		require.NoError(t, err)
		mockLocalFile.AssertExpectations(t)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()
		// 準備
		s := NewScreenshot(nil, nil)
		s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
			return "", failure.New(apperr.WailsError)
		}

		// テスト
		saved, err := s.SaveWithDialog(context.Background(), filename, base64Data)

		// アサーション
		assert.False(t, saved)
		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, apperr.WailsError, code)
	})

	t.Run("異常系_キャンセル", func(t *testing.T) {
		t.Parallel()
		// 準備
		s := NewScreenshot(nil, nil)
		s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
			return "", nil
		}

		// テスト
		saved, err := s.SaveWithDialog(context.Background(), filename, base64Data)

		// アサーション
		assert.False(t, saved)
		require.NoError(t, err)
	})
}
