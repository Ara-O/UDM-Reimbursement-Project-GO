package models

import (
	"encoding/json"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type FoapaDetails struct {
	FoapaName     string `json:"foapa_name" bson:"foapa_name"`
	Organization  string `json:"organization"`
	Account       string `json:"account"`
	Program       string `json:"program"`
	Activity      string `json:"activity"`
	CurrentAmount string `json:"current_amount" bson:"current_amount"`
	Fund          string `json:"fund"`
	InitialAmount string `json:"initial_amount" bson:"initial_amount"`
}

type UserDataPreVerification struct {
	FirstName        string `json:"first_name" validate:"required,max=50" mapstructure:"first_name"`
	LastName         string `json:"last_name" validate:"required,max=50" mapstructure:"last_name"`
	PhoneNumber      int64  `json:"phone_number" validate:"required,number" mapstructure:"phone_number"`
	WorkEmail        string `json:"work_email" validate:"required,alphanum" mapstructure:"work_email"`
	EmploymentNumber int64  `json:"employment_number" validate:"required,number" mapstructure:"employment_number"`
	Department       string `json:"department" validate:"required" mapstructure:"department"`
}

type UserDataPostVerification struct {
	FirstName        string         `json:"first_name" bson:"first_name" validate:"required,min=2,max=50"`
	LastName         string         `json:"last_name" bson:"last_name" validate:"required,min=2,max=50"`
	PhoneNumber      int64          `json:"phone_number" bson:"phone_number" validate:"required,number"`
	WorkEmail        string         `json:"work_email" bson:"work_email" validate:"required,email"`
	EmploymentNumber int64          `json:"employment_number" bson:"employment_number"  validate:"required,number"`
	Department       string         `json:"department" bson:"department" validate:"required"`
	City             string         `json:"city" bson:"city" validate:"required"`
	Country          string         `json:"country" bson:"country" validate:"required"`
	FoapaDetails     []FoapaDetails `json:"foapa_details" bson:"foapa_details" validate:"required"`
	MailingAddress   string         `json:"mailing_address" bson:"mailing_address"  validate:"required"`
	Password         string         `json:"password" bson:"password"  validate:"required"`
	PostalCode       string         `json:"postal_code" bson:"postal_code" validate:"required"`
	State            string         `json:"state" bson:"state"  validate:"required"`
}

func (u *UserDataPreVerification) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *UserDataPreVerification) ValidateStruct() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(u)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (u *UserDataPostVerification) ValidateStruct() error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(u)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
