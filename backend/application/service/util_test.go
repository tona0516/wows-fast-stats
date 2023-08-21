package service

import (
	"testing"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
	"github.com/stretchr/testify/assert"
)

func TestUtil_makeRange(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []int{1, 2, 3, 4}, makeRange(1, 5))
	assert.Equal(t, []int{-5, -4, -3, -2, -1}, makeRange(-5, 0))
	assert.Equal(t, []int{}, makeRange(0, 0))
	assert.Equal(t, []int{}, makeRange(0, -1))
}

func TestUtil_doParallel_正常系(t *testing.T) {
	t.Parallel()

	values := makeRange(1, 5)

	var calls int
	err := doParallel(2, values, func(value int) error {
		calls++
		return nil
	})

	assert.Nil(t, err)
	assert.Equal(t, len(values), calls)
}

func TestUtil_doParallel_異常系(t *testing.T) {
	t.Parallel()

	values := makeRange(1, 5)

	expected := apperr.HTTPRequestError
	err := doParallel(2, values, func(value int) error {
		if value == values[len(values)-1] {
			return failure.New(expected)
		}
		return nil
	})

	code, ok := failure.CodeOf(err)
	assert.True(t, ok)
	assert.Equal(t, code, expected)
}
