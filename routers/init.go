package routers

import (
	"afiqo-location/api"
	"afiqo-location/helpers"
	"database/sql"
	"github.com/gomodule/redigo/redis"
)

var (
	dbPool                  *sql.DB
	cachePool               *redis.Pool
	logger                  *helpers.Logger
	customerService         *api.CustomerModule
	adminService            *api.AdminModule
	supplierService         *api.SupplierModule
	courierService          *api.CourierModule
	categoryService         *api.CategoryModule
	productService          *api.ProductModule
	orderService            *api.OrderModule
	orderProductService     *api.OrderProductModule
	paymentService          *api.PaymentModule
	warehouseService        *api.WarehouseModule
	warehouseProductService *api.WarehouseProductModule
	shipmentService         *api.ShipmentModule
	configurationService    *api.ConfigurationModule
)

func Init(db *sql.DB, cache *redis.Pool, log *helpers.Logger) {
	dbPool = db
	cachePool = cache
	logger = log
	customerService = api.NewCustomerModule(dbPool, cachePool, logger)
	adminService = api.NewAdminModule(dbPool, cachePool, logger)
	supplierService = api.NewSupplierModule(dbPool, cachePool, logger)
	courierService = api.NewCourierModule(dbPool, cachePool, logger)
	categoryService = api.NewCategoryModule(dbPool, cachePool, logger)
	productService = api.NewProductModule(dbPool, cachePool, logger)
	orderService = api.NewOrderModule(dbPool, cachePool, logger)
	orderProductService = api.NewOrderProductModule(dbPool, cachePool, logger)
	paymentService = api.NewPaymentModule(dbPool, cachePool, logger)
	warehouseService = api.NewWarehouseModule(dbPool, cachePool, logger)
	warehouseProductService = api.NewWarehouseProductModule(dbPool, cachePool, logger)
	shipmentService = api.NewShipmentModule(dbPool, cachePool, logger)
	configurationService = api.NewConfigurationModule(dbPool, cachePool, logger)
}
