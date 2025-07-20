package models

type Article struct {
	BaseModel
	Title   string `gorm:"type:varchar(128);not null" json:"title"`
	Content string `gorm:"type:text;not null" json:"content"`
	UserID string `gorm:"not null" json:"user_id"`
	User   User `gorm:"foreignKey:UserID" json:"user"`
}