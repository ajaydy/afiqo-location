package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerCourierList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCourierList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return courierService.List(ctx, filter)
}

func HandlerCourierDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	courierID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCourierDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.CourierDetailParam{ID: courierID}

	return courierService.Detail(ctx, param)
}

func HandlerCourierAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.CourierAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCourierAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return courierService.Add(ctx, param)
}

func HandlerCourierUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	courierID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCourierUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.CourierUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerCourierUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = courierID

	return courierService.Update(ctx, param)
}

func HandlerCourierLogin(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.CourierLoginParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCourierLogin/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}
	return courierService.Login(ctx, param)
}

func HandlerCourierDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	courierID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCourierDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.CourierDeleteParam

	param.ID = courierID

	return courierService.Delete(ctx, param)
}

func HandlerCourierPasswordUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	courierID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	var param api.PasswordUpdateParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler",
			"HandlerCourierPasswordUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = courierID

	return courierService.PasswordUpdate(ctx, param)
}
