package schemas

type ArticleCreateRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=128"`
	Content     string `json:"content" validate:"required,min=1,max=10000"`
	Media       []string `json:"media" validate:"omitempty,dive,required,min=1,max=128"` // Filenames list
}

type MediaCreateResponse struct {
	FileName string `json:"file_name"`
	S3Key    string `json:"s3_key"`
}

type ArticleCreateResponse struct {
	ID             string              `json:"id"`
	Title          string              `json:"title"`
	Content        string              `json:"content"`
	Media          []MediaCreateResponse `json:"media"`
	AuthorUsername string              `json:"author_username"`
}