package util

import (
	"context"
	"reflect"
	"strings"
	"unicode"

	"golang.org/x/sync/errgroup"
)

func FieldQuery(target reflect.Type) string {
	fields := []string{}
	fieldsRecursive([]string{}, target, &fields)

	return strings.Join(fields, ",")
}

func ToSnakeCase(s string) string {
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

func MakeRange(min, max int) []int {
	if min > max {
		return []int{}
	}

	a := make([]int, max-min)
	for i := range a {
		a[i] = min + i
	}

	return a
}

func DoParallel[T any](values []T, fn func(value T) error) error {
	eg, _ := errgroup.WithContext(context.Background())

	for _, v := range values {
		eg.Go(func() error {
			return fn(v)
		})
	}

	return eg.Wait()
}

func fieldsRecursive(parentNames []string, t reflect.Type, result *[]string) {
	for i := range t.NumField() {
		field := t.Field(i)

		if field.Type.Kind() == reflect.Struct {
			name := ToSnakeCase(field.Name)
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

func SafeDivide(numerator float64, denominator uint) float64 {
	if denominator == 0 {
		return 0
	}
	return numerator / float64(denominator)
}
