package models

type User struct{
	Id int `gorm:"primaryKey;autoIncrement" json:"id"`
	Username string `json:"username"`
	Email string `gorm:"unique" json:"email"`
	Password []byte ` json:"-"`
	
}