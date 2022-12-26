package controllers

import (
	"errors"
	"gl-farming/app/models"
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

func (ctrl UIDController) Login(c echo.Context) error {

	var authData models.UserCredentials

	if err := c.Bind(&authData); err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	uid, err := ctrl.Service.Login(authData)
	if err != nil {
		badCredentials := errors.New("incorrect user credentials")
		return c.JSON(http.StatusBadRequest, badCredentials)
	}

	return c.JSON(http.StatusOK, uid)
}

func (ctrl UIDController) Check(c echo.Context) error {

	if err := ctrl.Service.Check(c); err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}
	return c.String(http.StatusOK, "success")
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
