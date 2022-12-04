package models

import (
	"gl-farming/app/constants/requestStatus"

	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type TableData struct {
	ID       string      `json:"_id,omitempty" bson:"_id"`
	Type     AccountType `json:"type,omitempty" bson:"type"`
	Location Location    `json:"location,omitempty" bson:"location"`
	Status   int         `json:"status,omitempty" bson:"status"`
	Team     Team        `json:"team,omitempty" bson:"team"`

	Quantity uint `json:"quantity" bson:"quantity"`
	Valid    uint `json:"valid,omitempty" bson:"valid"`

	Currency Currency `json:"currency,omitempty" bson:"currency,omitempty"`
	Rate     float64  `json:"rate,omitempty" bson:"rate,omitempty"`
	Price    float64  `json:"price,omitempty" bson:"price,omitempty"`
	Total    float64  `json:"total,omitempty" bson:"total,omitempty"`

	CrossRate float64 `json:"crossRate" bson:"crossRate,omitempty"`

	BaseCurrency Currency `json:"baseCurrency,omitempty" bson:"baseCurrency,omitempty"`
	BaseRate     float64  `json:"baseRate,omitempty" bson:"baseRate,omitempty"`
	BasePrice    float64  `json:"basePrice,omitempty" bson:"basePrice,omitempty"`
	BaseTotal    float64  `json:"baseTotal,omitempty" bson:"baseTotal,omitempty"`

	Buyer  Employee `json:"buyer,omitempty" bson:"buyer,omitempty"`
	Farmer Employee `json:"farmer,omitempty" bson:"farmer,omitempty"`

	TakenBy     Employee `json:"takenBy,omitempty" bson:"takenBy,omitempty"`
	UpdatedBy   Employee `json:"updatedBy,omitempty" bson:"updatedBy,omitempty"`
	CancelledBy Employee `json:"cancelledBy,omitempty" bson:"cancelledBy,omitempty"`
	CompletedBy Employee `json:"completedBy,omitempty" bson:"completedBy,omitempty"`
	ReturnedBy  Employee `json:"returnedBy,omitempty" bson:"returnedBy,omitempty"`

	Description       string `json:"description,omitempty" bson:"description,omitempty"`
	CancellationCause string `json:"cancellationCause,omitempty" bson:"cancellationCause,omitempty"`

	FileName string `json:"fileName,omitempty" bson:"fileName,omitempty"`

	DateCreated   int64 `json:"dateCreated,omitempty" bson:"dateCreated,omitempty"`
	DateTaken     int64 `json:"dateTaken,omitempty" bson:"dateTaken,omitempty"`
	DateUpdated   int64 `json:"dateUpdated,omitempty" bson:"dateUpdated,omitempty"`
	DateCancelled int64 `json:"dateCancelled,omitempty" bson:"dateCancelled,omitempty"`
	DateCompleted int64 `json:"dateCompleted,omitempty" bson:"dateCompleted,omitempty"`
	DateReturned  int64 `json:"dateReturned,omitempty" bson:"dateReturned,omitempty"`
}

type TableDataRequest struct {
	DataSlice  []TableData
	Filter     bson.D
	Projection bson.D
}

type Period struct {
	StartISO  string    `json:"startDate"`
	EndISO    string    `json:"endDate"`
	StartDate time.Time `json:"unixStart"`
	EndDate   time.Time `json:"unixEnd"`
}

func (p *Period) Convert() {
	date_format := "2006-01-02"

	if p.StartISO == "" {
		p.StartISO = "1970-01-01"
	}

	p.StartDate, _ = time.Parse(date_format, p.StartISO)

	if p.EndISO != "" {
		p.EndDate, _ = time.Parse(date_format, p.EndISO)
		return
	}

	p.EndDate = time.Now()
}

func (t *TableDataRequest) GetAll() (filter bson.D) {
	t.Filter = bson.D{bson.E{}}
	return
}

func (t *TableDataRequest) GetBuyerRequests(uid UID, period Period, status int) {

	t.Filter = bson.D{
		bson.E{Key: "$and", Value: bson.A{
			bson.D{
				bson.E{Key: "buyer.id", Value: uid.UserID},
				bson.E{Key: "status", Value: requestStatus.Pending},
				bson.E{Key: "dateCreated", Value: bson.M{"$gte": period.StartDate.Unix()}},
				bson.E{Key: "dateCreated", Value: bson.M{"$lte": period.EndDate.Unix()}},
			}}}}

	switch status {

	case requestStatus.Pending:

		t.Projection = bson.D{
			bson.E{Key: "_id", Value: 1},
			bson.E{Key: "type", Value: 1},
			bson.E{Key: "location", Value: 1},
			bson.E{Key: "team", Value: 1},
			bson.E{Key: "quantity", Value: 1},
			bson.E{Key: "farmer", Value: 1},
			bson.E{Key: "description", Value: 1},
			bson.E{Key: "dateCreated", Value: 1},
			bson.E{Key: "fileName", Value: 1},
			bson.E{Key: "dateUpdated", Value: 1},
			bson.E{Key: "updatedBy", Value: 1},
		}
	case requestStatus.Inwork:
		t.Projection = bson.D{
			bson.E{Key: "_id", Value: 1},
			bson.E{Key: "type", Value: 1},
			bson.E{Key: "location", Value: 1},
			bson.E{Key: "farmer", Value: 1},
			bson.E{Key: "team", Value: 1},
			bson.E{Key: "quantity", Value: 1},
			bson.E{Key: "currency", Value: 1},
			bson.E{Key: "dateCreated", Value: 1},
			bson.E{Key: "dateUpdated", Value: 1},
			bson.E{Key: "dateTaken", Value: 1},
			bson.E{Key: "description", Value: 1},
			bson.E{Key: "fileName", Value: 1},
			bson.E{Key: "takenBy", Value: 1},
			bson.E{Key: "updatedBy", Value: 1},
		}
	case requestStatus.Complete:
		t.Projection = bson.D{
			bson.E{Key: "_id", Value: 1},
			bson.E{Key: "type", Value: 1},
			bson.E{Key: "location", Value: 1},
			bson.E{Key: "team", Value: 1},
			bson.E{Key: "farmer", Value: 1},
			bson.E{Key: "quantity", Value: 1},
			bson.E{Key: "valid", Value: 0},
			bson.E{Key: "currency", Value: 1},
			bson.E{Key: "rate", Value: 1},
			bson.E{Key: "price", Value: 1},
			bson.E{Key: "total", Value: 1},
			bson.E{Key: "crossRate", Value: 1},
			bson.E{Key: "baseCurrency", Value: 1},
			bson.E{Key: "baseRate", Value: 1},
			bson.E{Key: "basePrice", Value: 1},
			bson.E{Key: "baseTotal", Value: 1},
			bson.E{Key: "description", Value: 1},
			bson.E{Key: "fileName", Value: 1},
			bson.E{Key: "dateCreated", Value: 1},
			bson.E{Key: "dateTaken", Value: 1},
			bson.E{Key: "dateUpdated", Value: 1},
			bson.E{Key: "dateCompleted", Value: 1},
			bson.E{Key: "completedBy", Value: 1},
			bson.E{Key: "updatedBy", Value: 1},
		}

	case requestStatus.Cancelled:
		t.Projection = bson.D{
			bson.E{Key: "_id", Value: 1},
			bson.E{Key: "type", Value: 1},
			bson.E{Key: "location", Value: 1},
			bson.E{Key: "team", Value: 1},
			bson.E{Key: "farmer", Value: 1},
			bson.E{Key: "quantity", Value: 1},
			bson.E{Key: "valid", Value: 0},
			bson.E{Key: "currency", Value: 1},
			bson.E{Key: "rate", Value: 1},
			bson.E{Key: "price", Value: 1},
			bson.E{Key: "total", Value: 1},
			bson.E{Key: "crossRate", Value: 1},
			bson.E{Key: "baseCurrency", Value: 1},
			bson.E{Key: "baseRate", Value: 1},
			bson.E{Key: "basePrice", Value: 1},
			bson.E{Key: "baseTotal", Value: 1},
			bson.E{Key: "description", Value: 1},
			bson.E{Key: "fileName", Value: 1},
			bson.E{Key: "dateCreated", Value: 1},
			bson.E{Key: "dateTaken", Value: 1},
			bson.E{Key: "dateUpdated", Value: 1},
			bson.E{Key: "cancellationCause", Value: 1},
			bson.E{Key: "dateCancelled", Value: 1},
			bson.E{Key: "cancelledBy", Value: 1},
		}
	}

	switch uid.RoleID {
	case 6:

	}

}

func (t *TableDataRequest) GetFarmerRequests(farmerAccess FarmerAccessList, period Period, status int) {

	t.Filter = bson.D{
		bson.E{Key: "$and", Value: bson.A{
			bson.D{
				bson.E{Key: "farmer.id", Value: farmerAccess.Farmer.ID},
				bson.E{Key: "team.number", Value: bson.D{{Key: "$in", Value: farmerAccess.Teams}}},
				bson.E{Key: "status", Value: status},
				bson.E{Key: "dateCreated", Value: bson.M{"$gte": period.StartDate.Unix()}},
				bson.E{Key: "dateCreated", Value: bson.M{"$lte": period.EndDate.Unix()}},
			},
		}},
	}

	switch status {
	case requestStatus.Pending:
		t.Filter = bson.D{
			bson.E{Key: "$and", Value: bson.A{
				bson.D{
					bson.E{Key: "team", Value: bson.D{{Key: "$in", Value: farmerAccess.Teams}}},
					bson.E{Key: "status", Value: requestStatus.Pending},
					bson.E{Key: "dateCreated", Value: bson.M{"$gte": period.StartDate.Unix()}},
					bson.E{Key: "dateCreated", Value: bson.M{"$lte": period.EndDate.Unix()}},
				},
			}}}

		t.Projection = bson.D{
			bson.E{Key: "_id", Value: 1},
			bson.E{Key: "type", Value: 1},
			bson.E{Key: "location", Value: 1},
			bson.E{Key: "team", Value: 1},
			bson.E{Key: "buyer", Value: 1},
			bson.E{Key: "quantity", Value: 1},
			bson.E{Key: "description", Value: 1},
			bson.E{Key: "dateCreated", Value: 1},
			bson.E{Key: "fileName", Value: 1},
			bson.E{Key: "dateUpdated", Value: 1},
			bson.E{Key: "updatedBy", Value: 1},
		}
	case requestStatus.Inwork:

		t.Projection = bson.D{
			bson.E{Key: "_id", Value: 1},
			bson.E{Key: "type", Value: 1},
			bson.E{Key: "location", Value: 1},
			bson.E{Key: "buyer", Value: 1},
			bson.E{Key: "team", Value: 1},
			bson.E{Key: "quantity", Value: 1},
			bson.E{Key: "currency", Value: 1},
			bson.E{Key: "dateCreated", Value: 1},
			bson.E{Key: "dateUpdated", Value: 1},
			bson.E{Key: "dateTaken", Value: 1},
			bson.E{Key: "description", Value: 1},
			bson.E{Key: "fileName", Value: 1},
			bson.E{Key: "takenBy", Value: 1},
			bson.E{Key: "updatedBy", Value: 1},
		}

	case requestStatus.Complete:
		t.Projection = bson.D{
			bson.E{Key: "_id", Value: 1},
			bson.E{Key: "type", Value: 1},
			bson.E{Key: "location", Value: 1},
			bson.E{Key: "buyer", Value: 1},
			bson.E{Key: "team", Value: 1},
			bson.E{Key: "quantity", Value: 1},
			bson.E{Key: "valid", Value: 0},
			bson.E{Key: "currency", Value: 1},
			bson.E{Key: "rate", Value: 1},
			bson.E{Key: "price", Value: 1},
			bson.E{Key: "total", Value: 1},
			bson.E{Key: "crossRate", Value: 1},
			bson.E{Key: "baseCurrency", Value: 1},
			bson.E{Key: "baseRate", Value: 1},
			bson.E{Key: "basePrice", Value: 1},
			bson.E{Key: "baseTotal", Value: 1},
			bson.E{Key: "description", Value: 1},
			bson.E{Key: "fileName", Value: 1},
			bson.E{Key: "dateCreated", Value: 1},
			bson.E{Key: "dateTaken", Value: 1},
			bson.E{Key: "dateUpdated", Value: 1},
			bson.E{Key: "dateCompleted", Value: 1},
			bson.E{Key: "completedBy", Value: 1},
			bson.E{Key: "updatedBy", Value: 1},
		}
	case requestStatus.Cancelled:
		t.Projection = bson.D{
			bson.E{Key: "_id", Value: 1},
			bson.E{Key: "type", Value: 1},
			bson.E{Key: "location", Value: 1},
			bson.E{Key: "team", Value: 1},
			bson.E{Key: "buyer", Value: 1},
			bson.E{Key: "quantity", Value: 1},
			bson.E{Key: "valid", Value: 0},
			bson.E{Key: "currency", Value: 1},
			bson.E{Key: "rate", Value: 1},
			bson.E{Key: "price", Value: 1},
			bson.E{Key: "total", Value: 1},
			bson.E{Key: "crossRate", Value: 1},
			bson.E{Key: "baseCurrency", Value: 1},
			bson.E{Key: "baseRate", Value: 1},
			bson.E{Key: "basePrice", Value: 1},
			bson.E{Key: "baseTotal", Value: 1},
			bson.E{Key: "description", Value: 1},
			bson.E{Key: "fileName", Value: 1},
			bson.E{Key: "dateCreated", Value: 1},
			bson.E{Key: "dateTaken", Value: 1},
			bson.E{Key: "dateUpdated", Value: 1},
			bson.E{Key: "cancellationCause", Value: 1},
			bson.E{Key: "dateCancelled", Value: 1},
			bson.E{Key: "cancelledBy", Value: 1},
		}
	}
}

// bson.D{
// 	bson.E{Key: "_id", Value: 1},
// 	bson.E{Key: "type", Value: 1},
// 	bson.E{Key: "location", Value: 1},
// 	bson.E{Key: "status", Value: 1},
// 	bson.E{Key: "team", Value: 1},
// 	bson.E{Key: "quantity", Value: 1},
// bson.E{Key: "valid", Value: 0},
// bson.E{Key: "currency", Value: 1},
// bson.E{Key: "rate", Value: 1},
// bson.E{Key: "price", Value: 1},
// bson.E{Key: "total", Value: 1},
// bson.E{Key: "crossRate", Value: 0},
// bson.E{Key: "baseCurrency", Value: 0},
// bson.E{Key: "baseRate", Value: 0},
// bson.E{Key: "basePrice", Value: 0},
// bson.E{Key: "baseTotal", Value: 0},
// bson.E{Key: "buyer", Value: 1},
// bson.E{Key: "farmer", Value: 0},
// bson.E{Key: "takenBy", Value: 0},
// bson.E{Key: "updatedBy", Value: 0},
// bson.E{Key: "cancelledBy", Value: 0},
// bson.E{Key: "completedBy", Value: 0},
// bson.E{Key: "returnedBy", Value: 0},
// bson.E{Key: "description", Value: 1},
// bson.E{Key: "cancellationCause", Value: 0},
// bson.E{Key: "fileName", Value: 0},
// bson.E{Key: "dateCreated", Value: 0},
// bson.E{Key: "dateTaken", Value: 0},
// bson.E{Key: "dateUpdated", Value: 0},
// bson.E{Key: "dateCancelled", Value: 0},
// bson.E{Key: "dateCompleted", Value: 0},
// bson.E{Key: "dateReturned", Value: 0},
// }
