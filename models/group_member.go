package models

type GroupMember struct {
	BaseModel
	GroupId int `json:"-" form:"group_id" binding:"required"`
	User User `json:"member" binding:"-"`
	UserId int `json:"-" form:"user_id" binding:"required"`
}