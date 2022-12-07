package services

import (
	"encoding/json"
	"gl-farming/app/constants/gipsyUI"
	userRole "gl-farming/app/constants/roles"
	"gl-farming/app/helper"
	"gl-farming/app/models"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type TeamService interface {
	GetTeamleadByTeamID(*int, *string) (models.Employee, error)
	SetTeamlead(*int, *string, *models.Team) error
	GetTeams(*string) ([]int, error)
	GetFarmers(string) ([]models.Employee, error)

	FullAccess(echo.Context, string, *models.FarmerAccessList) error
	FullRevoke(echo.Context, *models.FarmerAccessList) error
	AddAccess(echo.Context, *models.AccessRequest) error
	RevokeAccess(echo.Context, *models.AccessRequest) error
	UpdateAccess(echo.Context, *models.AccessRequest) error

	GetAccess(echo.Context, *models.FarmerAccessList) ([]int, error)
	GetAllAccesses(echo.Context, *[]models.FarmerAccess) error
	DeleteAllAccesses(echo.Context) (int, error)
}

type TeamServiceImpl struct {
	collection *mongo.Collection
}

func NewTeamService(collection *mongo.Collection) TeamService {
	return &TeamServiceImpl{
		collection: collection,
	}
}

type responseUser struct {
	ID       int    `json:"ID"`
	Username string `json:"Username"`
	TeamID   int    `json:"Teamid"`
	IsActive bool   `json:"isActive"`
}

func (s TeamServiceImpl) GetTeamleadByTeamID(teamID *int, adminToken *string) (models.Employee, error) {

	var tl models.Employee
	var teamLeads []responseUser

	if err := s.GetUsersByRole(&teamLeads, adminToken, 2); err != nil {
		return tl, err
	}

	for _, teamlead := range teamLeads {
		if teamlead.TeamID == *teamID {
			tl.ID = teamlead.ID
			tl.FullName = teamlead.Username
			tl.Role = userRole.TeamLead
			break
		}
	}

	return tl, nil
}

func (s TeamServiceImpl) GetTeams(adminToken *string) ([]int, error) {

	var teamleads []responseUser
	teams := make([]int, 1)

	if err := s.GetUsersByRole(&teamleads, adminToken, 2); err != nil {
		return teams, nil
	}

	for _, teamlead := range teamleads {
		if teamlead.TeamID == 0 {
			continue
		}
		teams = append(teams, teamlead.TeamID)
	}

	teams = helper.Unique(teams)
	teams = helper.BubbleSort(teams)

	return teams, nil
}

func (s TeamServiceImpl) SetTeamlead(teamID *int, adminToken *string, team *models.Team) error {

	var teamLeads []responseUser

	if err := s.GetUsersByRole(&teamLeads, adminToken, 2); err != nil {
		return err
	}

	for _, teamlead := range teamLeads {
		if teamlead.TeamID == *teamID {
			team.Teamlead.ID = teamlead.ID
			team.Teamlead.FullName = teamlead.Username
			team.Teamlead.Role = userRole.TeamLead
			break
		}
	}

	return nil
}

func (s TeamServiceImpl) GetUsersByRole(users *[]responseUser, adminToken *string, roleID int) error {

	role := strconv.Itoa(roleID)
	url := gipsyUI.Basepath + gipsyUI.UsersByRoleEndpoint + role

	bearer := "BEARER " + *adminToken

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return err
	}

	request.Header.Add("Authorization", bearer)

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, _ := io.ReadAll(response.Body)

	if err := json.Unmarshal([]byte(body), users); err != nil {
		return err
	}

	return nil
}
