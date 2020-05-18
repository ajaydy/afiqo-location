package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerShipmentList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerShipmentList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return shipmentService.List(ctx, filter)
}

func HandlerShipmentListByCourierID(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerShipmentListByCourierID/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	courierID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	param := api.CourierDataParam{ID: courierID}

	return shipmentService.ListByCourierID(ctx, filter, param)
}

func HandlerShipmentListByCustomerID(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerShipmentListByCustomerID/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	customerID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	param := api.CustomerDataParam{ID: customerID}

	return shipmentService.ListByCustomerID(ctx, filter, param)
}

func HandlerShipmentDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	shipmentID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerShipmentDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.ShipmentDetailParam{ID: shipmentID}

	return shipmentService.Detail(ctx, param)
}

func HandlerShipmentAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.ShipmentAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerShipmentAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	return shipmentService.Add(ctx, param)
}
