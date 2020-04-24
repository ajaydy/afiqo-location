package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

func HandlerAdminLogin(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	var param api.AdminLoginParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {
		return nil, helpers.ErrorWrap(err, "handler", "HandlerAdminLogin/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}
	return adminService.Login(ctx, param)
}

func HandlerAdminPasswordUpdate(w http.ResponseWriter, r *http.Request) (interface{}, *helpers.Error) {

	ctx := r.Context()

	adminID := uuid.FromStringOrNil(ctx.Value("user_id").(string))

	var param api.PasswordUpdateParam

	err := helpers.ParseBodyRequestData(ctx, r, &param)
	if err != nil {

		return nil, helpers.ErrorWrap(err, "handler",
			"HandlerAdminPasswordUpdate/ParseBodyRequestData",
			helpers.BadRequestMessage, http.StatusBadRequest)

	}

	param.ID = adminID

	return adminService.PasswordUpdate(ctx, param)
}
