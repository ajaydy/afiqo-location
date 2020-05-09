package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"html"
	"net/http"
	"reflect"
	"strings"
)

var decoder = schema.NewDecoder()
var validate *validator.Validate

type (
	FilterOption struct {
		Limit  int    `json:"limit" schema:"limit"`
		Offset int    `json:"offset" schema:"offset"`
		Search string `json:"search" schema:"search"`
		Dir    string `json:"dir" schema:"dir"`
	}

	Filter struct {
		FilterOption `json:"filter,omitempty"`
		ID           uuid.UUID       `json:"id" schema:"id"`
		Longitude    decimal.Decimal `json:"longitude" schema:"longitude"`
		Latitude     decimal.Decimal `json:"latitude" schema:"latitude"`
		SupplierID   uuid.UUID       `json:"supplier_id" schema:"supplier_id"`
		CategoryID   uuid.UUID       `json:"category_id" schema:"category_id"`
		CustomerID   uuid.UUID       `json:"customer_id" schema:"customer_id"`
		CourierID    uuid.UUID       `json:"courier_id" schema:"courier_id"`
		OrderID      uuid.UUID       `json:"order_id" schema:"order_id"`
		ProductID    uuid.UUID       `json:"product_id" schema:"product_id"`
		WarehouseID  uuid.UUID       `json:"warehouse_id" schema:"warehouse_id"`
	}
)

func ParseBodyRequestData(ctx context.Context, r *http.Request, data interface{}) error {

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return err
	}

	value := reflect.ValueOf(data).Elem()
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		if field.Type() != reflect.TypeOf("") {
			continue
		}
		str := field.Interface().(string)
		field.SetString(html.EscapeString(str))

	}
	validate = validator.New()
	err = validate.Struct(data)
	//validationErrors := err.(validator.ValidationErrors)

	if err != nil {
		return err
	}

	return nil

}

func ParseFilter(ctx context.Context, r *http.Request) (Filter, error) {
	marshal, _ := json.Marshal(r.URL.Query())
	fmt.Println(string(marshal))
	var filter Filter
	err := decoder.Decode(&filter, r.URL.Query())
	if err != nil {
		return filter, nil
	}

	if strings.ToLower(filter.Dir) != "asc" && strings.ToLower(filter.Dir) != "desc" {
		filter.Dir = "ASC"
	}

	return filter, nil
}
