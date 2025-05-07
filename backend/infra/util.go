package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
	"golang.org/x/sync/errgroup"
)

const dirPerm = os.FileMode(0o755)

func readJSON[T any](name string, defaulValue T) (T, error) {
	errCtx := failure.Context{"path": name}

	f, err := os.ReadFile(path.Clean(name))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return defaulValue, failure.New(apperr.FileNotExist, errCtx)
		}

		return defaulValue, failure.Wrap(err, errCtx)
	}

	errCtx["target"] = string(f)

	err = json.Unmarshal(f, &defaulValue)

	return defaulValue, failure.Wrap(err, errCtx)
}

func writeJSON[T any](name string, target T) error {
	//nolint:errchkjson
	marshaled, _ := json.Marshal(target)
	errCtx := failure.Context{"path": name, "target": string(marshaled)}

	_ = os.MkdirAll(filepath.Dir(name), dirPerm)

	f, err := os.Create(path.Clean(name))
	if err != nil {
		return failure.Wrap(err, errCtx)
	}

	//nolint: errcheck
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(target)

	return failure.Wrap(err, errCtx)
}

func prettyJSON(str string) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, []byte(str), "", "    "); err != nil {
		return str
	}

	return buf.String()
}

func writeData(name string, target []byte) error {
	_ = os.MkdirAll(filepath.Dir(name), dirPerm)

	f, err := os.Create(path.Clean(name))
	if err != nil {
		return failure.Wrap(err)
	}
	//nolint: errcheck
	defer f.Close()

	_, err = f.Write(target)

	return failure.Wrap(err)
}

func makeRange(start, end int) []int {
	if start > end {
		return []int{}
	}

	a := make([]int, end-start)
	for i := range a {
		a[i] = start + i
	}

	return a
}

func doParallel[T any](values []T, fn func(value T) error) error {
	eg, _ := errgroup.WithContext(context.Background())

	for _, v := range values {
		eg.Go(func() error {
			return fn(v)
		})
	}

	return eg.Wait()
}

func fieldQuery(target reflect.Type) string {
	fields := []string{}
	fieldsRecursive([]string{}, target, &fields)

	return strings.Join(fields, ",")
}

func fieldsRecursive(parentNames []string, t reflect.Type, result *[]string) {
	for i := range t.NumField() {
		field := t.Field(i)

		if field.Type.Kind() == reflect.Struct {
			name := toSnakeCase(field.Name)
			fieldsRecursive(append(parentNames, name), field.Type, result)

			continue
		}

		column := field.Tag.Get("json")
		if len(parentNames) > 0 {
			*result = append(*result, strings.Join(parentNames, ".")+"."+column)
		} else {
			*result = append(*result, column)
		}
	}
}

func toSnakeCase(s string) string {
	runes := []rune(s)
	result := make([]rune, 0)

	for i, r := range runes {
		if !unicode.IsUpper(r) {
			result = append(result, r)

			continue
		}

		if i > 0 && unicode.IsLower(runes[i-1]) {
			result = append(result, '_')
		}

		result = append(result, unicode.ToLower(r))
	}

	return string(result)
}
