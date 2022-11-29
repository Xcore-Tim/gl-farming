package main

import (
	"gl-farming/app/controllers"
	"gl-farming/app/services"
	"gl-farming/database"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type AppControllers struct {
	Locations controllers.LocationController
}

type AppServices struct {
	Locations services.LocationService
}

var (
	appControllers = AppControllers{
		Locations: controllers.NewLocationController(),
	}
)

func (ac *AppControllers) InitServices(collections *database.Collections) {
	ac.Locations.Service = services.NewLocationService(collections.Locations)
}

// @title Echo Swagger Example API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {

	dbCollections, err := database.Init()

	if err != nil {
		log.Fatal("asdasd")
	}

	appControllers.InitServices(dbCollections)

	//New Echo instance
	e := echo.New()

	//Middleware usage
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	SetRoutes(e)

	e.Logger.Fatal(e.Start(":80"))

}

func SetRoutes(e *echo.Echo) {

	root := e.Group("/v2")

	locations := root.Group("/locations")
	locations.POST("/create", appControllers.Locations.Create)

	locationsGet := locations.Group("/get")
	locationsGet.GET("/all", appControllers.Locations.GetAll)

	locationsDelete := locations.Group("/delete")
	locationsDelete.DELETE("/all", appControllers.Locations.DeleteAll)
}

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}
