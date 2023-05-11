package vo

import (
	"reflect"
	"strings"
	"unicode"
)

func fieldQuery(target reflect.Type) string {
	fields := []string{}
	fieldsRecursive([]string{}, target, &fields)

	return strings.Join(fields, ",")
}

func fieldsRecursive(parentNames []string, t reflect.Type, result *[]string) {
	for i := 0; i < t.NumField(); i++ {
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
	b := &strings.Builder{}
	for i, r := range s {
		if i == 0 {
			b.WriteRune(unicode.ToLower(r))

			continue
		}
		if unicode.IsUpper(r) {
			b.WriteRune('_')
			b.WriteRune(unicode.ToLower(r))

			continue
		}
		b.WriteRune(r)
	}

	return b.String()
}
