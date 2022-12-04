package controllers

import "gl-farming/app/services"

type AppControllers struct {
	AccountRequests        AccountRequestController
	Tables                 TableController
	FarmerAccessController FarmerAccessController
	AccountTypes           AccountTypeController
	Currency               CurrencyController
	Locations              LocationController
	UID                    UIDController
}

func NewAppControllers(appServices services.AppServices) AppControllers {

	var ac = AppControllers{
		AccountRequests:        NewAccountRequestController(appServices),
		Tables:                 NewTableController(appServices),
		FarmerAccessController: NewFarmerAccessController(appServices),
		AccountTypes:           NewAccountTypeController(),
		Currency:               NewCurrencyController(),
		Locations:              NewLocationController(),
		UID:                    NewUIDController(),
	}

	ac.AccountRequests.Services = appServices

	ac.AccountTypes.Service = appServices.AccountTypes
	ac.Currency.Service = appServices.Currency
	ac.Locations.Service = appServices.Locations
	ac.UID.Service = appServices.UID

	return ac
}
