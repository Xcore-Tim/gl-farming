package controllers

import (
	"gl-farming/app/constants/requestStatus"
	"gl-farming/app/helper"
	"gl-farming/app/models"
	"gl-farming/app/services"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AccountRequestController struct {
	Services services.AppServices
}

func NewAccountRequestController(appServices services.AppServices) AccountRequestController {
	return AccountRequestController{
		Services: appServices,
	}
}

// Create godoc
// @Summary      Create account request
// @Description  creates new account request
// @Tags         Account requests
// @Accept       json
// @Produce      json
// @Param        createRequest    body     models.CreateAccountRequest  false  "create request body info"
// @Success      200  {array}  models.AccountRequest
// @Router       /v2/accountRequests/create [post]
func (ctrl AccountRequestController) Create(c echo.Context) error {

	uid, err := ctrl.Services.UID.GetUID(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var createAccountRequest models.CreateAccountRequest
	if err := c.Bind(&createAccountRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	adminToken, err := ctrl.Services.UID.GetAdminToken()
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var accountRequest models.AccountRequest

	if accountRequest.Type, err = ctrl.Services.AccountTypes.Get(c, createAccountRequest.TypeID); err != nil {
		return c.String(http.StatusBadRequest, "bad type")
	}

	if accountRequest.Location, err = ctrl.Services.Locations.Get(c, createAccountRequest.LocationID); err != nil {
		return c.String(http.StatusBadRequest, "bad location")
	}

	accountRequest.Status = requestStatus.Pending
	accountRequest.Buyer.FillWithUID(&uid)

	accountRequest.Team = models.Team{ID: uid.TeamID}
	if err := ctrl.Services.Teams.SetTeamlead(&uid.TeamID, &adminToken, &accountRequest.Team); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	accountRequest.Quantity = createAccountRequest.Quantity

	if createAccountRequest.Price != 0 {
		accountRequest.Price = createAccountRequest.Price
		accountRequest.Total = helper.CalculateTotal(createAccountRequest.Quantity, createAccountRequest.Price)
		if accountRequest.Currency, err = ctrl.Services.Currency.Get(c, createAccountRequest.CurrencyID); err != nil {
			return c.String(http.StatusBadRequest, "bad currency")
		}

		ctrl.setCurrency(c, &accountRequest)
	}

	accountRequest.Description = createAccountRequest.Description
	accountRequest.DateCreated = time.Now().Unix()

	if err := ctrl.Services.AccountRequests.Create(c, &accountRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// Cancel godoc
// @Summary      Cancel account request
// @Description  cancels account request
// @Tags         Account requests
// @Accept       json
// @Produce      json
// @Param        cancelRequest    body     models.CancelAccountRequest  false  "cancel request body info"
// @Success      200  {string}  string
// @Router       /v2/accountRequests/cancel [put]
func (ctrl AccountRequestController) Cancel(c echo.Context) error {

	var cancelRequest models.CancelAccountRequest

	if err := c.Bind(&cancelRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	uid, err := ctrl.Services.UID.GetUID(c)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	cancelRequest.CancelledBy.FillWithUID(&uid)
	cancelRequest.DateCancelled = time.Now().Unix()

	if err := ctrl.Services.AccountRequests.Cancel(c, &cancelRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// Take godoc
// @Summary      Take account request in work
// @Description  Takes account request in work
// @Tags         Account requests
// @Accept       json
// @Produce      json
// @Param        requestID    query     string  false  "taken request id"
// @Success      200  {string}  string
// @Router       /v2/accountRequests/take [put]
func (ctrl AccountRequestController) Take(c echo.Context) error {

	var farmer models.Employee

	uid, err := ctrl.Services.UID.GetUID(c)
	farmer.FillWithUID(&uid)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	ctrl.Services.AccountRequests.Take(c, &farmer)
	return c.String(http.StatusOK, "success")
}

// Update godoc
// @Summary      Update account request
// @Description  Updates account request
// @Tags         Account requests
// @Accept       json
// @Produce      json
// @Param        updateRequest    body    models.UpdateAccountRequest  false  "update request body"
// @Success      200  {string}  string
// @Router       /v2/accountRequests/update [put]
func (ctrl AccountRequestController) Update(c echo.Context) error {

	var updateAccountRequest models.UpdateAccountRequest
	if err := c.Bind(&updateAccountRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	uid, err := ctrl.Services.UID.GetUID(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	accountRequest, err := ctrl.Services.AccountRequests.Get(c, updateAccountRequest.RequestID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	accountRequest.MongoID, err = primitive.ObjectIDFromHex(updateAccountRequest.RequestID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	accountRequest.UpdatedBy.FillWithUID(&uid)

	if updateAccountRequest.TypeID != accountRequest.Type.ID {
		if accountRequest.Type, err = ctrl.Services.AccountTypes.Get(c, updateAccountRequest.TypeID); err != nil {
			return c.String(http.StatusBadRequest, "bad account type")
		}
	}

	if updateAccountRequest.LocationID != accountRequest.Location.ID {
		if accountRequest.Location, err = ctrl.Services.Locations.Get(c, updateAccountRequest.LocationID); err != nil {
			return c.String(http.StatusBadRequest, "bad location")
		}
	}

	accountRequest.Quantity = updateAccountRequest.Quantity

	if updateAccountRequest.Price != 0 {
		accountRequest.Price = updateAccountRequest.Price
		accountRequest.Total = helper.CalculateTotal(updateAccountRequest.Quantity, updateAccountRequest.Price)
		if accountRequest.Currency.ID != updateAccountRequest.CurrencyID {
			if accountRequest.Currency, err = ctrl.Services.Currency.Get(c, updateAccountRequest.CurrencyID); err != nil {
				return c.String(http.StatusBadRequest, "bad currency")
			}
			ctrl.setCurrency(c, &accountRequest)
		}
	}

	accountRequest.DateUpdated = time.Now().Unix()
	accountRequest.UpdatedBy.FillWithUID(&uid)

	if err := ctrl.Services.AccountRequests.Update(c, &accountRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, accountRequest)
}

// Complete godoc
// @Summary      Complete account request
// @Description  Completes account request
// @Tags         Account requests
// @Accept       json
// @Produce      json
// @Param        updateRequest    body    models.CompleteAccountRequest  false  "complete request body"
// @Success      200  {string}  string
// @Router       /v2/accountRequests/complete [put]
func (ctrl AccountRequestController) Complete(c echo.Context) error {

	uid, err := ctrl.Services.UID.GetUID(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	var completeAccountRequest models.CompleteAccountRequest
	if err := c.Bind(&completeAccountRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	accountRequest, err := ctrl.Services.AccountRequests.Get(c, completeAccountRequest.RequestID)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	accountRequest.Valid = completeAccountRequest.Valid

	if accountRequest.Price != completeAccountRequest.Price || completeAccountRequest.Price != 0 {
		accountRequest.Price = completeAccountRequest.Price
		accountRequest.Total = helper.CalculateTotal(accountRequest.Quantity, accountRequest.Price)
		if accountRequest.Currency.ID != completeAccountRequest.CurrencyID {
			if accountRequest.Currency, err = ctrl.Services.Currency.Get(c, completeAccountRequest.CurrencyID); err != nil {
				return c.String(http.StatusBadRequest, "bad currency")
			}
			ctrl.setCurrency(c, &accountRequest)
		}
	}

	accountRequest.DateCompleted = time.Now().Unix()
	accountRequest.CompletedBy.FillWithUID(&uid)

	if err := ctrl.Services.AccountRequests.Complete(c, &accountRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "success")
}

func (ctrl AccountRequestController) Return(c echo.Context) error {

	var returnAccountRequest models.ReturnAccountRequest
	uid, err := ctrl.Services.UID.GetUID(c)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	returnAccountRequest.CancelledBy.FillWithUID(&uid)

	if returnAccountRequest.RequestID, err = primitive.ObjectIDFromHex(c.QueryParam("requestID")); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	if err := ctrl.Services.AccountRequests.Return(c, &returnAccountRequest); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.String(http.StatusOK, "success")
}

// DeleteAll godoc
// @Summary      Delete all account requests
// @Description  deletes all account requests
// @Tags         Account requests
// @Accept       json
// @Produce      json
// @Success      200
// @Router       /v2/accountRequests/delete/all [delete]
func (ctrl AccountRequestController) DeleteAll(c echo.Context) error {

	count, err := ctrl.Services.AccountRequests.DeleteAll(c)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, count)
}

func (ctrl AccountRequestController) setCurrency(c echo.Context, accountRequest *models.AccountRequest) error {

	currencyRates, err := ctrl.Services.Currency.GetCurrencyRates()

	if err != nil {
		return err
	}

	accountRequest.Rate = currencyRates[accountRequest.Currency.ISO]
	if accountRequest.Rate == 0 {
		accountRequest.Rate = 1
	}

	if accountRequest.Currency.ISO == "USD" {
		accountRequest.BaseCurrency = accountRequest.Currency
		accountRequest.BaseRate = accountRequest.Rate
	} else {
		accountRequest.BaseCurrency, _ = ctrl.Services.Currency.GetByISO(c, "USD")
		accountRequest.BaseRate = currencyRates[accountRequest.BaseCurrency.ISO]
	}

	accountRequest.CrossRate = (accountRequest.Rate / accountRequest.BaseRate)
	accountRequest.CrossRate = helper.RoundFloat(accountRequest.CrossRate, 2)

	accountRequest.BasePrice = accountRequest.CrossRate * accountRequest.Price
	accountRequest.BasePrice = helper.RoundFloat(accountRequest.BasePrice, 2)
	accountRequest.BaseTotal = helper.CalculateTotal(accountRequest.Quantity, accountRequest.BasePrice)

	return nil
}
