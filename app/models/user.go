package models

import (
	"fmt"
	"github.com/robfig/revel"
	"regexp"
)

type User struct {
	Email    string
	Password string
	Words    []string
}

func (u *User) String() string {
	return fmt.Sprintf("User(%s)", u.Email)
}

var emailPattern = regexp.MustCompile("[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?")

func (user *User) Validate(v *revel.Validation) {
	v.Required(user.Email)
	v.Match(user.Email, emailPattern).Message("We're sorry, but this email address seems to be invalid. Please enter a valid email address.")
	ValidatePassword(v, user.Password).Key("user.Password")
}

func ValidatePassword(v *revel.Validation, password string) *revel.ValidationResult {
	return v.Check(password,
		revel.Required{},
		revel.MinSize{5},
	)
}
