package api

import (
	"afiqo-location/email"
	"afiqo-location/helpers"
	"afiqo-location/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type (
	PaymentModule struct {
		db     *sql.DB
		cache  *redis.Pool
		logger *helpers.Logger
		name   string
	}

	PaymentUpdateParam struct {
		ID uuid.UUID `json:"id"`
	}

	PaymentDetailParam struct {
		ID uuid.UUID `json:"id"`
	}
)

func NewPaymentModule(db *sql.DB, cache *redis.Pool, logger *helpers.Logger) *PaymentModule {
	return &PaymentModule{
		db:     db,
		cache:  cache,
		logger: logger,
		name:   "module/payment",
	}
}
func (s PaymentModule) Detail(ctx context.Context, param PaymentDetailParam) (interface{}, *helpers.Error) {
	payment, err := models.GetOnePayment(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/GetOnePayment", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := payment.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Detail/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil
}

func (s PaymentModule) List(ctx context.Context, filter helpers.Filter) (interface{}, *helpers.Error) {

	payments, err := models.GetAllPayment(ctx, s.db, filter)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "List/GetAllPayment", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var paymentResponse []models.PaymentResponse
	for _, payment := range payments {
		response, err := payment.Response(ctx, s.db, s.logger)

		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "List/Response", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		paymentResponse = append(paymentResponse, response)
	}

	return paymentResponse, nil
}

func (s PaymentModule) Update(ctx context.Context, param PaymentUpdateParam) (interface{}, *helpers.Error) {

	payment, err := models.GetOnePayment(ctx, s.db, param.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetOnePayment", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	order, err := models.GetOneOrder(ctx, s.db, payment.OrderID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetOneOrder", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	if order.CustomerID != uuid.FromStringOrNil(ctx.Value("user_id").(string)) {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetOneOrder", helpers.OrderErrorMessage,
			http.StatusInternalServerError)
	}

	payment = models.PaymentModel{
		ID:     payment.ID,
		Status: 1,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err = payment.Update(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Update", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	orderUpdate := models.OrderModel{
		ID:     order.ID,
		Status: 1,
		UpdatedBy: uuid.NullUUID{
			UUID:  uuid.FromStringOrNil(ctx.Value("user_id").(string)),
			Valid: true,
		},
	}

	err = orderUpdate.UpdateStatus(ctx, s.db)
	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/UpdateStatus", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	customer, err := models.GetOneCustomer(ctx, s.db, order.CustomerID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetOneCustomer", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	orderProducts, err := models.GetAllOrderProductByOrderID(ctx, s.db, order.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetAllOrderProductByOrderID", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	var entries [][]email.Entry
	for _, orderProduct := range orderProducts {
		var columns []email.Entry
		product, err := models.GetOneProduct(ctx, s.db, orderProduct.ProductID)
		if err != nil {
			return nil, helpers.ErrorWrap(err, s.name, "Update/GetOneProduct", helpers.InternalServerError,
				http.StatusInternalServerError)
		}
		column := email.Entry{
			Key:   "Item",
			Value: product.Name,
		}
		column2 := email.Entry{
			Key:   "Description",
			Value: product.Description,
		}
		column3 := email.Entry{
			Key:   "Quantity",
			Value: fmt.Sprintf("%d", orderProduct.Quantity),
		}
		column4 := email.Entry{
			Key:   "Price",
			Value: fmt.Sprintf("RM %s", product.Price),
		}

		column5 := email.Entry{
			Key:   "Subtotal",
			Value: fmt.Sprintf("RM %s", orderProduct.SubTotal),
		}

		columns = append(columns, column, column2, column3, column4, column5)
		entries = append(entries, columns)
	}

	data := email.MailData{
		Name:  customer.Name,
		Entry: entries,
	}

	body, err := data.GenerateForReceipt()

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GenerateForReceipt", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	mail := email.Mail{
		Subject: "Order Processing",
		Body:    body,
		To:      "mohdjamilafiq@gmail.com",
	}

	go func() {
		mail.SendEmail()
	}()

	payment, err = models.GetOnePayment(ctx, s.db, payment.ID)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/GetOnePayment", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	response, err := payment.Response(ctx, s.db, s.logger)

	if err != nil {
		return nil, helpers.ErrorWrap(err, s.name, "Update/Response", helpers.InternalServerError,
			http.StatusInternalServerError)
	}

	return response, nil

}
