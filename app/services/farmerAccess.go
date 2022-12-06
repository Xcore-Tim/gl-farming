package services

import (
	"errors"
	userRole "gl-farming/app/constants/roles"
	"gl-farming/app/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
)

func (s TeamServiceImpl) GetFarmers(adminToken string) ([]models.Employee, error) {
	var farmerAccess []responseUser
	var farmers []models.Employee

	if err := s.GetUsersByRole(&farmerAccess, &adminToken, 6); err != nil {
		return farmers, err
	}

	for _, farmerUser := range farmerAccess {
		if farmerUser.IsActive {
			var farmer models.Employee
			farmer.ID = farmerUser.ID
			farmer.FullName = farmerUser.Username
			farmer.Role = userRole.Farmer
			farmers = append(farmers, farmer)
		}
	}

	return farmers, nil
}

func (s TeamServiceImpl) AddAccess(c echo.Context, farmerAccess *models.AccessRequest) error {

	if _, err := s.collection.InsertOne(c.Request().Context(), farmerAccess); err != nil {
		return err
	}

	return nil
}

func (s TeamServiceImpl) RevokeAccess(c echo.Context, farmerAccess *models.AccessRequest) error {

	filter := bson.D{
		bson.E{Key: "farmer", Value: farmerAccess.Farmer},
		bson.E{Key: "team", Value: farmerAccess.TeamID},
	}
	result, _ := s.collection.DeleteOne(c.Request().Context(), filter)

	if result.DeletedCount != 1 {
		return errors.New("no farmer access found")
	}

	return nil
}

func (s TeamServiceImpl) UpdateAccess(c echo.Context, farmerAccess *models.AccessRequest) error {

	filter := bson.D{
		bson.E{Key: "farmer", Value: farmerAccess.Farmer},
		bson.E{Key: "team", Value: farmerAccess.TeamID},
	}

	update := bson.D{bson.E{Key: "$set", Value: bson.D{
		bson.E{Key: "farmer", Value: farmerAccess.Farmer},
		bson.E{Key: "team", Value: farmerAccess.TeamID},
	}}}

	result, _ := s.collection.UpdateOne(c.Request().Context(), filter, update)

	if result.MatchedCount != 1 {
		return errors.New("no farmer access found")
	}

	return nil
}

func (s TeamServiceImpl) GetAccess(c echo.Context, farmerAccess *models.FarmerAccessList) ([]int, error) {

	var accessList []int

	filter := bson.D{
		bson.E{Key: "farmer", Value: farmerAccess.Farmer},
	}

	cursor, err := s.collection.Find(c.Request().Context(), filter)

	if err != nil {
		return accessList, err
	}

	defer cursor.Close(c.Request().Context())

	for cursor.Next(c.Request().Context()) {
		var farmerAccess models.FarmerAccess
		err := cursor.Decode(&farmerAccess)
		if err != nil {
			return accessList, err
		}
		farmerAccess.ConvertID()
		accessList = append(accessList, farmerAccess.Team)
	}

	if err := cursor.Err(); err != nil {
		return accessList, err
	}

	return accessList, err
}

func (s TeamServiceImpl) GetAllAccesses(c echo.Context, accessList *[]models.FarmerAccess) error {

	cursor, err := s.collection.Find(c.Request().Context(), bson.D{{}})

	if err != nil {
		return err
	}

	defer cursor.Close(c.Request().Context())

	for cursor.Next(c.Request().Context()) {
		var farmerAccess models.FarmerAccess
		err := cursor.Decode(&farmerAccess)
		if err != nil {
			return err
		}
		farmerAccess.ConvertID()
		*accessList = append(*accessList, farmerAccess)
	}

	if err := cursor.Err(); err != nil {
		return err
	}

	return nil
}

func (s TeamServiceImpl) DeleteAllAccesses(c echo.Context) (int, error) {

	filter := bson.D{bson.E{}}
	result, err := s.collection.DeleteMany(c.Request().Context(), filter)

	if err != nil {
		return 0, err
	}

	return int(result.DeletedCount), nil
}
