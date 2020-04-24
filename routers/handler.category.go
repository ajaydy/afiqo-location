package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerCategoryList(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	filter, err := helpers.ParseFilter(ctx, r)

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCategoryList/parseFilter",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}
	return categoryService.List(ctx, filter)
}

func HandlerCategoryDetail(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	categoryID, err := uuid.FromString(params["id"])
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCategoryDetail/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	param := api.CategoryDetailParam{ID: categoryID}

	return categoryService.Detail(ctx, param)
}

func HandlerCategoryAdd(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.CategoryAddParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCategoryAdd/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return categoryService.Add(ctx, param)
}

func HandlerCategoryUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	categoryID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCategoryUpdate/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.CategoryUpdateParam

	err = helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerCategoryUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = categoryID

	return categoryService.Update(ctx, param)
}

func HandlerCategoryDelete(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	params := mux.Vars(r)

	categoryID, err := uuid.FromString(params["id"])

	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerCategoryDelete/parseID",
			helpers.BadRequestMessage, http.StatusBadRequest)
	}

	var param api.CategoryDeleteParam

	param.ID = categoryID

	return categoryService.Delete(ctx, param)
}
