package api

import (
	"afiqo-location/helpers"
	"afiqo-location/models"
	"afiqo-location/session"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type (
	CustomerModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	CustomerWithSession struct {
		Customer models.CustomerResponse `json:"customer"`
		Session  string                  `json:"session"`
	}

	CustomerRegisterParam struct {
		Name            string    `json:"name" validate:"max=20,min=4,required"`
		Address         string    `json:"address" validate:"omitempty"`
		DateOfBirth     time.Time `json:"date_of_birth" validate:"lt,required"`
		Gender          int       `json:"gender" validate:"max=2,min=1,required"`
		PhoneNo         string    `json:"phone_no" validate:"required"`
		Email           string    `json:"email" validate:"email,required"`
		Password        string    `json:"password" validate:"required"`
		ConfirmPassword string    `json:"confirm_password" validate:"required"`
	}

	CustomerLoginParam struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	CustomerDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	CustomerAddParam struct {
		Name        string    `json:"name" validate:"max=20,min=4,required"`
		Address     string    `json:"address" validate:"omitempty"`
		DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
		Gender      int       `json:"gender" validate:"max=2,min=1,required"`
		PhoneNo     string    `json:"phone_no" validate:"required"`
		Email       string    `json:"email" validate:"email,required"`
		Password    string    `json:"password" validate:"required"`
	}

	CustomerUpdateParam struct {
		ID          uuid.UUID
		Name        string    `json:"name" validate:"max=20,min=4,required"`
		Address     string    `json:"address" validate:"omitempty"`
		DateOfBirth time.Time `json:"date_of_birth" validate:"required"`
		Gender      int       `json:"gender" validate:"max=1,min=0,required"`
		PhoneNo     string    `json:"phone_no" validate:"required"`
	}

	CustomerDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewCustomerModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *CustomerModule {
	return &CustomerModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/customer",
	}
}

func (s CustomerModule) Register(ctx context.Context, param CustomerRegisterParam) (interface{}, *helpers.Error) {

	if param.Password != param.ConfirmPassword {
		return nil, helpers.ErrorWrap(errors.New("Password Does Not Match !"), s.name,
			"Register/Password", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	customer := models.CustomerModel{
		Name:        param.Name,
		Address:     param.Address,
		DateOfBirth: param.DateOfBirth,
		Gender:      param.Gender,
		PhoneNo:     param.PhoneNo,
		Email:       param.Email,
		Password:    param.Password,
		CreatedBy:   uuid.NewV4(),
	}

	err := customer.Insert(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Register/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	session := session.Session{
		UserID:     customer.ID,
		SessionKey: fmt.Sprintf(`%s:%s`, session.USER_SESSION, uuid.NewV4()),
		Expiry:     86400,
		Role:       session.CUSTOMER_ROLE,
	}

	err = session.Store(ctx)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Register/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	customerSession := CustomerWithSession{
		Customer: customer.Response(),
		Session:  session.SessionKey,
	}

	return customerSession, nil
}

func (s CustomerModule) Login(ctx context.Context, param CustomerLoginParam) (interface{}, *helpers.Error) {

	customer, err := models.GetOneCustomerByEmail(ctx, s.db, param.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, helpers.ErrorWrap(err, s.name, "Login/Email", helpers.IncorrectEmailMessage,
				http.StatusInternalServerError)
		}
		return nil, helpers.ErrorWrap(err, s.name, "Login/GetOneCustomerByEmail", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(param.Password))
	if err != nil {
		return nil, helpers.ErrorWrap(errors.New("Invalid Password"), s.name, "Login/CompareHashAndPassword",
			helpers.IncorrectPasswordMessage,
			http.StatusInternalServerError)
	}

	session := session.Session{
		UserID:     customer.ID,
		SessionKey: fmt.Sprintf(`%s:%s`, session.USER_SESSION, uuid.NewV4()),
		Expiry:     86400,
		Role:       session.CUSTOMER_ROLE,
	}

	err = session.Store(ctx)

	customerSession := CustomerWithSession{
		Customer: customer.Response(),
		Session:  session.SessionKey,
	}

	return customerSession, nil

}

func (s CustomerModule) PasswordUpdate(ctx context.Context, param PasswordUpdateParam) (interface{}, *helpers.Error) {

	customer, err := models.GetOneCustomer(ctx, s.db, param.ID)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "PasswordUpdate/GetOneCustomer", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(param.CurrentPassword))

	if err != nil {
		return nil, helpers.ErrorWrap(errors.New("Current Password Is Incorrect!"), s.name,
			"PasswordUpdate/CompareHashAndPassword",
			helpers.IncorrectPasswordMessage,
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

	customer = models.CustomerModel{
		ID:       param.ID,
		Password: param.NewPassword,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}
	err = customer.PasswordUpdate(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "PasswordUpdate/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	updatePasswordResponse := models.UpdatePasswordResponse{
		Message: "Password Successfully Changed",
	}

	return updatePasswordResponse, nil

}

func (s CustomerModule) Detail(ctx context.Context, param CustomerDetailParam) (interface{}, *helpers.Error) {
	customer, err := models.GetOneCustomer(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneCustomer", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return customer.Response(), nil
}

func (s CustomerModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	customers, err := models.GetAllCustomer(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllCustomer", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var customerResponse []models.CustomerResponse
	for _, customer := range customers {
		customerResponse = append(customerResponse, customer.Response())
	}

	return customerResponse, nil
}

func (s CustomerModule) Add(ctx context.Context, param CustomerAddParam) (interface{}, *helpers.Error) {

	customer := models.CustomerModel{
		Name:        param.Name,
		Gender:      param.Gender,
		DateOfBirth: param.DateOfBirth,
		Address:     param.Address,
		PhoneNo:     param.PhoneNo,
		Email:       param.Email,
		Password:    param.Password,
		CreatedBy:   uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := customer.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return customer.Response(), nil
}

func (s CustomerModule) Update(ctx context.Context, param CustomerUpdateParam) (interface{}, *helpers.Error) {

	customer := models.CustomerModel{
		ID:          uuid.FromStringOrNil(ctx.Value("user_id").(string)),
		Name:        param.Name,
		Gender:      param.Gender,
		DateOfBirth: param.DateOfBirth,
		Address:     param.Address,
		PhoneNo:     param.PhoneNo,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := customer.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return customer.Response(), nil

}

func (s CustomerModule) Delete(ctx context.Context, param CustomerDeleteParam) (interface{}, *helpers.Error) {

	customer := models.CustomerModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := customer.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
