package infra

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringRoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")
		expectedContent := "Hello, World!"

		err := writeString(filePath, expectedContent)
		assert.NoError(t, err)

		result, err := readString(filePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedContent, result)
	})

	t.Run("異常系_ファイルが存在しない", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		nonExistentPath := filepath.Join(tmpDir, "nonexistent.txt")

		_, err := readString(nonExistentPath)
		assert.Error(t, err)
	})

	t.Run("正常系_空のファイル", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "empty.txt")

		err := writeString(filePath, "")
		assert.NoError(t, err)

		result, err := readString(filePath)
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("正常系_ネストされたディレクトリを作成", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "dir1", "dir2", "test.txt")
		content := "Nested content"

		err := writeString(filePath, content)
		assert.NoError(t, err)

		result, err := readString(filePath)
		assert.NoError(t, err)
		assert.Equal(t, content, result)
	})

	t.Run("正常系_既存ファイルを上書き", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.txt")

		err := writeString(filePath, "Original content")
		assert.NoError(t, err)

		newContent := "Updated content"
		err = writeString(filePath, newContent)
		assert.NoError(t, err)

		result, err := readString(filePath)
		assert.NoError(t, err)
		assert.Equal(t, newContent, result)
	})
}

func TestJSONRoundTrip(t *testing.T) {
	t.Parallel()

	t.Run("正常系", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.json")

		type TestData struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		}

		originalData := TestData{Name: "John", Age: 30}

		err := writeJSON(filePath, originalData)
		assert.NoError(t, err)

		result, err := readJSON[TestData](filePath)
		assert.NoError(t, err)
		assert.Equal(t, originalData, result)
	})

	t.Run("異常系_ファイルが存在しない", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		nonExistentPath := filepath.Join(tmpDir, "nonexistent.json")

		type TestData struct {
			Name string `json:"name"`
		}

		_, err := readJSON[TestData](nonExistentPath)
		assert.Error(t, err)
	})

	t.Run("異常系_不正なJSON形式", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "invalid.json")

		err := writeString(filePath, "{ invalid json }")
		assert.NoError(t, err)

		type TestData struct {
			Name string `json:"name"`
		}

		_, err = readJSON[TestData](filePath)
		assert.Error(t, err)
	})

	t.Run("正常系_複雑な構造体", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "complex.json")

		type Address struct {
			Street string `json:"street"`
			City   string `json:"city"`
		}

		type Person struct {
			Name    string   `json:"name"`
			Age     int      `json:"age"`
			Address Address  `json:"address"`
			Tags    []string `json:"tags"`
		}

		originalData := Person{
			Name: "Jane",
			Age:  25,
			Address: Address{
				Street: "123 Main St",
				City:   "Springfield",
			},
			Tags: []string{"developer", "gamer"},
		}

		err := writeJSON(filePath, originalData)
		assert.NoError(t, err)

		result, err := readJSON[Person](filePath)
		assert.NoError(t, err)
		assert.Equal(t, originalData, result)
	})

	t.Run("正常系_ネストされたディレクトリを作成", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "dir1", "dir2", "test.json")

		type TestData struct {
			Value string `json:"value"`
		}

		data := TestData{Value: "nested"}

		err := writeJSON(filePath, data)
		assert.NoError(t, err)

		result, err := readJSON[TestData](filePath)
		assert.NoError(t, err)
		assert.Equal(t, data, result)
	})

	t.Run("正常系_既存ファイルを上書き", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "test.json")

		type TestData struct {
			Name string `json:"name"`
		}

		originalData := TestData{Name: "Bob"}
		err := writeJSON(filePath, originalData)
		assert.NoError(t, err)

		updatedData := TestData{Name: "Charlie"}
		err = writeJSON(filePath, updatedData)
		assert.NoError(t, err)

		result, err := readJSON[TestData](filePath)
		assert.NoError(t, err)
		assert.Equal(t, updatedData, result)
	})

	t.Run("正常系_空のスライス", func(t *testing.T) {
		t.Parallel()

		tmpDir := t.TempDir()
		filePath := filepath.Join(tmpDir, "empty_slice.json")

		type TestData struct {
			Items []string `json:"items"`
		}

		data := TestData{Items: []string{}}

		err := writeJSON(filePath, data)
		assert.NoError(t, err)

		result, err := readJSON[TestData](filePath)
		assert.NoError(t, err)
		assert.NotNil(t, result.Items)
		assert.Len(t, result.Items, 0)
	})
}
