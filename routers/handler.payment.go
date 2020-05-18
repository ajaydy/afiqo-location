package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerPaymentList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerPaymentList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return paymentService.List(ctx, filter)
}

func HandlerPaymentDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	paymentID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerPaymentDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.PaymentDetailParam{ID: paymentID}

	return paymentService.Detail(ctx, param)
}

func HandlerPaymentUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	paymentID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerPaymentUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.PaymentUpdateParam{ID: paymentID}

	return paymentService.Update(ctx, param)
}
