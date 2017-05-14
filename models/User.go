package models

import (
	"crypto/md5"
	"fmt"
)

type User struct {
	BaseModel
	FirstName  string `json:"first_name" form:"first_name" binding:"required"`
	LastName  string `json:"last_name" form:"last_name" binding:"required"`	
	MiddleName string `json:"middle_name" form:"middle_name"`
	Gender string `json:"gender" form:"gender"`
	Email string `json:"email" form:"email" sql:"type:varchar(255);index;not null;unique"`
	Address string `json:"address" form:"address"`
	ContactNo string `json:"contact_no" form:"contact_no" binding:"required"`
	Status string `json:"status"`
	UserRole string `json:"user_role" form:"user_role" binding:"required"`
	Password string `json:"-" form:"password" binding:"required"`
	PicUrl string `json:"pic_url" form:"pic_url"`
}

func (u *User) BeforeCreate() (err error) {
	u.Status = "active"
	defaultPic := fmt.Sprintf("%x", md5.Sum([]byte(u.Email)))
	u.PicUrl = fmt.Sprintf("http://www.gravatar.com/avatar/%s?d=identicon", defaultPic)
	return
}