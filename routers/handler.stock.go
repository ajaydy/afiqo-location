package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerStockList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStockList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return stockService.List(ctx, filter)
}

func HandlerStockListBySupplierID(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStockListBySupplierID/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	supplierID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	param := api.SupplierDataParam{ID: supplierID}

	return stockService.ListBySupplierID(ctx, filter, param)
}

func HandlerStockDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	stockID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStockDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.StockDetailParam{ID: stockID}

	return stockService.Detail(ctx, param)
}

func HandlerStockAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.StockAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStockAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return stockService.Add(ctx, param)
}

func HandlerStockUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	stockID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStockUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.StockUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerStockUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = stockID

	return stockService.Update(ctx, param)
}
