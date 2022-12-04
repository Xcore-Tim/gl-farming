package controllers

import (
	"gl-farming/app/models"
	"gl-farming/app/services"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type TableController struct {
	Services services.AppServices
}

func NewTableController(appServices services.AppServices) TableController {
	return TableController{
		Services: appServices,
	}
}

func (ctrl TableController) Get(c echo.Context) error {

	uid, err := ctrl.Services.UID.GetUID(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var tableData models.TableDataRequest

	status, _ := strconv.Atoi(c.QueryParam("status"))

	switch uid.RoleID {
	case 2, 3, 4, 7:
		tableData, err = ctrl.getBuyerRequests(c, status)
	case 6:
		tableData, err = ctrl.getFarmerRequests(c, status)
	}

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, tableData)
}

func (ctrl TableController) getBuyerRequests(c echo.Context, status int) (models.TableDataRequest, error) {

	var tableDataRequest models.TableDataRequest

	uid, err := ctrl.Services.UID.GetUID(c)
	if err != nil {
		return tableDataRequest, err
	}

	var period models.Period
	if err := c.Bind(&period); err != nil {
		return tableDataRequest, err
	}
	period.Convert()

	tableDataRequest.GetBuyerRequests(uid, period, status)

	if err := ctrl.Services.Tables.Get(c, &tableDataRequest); err != nil {
		return tableDataRequest, err
	}

	return tableDataRequest, nil
}

func (ctrl TableController) getFarmerRequests(c echo.Context, status int) (models.TableDataRequest, error) {

	var tableDataRequest models.TableDataRequest

	uid, err := ctrl.Services.UID.GetUID(c)
	if err != nil {
		return tableDataRequest, err
	}

	var period models.Period
	if err := c.Bind(&period); err != nil {
		return tableDataRequest, err
	}
	period.Convert()

	var farmerAccess models.FarmerAccessList
	farmerAccess.Farmer.FillWithUID(&uid)
	farmerAccess.Teams = make([]int, 1)

	teams, err := ctrl.Services.Teams.GetAccess(c, &farmerAccess)
	if err != nil {
		return tableDataRequest, err
	}
	farmerAccess.Teams = append(farmerAccess.Teams, teams...)

	tableDataRequest.GetFarmerRequests(farmerAccess, period, status)

	return tableDataRequest, nil
}

// GetAll godoc
// @Summary      Get all account requests
// @Description  returns all account requests
// @Tags         Account requests
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.AccountRequest
// @Router       /v2/accountRequests/get/all [get]
func (ctrl TableController) GetAll(c echo.Context) error {

	var tableDataRequest models.TableDataRequest
	tableDataRequest.GetAll()

	err := ctrl.Services.Tables.GetAll(c, &tableDataRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, tableDataRequest.DataSlice)
}
