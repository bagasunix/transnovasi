package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Customer struct {
	Name     string `json:"name"`
	Sex      int    `json:"sex"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
}

func (c Customer) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(6, 0)),
		validation.Field(&c.Sex, validation.Required, validation.In(1, 2)),
		validation.Field(&c.Name, validation.Required),
	)
}
