package schemas

// MediaUploadResponse is the response schema for generating temporary upload URLs
type MediaUploadResponse struct {
	TempURL string `json:"temp_url"`
	FileID  string `json:"file_id"`
}
