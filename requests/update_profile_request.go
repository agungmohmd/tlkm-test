package request

type UpdateProfileRequest struct {
	Name     string `json:"name" validate:"required"`
	Username string `json:"username" validate:"required"`
	Address  string `json:"address" validate:"required"`
}
