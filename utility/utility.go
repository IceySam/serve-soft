package utility

import (
	"encoding/json"
	"errors"
	"fmt"
)

/*
* validate request data
 */
func Validate(s any) error {
	var m map[string]string
	data, _ := json.Marshal(s)
	json.Unmarshal(data, &m)

	errs := ""

	for k, v := range m {
		if v == "" || len(v) == 0 {
			errs = fmt.Sprintf("%s, %s is required", errs, k)
		}
	}

	if len(errs) > 0 {
		return errors.New(errs)
	}
	return nil
}
