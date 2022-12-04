package controllers

import (
	"gl-farming/app/models"
	"gl-farming/app/services"
	"net/http"
	"strconv"

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

func (ctrl FarmerAccessController) Add(c echo.Context) error {

	var farmer models.Employee

	if err := c.Bind(&farmer); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var farmerAccess models.FarmerAccess
	team, _ := strconv.Atoi(c.QueryParam("teamID"))
	farmerAccess.Farmer = farmer
	farmerAccess.Team = team

	if err := ctrl.Services.Teams.UpdateAccess(c, &farmerAccess); err != nil {
		if err := ctrl.Services.Teams.AddAccess(c, &farmerAccess); err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
	}

	return c.String(http.StatusOK, "success")
}

func (ctrl FarmerAccessController) Revoke(c echo.Context) error {

	var farmer models.Employee

	if err := c.Bind(&farmer); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var farmerAccess models.FarmerAccess
	team, _ := strconv.Atoi(c.QueryParam("teamID"))
	farmerAccess.Farmer = farmer
	farmerAccess.Team = team

	if err := ctrl.Services.Teams.RevokeAccess(c, &farmerAccess); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

func (ctrl FarmerAccessController) GetAccess(c echo.Context) error {

	// var farmer models.Employee

	// if err := c.Bind(&farmer); err != nil {
	// 	return c.String(http.StatusBadRequest, err.Error())
	// }

	// var farmerAccess models.FarmerAccess

	// farmerAccess.Farmer = farmer

	// accessList, err := ctrl.Services.Teams.GetAccess(c, &farmerAccess)
	// if err != nil {
	// 	return c.String(http.StatusBadRequest, err.Error())
	// }

	return c.JSON(http.StatusOK, 1)
}

func (ctrl FarmerAccessController) GetAll(c echo.Context) error {

	var accessList []models.FarmerAccess

	if err := ctrl.Services.Teams.GetAllAccesses(c, &accessList); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, accessList)
}

func (ctrl FarmerAccessController) GetTeams(c echo.Context) error {

	adminToken, err := ctrl.Services.UID.GetAdminToken()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	teams, err := ctrl.Services.Teams.GetTeams(&adminToken)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, teams)

	return nil
}

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
		var farmerAccess models.FarmerAccessList
		farmerAccess.Farmer = farmer
		farmerAccess.Teams = make([]int, 1)

		teams, err := ctrl.Services.Teams.GetAccess(c, &farmerAccess)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		farmerAccess.Teams = append(farmerAccess.Teams, teams...)

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
