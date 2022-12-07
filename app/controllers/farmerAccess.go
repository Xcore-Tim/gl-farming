package controllers

import (
	"gl-farming/app/helper"
	"gl-farming/app/models"
	"gl-farming/app/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type FarmerAccessController struct {
	Services services.AppServices
}

func NewFarmerAccessController(appServices services.AppServices) FarmerAccessController {
	return FarmerAccessController{
		Services: appServices,
	}
}

// Add godoc
// @Summary      Add access
// @Description  adds access to farmer
// @Tags         Farmer Access
// @Accept       json
// @Produce      json
// @Param        farmer    body     models.AccessRequest  false  "farmer uid"
// @Success      200  {array}  models.AccessRequest
// @Router       /v2/farmerAccess/add [post]
func (ctrl FarmerAccessController) Add(c echo.Context) error {

	var addAccessRequest models.AccessRequest

	if err := c.Bind(&addAccessRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := ctrl.Services.Teams.UpdateAccess(c, &addAccessRequest); err != nil {
		if err := ctrl.Services.Teams.AddAccess(c, &addAccessRequest); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	}

	return c.String(http.StatusOK, "success")
}

// FullAccess godoc
// @Summary      Full access
// @Description  Sets full access to all teams for farmer
// @Tags         Farmer Access
// @Accept       json
// @Produce      json
// @Param        fullAccessRequest    body     models.FullAccessRequest  false  "farmer uid"
// @Success      200  {string}  string	"success"
// @Router       /v2/farmerAccess/add/all [put]
func (ctrl FarmerAccessController) FullAccess(c echo.Context) error {

	adminToken, err := ctrl.Services.UID.GetAdminToken()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var fullAccessRequest models.FullAccessRequest
	if err := c.Bind(&fullAccessRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	farmerAccessList := models.FarmerAccessList{
		Farmer: fullAccessRequest.Farmer,
		Teams:  make([]int, 0, 1),
	}
	if err = ctrl.Services.Teams.FullAccess(c, adminToken, &farmerAccessList); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// FullAccess godoc
// @Summary      Full access
// @Description  Sets full access to all teams for farmer
// @Tags         Farmer Access
// @Accept       json
// @Produce      json
// @Param        fullAccessRequest    body     models.FullAccessRequest  false  "farmer uid"
// @Success      200  {string}  string	"success"
// @Router       /v2/farmerAccess/revoke/all [put]
func (ctrl FarmerAccessController) FullRevoke(c echo.Context) error {

	var fullAccessRequest models.FullAccessRequest
	if err := c.Bind(&fullAccessRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var farmerAccessList = models.FarmerAccessList{
		Farmer: fullAccessRequest.Farmer,
		Teams:  make([]int, 1),
	}

	teams, err := ctrl.Services.Teams.GetAccess(c, &farmerAccessList)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	farmerAccessList.Teams = append(farmerAccessList.Teams, teams...)

	if err := ctrl.Services.Teams.FullRevoke(c, &farmerAccessList); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// Revoke godoc
// @Summary      revoke access
// @Description  revokes access to farmer
// @Tags         Farmer Access
// @Accept       json
// @Produce      json
// @Param        accessRequest    body     models.AccessRequest  false  "farmer uid"
// @Success      200  {array}  models.AccessRequest
// @Router       /v2/farmerAccess/revoke [put]
func (ctrl FarmerAccessController) Revoke(c echo.Context) error {

	var revokeAccessRequest models.AccessRequest

	if err := c.Bind(&revokeAccessRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := ctrl.Services.Teams.RevokeAccess(c, &revokeAccessRequest); err != nil {
		c.JSON(http.StatusOK, revokeAccessRequest)
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// GetAll godoc
// @Summary      Get all accesses
// @Description  returns all accesses
// @Tags         Table data
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.FarmerAccess
// @Router       /v2/farmerAccess/get/all [get]
func (ctrl FarmerAccessController) GetAll(c echo.Context) error {

	var accessList []models.FarmerAccess

	if err := ctrl.Services.Teams.GetAllAccesses(c, &accessList); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, accessList)
}

// GetTeams godoc
// @Summary      Get teams
// @Description  returns all teams
// @Tags         Farmer Access
// @Accept       json
// @Produce      json
// @Success      200  {array}  int
// @Router       /v2/farmerAccess/get/teams [get]
func (ctrl FarmerAccessController) GetTeams(c echo.Context) error {

	adminToken, err := ctrl.Services.UID.GetAdminToken()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	teams, err := ctrl.Services.Teams.GetTeams(&adminToken)

	for i, v := range teams {
		if v == 0 {
			teams = helper.RemoveElementFromSlice(teams, i)
			break
		}
	}

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, teams)
}

// GetFarmers godoc
// @Summary      Get farmer access
// @Description  returns all farmers accesses
// @Tags         Farmer Access
// @Accept       json
// @Produce      json
// @Success      200  {array}  models.FarmerAccessList
// @Router       /v2/farmerAccess/get/farmers [get]
func (ctrl FarmerAccessController) GetFarmers(c echo.Context) error {

	adminToken, err := ctrl.Services.UID.GetAdminToken()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	farmers, err := ctrl.Services.Teams.GetFarmers(adminToken)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var accessList []models.FarmerAccessList
	for _, farmer := range farmers {

		var farmerAccess = models.FarmerAccessList{
			Farmer: farmer,
			Teams:  make([]int, 1),
		}

		teams, err := ctrl.Services.Teams.GetAccess(c, &farmerAccess)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		farmerAccess.Teams = append(farmerAccess.Teams, teams...)

		for i, v := range farmerAccess.Teams {
			if v == 0 {
				farmerAccess.Teams = helper.RemoveElementFromSlice(farmerAccess.Teams, i)
			}
		}

		accessList = append(accessList, farmerAccess)
	}

	return c.JSON(http.StatusOK, accessList)
}

func (ctrl FarmerAccessController) DeleteAll(c echo.Context) error {
	count, err := ctrl.Services.Teams.DeleteAllAccesses(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusBadRequest, count)
}
