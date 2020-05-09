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

//func HandlerShipmentUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
//
//	ctx := r.Context()
//
//	params := mux.Vars(r)
//
//	shipmentID, err := uuid.FromString(params["id"])
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerShipmentUpdate/parseID",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//	}
//
//	var param api.ShipmentUpdateParam
//
//	err = helpers.ParseBodyRequestData(ctx, r, &param)
//	if err != nil {
//
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerShipmentUpdate/ParseBodyRequestData",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//
//	}
//
//	param.ID = shipmentID
//
//	return shipmentService.Update(ctx, param)
//}
//
//func HandlerShipmentDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {
//
//	ctx := r.Context()
//
//	params := mux.Vars(r)
//
//	shipmentID, err := uuid.FromString(params["id"])
//
//	if err != nil {
//		return nil, helpers.ErrorWrap(err, "handler", "HandlerShipmentDelete/parseID",
//			helpers.BadRequestMessage, http.StatusBadRequest)
//	}
//
//	var param api.ShipmentDeleteParam
//
//	param.ID = shipmentID
//
//	return shipmentService.Delete(ctx, param)
//}
