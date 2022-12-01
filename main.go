package main

import (
	"gl-farming/app/controllers"
	"gl-farming/app/services"
	"gl-farming/database"
	"log"

	_ "gl-farming/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type AppControllers struct {
	Locations    controllers.LocationController
	AccountTypes controllers.AccountTypeController
	Currency     controllers.CurrencyController
}

var appControllers = AppControllers{
	Locations:    controllers.NewLocationController(),
	AccountTypes: controllers.NewAccountTypeController(),
	Currency:     controllers.NewCurrencyController(),
}

func (ac *AppControllers) InitServices(collections *database.Collections) {
	ac.Locations.Service = services.NewLocationService(collections.Locations)
	ac.AccountTypes.Service = services.NewAccountTypeService(collections.AccountTypes)
	ac.Currency.Service = services.NewCurrencyService(collections.Currency)
}

// @title Gipsyland Farming
// @version 2.0
// @description Farming service API description.
func main() {

	dbCollections, err := database.Init()

	if err != nil {
		log.Fatal("asdasd")
	}

	appControllers.InitServices(dbCollections)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	SetRoutes(e)

	e.Logger.Fatal(e.Start(":80"))

}

func SetRoutes(e *echo.Echo) {

	root := e.Group("/v2")
	root.GET("/swagger/*", echoSwagger.WrapHandler)

	locations := root.Group("/locations")
	locations.POST("/create", appControllers.Locations.Create)
	locationsGet := locations.Group("/get")
	locationsGet.GET("", appControllers.Locations.Get)
	locationsGet.GET("/all", appControllers.Locations.GetAll)
	locationsDelete := locations.Group("/delete")
	locationsDelete.DELETE("", appControllers.Locations.Delete)
	locationsDelete.DELETE("/all", appControllers.Locations.DeleteAll)

	accountTypes := root.Group("/accountTypes")
	accountTypes.POST("/create", appControllers.AccountTypes.Create)
	accountTypesGet := accountTypes.Group("/get")
	accountTypesGet.GET("", appControllers.AccountTypes.Get)
	accountTypesGet.GET("/all", appControllers.AccountTypes.GetAll)
	accountTypesDelete := accountTypes.Group("/delete")
	accountTypesDelete.DELETE("", appControllers.AccountTypes.Delete)
	accountTypesDelete.DELETE("/all", appControllers.AccountTypes.DeleteAll)

	currency := root.Group("/currency")
	currency.POST("/create", appControllers.Currency.Create)
	currencyGet := currency.Group("/get")
	currencyGet.GET("", appControllers.Currency.Get)
	currencyGet.GET("/all", appControllers.Currency.GetAll)
	currencyDelete := currency.Group("/delete")
	currencyDelete.DELETE("", appControllers.Currency.Delete)
	currencyDelete.DELETE("/all", appControllers.Currency.DeleteAll)

}
