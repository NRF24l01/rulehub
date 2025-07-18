package models

type Media struct {
	BaseModel
	FileName string `gorm:"type:varchar(128);not null" json:"file_name"`
	S3Key  string `gorm:"type:varchar(256);not null" json:"s3_key"`
	ArticleID string `gorm:"not null" json:"article_id"`
	Article   Article `gorm:"foreignKey:ArticleID" json:"article"`
}