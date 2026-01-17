package infra

import (
	"errors"
	"testing"
	"wfs/backend/adapter"
	"wfs/backend/config"
	"wfs/backend/core"
	"wfs/backend/mock"

	"github.com/morikuni/failure"
	"github.com/rs/zerolog"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

			injector := do.New()
			do.ProvideValue(injector, config.Config{
				Basic: config.BasicConfig{
					Name: "test-app",
				},
				Logger: config.LoggerConfig{
					Level: zerolog.DebugLevel,
				},
				LocalFile: config.LocalFileConfig{
					RootDir: t.TempDir(),
				},
			})
			do.ProvideNamed(injector, "alert-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
				return (*mock.MockDiscordClient)(nil), nil
			})
			do.ProvideNamed(injector, "info-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
				return (*mock.MockDiscordClient)(nil), nil
			})
			instance, err := NewLogger(injector)
			require.NoError(t, err)

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
			mockInfoDiscord := mock.NewMockDiscordClient(ctrl)
			mockInfoDiscord.EXPECT().Comment(gomock.Any(), gomock.Any()).Return(nil).Times(1)

			injector := do.New()
			do.ProvideValue(injector, config.Config{
				Basic: config.BasicConfig{
					Name: "test-app",
				},
				Logger: config.LoggerConfig{
					Level: zerolog.InfoLevel,
				},
				LocalFile: config.LocalFileConfig{
					RootDir: t.TempDir(),
				},
			})
			do.ProvideNamed(injector, "alert-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
				return (*mock.MockDiscordClient)(nil), nil
			})
			do.ProvideNamed(injector, "info-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
				return mockInfoDiscord, nil
			})
			instance, err := NewLogger(injector)
			require.NoError(t, err)

			assert.NotPanics(t, func() {
				instance.Info(tc.message, tc.contexts)
			})
		})
	}
}

func TestLogger_Info_DicordError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockInfoDiscord := mock.NewMockDiscordClient(ctrl)
	expectedErr := failure.New(core.ErrDiscordAPI)
	mockInfoDiscord.EXPECT().Comment(gomock.Any(), gomock.Any()).Return(expectedErr)

	injector := do.New()
	do.ProvideValue(injector, config.Config{
		Basic: config.BasicConfig{
			Name: "test-app",
		},
		Logger: config.LoggerConfig{
			Level: zerolog.InfoLevel,
		},
		LocalFile: config.LocalFileConfig{
			RootDir: t.TempDir(),
		},
	})
	do.ProvideNamed(injector, "alert-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
		return (*mock.MockDiscordClient)(nil), nil
	})
	do.ProvideNamed(injector, "info-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
		return mockInfoDiscord, nil
	})
	instance, err := NewLogger(injector)
	require.NoError(t, err)

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
			mockAlertDiscord := mock.NewMockDiscordClient(ctrl)
			mockAlertDiscord.EXPECT().Comment(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

			injector := do.New()
			do.ProvideValue(injector, config.Config{
				Basic: config.BasicConfig{
					Name: "test-app",
				},
				Logger: config.LoggerConfig{
					Level: zerolog.ErrorLevel,
				},
				LocalFile: config.LocalFileConfig{
					RootDir: t.TempDir(),
				},
			})
			do.ProvideNamed(injector, "alert-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
				return mockAlertDiscord, nil
			})
			do.ProvideNamed(injector, "info-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
				return (*mock.MockDiscordClient)(nil), nil
			})
			instance, err := NewLogger(injector)
			require.NoError(t, err)

			assert.NotPanics(t, func() {
				instance.Error(tc.err, tc.contexts)
			})
		})
	}
}

func TestLogger_Error_DicordError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	mockAlertDiscord := mock.NewMockDiscordClient(ctrl)
	expectedErr := failure.New(core.ErrDiscordAPI)
	mockAlertDiscord.EXPECT().Comment(gomock.Any(), gomock.Any()).Return(expectedErr)

	injector := do.New()
	do.ProvideValue(injector, config.Config{
		Basic: config.BasicConfig{
			Name: "test-app",
		},
		Logger: config.LoggerConfig{
			Level: zerolog.ErrorLevel,
		},
		LocalFile: config.LocalFileConfig{
			RootDir: t.TempDir(),
		},
	})
	do.ProvideNamed(injector, "alert-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
		return mockAlertDiscord, nil
	})
	do.ProvideNamed(injector, "info-discord-client", func(i do.Injector) (adapter.DiscordClient, error) {
		return (*mock.MockDiscordClient)(nil), nil
	})
	instance, err := NewLogger(injector)
	require.NoError(t, err)

	assert.NotPanics(t, func() {
		instance.Error(errors.New("some error"), nil)
	})
}
