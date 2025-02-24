package validator

import (
	"go-rest-api/model"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type IUserValidator interface {
	UserValidate(user model.User) error
}

type userValidator struct{}

func NewUserValidator() IUserValidator {
	return &userValidator{}
}

func (uv *userValidator) UserValidate(user model.User) error {
	return validation.ValidateStruct(&user,
		validation.Field(&user.Email,
			validation.Required.Error("email is required"),
			validation.Length(1, 30).Error("email must be 1 to 30 characters"),
			is.Email.Error("invalid email address"),
		),
		validation.Field(&user.Password,
			validation.Required.Error("password is required"),
			validation.Length(6, 30).Error("password must be 6 to 30 characters"),
		),
	)
}
