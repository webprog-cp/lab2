package handlers

import (
	"errors"
	"net/http"

	"cp_lab1/internal/const_errors"
	"cp_lab1/internal/dto"
	"cp_lab1/internal/mapper"
	"cp_lab1/internal/repositories"

	"github.com/gin-gonic/gin"
)

type CatHandler struct {
	repo        repositories.CatRepo
	shelterRepo repositories.ShelterRepo
}

func NewCatHandler(repo repositories.CatRepo, shelterRepo repositories.ShelterRepo) *CatHandler {
	return &CatHandler{
		repo:        repo,
		shelterRepo: shelterRepo,
	}
}

// GetAll godoc
// @Summary Get all cats
// @Tags cats
// @Produce json
// @Success 200 {array} dto.CatResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cats [get]
func (ch *CatHandler) GetAll(c *gin.Context) {
	cats, err := ch.repo.GetAll()
	if err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.CatResponses(cats))
}

// GetByID godoc
// @Summary Get cat by ID
// @Tags cats
// @Produce json
// @Param id path int true "Cat ID"
// @Success 200 {object} dto.CatResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cats/{id} [get]
func (ch *CatHandler) GetByID(c *gin.Context) {
	id, valid := idParam(c)
	if !valid {
		return
	}

	cat, err := ch.repo.GetByID(id)
	if err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.CatResponse(cat))
}

// Create godoc
// @Summary Create cat
// @Tags cats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param cat body dto.CatRequest true "Cat"
// @Success 201 {object} dto.CatResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cats [post]
func (ch *CatHandler) Create(c *gin.Context) {
	var request dto.CatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	if err := ch.shelterHasCapacity(request.ShelterID); err != nil {
		respondRepositoryError(c, err)
		return
	}

	cat, err := ch.repo.Create(mapper.CatFromRequest(request))
	if err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.CatResponse(cat))
}

// Update godoc
// @Summary Update cat
// @Tags cats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Param cat body dto.CatRequest true "Cat"
// @Success 200 {object} dto.CatResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cats/{id} [put]
func (ch *CatHandler) Update(c *gin.Context) {
	id, valid := idParam(c)
	if !valid {
		return
	}

	var request dto.CatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "invalid request body"})
		return
	}

	existing, err := ch.repo.GetByID(id)
	if err != nil {
		respondRepositoryError(c, err)
		return
	}
	if existing.ShelterID != request.ShelterID {
		if err := ch.shelterHasCapacity(request.ShelterID); err != nil {
			respondRepositoryError(c, err)
			return
		}
	}

	cat := mapper.CatFromRequest(request)
	if err := ch.repo.Update(id, cat); err != nil {
		respondRepositoryError(c, err)
		return
	}

	cat.ID = id
	c.JSON(http.StatusOK, mapper.CatResponse(cat))
}

func (ch *CatHandler) shelterHasCapacity(shelterID uint64) error {
	shelter, err := ch.shelterRepo.GetByID(shelterID)
	if errors.Is(err, const_errors.NoEntityByID) {
		return const_errors.InvalidReference
	}
	if err != nil {
		return err
	}

	cats, err := ch.repo.GetCatsByShelterID(shelterID)
	if err != nil {
		return err
	}
	if uint64(len(cats)) >= shelter.Capacity {
		return const_errors.ShelterAtCapacity
	}

	return nil
}

// Delete godoc
// @Summary Delete cat
// @Tags cats
// @Security ApiKeyAuth
// @Param id path int true "Cat ID"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cats/{id} [delete]
func (ch *CatHandler) Delete(c *gin.Context) {
	id, valid := idParam(c)
	if !valid {
		return
	}

	if err := ch.repo.Delete(id); err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteAll godoc
// @Summary Delete all cats
// @Tags cats
// @Security ApiKeyAuth
// @Success 204
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /cats [delete]
func (ch *CatHandler) DeleteAll(c *gin.Context) {
	if err := ch.repo.DeleteAll(); err != nil {
		respondRepositoryError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
