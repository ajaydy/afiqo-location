package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerOrderProductList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerOrderProductList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return orderProductService.List(ctx, filter)
}

func HandlerOrderProductDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	orderProductID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerOrderProductDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.OrderProductDetailParam{ID: orderProductID}

	return orderProductService.Detail(ctx, param)
}
