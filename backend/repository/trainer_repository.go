package repository

import (
	"di/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

type trainerRepositoryImpl struct {
	DB *gorm.DB
}

func NewTrainerRepository(gormDB *gorm.DB) TrainerRepository {
	return &trainerRepositoryImpl{
		DB: gormDB,
	}
}

func (repo *trainerRepositoryImpl) FindByID(id uint) (*model.Trainer, error) {
	var trainer = model.Trainer{}

	result := repo.DB.Preload("User").First(&trainer, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &trainer, nil
}

func (repo *trainerRepositoryImpl) FindByOwner(ownerId uint) ([]model.Trainer, error) {

	var trainers []model.Trainer

	result := repo.DB.Preload("User").Where("user_id = ?", ownerId).Order("id desc").Find(&trainers)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return trainers, nil
}

func (repo *trainerRepositoryImpl) Create(trainer *model.Trainer) error {
	result := repo.DB.Create(trainer)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *trainerRepositoryImpl) Update(trainer *model.Trainer) error {
	result := repo.DB.Save(trainer)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		repo.Create(trainer)
	}

	return nil
}

func (repo *trainerRepositoryImpl) Delete(id uint) error {
	result := repo.DB.Delete(&model.Trainer{}, id)

	if result.Error != nil {
		log.Printf("Failed to delete pipeline. Reason: %v\n", result.Error)
		return result.Error
	}

	return nil
}
