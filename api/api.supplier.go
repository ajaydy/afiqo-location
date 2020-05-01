package api

import (
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
	SupplierModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	SupplierWithSession struct {
		Supplier models.SupplierResponse `json:"supplier"`
		Session  string                  `json:"session"`
	}

	SupplierDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	SupplierLoginParam struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	SupplierAddParam struct {
		Name    string `json:"name"`
		PhoneNo string `json:"phone_no"`
		Email   string `json:"email"`
	}

	SupplierUpdateParam struct {
		Name    string `json:"name" validate:"max=20,min=4,required"`
		PhoneNo string `json:"phone_no" validate:"required"`
	}

	SupplierDeleteParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewSupplierModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *SupplierModule {
	return &SupplierModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/supplier",
	}
}

func (s SupplierModule) Login(ctx context.Context, param SupplierLoginParam) (interface{}, *helpers.Error) {

	supplier, err := models.GetOneSupplierByEmail(ctx, s.db, param.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, helpers.ErrorWrap(err, s.name, "Login/Email", helpers.IncorrectEmailMessage,
				http.StatusInternalServerError)
		}
		return nil, helpers.ErrorWrap(err, s.name, "Login/GetOneSupplierByEmail", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(supplier.Password), []byte(param.Password))
	if err != nil {
		return nil, helpers.ErrorWrap(errors.New("Invalid Password"), s.name, "Login/CompareHashAndPassword",
			helpers.IncorrectPasswordMessage,
			http.StatusInternalServerError)
	}

	session := session.Session{
		UserID:     supplier.ID,
		SessionKey: fmt.Sprintf(`%s:%s`, session.USER_SESSION, uuid.NewV4()),
		Expiry:     86400,
		Role:       session.SUPPLIER_ROLE,
	}

	err = session.Store(ctx)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Login/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	supplierSession := SupplierWithSession{
		Supplier: supplier.Response(),
		Session:  session.SessionKey,
	}

	return supplierSession, nil

}

func (s SupplierModule) PasswordUpdate(ctx context.Context, param PasswordUpdateParam) (interface{}, *helpers.Error) {

	supplier, err := models.GetOneSupplier(ctx, s.db, param.ID)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "PasswordUpdate/GetOneSupplier", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	err = bcrypt.CompareHashAndPassword([]byte(supplier.Password), []byte(param.CurrentPassword))

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

	supplier = models.SupplierModel{
		ID:       param.ID,
		Password: param.NewPassword,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}
	err = supplier.PasswordUpdate(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "PasswordUpdate/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	updatePasswordResponse := models.UpdatePasswordResponse{
		Message: "Password Successfully Changed",
	}

	return updatePasswordResponse, nil

}

func (s SupplierModule) Detail(ctx context.Context, param SupplierDetailParam) (interface{}, *helpers.Error) {
	supplier, err := models.GetOneSupplier(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneSupplier", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return supplier.Response(), nil
}

func (s SupplierModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	suppliers, err := models.GetAllSupplier(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllSupplier", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var supplierResponse []models.SupplierResponse
	for _, supplier := range suppliers {
		supplierResponse = append(supplierResponse, supplier.Response())
	}

	return supplierResponse, nil
}

func (s SupplierModule) Add(ctx context.Context, param SupplierAddParam) (interface{}, *helpers.Error) {

	password := util.RandomString(12)

	supplier := models.SupplierModel{
		Name:      param.Name,
		PhoneNo:   param.PhoneNo,
		Email:     param.Email,
		Password:  password,
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := supplier.Insert(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return supplier.Response(), nil
}

func (s SupplierModule) Update(ctx context.Context, param SupplierUpdateParam) (interface{}, *helpers.Error) {

	supplier := models.SupplierModel{
		ID:      uuid.FromStringOrNil(ctx.Value("user_id").(string)),
		Name:    param.Name,
		PhoneNo: param.PhoneNo,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := supplier.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return supplier.Response(), nil

}

func (s SupplierModule) Delete(ctx context.Context, param SupplierDeleteParam) (interface{}, *helpers.Error) {

	supplier := models.SupplierModel{
		ID: param.ID,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err := supplier.Delete(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Delete/Delete", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return nil, nil

}
