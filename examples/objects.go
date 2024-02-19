package examples

import "github.com/IceySam/serve-soft/utility"

type Food struct {
	Name string `json:"name"`
}

func (f *Food) validate() error {
	return utility.Validate(f)
}
