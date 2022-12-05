package controllers

import (
	"errors"
	"gl-farming/app/models"
	"gl-farming/app/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type LocationController struct {
	Service services.LocationService
}

func NewLocationController() LocationController {
	return LocationController{}
}

// Create godoc
// @Summary      Create location
// @Description  creates location
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        id    body     models.Location  true  "location body json"
// @Success      200  {string}   string
// @Router       /v2/locations/create [post]
func (ctrl LocationController) Create(c echo.Context) error {

	var location models.Location

	if err := c.Bind(&location); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := ctrl.Service.Create(c, location); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// Get godoc
// @Summary      Update location
// @Description  Updates location by id
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        id    body     models.Location  true  "location body"
// @Success      200  {string}  string
// @Router       /v2/locations/update [patch]
func (ctrl LocationController) Update(c echo.Context) error {
	var location models.Location

	if err := c.Bind(&location); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	location.ConvertToMongoID()

	if err := ctrl.Service.Update(c, &location); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.String(http.StatusBadRequest, "success")
}

// Get godoc
// @Summary      Delete location
// @Description  Deletes location by id
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "location id"
// @Success      200  {string}  string
// @Router       /v2/locations/delete [delete]
func (ctrl LocationController) Delete(c echo.Context) error {

	id := c.QueryParam("id")

	if err := ctrl.Service.Delete(c, id); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "successfuly deleted location")
}

// Get godoc
// @Summary      Get location
// @Description  get location by id
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "location id"
// @Success      200  {object}   models.Location
// @Router       /v2/locations/get [get]
func (ctrl LocationController) Get(c echo.Context) error {

	id := c.QueryParam("id")

	location, err := ctrl.Service.Get(c, id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, location)
}

// GetAll godoc
// @Summary      Get all locations
// @Description  returns all locations
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Location
// @Router       /v2/locations/get/all [get]
func (ctrl LocationController) GetAll(c echo.Context) error {

	var locations []models.Location

	if err := ctrl.Service.GetAll(c, &locations); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(locations) == 0 {
		return c.JSON(http.StatusOK, errors.New("no locations found"))
	}

	return c.JSON(http.StatusOK, locations)
}

// DeleteAll godoc
// @Summary      Delete all locations
// @Description  deletes all locations
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /v2/locations/delete/all [delete]
func (ctrl LocationController) DeleteAll(c echo.Context) error {

	count, err := ctrl.Service.DeleteAll(c)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, count)
}
