package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerOrderList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerOrderList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return orderService.List(ctx, filter)
}

func HandlerOrderListByCustomerID(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerOrderListByCustomerID/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	customerID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	param := api.CustomerDataParam{ID: customerID}

	return orderService.ListByCustomerID(ctx, filter, param)
}

func HandlerOrderDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	orderID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerOrderDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.OrderDetailParam{ID: orderID}

	return orderService.Detail(ctx, param)
}

func HandlerOrder(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.OrderParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerOrder/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return orderService.Order(ctx, param)
}

func HandlerOrderDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	orderID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerOrderDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.OrderDeleteParam

	param.ID = orderID

	return orderService.Delete(ctx, param)
}
