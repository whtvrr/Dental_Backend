package requests

type CreateComplaintRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
}

type UpdateComplaintRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description,omitempty"`
	Category    *string `json:"category,omitempty"`
}