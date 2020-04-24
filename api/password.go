package api

import uuid "github.com/satori/go.uuid"

type (
	PasswordUpdateParam struct {
		ID                 uuid.UUID `json:"id"`
		CurrentPassword    string    `json:"current_password"`
		NewPassword        string    `json:"new_password"`
		ConfirmNewPassword string    `json:"confirm_new_password"`
	}
)
