package response

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
