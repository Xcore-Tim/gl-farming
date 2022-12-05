package services

import (
	"bytes"
	"encoding/json"
	"gl-farming/app/constants/gipsyUI"
	"gl-farming/app/models"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

type UIDService interface {
	GetUID(echo.Context) (models.UID, error)
	GetAdminToken() (string, error)
	Login(models.UserCredentials) (models.UID, error)
	// GetAuthToken(*echo.Context) string
}

type UIDServiceImpl struct {
}

func NewUIDService() UIDService {
	return &UIDServiceImpl{}
}

func (s UIDServiceImpl) GetUID(ctx echo.Context) (models.UID, error) {

	var UID models.UID

	tokenString := s.GetAuthToken(ctx)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return UID, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return UID, err
	}

	UID.UserID, _ = strconv.Atoi(claims["UserId"].(string))
	UID.RoleID, _ = strconv.Atoi(claims["TeamId"].(string))
	UID.TeamID, _ = strconv.Atoi(claims["RoleId"].(string))
	UID.Token = tokenString

	url := gipsyUI.Basepath + gipsyUI.UsersEndpoint + claims["UserId"].(string)
	bearer := "BEARER " + tokenString

	request, err := http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return UID, err
	}

	request.Header.Add("Authorization", bearer)

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return UID, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return UID, err
	}

	if err := json.Unmarshal([]byte(body), &UID); err != nil {
		return UID, err
	}

	return UID, nil
}

func (s UIDServiceImpl) Login(userCredentials models.UserCredentials) (models.UID, error) {

	var uid models.UID

	urlPath := gipsyUI.Basepath + gipsyUI.AuthEndpoint

	requestBody, _ := json.Marshal(map[string]string{
		"email":    userCredentials.Email,
		"password": userCredentials.Password,
	})

	bodyReader := bytes.NewBuffer(requestBody)

	request, err := http.NewRequest(http.MethodPost, urlPath, bodyReader)

	if err != nil {
		return uid, err
	}

	request.Header.Set("Content-Type", gipsyUI.AuthContentType)

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return uid, err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return uid, err
	}

	if err = json.Unmarshal([]byte(body), &uid); err != nil {
		return uid, err
	}

	url := gipsyUI.Basepath + gipsyUI.UsersEndpoint + strconv.Itoa(uid.UserID)
	bearer := "BEARER " + uid.Token

	request, err = http.NewRequest(http.MethodGet, url, nil)

	if err != nil {
		return uid, err
	}

	request.Header.Add("Authorization", bearer)

	response, err = client.Do(request)

	if err != nil {
		return uid, err
	}

	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)

	if err != nil {
		return uid, err
	}

	if err := json.Unmarshal([]byte(body), &uid); err != nil {
		return uid, err
	}

	return uid, nil

}

func (s UIDServiceImpl) GetAdminToken() (string, error) {

	var uid models.UID

	urlPath := gipsyUI.Basepath + gipsyUI.AuthEndpoint

	requestBody, _ := json.Marshal(map[string]string{
		"email":    gipsyUI.ServiceAdminEmail,
		"password": gipsyUI.ServiceAdminPassword,
	})

	bodyReader := bytes.NewBuffer(requestBody)

	request, err := http.NewRequest(http.MethodPost, urlPath, bodyReader)

	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", gipsyUI.AuthContentType)

	client := &http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return "", err
	}

	if err = json.Unmarshal([]byte(body), &uid); err != nil {
		return "", err
	}

	return uid.Token, nil
}

func (UIDServiceImpl) GetAuthToken(ctx echo.Context) string {

	token := ctx.Request().Header.Get("Authorization")
	token = strings.ReplaceAll(token, "Bearer ", "")
	return token
}
