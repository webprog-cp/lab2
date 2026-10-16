package app

import (
	"fmt"

	"cp_lab1/internal/config"
	h "cp_lab1/internal/handlers"
	"cp_lab1/internal/models"
	r "cp_lab1/internal/repositories"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	CatHandler     *h.CatHandler
	ShelterHandler *h.ShelterHandler
}

func New(cfg *config.Config) (*App, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := db.AutoMigrate(&models.Shelter{}, &models.Cat{}); err != nil {
		return nil, fmt.Errorf("migrate database: %w", err)
	}

	catRepo := r.NewCatRepoPsql(db)
	shelterRepo := r.NewShelterRepoPsql(db)

	catHandler := h.NewCatHandler(catRepo, shelterRepo)
	shelterHandler := h.NewShelterHandler(shelterRepo, catRepo)

	return &App{
		CatHandler:     catHandler,
		ShelterHandler: shelterHandler,
	}, nil
}
