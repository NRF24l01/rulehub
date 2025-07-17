package schemas

type SignInRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=6,max=128,strongpwd"`
}

type SignInResponse struct {
	AccessJWT string `json:"access_token"`
}

type SignUpRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Password string `json:"password" validate:"required,min=6,max=128,strongpwd"`
}

type SignUpResponse struct {
	ID	     string `json:"id"`
	Username string `json:"username"`
}