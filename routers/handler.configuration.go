package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"net/http"
)

func HandlerConfigurationUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.ConfigurationUpdateParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler", "HandlerSupplierUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	return configurationService.Update(ctx, param)
}
