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

func (ctrl LocationController) Create(c echo.Context) error {

	var location models.Location

	if err := c.Bind(&location); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	locationDTO := location.ToDTO()

	if err := ctrl.Service.Create(c, locationDTO); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, location)
}

func (ctrl LocationController) Update(c echo.Context) {
}

func (ctrl LocationController) Delete(c echo.Context) {
}

func (ctrl LocationController) Get(c echo.Context) {
}

func (ctrl LocationController) GetAll(c echo.Context) error {

	var locations []models.LocationDTO

	if err := ctrl.Service.GetAll(c, &locations); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(locations) == 0 {
		return c.JSON(http.StatusOK, errors.New("no locations found"))
	}

	return c.JSON(http.StatusOK, locations)

}

func (ctrl LocationController) DeleteAll(c echo.Context) error {

	count, err := ctrl.Service.DeleteAll(c)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, count)
}
