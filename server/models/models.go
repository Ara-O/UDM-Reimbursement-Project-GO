package models

import "encoding/json"

type FoapaDetails struct {
}

type UserDataPreVerification struct {
	FirstName        string `json:"first_name" validate:"required,min=2,max=50" mapstructure:"first_name"`
	LastName         string `json:"last_name" validate:"required,min=2,max=50" mapstructure:"last_name"`
	PhoneNumber      int64  `json:"phone_number" validate:"required,number" mapstructure:"phone_number"`
	WorkEmail        string `json:"work_email" validate:"required,alphanum" mapstructure:"work_email"`
	EmploymentNumber int64  `json:"employment_number" validate:"required,number" mapstructure:"employment_number"`
	Department       string `json:"department" validate:"required" mapstructure:"department"`
}

type UserDataPostVerification struct {
	FirstName        string         `json:"first_name" validate:"required,min=2,max=50" mapstructure:"first_name"`
	LastName         string         `json:"last_name" validate:"required,min=2,max=50" mapstructure:"last_name"`
	PhoneNumber      int64          `json:"phone_number" validate:"required,number" mapstructure:"phone_number"`
	WorkEmail        string         `json:"work_email" validate:"required,alphanum" mapstructure:"work_email"`
	EmploymentNumber int64          `json:"employment_number" validate:"required,number" mapstructure:"employment_number"`
	Department       string         `json:"department" validate:"required" mapstructure:"department"`
	City             string         `json:"city" validate:"required"`
	Country          string         `json:"country" validate:"required"`
	FoapaDetails     []FoapaDetails `json:"foapa_details" validate:"required"`
	MailingAddress   string         `json:"mailing_address" validate:"required"`
	Password         string         `json:"password" validate:"required"`
	PostalCode       string         `json:"postal_code" validate:"required"`
	State            string         `json:"state" validate:"required"`
}

func (u *UserDataPreVerification) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}
