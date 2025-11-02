package requests

type CreateAnamnesisRequest struct {
	Text string `json:"text" binding:"required"`
}

type UpdateAnamnesisRequest struct {
	Text string `json:"text" binding:"required"`
}
