package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "exp/mngr/api/models"
)

type GroupMemberHandler struct {
	db *gorm.DB
}

func NewGroupMemberHandler(db *gorm.DB) *GroupMemberHandler {
	return &GroupMemberHandler{db}
}

//get all group members
func (handler GroupMemberHandler) Index(c *gin.Context) {
	groupMembers := []m.GroupMember{}		
	handler.db.Find(&groupMembers)
	c.JSON(http.StatusOK, groupMembers)
	return
}

//create new group member
func (handler GroupMemberHandler) Create(c *gin.Context) {
	var groupMember m.GroupMember
	err := c.Bind(&groupMember)
	if err == nil {
		existingGroup := m.Group{}
		existingGroupQuery := handler.db.Where("id = ?", groupMember.GroupId).First(&existingGroup)
		if existingGroupQuery.RowsAffected > 0 {
			existingMember := m.User{}
			existingMemberQuery := handler.db.Where("id = ?", groupMember.UserId).First(&existingMember)

			if existingMemberQuery.RowsAffected > 0 {
				existingGroupMember := m.GroupMember{}
				existingGroupMemberQuery := handler.db.Where("user_id = ? AND group_id = ?", groupMember.UserId, groupMember.GroupId).First(&existingGroupMember)

				if existingGroupMemberQuery.RowsAffected == 0 {
					saveResult := handler.db.Save(&groupMember)
					if saveResult.RowsAffected > 0 {
						respond(http.StatusCreated, "Member successfully added to group.", c, false)
					} else {
						respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
					}	
				} else {
					respond(http.StatusBadRequest, "Member already in group.", c, true)
				}
			} else {
				respond(http.StatusBadRequest, "Member record not found.", c, true)
			}
		} else {
			respond(http.StatusNotFound, "Group not found.", c, true)
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
}

//delete station pic
func (handler GroupMemberHandler) Delete(c *gin.Context) {
	userId, _ := strconv.Atoi(c.PostForm("user_id"))
	groupMember := m.GroupMember{}
	groupMemberQuery := handler.db.Where("user_id = ?", userId).First(&groupMember)

	if groupMemberQuery.RowsAffected > 0 {
		deleteResult := handler.db.Delete(&groupMember)
		if deleteResult.RowsAffected > 0 {
			respond(http.StatusOK, "Member successfully removed from group.", c, true)
		} else {
			respond(http.StatusBadRequest, deleteResult.Error.Error(), c, true)
		}	
	} else {
		respond(http.StatusNotFound, "Record not found.", c, true)
	}
}
