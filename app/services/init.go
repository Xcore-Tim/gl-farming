package services

import "gl-farming/database"

type AppServices struct {
	AccountRequests AccountRequestService
	Tables          TableService
	AccountTypes    AccountTypeService
	Teams           TeamService
	Currency        CurrencyService
	Locations       LocationService
	UID             UIDService
}

func (as *AppServices) Init(collections *database.Collections) {
	as.AccountRequests = NewAccountRequestService(collections.AccountRequests)
	as.Tables = NewTableService(collections.AccountRequests)
	as.AccountTypes = NewAccountTypeService(collections.AccountTypes)
	as.Currency = NewCurrencyService(collections.Currency)
	as.Locations = NewLocationService(collections.Locations)
	as.Teams = NewTeamService(collections.FarmerAccess)
	as.UID = NewUIDService()
}
