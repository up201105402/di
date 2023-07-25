package repository

import (
	"di/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

type trainedRepositoryImpl struct {
	DB *gorm.DB
}

func NewTrainedRepository(gormDB *gorm.DB) model.TrainedRepository {
	return &trainedRepositoryImpl{
		DB: gormDB,
	}
}

func (repo *trainedRepositoryImpl) FindByID(id uint) (*model.Trained, error) {
	var trained = model.Trained{}

	result := repo.DB.Preload("User").First(&trained, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &trained, nil
}

func (repo *trainedRepositoryImpl) FindByOwner(ownerId uint) ([]model.Trained, error) {

	var traineds []model.Trained

	result := repo.DB.Preload("User").Where("user_id = ?", ownerId).Order("id desc").Find(&traineds)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return traineds, nil
}

func (repo *trainedRepositoryImpl) Create(trained *model.Trained) error {
	result := repo.DB.Create(trained)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *trainedRepositoryImpl) Update(trained *model.Trained) error {
	result := repo.DB.Save(trained)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		repo.Create(trained)
	}

	return nil
}

func (repo *trainedRepositoryImpl) Delete(id uint) error {
	result := repo.DB.Delete(&model.Trained{}, id)

	if result.Error != nil {
		log.Printf("Failed to delete pipeline. Reason: %v\n", result.Error)
		return result.Error
	}

	return nil
}
