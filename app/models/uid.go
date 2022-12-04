package models

const (
	UI_API_BASEPATH string = "https://g-identity-test.azurewebsites.net"
	AuthEndpoint    string = "/v1/accounts/auth"
	ContentTypeAuth string = "application/json-patch+json"
)

type UID struct {
	UserID   int    `json:"userID"`
	RoleID   int    `json:"roleID"`
	TeamID   int    `json:"teamID"`
	Username string `json:"username"`
	Token    string `json:"token"`
}
