package api

import (
	"afiqo-location/helpers"
	"afiqo-location/models"
	"context"
	"database/sql"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type (
	ShipmentModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	ShipmentDetailParam struct {
		ID uuid.UUID `json:"id"`
	}

	CourierDataParam struct {
		ID uuid.UUID `json:"id"`
	}

	ShipmentAddParam struct {
		CourierID uuid.UUID `json:"courier_id" validate:"required"`
		OrderID   uuid.UUID `json:"order_id" validate:"required"`
		Status    int       `json:"status" validate:"required"`
	}

	ShipmentUpdateParam struct {
		ID     uuid.UUID `json:"id"`
		Status int       `json:"status"`
	}
)

func NewShipmentModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *ShipmentModule {
	return &ShipmentModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/shipment",
	}
}

func (s ShipmentModule) Detail(ctx context.Context, param ShipmentDetailParam) (interface{}, *helpers.Error) {
	shipment, err := models.GetOneShipment(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOneShipment", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := shipment.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s ShipmentModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	shipments, err := models.GetAllShipment(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllShipment", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var shipmentResponse []models.ShipmentResponse
	for _, shipment := range shipments {
		response, err := shipment.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		shipmentResponse = append(shipmentResponse, response)
	}

	return shipmentResponse, nil
}

func (s ShipmentModule) ListByCourierID(ctx context.Context, filter helpers.Filter, param CourierDataParam) (
	interface{}, *helpers.Error) {

	shipments, err := models.GetAllShipmentByCourierID(ctx, s.db, filter, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ListByCourierID/GetAllShipmentByCourierID", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var shipmentResponse []models.ShipmentResponse

	for _, shipment := range shipments {
		response, err := shipment.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "ListByCourierID/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		shipmentResponse = append(shipmentResponse, response)
	}

	return shipmentResponse, nil
}

func (s ShipmentModule) ListByCustomerID(ctx context.Context, filter helpers.Filter, param CustomerDataParam) (
	interface{}, *helpers.Error) {

	shipments, err := models.GetAllShipmentByCustomerID(ctx, s.db, filter, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "ListByCustomerID/GetAllShipmentByCustomerID", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var shipmentResponse []models.ShipmentResponse

	for _, shipment := range shipments {
		response, err := shipment.Response(ctx, s.db, s.logger)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "ListByCustomerID/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}

		shipmentResponse = append(shipmentResponse, response)
	}

	return shipmentResponse, nil
}

func (s ShipmentModule) Add(ctx context.Context, param ShipmentAddParam) (interface{}, *helpers.Error) {

	order, err := models.GetOneOrder(ctx, s.db, param.OrderID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneOrder", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	if order.Status != 1 {
		return nil, helpers.ErrorWrap(err, s.name, "Add/GetOneOrder", "Order Not Confirmed",
			http.StatusInternalServerError)
	}

	shipment := models.ShipmentModel{
		CourierID: param.CourierID,
		OrderID:   param.OrderID,
		Status:    param.Status,
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err = shipment.Insert(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := shipment.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}

func (s ShipmentModule) UpdateStatus(ctx context.Context, param ShipmentUpdateParam) (interface{}, *helpers.Error) {

	shipment := models.ShipmentModel{
		ID:        param.ID,
		Status:    param.Status,
		CreatedBy: uuid.FromStringOrNil(ctx.Value("user_id").(string)),
	}

	err := shipment.UpdateStatus(ctx, s.db)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Insert", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := shipment.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Add/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}
