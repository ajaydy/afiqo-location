package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerSupplierList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSupplierList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return supplierService.List(ctx, filter)
}

func HandlerSupplierDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	supplierID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSupplierDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.SupplierDetailParam{ID: supplierID}

	return supplierService.Detail(ctx, param)
}

func HandlerSupplierUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.SupplierUpdateParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerSupplierUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return supplierService.Update(ctx, param)
}

func HandlerSupplierAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.SupplierAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSupplierAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return supplierService.Add(ctx, param)
}

func HandlerSupplierLogin(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.SupplierLoginParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSupplierLogin/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}
	return supplierService.Login(ctx, param)
}

func HandlerSupplierDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	supplierID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerSupplierDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.SupplierDeleteParam

	param.ID = supplierID

	return supplierService.Delete(ctx, param)
}

func HandlerSupplierPasswordUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	supplierID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	var param api.PasswordUpdateParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler",
			"HandlerSupplierPasswordUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = supplierID

	return supplierService.PasswordUpdate(ctx, param)
}
