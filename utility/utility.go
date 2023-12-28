package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
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
func ToMap(s interface{}) (map[string]interface{}, string, error) {
	ty := strings.Split(reflect.TypeOf(s).String(), ".")[1]
	var m map[string]interface{}
	data, err := json.Marshal(s)
	if err != nil {
		return nil, ty, err
	}
	err = json.Unmarshal(data, &m)
	if err != nil {
		return nil, ty, err
	}

	return m, ty, nil
}

func ToStuct(m map[string]interface{}, i interface{})  {
	structVal := reflect.ValueOf(i).Elem()
    
	for k, v := range m {
		field := structVal.FieldByName(k)
		if !field.IsValid() {
			log.Fatalf("No such field: %s in obj", k)
		}
		if !field.CanSet() {
			log.Fatalf("Cannot set %s field value", k)
		}
		if field.Type() != reflect.TypeOf(v) {
			log.Fatalf("Cannot assign %s to %s",reflect.TypeOf(v),  field.Type())
		}
		field.Set(reflect.ValueOf(v))
	}
}

func MapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TypeEquals(data any, compare any) bool {
	if reflect.TypeOf(data).AssignableTo(reflect.TypeOf(compare)) {
		return true
	}
	return false
}
