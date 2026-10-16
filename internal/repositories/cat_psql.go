package repositories

import (
	"errors"
	"fmt"

	"cp_lab1/internal/const_errors"
	"cp_lab1/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CatRepoPsql struct {
	db *gorm.DB
}

func NewCatRepoPsql(db *gorm.DB) *CatRepoPsql {
	return &CatRepoPsql{
		db: db,
	}
}

func (cr *CatRepoPsql) GetAll() ([]models.Cat, error) {
	var cats []models.Cat
	if err := cr.db.Order("id").Find(&cats).Error; err != nil {
		return nil, fmt.Errorf("get all cats: %w", err)
	}

	return cats, nil
}

func (cr *CatRepoPsql) GetByID(id uint64) (models.Cat, error) {
	var cat models.Cat
	if err := cr.db.First(&cat, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Cat{}, fmt.Errorf("cat %d: %w", id, const_errors.NoEntityByID)
		}
		return models.Cat{}, fmt.Errorf("get cat %d: %w", id, err)
	}

	return cat, nil
}

func (cr *CatRepoPsql) GetCatsByShelterID(shelterID uint64) ([]models.Cat, error) {
	var cats []models.Cat
	if err := cr.db.Where("shelter_id = ?", shelterID).Order("id").Find(&cats).Error; err != nil {
		return nil, fmt.Errorf("get cats for shelter %d: %w", shelterID, err)
	}

	return cats, nil
}

func shelterHasCapacity(tx *gorm.DB, shelterID uint64) error {
	var shelter models.Shelter
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&shelter, shelterID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("shelter %d: %w", shelterID, const_errors.InvalidReference)
		}
		return fmt.Errorf("get shelter %d: %w", shelterID, err)
	}

	var count int64
	if err := tx.Model(&models.Cat{}).Where("shelter_id = ?", shelterID).Count(&count).Error; err != nil {
		return fmt.Errorf("count cats in shelter %d: %w", shelterID, err)
	}
	if uint64(count) >= shelter.Capacity {
		return fmt.Errorf("shelter %d: %w", shelterID, const_errors.ShelterAtCapacity)
	}

	return nil
}

func (cr *CatRepoPsql) Create(cat models.Cat) (models.Cat, error) {
	err := cr.db.Transaction(func(tx *gorm.DB) error {
		if err := shelterHasCapacity(tx, cat.ShelterID); err != nil {
			return err
		}
		return tx.Create(&cat).Error
	})
	if err != nil {
		return models.Cat{}, fmt.Errorf("create cat: %w", err)
	}

	return cat, nil
}

func (cr *CatRepoPsql) Update(id uint64, cat models.Cat) error {
	err := cr.db.Transaction(func(tx *gorm.DB) error {
		var existing models.Cat
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&existing, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("cat %d: %w", id, const_errors.NoEntityByID)
			}
			return err
		}

		if existing.ShelterID != cat.ShelterID {
			if err := shelterHasCapacity(tx, cat.ShelterID); err != nil {
				return err
			}
		}

		return tx.Model(&existing).Updates(map[string]any{
			"name":       cat.Name,
			"age_months": cat.AgeMonths,
			"gender":     cat.Gender,
			"is_adopted": cat.IsAdopted,
			"shelter_id": cat.ShelterID,
		}).Error
	})
	if err != nil {
		return fmt.Errorf("update cat %d: %w", id, err)
	}

	return nil
}

func (cr *CatRepoPsql) Delete(id uint64) error {
	result := cr.db.Delete(&models.Cat{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete cat %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("cat %d: %w", id, const_errors.NoEntityByID)
	}

	return nil
}

func (cr *CatRepoPsql) DeleteAll() error {
	if err := cr.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Cat{}).Error; err != nil {
		return fmt.Errorf("delete all cats: %w", err)
	}

	return nil
}
