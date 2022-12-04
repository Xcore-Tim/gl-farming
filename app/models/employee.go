package models

type Employee struct {
	ID       int    `json:"_id,omitempty"`
	FullName string `json:"fullName" bson:"fullName"`
	Role     int    `json:"role" bson:"role"`
}

func (e *Employee) FillWithUID(uid *UID) {
	e.ID = uid.UserID
	e.FullName = uid.Username
	e.Role = uid.RoleID
}
