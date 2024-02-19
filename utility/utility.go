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
	name := strings.Split(ty.String(), ".")[1]
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
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, i)
	if err != nil {
		return err
	}
	ty := reflect.ValueOf(i).Elem()
	if ty.Kind() == reflect.Struct && ty.NumField() != len(m) {
		return fmt.Errorf("private field")
	}
	return nil
}

func ToStructArray(m []map[string]interface{}, i interface{}) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, i)
	if err != nil {
		return err
	}
	ty := reflect.ValueOf(i).Elem()
	if ty.Kind() == reflect.Slice && len(m) > 0 && ty.Index(0).NumField() != len(m[0]) {
		return fmt.Errorf("private field")
	}
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
	if (data == nil && compare == nil) {
		return true
	} else if (data == nil && compare != nil) {
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
		return "NULL"
	} else {
		return str
	}
}
