package repository

import (
	"todo/internal/model"

	"gorm.io/gorm"
)

type SettingRepository struct {
	db *gorm.DB
}

func NewSettingRepository(db *gorm.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

func (r *SettingRepository) SetShowCompleted(showCompleted string) error {
	setting := model.Setting{Key: "show_completed"}
	return r.db.Where("key = ?", setting.Key).
		Assign(model.Setting{Value: showCompleted}).
		FirstOrCreate(&setting).Error
}

func (r *SettingRepository) GetShowCompleted() (string, error) {
	setting := model.Setting{Key: "show_completed"}
	if err := r.db.First(&setting, "key = ?", setting.Key).Error; err != nil {
		return "", err
	}
	return setting.Value, nil
}
