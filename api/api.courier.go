package api

import (
	"afiqo-location/email"
	"afiqo-location/helpers"
	"afiqo-location/models"
	"afiqo-location/session"
	"afiqo-location/util"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type (
	CourierModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	CourierWithSession struct {
		Courier models.CourierResponse `json:"courier"`
		Session string                 `json:"session"`
	}

	CourierLoginParam struct {
		Email    string `json:"email" validate:"email,required"`
		Password string `json:"password" validate:"required"`
	}

	CourierDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	CourierAddParam struct {
		Name    string `json:"name" validate:"max=20,min=4,required"`
		Address string `json:"address" validate:"required"`
		PhoneNo string `json:"phone_no" validate:"required"`
		Email   string `json:"email" validate:"email,required"`
	}

	CourierUpdateParam struct {
		ID      uuid.UUID `json:"id"`
		Name    string    `json:"name" validate:"max=20,min=4,required"`
		Address string    `json:"address" validate:"required"`
		PhoneNo string    `json:"phone_no" validate:"required"`
	}

	CourierDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewCourierModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *CourierModule {
	return &CourierModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/courier",
	}
}

func (s CourierModule) Login(ctx context.Context, param CourierLoginParam) (interface{}, *helpers.Error) {

	courier, err := models.GetOneCourierByEmail(ctx, s.db, param.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, helpers.ErrorWrap(err, s.name, "Login/Email", helpers.IncorrectEmailMessage,
				http.StatusInternalServerError)
		}
		return nil, helpers.ErrorWrap(err, s.name, "Login/GetOneCourierByEmail", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(courier.Password), []byte(param.Password))
	if err != nil {
		return nil, helpers.ErrorWrap(errors.New("Invalid Password"), s.name, "Login/CompareHashAndPassword",
			helpers.IncorrectPasswordMessage,
			http.StatusInternalServerError)
	}

	session := session.Session{
		UserID:     courier.ID,
		SessionKey: fmt.Sprintf(`%s:%s`, session.USER_SESSION, uuid.NewV4()),
		Expiry:     86400,
		Role:       session.COURIER_ROLE,
	}

	err = session.Store(ctx)

	courierSession := CourierWithSession{
		Courier: courier.Response(),
		Session: session.SessionKey,
	}

	return courierSession, nil

}

func (s CourierModule) PasswordUpdate(ctx context.Context, param PasswordUpdateParam) (interface{}, *helpers.Error) {

	courier, err := models.GetOneCourier(ctx, s.db, param.ID)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "PasswordUpdate/GetOneCourier", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(courier.Password), []byte(param.CurrentPassword))

	if err != nil {
		return nil, helpers.ErrorWrap(errors.New("Current Password Is Incorrect!"), s.name,
			"PasswordUpdate/CompareHashAndPassword",
			helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	if param.NewPassword == param.CurrentPassword {
		return nil, helpers.ErrorWrap(errors.New("New Password Cannot Be Same With Current Password"), s.name,
			"PasswordUpdate/CurrentPasswordComparison", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	if param.NewPassword != param.ConfirmNewPassword {
		return nil, helpers.ErrorWrap(errors.New("New Password Does Not Match"), s.name,
			"PasswordUpdate/NewPassword", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	courier = models.CourierModel{
		ID:       param.ID,
		Password: param.NewPassword,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}
	err = courier.PasswordUpdate(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "PasswordUpdate/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	updatePasswordResponse := models.UpdatePasswordResponse{
		Message: "Password Successfully Changed",
	}

	return updatePasswordResponse, nil

}

func (s CourierModule) Detail(ctx context.Context, param CourierDetailParam) (interface{}, *helpers.Error) {
	courier, err := models.GetOneCourier(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneCourier", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return courier.Response(), nil
}

func (s CourierModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	couriers, err := models.GetAllCourier(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllCourier", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var courierResponse []models.CourierResponse
	for _, courier := range couriers {
		courierResponse = append(courierResponse, courier.Response())
	}

	return courierResponse, nil
}

func (s CourierModule) Add(ctx context.Context, param CourierAddParam) (interface{}, *helpers.Error) {

	password := util.RandomString(10)

	courier := models.CourierModel{
		Name:      param.Name,
		PhoneNo:   param.PhoneNo,
		Email:     param.Email,
		Address:   param.Address,
		Password:  password,
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := courier.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	data := email.MailData{
		Name: param.Name,
		Actions: []email.Action{
			{
				Button: email.Button{
					Text: password,
				},
			},
		},
	}

	body, err := data.GenerateForPassword()

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GenerateForPassword", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	mail := email.Mail{
		Subject: "Password For Login",
		Body:    body,
		To:      courier.Email,
	}

	go func() {
		mail.SendEmail()
	}()

	return courier.Response(), nil
}

func (s CourierModule) Update(ctx context.Context, param CourierUpdateParam) (interface{}, *helpers.Error) {

	courier := models.CourierModel{
		ID:      param.ID,
		Name:    param.Name,
		Address: param.Address,
		PhoneNo: param.PhoneNo,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := courier.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return courier.Response(), nil

}

func (s CourierModule) Delete(ctx context.Context, param CourierDeleteParam) (interface{}, *helpers.Error) {

	courier := models.CourierModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := courier.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
