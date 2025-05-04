package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"
	"wfs/backend/apperr"

	"github.com/morikuni/failure"
	"golang.org/x/sync/errgroup"
)

func readJSON[T any](path string, defaulValue T) (T, error) {
	errCtx := failure.Context{"path": path}

	f, err := os.ReadFile(path)
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

func writeJSON[T any](path string, target T) error {
	//nolint:errchkjson
	marshaled, _ := json.Marshal(target)
	errCtx := failure.Context{"path": path, "target": string(marshaled)}

	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	f, err := os.Create(path)
	if err != nil {
		return failure.Wrap(err, errCtx)
	}
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

func writeData(path string, target []byte) error {
	_ = os.MkdirAll(filepath.Dir(path), 0o755)

	f, err := os.Create(path)
	if err != nil {
		return failure.Wrap(err)
	}
	defer f.Close()

	_, err = f.Write(target)
	return failure.Wrap(err)
}

func makeRange(min, max int) []int {
	if min > max {
		return []int{}
	}

	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
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
