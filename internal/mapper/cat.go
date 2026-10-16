package mapper

import (
	"cp_lab1/internal/dto"
	"cp_lab1/internal/models"
)

func CatFromRequest(request dto.CatRequest) models.Cat {
	return models.Cat{
		Name:      request.Name,
		AgeMonths: *request.AgeMonths,
		Gender:    request.Gender,
		IsAdopted: request.IsAdopted,
		ShelterID: request.ShelterID,
	}
}

func CatResponse(cat models.Cat) dto.CatResponse {
	return dto.CatResponse{
		ID:        cat.ID,
		Name:      cat.Name,
		AgeMonths: cat.AgeMonths,
		Gender:    cat.Gender,
		IsAdopted: cat.IsAdopted,
		ShelterID: cat.ShelterID,
	}
}

func CatResponses(cats []models.Cat) []dto.CatResponse {
	responses := make([]dto.CatResponse, len(cats))
	for i, cat := range cats {
		responses[i] = CatResponse(cat)
	}
	return responses
}
