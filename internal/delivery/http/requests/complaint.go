package requests

type CreateComplaintRequest struct {
	Title string `json:"title" binding:"required"`
}

type UpdateComplaintRequest struct {
	Title string `json:"title" binding:"required"`
}