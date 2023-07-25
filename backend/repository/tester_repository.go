package repository

import (
	"di/model"
	"errors"
	"log"

	"gorm.io/gorm"
)

type testerRepositoryImpl struct {
	DB *gorm.DB
}

func NewTesterRepository(gormDB *gorm.DB) model.TesterRepository {
	return &testerRepositoryImpl{
		DB: gormDB,
	}
}

func (repo *testerRepositoryImpl) FindByID(id uint) (*model.Tester, error) {
	var tester = model.Tester{}

	result := repo.DB.Preload("User").First(&tester, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &tester, nil
}

func (repo *testerRepositoryImpl) FindByOwner(ownerId uint) ([]model.Tester, error) {

	var testers []model.Tester

	result := repo.DB.Preload("User").Where("user_id = ?", ownerId).Order("id desc").Find(&testers)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return testers, nil
}

func (repo *testerRepositoryImpl) Create(tester *model.Tester) error {
	result := repo.DB.Create(tester)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *testerRepositoryImpl) Update(tester *model.Tester) error {
	result := repo.DB.Save(tester)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		repo.Create(tester)
	}

	return nil
}

func (repo *testerRepositoryImpl) Delete(id uint) error {
	result := repo.DB.Delete(&model.Tester{}, id)

	if result.Error != nil {
		log.Printf("Failed to delete pipeline. Reason: %v\n", result.Error)
		return result.Error
	}

	return nil
}
