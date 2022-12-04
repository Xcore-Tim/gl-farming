package controllers

import (
	"gl-farming/app/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

type UIDController struct {
	Service services.UIDService
}

func NewUIDController() UIDController {
	return UIDController{}
}

func (ctrl UIDController) GetUID(c echo.Context) error {

	if uid, err := ctrl.Service.GetUID(c); err == nil {
		return c.JSON(http.StatusOK, uid)
	} else {
		return c.String(http.StatusBadRequest, err.Error())
	}

}

func (ctrl UIDController) GetServiceAdminUID(c echo.Context) error {
	if token, err := ctrl.Service.GetAdminToken(); err == nil {
		return c.String(http.StatusOK, token)
	} else {
		return c.String(http.StatusBadRequest, err.Error())
	}
}
