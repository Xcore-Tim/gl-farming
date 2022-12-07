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

// Get godoc
// @Summary      Get account requests by period and employee
// @Description  returns all account requests by period and employee
// @Tags         Table data
// @Accept       json
// @Produce      json
// @Param        startDate    query     string  false  "period start date"
// @Param        endDate    query     string  false  "period end date"
// @Param        status    query     string  false  "status"
// @Success      200  {array}  models.TableData
// @Router       /v2/tableData/get [get]
func (ctrl TableController) Get(c echo.Context) error {

	uid, err := ctrl.Services.UID.GetUID(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var period models.Period
	period.StartISO = c.QueryParam("startDate")
	period.EndISO = c.QueryParam("endDate")
	period.Convert()

	var tableData models.TableDataRequest

	status, _ := strconv.Atoi(c.QueryParam("status"))

	switch uid.RoleID {
	case 2, 3, 4, 7:
		tableData, err = ctrl.getBuyerRequests(c, uid, period, status)
	case 6:
		tableData, err = ctrl.getFarmerRequests(c, uid, period, status)
	case 5:
		tableData, err = ctrl.getTlfRequests(c, uid, period, status)
	}

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, tableData.DataSlice)
}

// GetTeamleadTables godoc
// @Summary      Get account requests by period and employee
// @Description  returns all account requests by period and employee
// @Tags         Table data
// @Accept       json
// @Produce      json
// @Param        getTeamleadTables    body  models.TeamleadTableRequest   false  "status"
// @Success      200  {array}  models.TableData
// @Router       /v2/tableData/teamlead/get [post]
func (ctrl TableController) GetTeamleadTables(c echo.Context) error {

	var teamleadTableRequest models.TeamleadTableRequest

	err := c.Bind(&teamleadTableRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	teamleadTableRequest.Period.Convert()

	var tableData models.TableDataRequest

	switch teamleadTableRequest.UID.RoleID {
	case 2, 3, 4, 7:
		tableData, err = ctrl.getBuyerRequests(c, teamleadTableRequest.UID, teamleadTableRequest.Period, teamleadTableRequest.Status)
	case 6:
		tableData, err = ctrl.getFarmerRequests(c, teamleadTableRequest.UID, teamleadTableRequest.Period, teamleadTableRequest.Status)
	}

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, tableData.DataSlice)
}

// GetFarmerList godoc
// @Summary      Get farmer list
// @Description  returns farmer list
// @Tags         Table data
// @Accept       json
// @Produce      json
// @Param        startDate    query     string  false  "period start date"
// @Param        endDate    query     string  false  "period end date"
// @Success      200  {array}  models.EmployeePipeline
// @Router       /v2/tableData/aggregate/farmers [get]
func (ctrl TableController) FarmerPipeline(c echo.Context) error {

	var period models.Period
	period.StartISO = c.QueryParam("startDate")
	period.EndISO = c.QueryParam("endDate")
	period.Convert()

	var farmers models.EmployeePipeline
	matchStage, groupStage := farmers.FarmerPipeline(period)

	atd, err := ctrl.Services.Tables.AggregateDataByUID(c, matchStage, groupStage)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, atd)
}

// GetBuyerList godoc
// @Summary      Get buyer list
// @Description  returns buyer list
// @Tags         Table data
// @Accept       json
// @Produce      json
// @Param        startDate    query     string  false  "period start date"
// @Param        endDate    query     string  false  "period end date"
// @Param        teamleadID    query     string  false  "teamlead id"
// @Success      200  {array}  models.EmployeePipeline
// @Router       /v2/tableData/aggregate/buyers [get]
func (ctrl TableController) BuyerPipiline(c echo.Context) error {

	var period models.Period
	var uid models.UID

	period.StartISO = c.QueryParam("startDate")
	period.EndISO = c.QueryParam("endDate")
	period.Convert()

	teamleadID := c.QueryParam("teamleadID")
	if teamleadID == "" {
		var err error
		uid, err = ctrl.Services.UID.GetUID(c)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		teamleadID = strconv.Itoa(uid.UserID)
	}

	var buyers models.EmployeePipeline
	matchStage, groupStage := buyers.BuyerPipiline(period, teamleadID)

	atd, err := ctrl.Services.Tables.AggregateDataByUID(c, matchStage, groupStage)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, atd)
}

// GetTeamList godoc
// @Summary      Get team list
// @Description  returns team list
// @Tags         Table data
// @Accept       json
// @Produce      json
// @Param        startDate    query     string  false  "period start date"
// @Param        endDate    query     string  false  "period end date"
// @Success      200  {array}  models.EmployeePipeline
// @Router       /v2/tableData/aggregate/teamleads [get]
func (ctrl TableController) TeamleadPipiline(c echo.Context) error {

	var period models.Period
	period.StartISO = c.QueryParam("startDate")
	period.EndISO = c.QueryParam("endDate")
	period.Convert()

	var teamleads models.TeamPipiline
	matchStage, groupStage := teamleads.Pipeline(period)

	atd, err := ctrl.Services.Tables.AggregateDataByTeam(c, matchStage, groupStage)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, atd)
}

func (ctrl TableController) getBuyerRequests(c echo.Context, uid models.UID, period models.Period, status int) (models.TableDataRequest, error) {

	var tableDataRequest models.TableDataRequest

	tableDataRequest.GetBuyerRequests(uid, period, status)

	if err := ctrl.Services.Tables.Get(c, &tableDataRequest); err != nil {
		return tableDataRequest, err
	}

	return tableDataRequest, nil
}

func (ctrl TableController) getFarmerRequests(c echo.Context, uid models.UID, period models.Period, status int) (models.TableDataRequest, error) {

	var tableDataRequest models.TableDataRequest

	var farmerAccess models.FarmerAccessList
	farmerAccess.Farmer.FillWithUID(&uid)
	farmerAccess.Teams = make([]int, 1)

	teams, err := ctrl.Services.Teams.GetAccess(c, &farmerAccess)
	if err != nil {
		return tableDataRequest, err
	}
	farmerAccess.Teams = append(farmerAccess.Teams, teams...)

	tableDataRequest.GetFarmerRequests(farmerAccess, period, status)

	if err := ctrl.Services.Tables.Get(c, &tableDataRequest); err != nil {
		return tableDataRequest, err
	}

	return tableDataRequest, nil
}

func (ctrl TableController) getTlfRequests(c echo.Context, uid models.UID, period models.Period, status int) (models.TableDataRequest, error) {

	var tableDataRequest models.TableDataRequest

	tableDataRequest.GetTlfRequests(uid, period, status)

	if err := ctrl.Services.Tables.Get(c, &tableDataRequest); err != nil {
		return tableDataRequest, err
	}

	return tableDataRequest, nil
}

// GetAll godoc
// @Summary      Get all account requests
// @Description  returns all account requests
// @Tags         Table data
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.AccountRequest
// @Router       /v2/tableData/get/all [get]
func (ctrl TableController) GetAll(c echo.Context) error {

	var tableDataRequest models.TableDataRequest
	tableDataRequest.GetAll()

	err := ctrl.Services.Tables.GetAll(c, &tableDataRequest)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, tableDataRequest.DataSlice)
}
