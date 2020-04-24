package session

import (
	"afiqo-location/helpers"
	"context"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

const (
	USER_SESSION = "USER_SESSION"
)

const (
	CUSTOMER_ROLE = "customer"
	ADMIN_ROLE    = "admin"
	SUPPLIER_ROLE = "supplier"
	COURIER_ROLE  = "courier"
)

type (
	Session struct {
		UserID     uuid.UUID `json:"user_id"`
		SessionKey string    `json:"session_key"`
		Expiry     int       `json:"expiry"`
		Role       string    `json:"role"`
	}

	SessionData struct {
		UserID uuid.UUID `json:"user_id"`
		Role   string    `json:"role"`
	}
)

func (s Session) Store(ctx context.Context) error {
	sessionData := SessionData{
		UserID: s.UserID,
		Role:   s.Role,
	}

	sessionMarshall, err := json.Marshal(sessionData)

	if err != nil {
		return err
	}

	err = helpers.SetDataToCacheWithExpiry(ctx, s.SessionKey, string(sessionMarshall), s.Expiry)

	if err != nil {
		return err
	}

	return nil

}

func (s Session) Get(ctx context.Context) (SessionData, error) {
	session, err := helpers.GetDataFromCache(ctx, s.SessionKey)
	if err != nil {
		return SessionData{}, err
	}

	var sessionData SessionData

	err = json.Unmarshal([]byte(session), &sessionData)

	if err != nil {
		return SessionData{}, err
	}

	return sessionData, nil

}
