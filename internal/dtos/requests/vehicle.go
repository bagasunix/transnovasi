package requests

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Vehicle struct {
	PlateNo  string `json:"plate_no"`
	Model    string `json:"model"`
	Brand    string `json:"brand"`
	Color    string `json:"color"`
	Year     string `json:"year"`
	MaxSpeed int    `json:"max_speed"`
	FuelType string `json:"fuel_type"`
}

func (v Vehicle) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.PlateNo, validation.Required),
		validation.Field(&v.Model, validation.Required),
		validation.Field(&v.Brand, validation.Required),
		validation.Field(&v.Color, validation.Required),
		validation.Field(&v.Year, validation.Required),
		validation.Field(&v.FuelType, validation.Required),
	)
}
