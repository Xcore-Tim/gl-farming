package controllers

import (
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

// Revoke godoc
// @Summary      revoke access
// @Description  revokes access to farmer
// @Tags         Farmer Access
// @Accept       json
// @Produce      json
// @Param        accessRequest    body     models.AccessRequest  false  "farmer uid"
// @Success      200  {array}  models.AccessRequest
// @Router       /v2/farmerAccess/revoke [delete]
func (ctrl FarmerAccessController) Revoke(c echo.Context) error {

	var revokeAccessRequest models.AccessRequest

	if err := c.Bind(&revokeAccessRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := ctrl.Services.Teams.RevokeAccess(c, &revokeAccessRequest); err != nil {
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
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	c.JSON(http.StatusOK, teams)

	return nil
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
