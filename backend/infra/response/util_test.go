package response

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUtil_wgApiField(t *testing.T) {
	t.Parallel()

	type TestData struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Detail struct {
			Hoge string `json:"hoge"`
			Fuga string `json:"fuga"`
		} `json:"detail"`
	}

	result := wgAPIField(&TestData{})

	expectedResult := "id,name,detail.hoge,detail.fuga"
	assert.Equal(t, expectedResult, result)
}
