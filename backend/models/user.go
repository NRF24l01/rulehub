package models

type User struct {
	BaseModel
	Username string `gorm:"type:varchar(32);unique;not null" json:"username"`
	Password string `gorm:"type:varchar(45);not null" json:"password"`
}