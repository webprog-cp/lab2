package handlers

import (
	"net/http"

	"cp_lab1/internal/const_errors"
	"cp_lab1/internal/dto"
	"cp_lab1/internal/mapper"
	"cp_lab1/internal/repositories"

	"github.com/gin-gonic/gin"
)

type ShelterHandler struct {
	shelterRepo repositories.ShelterRepo
	catRepo     repositories.CatRepo
}

func NewShelterHandler(sr repositories.ShelterRepo, cr repositories.CatRepo) *ShelterHandler {
	return &ShelterHandler{
		shelterRepo: sr,
		catRepo:     cr,
	}
}

// GetAll godoc
// @Summary Get all shelters
// @Tags shelters
// @Produce json
// @Success 200 {array} dto.ShelterResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shelters [get]
func (sh *ShelterHandler) GetAll(c *gin.Context) {
	shelters, err := sh.shelterRepo.GetAll()
	if err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.ShelterResponses(shelters))
}

// GetByID godoc
// @Summary Get shelter by ID
// @Tags shelters
// @Produce json
// @Param id path int true "Shelter ID"
// @Success 200 {object} dto.ShelterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shelters/{id} [get]
func (sh *ShelterHandler) GetByID(c *gin.Context) {
	id, valid := idParam(c)
	if !valid {
		return
	}

	shelter, err := sh.shelterRepo.GetByID(id)
	if err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.ShelterResponse(shelter))
}

// Create godoc
// @Summary Create shelter
// @Tags shelters
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param shelter body dto.ShelterRequest true "Shelter"
// @Success 201 {object} dto.ShelterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shelters [post]
func (sh *ShelterHandler) Create(c *gin.Context) {
	var request dto.ShelterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	shelter, err := sh.shelterRepo.Create(mapper.ShelterFromRequest(request))
	if err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.ShelterResponse(shelter))
}

// Update godoc
// @Summary Update shelter
// @Tags shelters
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Shelter ID"
// @Param shelter body dto.ShelterRequest true "Shelter"
// @Success 200 {object} dto.ShelterResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shelters/{id} [put]
func (sh *ShelterHandler) Update(c *gin.Context) {
	id, valid := idParam(c)
	if !valid {
		return
	}

	var request dto.ShelterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	cats, err := sh.catRepo.GetCatsByShelterID(id)
	if err != nil {
		respondRepositoryError(c, err)
		return
	}
	if uint64(len(cats)) > request.Capacity {
		respondRepositoryError(c, const_errors.CapacityBelowUsage)
		return
	}

	shelter := mapper.ShelterFromRequest(request)
	if err := sh.shelterRepo.Update(id, shelter); err != nil {
		respondRepositoryError(c, err)
		return
	}

	shelter.ID = id
	c.JSON(http.StatusOK, mapper.ShelterResponse(shelter))
}

// Delete godoc
// @Summary Delete shelter
// @Tags shelters
// @Security ApiKeyAuth
// @Param id path int true "Shelter ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shelters/{id} [delete]
func (sh *ShelterHandler) Delete(c *gin.Context) {
	id, valid := idParam(c)
	if !valid {
		return
	}

	if err := sh.shelterRepo.Delete(id); err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteAll godoc
// @Summary Delete all shelters
// @Tags shelters
// @Security ApiKeyAuth
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shelters [delete]
func (sh *ShelterHandler) DeleteAll(c *gin.Context) {
	if err := sh.shelterRepo.DeleteAll(); err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// GetCats godoc
// @Summary Get shelter cats
// @Tags shelters
// @Produce json
// @Param id path int true "Shelter ID"
// @Success 200 {array} dto.CatResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /shelters/{id}/cats [get]
func (sh *ShelterHandler) GetCats(c *gin.Context) {
	id, valid := idParam(c)
	if !valid {
		return
	}

	if _, err := sh.shelterRepo.GetByID(id); err != nil {
		respondRepositoryError(c, err)
		return
	}

	cats, err := sh.catRepo.GetCatsByShelterID(id)
	if err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.CatResponses(cats))
}
