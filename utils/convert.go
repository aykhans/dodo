package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type TruncatedMarshaller struct {
	Value    interface{}
	MaxItems int
}

func (t TruncatedMarshaller) MarshalJSON() ([]byte, error) {
	val := reflect.ValueOf(t.Value)

	if val.Kind() != reflect.Slice && val.Kind() != reflect.Array {
		return json.Marshal(t.Value)
	}

	length := val.Len()
	if length <= t.MaxItems {
		return json.Marshal(t.Value)
	}

	truncated := make([]interface{}, t.MaxItems+1)

	for i := 0; i < t.MaxItems; i++ {
		truncated[i] = val.Index(i).Interface()
	}

	remaining := length - t.MaxItems
	truncated[t.MaxItems] = fmt.Sprintf("+%d", remaining)

	return json.Marshal(truncated)
}

func PrettyJSONMarshal(v interface{}, maxItems int, prefix, indent string) []byte {
	truncated := processValue(v, maxItems)
	d, _ := json.MarshalIndent(truncated, prefix, indent)
	return d
}

func processValue(v interface{}, maxItems int) interface{} {
	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Map:
		newMap := make(map[string]interface{})
		iter := val.MapRange()
		for iter.Next() {
			k := iter.Key().String()
			newMap[k] = processValue(iter.Value().Interface(), maxItems)
		}
		return newMap

	case reflect.Slice, reflect.Array:
		return TruncatedMarshaller{Value: v, MaxItems: maxItems}

	case reflect.Struct:
		newMap := make(map[string]interface{})
		t := val.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.IsExported() {
				jsonTag := field.Tag.Get("json")
				if jsonTag == "-" {
					continue
				}
				fieldName := field.Name
				if jsonTag != "" {
					fieldName = jsonTag
				}
				newMap[fieldName] = processValue(val.Field(i).Interface(), maxItems)
			}
		}
		return newMap

	default:
		return v
	}
}
