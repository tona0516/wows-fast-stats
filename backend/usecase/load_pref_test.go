package usecase

import (
	"testing"
	"wfs/backend/adapter"
	"wfs/backend/core"
	"wfs/backend/mock"

	"github.com/morikuni/failure"
	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestLoadPref_Invoke(t *testing.T) {
	t.Parallel()

	t.Run("正常系_既存のPrefが取得できる", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		expectedPref := core.Pref{
			Version:      1,
			InstallPath:  "/path/to/wows",
			ZoomRate:     150,
			StatsExtra:   "extra_value",
			IsSendReport: true,
		}

		mockPrefStore.EXPECT().
			Pref().
			Return(expectedPref, nil)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
			return mockPrefStore, nil
		})
		do.Provide(injector, NewLoadPref)
		lp := do.MustInvoke[*LoadPref](injector)

		result, err := lp.Invoke()

		assert.NoError(t, err)
		assert.Equal(t, expectedPref, result)
	})

	t.Run("正常系_Prefが見つからない場合、デフォルト値を返す", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		mockPrefStore.EXPECT().
			Pref().
			Return(core.Pref{}, failure.New(core.ErrJSONNotFound))

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
			return mockPrefStore, nil
		})
		do.Provide(injector, NewLoadPref)
		lp := do.MustInvoke[*LoadPref](injector)

		result, err := lp.Invoke()

		assert.NoError(t, err)
		// デフォルト値のPrefが返される
		assert.Equal(t, core.DefaultPref(), result)
	})

	t.Run("異常系_ファイル読み込みエラーが発生した場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		expectedErr := failure.New(core.ErrJSONRead)
		mockPrefStore.EXPECT().
			Pref().
			Return(core.Pref{}, expectedErr)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
			return mockPrefStore, nil
		})
		do.Provide(injector, NewLoadPref)
		lp := do.MustInvoke[*LoadPref](injector)

		result, err := lp.Invoke()

		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("異常系_JSONパースエラーが発生した場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		expectedErr := failure.New(core.ErrJSONRead)
		mockPrefStore.EXPECT().
			Pref().
			Return(core.Pref{}, expectedErr)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
			return mockPrefStore, nil
		})
		do.Provide(injector, NewLoadPref)
		lp := do.MustInvoke[*LoadPref](injector)

		result, err := lp.Invoke()

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}
