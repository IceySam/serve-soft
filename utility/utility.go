package utility

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

/*
 Validate request data
*/
func Validate(s interface{}) error {
	errs := ""

	var m map[string]interface{}
	data, err := json.Marshal(s)
	if err != nil {
		errs = err.Error()
	}
	json.Unmarshal(data, &m)

	for k, v := range m {
		if v == nil || v == "" || v == 0 {
			errs = fmt.Sprintf("%s, %s is required", errs, k)
		}
	}

	if len(errs) > 0 {
		return errors.New(errs)
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

	return m, ty, name, nil
}

func ToStuct(m map[string]interface{}, i interface{}) {
	data, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, i)
	if err != nil {
		log.Fatal(err)
	}
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
