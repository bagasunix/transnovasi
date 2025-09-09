package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l Login) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.Required, is.Email),
		validation.Field(&l.Password, validation.Required, validation.Length(6, 0)),
	)
}
