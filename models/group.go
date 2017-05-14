package models

type Group struct {
	BaseModel
	Name  string `json:"name" form:"name" binding:"required"`
	Description  string `json:"description" form:"description"`	
	Status string `json:"status"`
	Members []GroupMember `json:"members"`
}

func (g *Group) BeforeCreate() (err error) {
	g.Status = "active"
	return
}