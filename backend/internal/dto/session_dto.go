package dto

type SessionRegisterRequest struct {
	Name        string `json:"name" binding:"required,min=1"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password" binding:"required,min=6"`
}

type SessionLoginRequest struct {
	Name     string `json:"name" binding:"required,min=1"`
	Password string `json:"password" binding:"required,min=6"`
}
