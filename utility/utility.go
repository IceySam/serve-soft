package utility

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

/*
 Validate request data
*/
func Validate(s interface{}) error {
	var errs error

	var m map[string]interface{}
	data, err := json.Marshal(s)
	if err != nil {
		errs = err
	}
	json.Unmarshal(data, &m)

	for k, v := range m {
		if v == nil || v == "" || v == 0 {
			errs = fmt.Errorf("%v, %v is required", errs, k)
		}
	}
	if errs != nil {
		return errs
	}
	return nil
}

/*
 Convert struct to map
*/
func ToMap(s interface{}) (map[string]interface{}, reflect.Type, string, error) {
	ty := reflect.TypeOf(s)
	name := strings.ToLower(strings.Split(ty.String(), ".")[1])
	var m map[string]interface{}
	data, err := json.Marshal(s)
	if err != nil {
		return nil, ty, name, err
	}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, ty, name, err
	}
	if ty.Kind() == reflect.Struct && ty.NumField() != len(m) {
		return nil, ty, name, fmt.Errorf("private filed")
	}
	if ty.Kind() == reflect.Ptr {
		v := reflect.ValueOf(s).Elem()
		if v.NumField() != len(m) {
			return nil, ty, name, fmt.Errorf("private filed")
		}
	}
	return m, ty, name, nil
}

func ToStruct(m map[string]interface{}, i interface{}) error {
	ty := reflect.TypeOf(i).Elem()
	vy := reflect.ValueOf(i).Elem()
	if ty.Kind() != reflect.Struct {
		return fmt.Errorf("requires struct found, %s", ty.Kind())
	}

	conv := make(map[string]interface{}, len(m))
	mapKeys := make(map[string]interface{}, len(m))
	for k, v := range m {
		conv[strings.ToLower(k)] = v
		mapKeys[strings.ToLower(k)] = k
	}

	for x := 0; x < ty.NumField(); x++ {
		structField := ty.Field(x)
		structValue := vy.Field(x)

		value := conv[strings.ToLower(structField.Name)]

		if mapKeys[strings.ToLower(structField.Name)] == nil {
			return fmt.Errorf("field %s unavailable", structField.Name)
		}

		var res reflect.Value
		if structField.Type.Kind() == reflect.Bool {
			b, err := strconv.ParseBool(fmt.Sprintf("%v", value))
			if err != nil {
				return fmt.Errorf("%v is not assignable to %v %v", reflect.TypeOf(value), structField.Name, structField.Type)
			}
			res = reflect.ValueOf(b)

		} else {
			if !reflect.TypeOf(value).AssignableTo(structField.Type) {
				return fmt.Errorf("%v is not assignable to %v %v", reflect.TypeOf(value), structField.Name, structField.Type)
			}
			res = reflect.ValueOf(value)
		}
		if value != nil {
			structValue.Set(res)
		}
	}

	return nil
}

func ToStructArray(m []map[string]interface{}, i interface{}) error {
	sliceType := reflect.TypeOf(i).Elem()
	sliceValue := reflect.ValueOf(i).Elem()

	// check slice
	if sliceType.Kind() != reflect.Slice {
		return fmt.Errorf("requires struct found, %s", sliceType.Kind())
	}

	ty := sliceType.Elem()
	// check struct
	if ty.Kind() != reflect.Struct {
		return fmt.Errorf("requires struct found, %s", ty.Kind())
	}

	// new slice
	vy := reflect.MakeSlice(sliceType, len(m), len(m))

	conv := make([]map[string]interface{}, 0)
	mapKeys := make([]map[string]interface{}, 0)
	for _, m := range m {
		item := make(map[string]interface{}, len(m))
		keyItem := make(map[string]interface{}, len(m))
		for k, v := range m {
			item[strings.ToLower(k)] = v
			keyItem[strings.ToLower(k)] = k
		}
		conv = append(conv, item)
		mapKeys = append(mapKeys, keyItem)
	}

	// verify and set fields
	for x := 0; x < len(m); x++ {
		for y := 0; y < ty.NumField(); y++ {
			structField := ty.Field(y)
			structValue := vy.Index(x).Field(y)

			value := conv[x][strings.ToLower(structField.Name)]

			if x == 0 && mapKeys[x][strings.ToLower(structField.Name)] == nil {
				return fmt.Errorf("field %s unavailable", structField.Name)
			}

			var res any
			if structField.Type.Kind() == reflect.Bool {
				b, err := strconv.ParseBool(fmt.Sprintf("%v", value))
				if err != nil {
					return fmt.Errorf("%v is not assignable to %v %v", reflect.TypeOf(value), structField.Name, structField.Type)
				}
				res = b
			} else {
				if x == 0 && !reflect.TypeOf(value).AssignableTo(structField.Type) {
					return fmt.Errorf("%v is not assignable to %v %v", reflect.TypeOf(value), structField.Name, structField.Type)
				}
				res = value
			}
			if value != nil {
				structValue.Set(reflect.ValueOf(res))
			}
		}
	}

	sliceValue.Set(vy)

	return nil
}

func MakeStruct(m map[string]interface{}, ty reflect.Type) interface{} {
	val := make([]reflect.StructField, 0)
	for k, v := range m {
		val = append(val, reflect.StructField{
			Name:    k,
			Type:    reflect.TypeOf(v),
			PkgPath: ty.String(),
		})
	}
	res := reflect.New(reflect.StructOf(val)).Elem()
	return res
}

func MapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TypeEquals(data any, compare any) bool {
	if data == nil && compare == nil {
		return true
	} else if data == nil && compare != nil {
		return false
	}
	return reflect.TypeOf(data).AssignableTo(reflect.TypeOf(compare))
}

func ParseAny(byt sql.RawBytes) any {
	str := string(byt)

	if val, err := strconv.ParseInt(str, 10, 64); err == nil {
		return val
	} else if val, err := strconv.ParseFloat(str, 64); err == nil {
		return val
	} else if val, err := strconv.ParseBool(str); err == nil {
		return val
	} else if byt == nil {
		return nil
	} else {
		return str
	}
}
