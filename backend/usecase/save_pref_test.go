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

func TestSavePref_Invoke(t *testing.T) {
	t.Parallel()

	t.Run("正常系_Prefが正常に保存される", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		prefToSave := core.Pref{
			Version:        1,
			GameClientPath: "/path/to/wows",
			ZoomRate:       150,
			StatsExtra:     "extra_value",
			IsSendReport:   true,
		}

		mockPrefStore.EXPECT().
			SetPref(prefToSave).
			Return(nil)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
			return mockPrefStore, nil
		})
		do.Provide(injector, NewSavePref)
		sp := do.MustInvoke[*SavePref](injector)

		err := sp.Invoke(prefToSave)

		assert.NoError(t, err)
	})

	t.Run("正常系_複数回呼び出しでそれぞれ異なるPrefが保存される", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		pref1 := core.Pref{
			Version:        1,
			GameClientPath: "/path1",
			ZoomRate:       100,
		}
		pref2 := core.Pref{
			Version:        1,
			GameClientPath: "/path2",
			ZoomRate:       200,
		}

		gomock.InOrder(
			mockPrefStore.EXPECT().
				SetPref(pref1).
				Return(nil),
			mockPrefStore.EXPECT().
				SetPref(pref2).
				Return(nil),
		)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
			return mockPrefStore, nil
		})
		do.Provide(injector, NewSavePref)
		sp := do.MustInvoke[*SavePref](injector)

		err1 := sp.Invoke(pref1)
		assert.NoError(t, err1)

		err2 := sp.Invoke(pref2)
		assert.NoError(t, err2)
	})

	t.Run("異常系_ファイル書き込みエラーが発生した場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		prefToSave := core.Pref{
			Version:        1,
			GameClientPath: "/path/to/wows",
		}

		expectedErr := failure.New(core.ErrJSONWrite)
		mockPrefStore.EXPECT().
			SetPref(prefToSave).
			Return(expectedErr)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
			return mockPrefStore, nil
		})
		do.Provide(injector, NewSavePref)
		sp := do.MustInvoke[*SavePref](injector)

		err := sp.Invoke(prefToSave)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("異常系_空のPrefを保存する場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		emptyPref := core.Pref{}

		mockPrefStore.EXPECT().
			SetPref(emptyPref).
			Return(nil)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
			return mockPrefStore, nil
		})
		do.Provide(injector, NewSavePref)
		sp := do.MustInvoke[*SavePref](injector)

		err := sp.Invoke(emptyPref)

		assert.NoError(t, err)
	})

	t.Run("異常系_ディレクトリ作成に失敗した場合", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		mockPrefStore := mock.NewMockPrefStore(ctrl)

		prefToSave := core.Pref{
			Version:        1,
			GameClientPath: "/path/to/wows",
		}

		expectedErr := failure.New(core.ErrJSONWrite)
		mockPrefStore.EXPECT().
			SetPref(prefToSave).
			Return(expectedErr)

		injector := do.New()
		do.Provide(injector, func(i do.Injector) (adapter.PrefStore, error) {
			return mockPrefStore, nil
		})
		do.Provide(injector, NewSavePref)
		sp := do.MustInvoke[*SavePref](injector)

		err := sp.Invoke(prefToSave)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
