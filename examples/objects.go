package examples

import "github.com/IceySam/serve-soft/utility"
type User struct {
	Id int64
	FirstName string
	OtherNames string
	Role string
}

func (f *User) validate() error {
	return utility.Validate(f)
}

type Food struct {
	Name string `json:"name"`
}

func (f *Food) validate() error {
	return utility.Validate(f)
}
