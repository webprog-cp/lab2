package repositories

import (
	"errors"
	"fmt"

	"cp_lab1/internal/const_errors"
	"cp_lab1/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ShelterRepoPsql struct {
	db *gorm.DB
}

func NewShelterRepoPsql(db *gorm.DB) *ShelterRepoPsql {
	return &ShelterRepoPsql{db: db}
}

func (sr *ShelterRepoPsql) GetAll() ([]models.Shelter, error) {
	var shelters []models.Shelter
	if err := sr.db.Order("id").Find(&shelters).Error; err != nil {
		return nil, fmt.Errorf("get all shelters: %w", err)
	}

	return shelters, nil
}

func (sr *ShelterRepoPsql) GetByID(id uint64) (models.Shelter, error) {
	var shelter models.Shelter
	if err := sr.db.First(&shelter, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Shelter{}, fmt.Errorf("shelter %d: %w", id, const_errors.NoEntityByID)
		}
		return models.Shelter{}, fmt.Errorf("get shelter %d: %w", id, err)
	}

	return shelter, nil
}

func (sr *ShelterRepoPsql) Create(shelter models.Shelter) (models.Shelter, error) {
	if err := sr.db.Create(&shelter).Error; err != nil {
		return models.Shelter{}, fmt.Errorf("create shelter: %w", err)
	}

	return shelter, nil
}

func (sr *ShelterRepoPsql) Update(id uint64, shelter models.Shelter) error {
	err := sr.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Shelter
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&existing, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("shelter %d: %w", id, const_errors.NoEntityByID)
			}
			return err
		}

		var count int64
		if err := tx.Model(&models.Cat{}).Where("shelter_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if uint64(count) > shelter.Capacity {
			return const_errors.CapacityBelowUsage
		}

		return tx.Model(&existing).Updates(map[string]any{
			"name":     shelter.Name,
			"city":     shelter.City,
			"address":  shelter.Address,
			"capacity": shelter.Capacity,
		}).Error
	})
	if err != nil {
		return fmt.Errorf("update shelter %d: %w", id, err)
	}

	return nil
}

func (sr *ShelterRepoPsql) Delete(id uint64) error {
	result := sr.db.Delete(&models.Shelter{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrForeignKeyViolated) {
			return fmt.Errorf("shelter %d: %w", id, const_errors.EntityHasDependents)
		}
		return fmt.Errorf("delete shelter %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("shelter %d: %w", id, const_errors.NoEntityByID)
	}

	return nil
}

func (sr *ShelterRepoPsql) DeleteAll() error {
	if err := sr.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Shelter{}).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return fmt.Errorf("delete all shelters: %w", const_errors.EntityHasDependents)
		}
		return fmt.Errorf("delete all shelters: %w", err)
	}

	return nil
}
