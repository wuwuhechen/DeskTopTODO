package repository

import (
	"encoding/json"
	"todo/internal/model"

	"gorm.io/gorm"
)

type WindowRepository struct {
	db *gorm.DB
}

func NewWindowRepository(db *gorm.DB) *WindowRepository {
	return &WindowRepository{db: db}
}

func (r *WindowRepository) LoadWindowSettings() (*model.WindowState, error) {
	var windowSetting model.Setting
	var windowState model.WindowState

	var jsonData []byte

	if err := r.db.Where("key = ?", "window_setting").First(&windowSetting).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			windowState = model.WindowState{
				Width:  340,
				Height: 460,
				X:      100,
				Y:      100,
				Pinned: false,
				Locked: false,
			}
			jsonData, err = json.Marshal(windowState)
			if err != nil {
				return nil, err
			}
			windowSetting = model.Setting{
				Key:   "window_setting",
				Value: string(jsonData),
			}
			if err := r.db.Create(&windowSetting).Error; err != nil {
				return nil, err
			}

			return &windowState, nil
		} else {
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(windowSetting.Value), &windowState); err != nil {
		return nil, err
	}

	return &windowState, nil
}

func (r *WindowRepository) SaveWindowSettings(state *model.WindowState) error {
	jsonData, err := json.Marshal(state)
	if err != nil {
		return err
	}

	setting := model.Setting{
		Key:   "window_setting",
		Value: string(jsonData),
	}

	return r.db.Where("key = ?", setting.Key).
		Assign(model.Setting{Value: setting.Value}).
		FirstOrCreate(&setting).Error
}
