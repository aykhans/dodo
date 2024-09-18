package utils

import (
	"encoding/json"
	"reflect"
	"strings"
)

func MarshalJSON(v any, maxSliceSize uint) string {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Slice && rv.Len() == 0 {
		return "[]"
	}

	data, err := json.MarshalIndent(truncateLists(v, int(maxSliceSize)), "", "  ")
	if err != nil {
		return "{}"
	}

	return strings.Replace(string(data), `"..."`, "...", -1)
}

func truncateLists(v interface{}, maxItems int) interface{} {
	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		if rv.Len() > maxItems {
			newSlice := reflect.MakeSlice(rv.Type(), maxItems, maxItems)
			reflect.Copy(newSlice, rv.Slice(0, maxItems))
			newSlice = reflect.Append(newSlice, reflect.ValueOf("..."))
			return newSlice.Interface()
		}
	case reflect.Map:
		newMap := reflect.MakeMap(rv.Type())
		for _, key := range rv.MapKeys() {
			newMap.SetMapIndex(key, reflect.ValueOf(truncateLists(rv.MapIndex(key).Interface(), maxItems)))
		}
		return newMap.Interface()
	case reflect.Struct:
		newStruct := reflect.New(rv.Type()).Elem()
		for i := 0; i < rv.NumField(); i++ {
			newStruct.Field(i).Set(reflect.ValueOf(truncateLists(rv.Field(i).Interface(), maxItems)))
		}
		return newStruct.Interface()
	case reflect.Ptr:
		if rv.IsNil() {
			return nil
		}
		return truncateLists(rv.Elem().Interface(), maxItems)
	}

	return v
}
