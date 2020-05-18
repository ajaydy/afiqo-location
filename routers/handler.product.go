package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerProductList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProductList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return productService.List(ctx, filter)
}

func HandlerProductListBySupplierID(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProductListBySupplierID/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	supplierID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	param := api.SupplierDataParam{ID: supplierID}

	return productService.ListBySupplierID(ctx, filter, param)
}

func HandlerProductListForCustomer(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProductListForCustomer/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.ForCustomerParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProductListForCustomer/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return productService.ListForCustomer(ctx, filter, param)
}

func HandlerProductDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	productID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProductDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.ProductDetailParam{ID: productID}

	return productService.Detail(ctx, param)
}

func HandlerProductAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.ProductAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProductAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return productService.Add(ctx, param)
}

func HandlerProductUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	productID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProductUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.ProductUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerProductUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = productID

	return productService.Update(ctx, param)
}

func HandlerProductDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	productID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerProductDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.ProductDeleteParam

	param.ID = productID

	return productService.Delete(ctx, param)
}
