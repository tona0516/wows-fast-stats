package service

import (
	"context"
	"path/filepath"
	"testing"
	"wfs/backend/apperr"
	"wfs/backend/domain/mock/repository"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"go.uber.org/mock/gomock"
)

const (
	filename   = "example.png"
	base64Data = "abc123"
)

func TestScreenshot_SaveForAuto(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		screenshotPath := filepath.Join("screenshot", filename)
		mockLocalFile := repository.NewMockLocalFile(ctrl)
		mockLocalFile.EXPECT().SaveScreenshot(screenshotPath, base64Data).Return(nil).AnyTimes()

		s := NewScreenshot(mockLocalFile)
		s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
			return screenshotPath, nil
		}

		// テスト
		err := s.SaveForAuto(filename, base64Data)

		// アサーション
		assert.NoError(t, err)
	})
}

func TestScreenshot_SaveWithDialog(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		// 準備
		screenshotPath := filepath.Join("directory", filename)
		mockLocalFile := repository.NewMockLocalFile(ctrl)
		mockLocalFile.EXPECT().SaveScreenshot(screenshotPath, base64Data).Return(nil).AnyTimes()

		s := NewScreenshot(mockLocalFile)
		s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
			return screenshotPath, nil
		}

		// テスト
		saved, err := s.SaveWithDialog(t.Context(), filename, base64Data)

		// アサーション
		assert.True(t, saved)
		assert.NoError(t, err)
	})

	t.Run("異常系", func(t *testing.T) {
		t.Parallel()
		// 準備
		s := NewScreenshot(nil)
		s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
			return "", failure.New(apperr.WailsError)
		}

		// テスト
		saved, err := s.SaveWithDialog(t.Context(), filename, base64Data)

		// アサーション
		assert.False(t, saved)
		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, apperr.WailsError, code)
	})

	t.Run("異常系_キャンセル", func(t *testing.T) {
		t.Parallel()
		// 準備
		s := NewScreenshot(nil)
		s.SaveFileDialog = func(ctx context.Context, dialogOptions runtime.SaveDialogOptions) (string, error) {
			return "", nil
		}

		// テスト
		saved, err := s.SaveWithDialog(t.Context(), filename, base64Data)

		// アサーション
		assert.False(t, saved)
		assert.NoError(t, err)
	})
}
