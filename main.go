package main

import (
	"context"
	"fmt"
	"gl-farming/app/constants/files"
	mongoParams "gl-farming/app/constants/mongo"
	"gl-farming/app/controllers"
	"gl-farming/app/services"
	"gl-farming/database"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "gl-farming/docs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// @title Gipsyland Farming
// @version 2.0
// @description Farming service API description.
func main() {

	ctx := context.TODO()

	connectionAddress := mongoParams.GetConnectionString()
	// connectionAddress := mongoParams.AzureProdAddress
	mongoConnection := options.Client().ApplyURI(connectionAddress)
	client, err := mongo.Connect(ctx, mongoConnection)
	defer client.Disconnect(ctx)

	if err != nil {
		log.Fatal(err.Error())
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println("Connection established")

	var Collections = database.Collections{
		AccountRequests: client.Database("gypsyland").Collection("accountRequests"),
		AccountTypes:    client.Database("gypsyland").Collection("accountTypes"),
		Locations:       client.Database("gypsyland").Collection("locations"),
		Currency:        client.Database("gypsyland").Collection("currency"),
		FarmerAccess:    client.Database("gypsyland").Collection("farmerAccess"),
	}

	var appServices services.AppServices
	appServices.Init(&Collections)
	var appControllers = controllers.NewAppControllers(appServices)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{"*"},
		Skipper:          middleware.DefaultSkipper,
		AllowCredentials: true,
	}))
	SetRoutes(e, appControllers)
	e.Static(files.Static, "static")

	e.Logger.Fatal(e.Start(":80"))

}

func checkRoot(c echo.Context) error {
	curDir, _ := os.Getwd()
	files, err := ioutil.ReadDir(curDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		c.String(http.StatusAccepted, file.Name()+"\n")
	}

	return c.String(http.StatusAccepted, "complete")

}

func checkRootFiles(c echo.Context) error {
	curDir, _ := os.Getwd()
	files, err := ioutil.ReadDir(curDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		c.String(http.StatusAccepted, file.Name()+"\n")
	}

	return c.String(http.StatusAccepted, "complete")

}

func SetRoutes(e *echo.Echo, appControllers controllers.AppControllers) {

	root := e.Group("/v2")
	root.GET("/swagger/*", echoSwagger.WrapHandler)
	root.GET("/", appControllers.Files.DownloadPage)
	root.GET("/checkRoot", checkRoot)

	auth := root.Group("/auth")
	auth.POST("/login", appControllers.UID.Login)
	auth.GET("/uid", appControllers.UID.GetUID)
	auth.GET("/uid/admin", appControllers.UID.GetServiceAdminUID)
	auth.GET("/check", appControllers.UID.Check)

	files := root.Group("/files")
	files.POST("/upload", appControllers.Files.Upload)
	filesDownload := files.Group("/download")
	filesDownload.GET("/file", appControllers.Files.Download)
	filesDownload.GET("/inline", appControllers.Files.DownloadInline)
	filesDownload.GET("/attachment", appControllers.Files.DownloadAttachment)

	accountTypes := root.Group("/accountTypes")
	accountTypes.POST("/create", appControllers.AccountTypes.Create)
	accountTypesGet := accountTypes.Group("/get")
	accountTypesGet.GET("", appControllers.AccountTypes.Get)
	accountTypesGet.GET("/all", appControllers.AccountTypes.GetAll)
	accountTypesDelete := accountTypes.Group("/delete")
	accountTypesDelete.DELETE("", appControllers.AccountTypes.Delete)
	accountTypesDelete.DELETE("/all", appControllers.AccountTypes.DeleteAll)

	farmerAccess := root.Group("/farmerAccess")
	farmerAccess.POST("/add", appControllers.FarmerAccessController.Add)
	farmerAccess.PUT("/revoke", appControllers.FarmerAccessController.Revoke)
	farmerAccess.PUT("/add/all", appControllers.FarmerAccessController.FullAccess)
	farmerAccess.PUT("/revoke/all", appControllers.FarmerAccessController.FullRevoke)
	farmerAccessGet := farmerAccess.Group("/get")
	farmerAccessGet.GET("/teams", appControllers.FarmerAccessController.GetTeams)
	farmerAccessGet.GET("/farmers", appControllers.FarmerAccessController.GetFarmers)
	farmerAccessGet.GET("/all", appControllers.FarmerAccessController.GetAll)
	farmerAccessDelete := farmerAccess.Group("/delete")
	farmerAccessDelete.DELETE("/all", appControllers.FarmerAccessController.DeleteAll)

	tableData := root.Group("/tableData")
	tableDataGet := tableData.Group("/get")
	tableDataGet.GET("", appControllers.Tables.Get)
	tableDataGet.GET("/all", appControllers.Tables.GetAll)
	tableDataAggregate := tableData.Group("/aggregate")
	tableDataAggregate.GET("/farmers", appControllers.Tables.FarmerPipeline)
	tableDataAggregate.GET("/buyers", appControllers.Tables.BuyerPipiline)
	tableDataAggregate.GET("/teamleads", appControllers.Tables.TeamleadPipiline)
	tableDataTeamlead := tableData.Group("/teamlead")
	tableDataTeamlead.POST("/get", appControllers.Tables.GetTeamleadTables)

	currency := root.Group("/currency")
	currency.POST("/create", appControllers.Currency.Create)
	currencyGet := currency.Group("/get")
	currencyGet.GET("", appControllers.Currency.Get)
	currencyGet.GET("/all", appControllers.Currency.GetAll)
	currencyDelete := currency.Group("/delete")
	currencyDelete.DELETE("", appControllers.Currency.Delete)
	currencyDelete.DELETE("/all", appControllers.Currency.DeleteAll)

	locations := root.Group("/locations")
	locations.POST("/create", appControllers.Locations.Create)
	locationsGet := locations.Group("/get")
	locationsGet.GET("", appControllers.Locations.Get)
	locationsGet.GET("/all", appControllers.Locations.GetAll)
	locationsDelete := locations.Group("/delete")
	locationsDelete.DELETE("", appControllers.Locations.Delete)
	locationsDelete.DELETE("/all", appControllers.Locations.DeleteAll)

	accountRequests := root.Group("/accountRequests")
	accountRequests.POST("/create", appControllers.AccountRequests.Create)
	accountRequests.PUT("/take", appControllers.AccountRequests.Take)
	accountRequests.PUT("/cancel", appControllers.AccountRequests.Cancel)
	accountRequests.PUT("/update", appControllers.AccountRequests.Update)
	accountRequests.PUT("/complete", appControllers.AccountRequests.Complete)
	accountRequests.GET("/get", appControllers.AccountRequests.Get)
	accountRequests.PUT("/return", appControllers.AccountRequests.Return)
	accountRequestsDelete := accountRequests.Group("/delete")
	accountRequestsDelete.DELETE("/all", appControllers.AccountRequests.DeleteAll)

}
