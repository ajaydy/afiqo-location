package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerCustomerList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCustomerList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return customerService.List(ctx, filter)
}

func HandlerCustomerDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	customerID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCustomerDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.CustomerDetailParam{ID: customerID}

	return customerService.Detail(ctx, param)
}

func HandlerCustomerAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.CustomerAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCustomerAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return customerService.Add(ctx, param)
}

func HandlerCustomerUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	customerID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCustomerUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.CustomerUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerCustomerUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = customerID

	return customerService.Update(ctx, param)
}

func HandlerCustomerRegister(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.CustomerRegisterParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCustomerRegister/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return customerService.Register(ctx, param)
}

func HandlerCustomerLogout(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	return customerService.Logout(ctx, r.Header.Get("session"))
}

func HandlerCustomerLogin(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.CustomerLoginParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCustomerLogin/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}
	return customerService.Login(ctx, param)
}

func HandlerCustomerDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	customerID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCustomerDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.CustomerDeleteParam

	param.ID = customerID

	return customerService.Delete(ctx, param)
}

func HandlerCustomerPasswordUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	customerID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	var param api.PasswordUpdateParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler",
			"HandlerCustomerPasswordUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = customerID

	return customerService.PasswordUpdate(ctx, param)
}
