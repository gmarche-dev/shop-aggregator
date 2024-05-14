package request

type UpdatePassword struct {
	Password    string `json:"password" binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
}
