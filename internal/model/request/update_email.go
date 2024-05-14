package request

type UpdateEmail struct {
	Email string `json:"email" binding:"required"`
}
