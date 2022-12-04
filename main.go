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

// @title Gipsyland Farming
// @version 2.0
// @description Farming service API description.
func main() {

	dbCollections, err := database.Init()

	if err != nil {
		log.Fatal("asdasd")
	}

	var appServices services.AppServices
	appServices.Init(dbCollections)
	var appControllers = controllers.NewAppControllers(appServices)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	SetRoutes(e, appControllers)

	e.Logger.Fatal(e.Start(":80"))

}

func SetRoutes(e *echo.Echo, appControllers controllers.AppControllers) {

	root := e.Group("/v2")
	root.GET("/swagger/*", echoSwagger.WrapHandler)

	auth := root.Group("/auth")
	auth.GET("/uid", appControllers.UID.GetUID)
	auth.GET("/uid/admin", appControllers.UID.GetServiceAdminUID)

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

	accountRequests := root.Group("/accountRequests")
	accountRequests.POST("/create", appControllers.AccountRequests.Create)
	accountRequests.PUT("/take", appControllers.AccountRequests.Take)
	accountRequests.PUT("/cancel", appControllers.AccountRequests.Cancel)
	accountRequests.PUT("/update", appControllers.AccountRequests.Update)
	accountRequests.PUT("/complete", appControllers.AccountRequests.Complete)
	accountRequestsDelete := accountRequests.Group("/delete")
	accountRequestsDelete.DELETE("/all", appControllers.AccountRequests.DeleteAll)

	tableData := root.Group("/tableData")
	tableDataGet := tableData.Group("/get")
	tableDataGet.GET("", appControllers.Tables.Get)
	tableDataGet.GET("/all", appControllers.Tables.GetAll)

	farmerAccess := root.Group("/farmerAccess")
	farmerAccess.POST("/add", appControllers.FarmerAccessController.Add)
	farmerAccess.DELETE("/revoke", appControllers.FarmerAccessController.Revoke)
	farmerAccessGet := farmerAccess.Group("/get")
	farmerAccessGet.GET("/teams", appControllers.FarmerAccessController.GetTeams)
	farmerAccessGet.GET("/farmers", appControllers.FarmerAccessController.GetFarmers)
	farmerAccessGet.GET("/all", appControllers.FarmerAccessController.GetAll)
	farmerAccessGet.GET("/access", appControllers.FarmerAccessController.GetAccess)
	farmerAccessDelete := farmerAccess.Group("/delete")
	farmerAccessDelete.DELETE("/all", appControllers.FarmerAccessController.DeleteAll)
}
