package dto

type CatRequest struct {
	Name      string  `json:"name" binding:"required,min=1,max=100"`
	AgeMonths *uint16 `json:"ageMonths" binding:"required,lte=360"`
	Gender    string  `json:"gender" binding:"required,oneof=male female"`
	IsAdopted bool    `json:"adopted"`
	ShelterID uint64  `json:"shelterID" binding:"required" minimum:"1"`
}

type CatResponse struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	AgeMonths uint16 `json:"ageMonths"`
	Gender    string `json:"gender"`
	IsAdopted bool   `json:"adopted"`
	ShelterID uint64 `json:"shelterID"`
}
