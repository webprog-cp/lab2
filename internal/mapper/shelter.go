package mapper

import (
	"cp_lab1/internal/dto"
	"cp_lab1/internal/models"
)

func ShelterFromRequest(request dto.ShelterRequest) models.Shelter {
	return models.Shelter{
		Name:     request.Name,
		City:     request.City,
		Address:  request.Address,
		Capacity: request.Capacity,
	}
}

func ShelterResponse(shelter models.Shelter) dto.ShelterResponse {
	return dto.ShelterResponse{
		ID:       shelter.ID,
		Name:     shelter.Name,
		City:     shelter.City,
		Address:  shelter.Address,
		Capacity: shelter.Capacity,
	}
}

func ShelterResponses(shelters []models.Shelter) []dto.ShelterResponse {
	responses := make([]dto.ShelterResponse, len(shelters))
	for i, shelter := range shelters {
		responses[i] = ShelterResponse(shelter)
	}
	return responses
}
