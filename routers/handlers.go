package routers

import (
	"afiqo-location/helpers"
	"afiqo-location/middleware"
	"afiqo-location/session"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type (
	HandlerFunc func(http.ResponseWriter, *http.Request) (interface{}, *helpers.Error)
)

func (fn HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var errs []string
	r.ParseForm()
	data, err := fn(w, r)
	if err != nil {
		errs = append(errs, err.Error())
		w.WriteHeader(err.StatusCode)
	}
	resp := helpers.Response{
		Data: data,
		BaseResponse: helpers.BaseResponse{
			Errors: errs,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		return
	}
}

func InitHandlers() *mux.Router {
	r := mux.NewRouter()

	http.Handle("/", r)

	apiV1 := r.PathPrefix("/api/v1").Subrouter()
	apiV1.Use(middleware.LoggingMiddleware)

	apiV1.Handle("/customers", middleware.SessionMiddleware(
		HandlerFunc(HandlerCustomerList))).Methods(http.MethodGet)
	apiV1.Handle("/customers/{id}", middleware.SessionMiddleware(
		HandlerFunc(HandlerCustomerDetail))).Methods(http.MethodGet)
	apiV1.Handle("/customers/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerCustomerUpdate), session.CUSTOMER_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/customers/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerCustomerDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)
	apiV1.Handle("/customer/password-update", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerCustomerPasswordUpdate), session.CUSTOMER_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/customer/login", HandlerFunc(HandlerCustomerLogin)).Methods(http.MethodPost)
	apiV1.Handle("/customer/register", HandlerFunc(HandlerCustomerRegister)).Methods(http.MethodPost)

	apiV1.Handle("/suppliers", middleware.SessionMiddleware(
		HandlerFunc(HandlerSupplierList))).Methods(http.MethodGet)
	apiV1.Handle("/suppliers/{id}", middleware.SessionMiddleware(
		HandlerFunc(HandlerSupplierDetail))).Methods(http.MethodGet)
	apiV1.Handle("/suppliers/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerSupplierUpdate), session.SUPPLIER_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/suppliers/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerSupplierDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)
	apiV1.Handle("/suppliers/password-update", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerSupplierPasswordUpdate), session.SUPPLIER_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/supplier/login", HandlerFunc(HandlerSupplierLogin)).Methods(http.MethodPost)
	apiV1.Handle("/supplier/register", HandlerFunc(HandlerSupplierRegister)).Methods(http.MethodPost)

	apiV1.Handle("/couriers", middleware.SessionMiddleware(
		HandlerFunc(HandlerCourierList))).Methods(http.MethodGet)
	apiV1.Handle("/couriers/{id}", middleware.SessionMiddleware(
		HandlerFunc(HandlerCourierDetail))).Methods(http.MethodGet)
	apiV1.Handle("/couriers", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerCourierAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/couriers/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerCourierUpdate), session.COURIER_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/couriers/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerCourierDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)
	apiV1.Handle("/courier/password-update", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerCourierPasswordUpdate), session.COURIER_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/courier/login", HandlerFunc(HandlerCourierLogin)).Methods(http.MethodPost)

	apiV1.Handle("/categories", middleware.SessionMiddleware(
		HandlerFunc(HandlerCategoryList))).Methods(http.MethodGet)
	apiV1.Handle("/categories/{id}", middleware.SessionMiddleware(
		HandlerFunc(HandlerCategoryDetail))).Methods(http.MethodGet)
	apiV1.Handle("/categories", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerCategoryAdd), session.ADMIN_ROLE))).Methods(http.MethodPost)
	apiV1.Handle("/categories/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerCategoryUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/categories/{id}", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerCategoryDelete), session.ADMIN_ROLE))).Methods(http.MethodDelete)

	apiV1.Handle("/admin/password-update", middleware.SessionMiddleware(middleware.RolesMiddleware(
		HandlerFunc(HandlerAdminPasswordUpdate), session.ADMIN_ROLE))).Methods(http.MethodPut)
	apiV1.Handle("/admin/login", HandlerFunc(HandlerAdminLogin)).Methods(http.MethodPost)

	return r
}
