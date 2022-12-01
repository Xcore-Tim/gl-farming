package controllers

import (
	"errors"
	"gl-farming/app/models"
	"gl-farming/app/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AccountTypeController struct {
	Service services.AccountTypeService
}

func NewAccountTypeController() AccountTypeController {
	return AccountTypeController{}
}

// Create godoc
// @Summary      Create account type
// @Description  creates account type
// @Tags         Account types
// @Accept       json
// @Produce      json
// @Param        id    body     models.AccountType  true  "account type body json"
// @Success      200  {string}   string
// @Router       /v2/accountTypes/create [post]
func (ctrl AccountTypeController) Create(c echo.Context) error {

	var accountType models.AccountType

	if err := c.Bind(&accountType); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := ctrl.Service.Create(c, accountType); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

func (ctrl AccountTypeController) Update(c echo.Context) {
}

// Get godoc
// @Summary      Delete account type
// @Description  Deletes account type by id
// @Tags         Account types
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "account type id"
// @Success      200  {string}  string
// @Router       /v2/accountTypes/delete [delete]
func (ctrl AccountTypeController) Delete(c echo.Context) error {

	id := c.QueryParam("id")

	if err := ctrl.Service.Delete(c, id); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "successfuly deleted account type")
}

// Get godoc
// @Summary      Get account types
// @Description  get account types by id
// @Tags         Account types
// @Accept       json
// @Produce      json
// @Param        id    query     string  true  "account type id"
// @Success      200  {object}   models.AccountType
// @Router       /v2/accountTypes/get [get]
func (ctrl AccountTypeController) Get(c echo.Context) error {

	id := c.QueryParam("id")

	accountType, err := ctrl.Service.Get(c, id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, accountType)
}

// GetAll godoc
// @Summary      Get all account types
// @Description  returns all account types
// @Tags         Account types
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.AccountType
// @Router       /v2/accountTypes/get/all [get]
func (ctrl AccountTypeController) GetAll(c echo.Context) error {

	var accountTypes []models.AccountType

	if err := ctrl.Service.GetAll(c, &accountTypes); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if len(accountTypes) == 0 {
		return c.JSON(http.StatusOK, errors.New("no account types found"))
	}

	return c.JSON(http.StatusOK, accountTypes)
}

// DeleteAll godoc
// @Summary      Delete all account types
// @Description  deletes all account types
// @Tags         Account types
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /v2/accountTypes/delete/all [delete]
func (ctrl AccountTypeController) DeleteAll(c echo.Context) error {

	count, err := ctrl.Service.DeleteAll(c)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, count)
}
