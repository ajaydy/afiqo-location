package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerWarehouseList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerWarehouseList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return warehouseService.List(ctx, filter)
}

func HandlerWarehouseDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	warehouseID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerWarehouseDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.WarehouseDetailParam{ID: warehouseID}

	return warehouseService.Detail(ctx, param)
}

func HandlerWarehouseDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	warehouseID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerWarehouseDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.WarehouseDeleteParam

	param.ID = warehouseID

	return warehouseService.Delete(ctx, param)
}
