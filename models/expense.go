package models

type Expense struct {
	BaseModel
	ExpenseName string `json:"name" form:"expense_name" binding:"required"`
	ExpenseDescription string `json:"description" form:"expense_description" binding:"required"`
	ExpenseType string `json:"type" form:"type" binding:"required"`
	Amount float64 `json:"amount" form:"amount" binding:"required"` 
	DueDate string `json:"due_date" form:"due_date" binding:"required"`
	Group Group `json:"group" binding:"-"`
	GroupId int `json:"-" form:"group" binding:"required"`
	Overdue int `json:"overdue"`
	Status string `json:"status"`
	User User `json:"created_by" binding:"-"`
	UserId int `json:"-" form:"created_by" binding:"required"`
}

func (e *Expense) BeforeCreate() (err error) {
	e.Status = "Pending"
	return
}
