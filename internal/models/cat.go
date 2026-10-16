package models

type Cat struct {
	ID        uint64  `json:"id"`
	Name      string  `json:"name"`
	AgeMonths uint16  `json:"ageMonths"`
	Gender    string  `json:"gender"`
	IsAdopted bool    `json:"adopted"`
	ShelterID uint64  `json:"shelterID" gorm:"not null;index"`
	Shelter   Shelter `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`
}
