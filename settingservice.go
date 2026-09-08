package main

import "todo/internal/repository"

type SettingService struct {
	settings *repository.SettingRepository
}

func NewSettingService(settings *repository.SettingRepository) *SettingService {
	return &SettingService{settings: settings}
}

func (s *SettingService) SetShowCompleted(showCompleted string) error {
	return s.settings.SetShowCompleted(showCompleted)
}

func (s *SettingService) GetShowCompleted() (string, error) {
	return s.settings.GetShowCompleted()
}
