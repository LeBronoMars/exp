package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	m "exp/mngr/api/models"
)

type ExpenseHandler struct {
	db *gorm.DB
}

func NewExpenseHandler(db *gorm.DB) *ExpenseHandler {
	return &ExpenseHandler{db}
}

//get all expense
func (handler ExpenseHandler) Index(c *gin.Context) {
	expenses := []m.Expense{}		
	//handler.db.Preload("CreatedBy").Preload("Group").Preload("Group.Members").Preload("Group.Members.User").Find(&expenses)
	handler.db.Preload("Group").Preload("User").Find(&expenses)
	c.JSON(http.StatusOK, expenses)
	return
}

//create new expense
func (handler ExpenseHandler) Create(c *gin.Context) {
	var expense m.Expense
	err := c.Bind(&expense)

	if err == nil {
		//check for existing group
		existingGroup := m.Group{}
		existingGroupQuery := handler.db.Where("id = ?", expense.GroupId).First(&existingGroup)
		
		if existingGroupQuery.RowsAffected == 0 {
			respond(http.StatusNotFound, "Group not found.", c, true)
		} else {
			//check for existing user
			existingUser := m.User{}
			existingUserQuery := handler.db.Where("id = ?", expense.UserId).First(&existingUser)

			if existingUserQuery.RowsAffected == 0 {
				respond(http.StatusNotFound, "User not found.", c, true)
			} else {
				saveResult := handler.db.Save(&expense)
				if saveResult.RowsAffected > 0 {
					respond(http.StatusCreated, "New expense successfully created.", c, false)
				} else {
					respond(http.StatusBadRequest, saveResult.Error.Error(), c, true)
				}
			}
		}
	} else {
		respond(http.StatusBadRequest, err.Error(), c, true)
	}
}

//search for group
// func (handler ExpenseHandler) Search(c *gin.Context) {
// 	groups := []m.Group{}		
// 	res := handler.db.Preload("Members").Preload("Members.User").Where("name LIKE ?", "%"+ c.Param("group_name") +"%").Find(&groups)
// 	if res.RowsAffected > 0 {
// 		c.JSON(http.StatusOK, groups)	
// 	} else {
// 		respond(http.StatusNotFound, "No records found.", c, true)
// 	}
// }