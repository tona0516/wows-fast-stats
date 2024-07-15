package response

import (
	"reflect"
	"strings"
)

func wgAPIField(i any) string {
	parents := make([]string, 0)
	elem := reflect.TypeOf(i).Elem()
	fields := jsonTagRecursive(parents, elem)
	joined := strings.Join(fields, ",")
	return joined
}

func jsonTagRecursive(parents []string, t reflect.Type) []string {
	tags := make([]string, 0)

	for i := range t.NumField() {
		f := t.Field(i)
		tag := f.Tag.Get("json")

		if f.Type.Kind() == reflect.Struct {
			childTags := jsonTagRecursive(append(parents, tag), f.Type)
			tags = append(tags, childTags...)
			continue
		}

		if len(parents) > 0 {
			tag = strings.Join(parents, ".") + "." + tag
		}
		tags = append(tags, tag)
	}

	return tags
}
