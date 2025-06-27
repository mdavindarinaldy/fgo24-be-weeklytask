package models

type User struct {
	Id          int    `db:"id"`
	Name        string `form:"name" json:"name" db:"name" binding:"required"`
	Email       string `form:"email" json:"email" db:"email" binding:"required,email"`
	PhoneNumber string `form:"phoneNumber" json:"phoneNumber" db:"phone_number" binding:"required"`
	Password    string `form:"password" json:"password" db:"password" binding:"required"`
	Pin         string `form:"pin" json:"pin" db:"pin" binding:"required"`
}

// var Token *string = new(string)
