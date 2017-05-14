package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "exp/mngr/api/models"
)

type GroupHandler struct {
	db *gorm.DB
}

func NewGroupHandler(db *gorm.DB) *GroupHandler {
	return &GroupHandler{db}
}

//get all groups
func (handler GroupHandler) Index(c *gin.Context) {
	groups := []m.Group{}		
	handler.db.Preload("Members").Preload("Members.User").Find(&groups)
	c.JSON(http.StatusOK, groups)
	return
}

//create new group
func (handler GroupHandler) Create(c *gin.Context) {
	var group m.Group
	err := c.Bind(&group)
	if err == nil {
		//check for group name
		existingGroup := m.Group{}
		existingGroupQuery := handler.db.Where("name = ?", group.Name).First(&existingGroup)
		if existingGroupQuery.RowsAffected > 0 {
			respond(http.StatusBadRequest, "Group name already taken.", c, true)
		} else {
			saveResult := handler.db.Save(&group)
			if saveResult.RowsAffected > 0 {
				c.JSON(http.StatusCreated, group)
			} else {
				respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
			}
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
}