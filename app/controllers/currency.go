package controllers

import (
	"errors"
	"gl-farming/app/models"
	"gl-farming/app/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CurrencyController struct {
	Service services.CurrencyService
}

func NewCurrencyController() CurrencyController {
	return CurrencyController{}
}

// Create godoc
// @Summary      Create currency
// @Description  creates currency
// @Tags         Currency
// @Accept       json
// @Produce      json
// @Param        id    body     models.Currency  true  "currency body json"
// @Success      200  {string}   string
// @Router       /v2/currency/create [post]
func (ctrl CurrencyController) Create(c echo.Context) error {

	var currency models.Currency

	if err := c.Bind(&currency); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := ctrl.Service.Create(c, currency); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

func (ctrl CurrencyController) Update(c echo.Context) {
}

// Get godoc
// @Summary      Delete currency
// @Description  Deletes currency by id
// @Tags         Currency
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "currency id"
// @Success      200  {string}  string
// @Router       /v2/currency/delete [delete]
func (ctrl CurrencyController) Delete(c echo.Context) error {

	id := c.QueryParam("id")

	if err := ctrl.Service.Delete(c, id); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "successfuly deleted currency")
}

// Get godoc
// @Summary      Get currency
// @Description  get currency by id
// @Tags         Currency
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "currency id"
// @Success      200  {object}   models.Currency
// @Router       /v2/currency/get [get]
func (ctrl CurrencyController) Get(c echo.Context) error {

	id := c.QueryParam("id")

	currency, err := ctrl.Service.Get(c, id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, currency)
}

// GetAll godoc
// @Summary      Get all currency
// @Description  returns all currency
// @Tags         Currency
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.Currency
// @Router       /v2/currency/get/all [get]
func (ctrl CurrencyController) GetAll(c echo.Context) error {

	var currency []models.Currency

	if err := ctrl.Service.GetAll(c, &currency); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(currency) == 0 {
		return c.JSON(http.StatusOK, errors.New("no currency found"))
	}

	return c.JSON(http.StatusOK, currency)
}

// DeleteAll godoc
// @Summary      Delete all currency
// @Description  deletes all currency
// @Tags         Currency
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /v2/currency/delete/all [delete]
func (ctrl CurrencyController) DeleteAll(c echo.Context) error {

	count, err := ctrl.Service.DeleteAll(c)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, count)
}
