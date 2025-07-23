package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"` // This will be the plain password from request, not hashed
}
