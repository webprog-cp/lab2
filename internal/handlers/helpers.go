package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"cp_lab1/internal/const_errors"
	"cp_lab1/internal/dto"

	"github.com/gin-gonic/gin"
)

func idParam(c *gin.Context) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "id must be an unsigned integer"})
		return 0, false
	}
	return id, true
}

func respondRepositoryError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, const_errors.NoEntityByID):
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "resource not found"})
	case errors.Is(err, const_errors.InvalidReference):
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "related resource not found"})
	case errors.Is(err, const_errors.EntityHasDependents):
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "resource has dependent records"})
	case errors.Is(err, const_errors.ShelterAtCapacity):
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "shelter is at full capacity"})
	case errors.Is(err, const_errors.CapacityBelowUsage):
		c.JSON(http.StatusConflict, dto.ErrorResponse{Error: "capacity is below current occupancy"})
	default:
		log.Printf("%s %s: %v", c.Request.Method, c.Request.URL.Path, err)
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "internal server error"})
	}
}
