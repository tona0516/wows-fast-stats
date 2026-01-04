package infra

import (
	"errors"
	"testing"
	"wfs/internal/mock"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLogger_Debug(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		message  string
		contexts map[string]string
	}{
		{
			name:     "正常系_メッセージのみでDebugログ出力",
			message:  "debug message",
			contexts: nil,
		},
		{
			name:    "正常系_コンテキスト付きでDebugログ出力",
			message: "debug with context",
			contexts: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:     "正常系_空のコンテキストでDebugログ出力",
			message:  "debug with empty context",
			contexts: map[string]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			instance := NewLogger(
				"test-app",
				"1.0.0",
				t.TempDir(),
				zerolog.DebugLevel,
				nil,
				nil,
			)

			assert.NotPanics(t, func() {
				instance.Debug(tc.message, tc.contexts)
			})
		})
	}
}

func TestLogger_Info(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		message  string
		contexts map[string]string
	}{
		{
			name:     "正常系_メッセージのみでInfoログ出力",
			message:  "info message",
			contexts: nil,
		},
		{
			name:    "正常系_コンテキスト付きでInfoログ出力",
			message: "info with context",
			contexts: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:     "正常系_空のコンテキストでInfoログ出力",
			message:  "info with empty context",
			contexts: map[string]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockInfoDiscord := mock.NewMockDiscordApiClient(ctrl)
			mockInfoDiscord.EXPECT().Comment(gomock.Any()).Return(nil)
			instance := NewLogger(
				"test-app",
				"1.0.0",
				t.TempDir(),
				zerolog.InfoLevel,
				nil,
				mockInfoDiscord,
			)

			assert.NotPanics(t, func() {
				instance.Info(tc.message, tc.contexts)
			})
		})
	}
}

func TestLogger_Info_DicordError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockInfoDiscord := mock.NewMockDiscordApiClient(ctrl)
	mockInfoDiscord.EXPECT().Comment(gomock.Any()).Return(errors.New("discord error"))
	instance := NewLogger(
		"test-app",
		"1.0.0",
		t.TempDir(),
		zerolog.InfoLevel,
		nil,
		mockInfoDiscord,
	)

	assert.NotPanics(t, func() {
		instance.Info("info message", nil)
	})
}

func TestLogger_Error(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		err      error
		contexts map[string]string
	}{
		{
			name:     "正常系_エラーのみでErrorログ出力",
			err:      errors.New("error occurred"),
			contexts: nil,
		},
		{
			name: "正常系_コンテキスト付きでErrorログ出力",
			err:  errors.New("api error"),
			contexts: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:     "正常系_空のコンテキストでErrorログ出力",
			err:      errors.New("general error"),
			contexts: map[string]string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			mockAlertDiscord := mock.NewMockDiscordApiClient(ctrl)
			mockAlertDiscord.EXPECT().Comment(gomock.Any()).Return(nil).AnyTimes()
			instance := NewLogger(
				"test-app",
				"1.0.0",
				t.TempDir(),
				zerolog.ErrorLevel,
				mockAlertDiscord,
				nil,
			)

			assert.NotPanics(t, func() {
				instance.Error(tc.err, tc.contexts)
			})
		})
	}
}

func TestLogger_Error_DicordError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockAlertDiscord := mock.NewMockDiscordApiClient(ctrl)
	mockAlertDiscord.EXPECT().Comment(gomock.Any()).Return(errors.New("discord error"))
	instance := NewLogger(
		"test-app",
		"1.0.0",
		t.TempDir(),
		zerolog.InfoLevel,
		mockAlertDiscord,
		nil,
	)

	assert.NotPanics(t, func() {
		instance.Error(errors.New("some error"), nil)
	})
}
