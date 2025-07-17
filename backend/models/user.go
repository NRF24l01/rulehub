package models

type User struct {
	BaseModel
	Username string `gorm:"type:varchar(32);unique;not null" json:"username"`
	Password string `gorm:"not null" json:"password"`
}