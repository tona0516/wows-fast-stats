package service

import (
	"testing"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUtil_makeRange(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []int{1, 2, 3, 4}, makeRange(1, 5))
	assert.Equal(t, []int{-5, -4, -3, -2, -1}, makeRange(-5, 0))
	assert.Equal(t, []int{}, makeRange(0, 0))
	assert.Equal(t, []int{}, makeRange(0, -1))
}

func TestUtil_doParallel(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		values := makeRange(1, 5)

		var calls int
		err := doParallel(values, func(value int) error {
			calls++
			return nil
		})

		require.NoError(t, err)
		assert.Len(t, values, calls)
	})
	t.Run("異常系", func(t *testing.T) {
		t.Parallel()

		values := makeRange(1, 5)

		expected := apperr.HTTPRequestError
		err := doParallel(values, func(value int) error {
			if value == values[len(values)-1] {
				return failure.New(expected)
			}
			return nil
		})

		code, ok := failure.CodeOf(err)
		assert.True(t, ok)
		assert.Equal(t, expected, code)
	})
}
