package dto

type ShelterRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=100"`
	City     string `json:"city" binding:"required,min=1,max=100"`
	Address  string `json:"address" binding:"required,min=1,max=200"`
	Capacity uint64 `json:"capacity" binding:"required,gt=0" minimum:"1"`
}

type ShelterResponse struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Capacity uint64 `json:"capacity"`
}
