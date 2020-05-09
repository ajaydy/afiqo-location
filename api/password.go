package api

import uuid "github.com/satori/go.uuid"

type (
	PasswordUpdateParam struct {
		ID                 uuid.UUID `json:"id"`
		CurrentPassword    string    `json:"current_password" validate:"required"`
		NewPassword        string    `json:"new_password" validate:"gt=6,required"`
		ConfirmNewPassword string    `json:"confirm_new_password" validate:"required"`
	}
)
