package main

import (
	"errors"
	"sync"
	"time"

	"todo/internal/model"
	"todo/internal/repository"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

const (
	minWindowWidth  = 320
	minWindowHeight = 320
	maxWindowWidth  = 900
	maxWindowHeight = 1200
)

type WindowService struct {
	app      *application.App
	window   *application.WebviewWindow
	settings *repository.WindowRepository

	timerMu   sync.Mutex
	saveMu    sync.Mutex
	saveTimer *time.Timer
}

func NewWindowService(settings *repository.WindowRepository) *WindowService {
	return &WindowService{
		settings: settings,
	}
}

// TODO
func (s *WindowService) AttachWindow(
	app *application.App,
	window *application.WebviewWindow,
) {
	s.app = app
	s.window = window

	window.OnWindowEvent(events.Common.WindowRuntimeReady, func(
		_ *application.WindowEvent,
	) {
		if err := s.RestoreState(); err != nil {
			s.app.Logger.Error(
				"restore window state failed",
				"error", err,
			)
			s.window.Center()
		}
	})

	window.OnWindowEvent(events.Common.WindowDidMove, func(
		_ *application.WindowEvent,
	) {
		s.ScheduleSave()
	})

	window.OnWindowEvent(events.Common.WindowDidResize, func(
		_ *application.WindowEvent,
	) {
		s.ScheduleSave()
	})
}

func (s *WindowService) LoadWindowSettings() (*model.WindowState, error) {
	return s.settings.LoadWindowSettings()
}

func (s *WindowService) SetPinned(pinned bool) error {
	s.saveMu.Lock()
	defer s.saveMu.Unlock()

	if s.window == nil {
		return errors.New("window is not attached")
	}

	state, err := s.settings.LoadWindowSettings()
	if err != nil {
		return err
	}

	state.Pinned = pinned
	if err := s.settings.SaveWindowSettings(state); err != nil {
		return err
	}

	s.window.SetAlwaysOnTop(pinned)
	return nil
}

func (s *WindowService) SetLocked(locked bool) error {
	s.saveMu.Lock()
	defer s.saveMu.Unlock()

	if s.window == nil {
		return errors.New("window is not attached")
	}

	state, err := s.settings.LoadWindowSettings()
	if err != nil {
		return err
	}

	state.Locked = locked
	if err := s.settings.SaveWindowSettings(state); err != nil {
		return err
	}

	s.window.SetResizable(!locked)
	return nil
}

func (s *WindowService) RestoreState() error {
	state, err := s.settings.LoadWindowSettings()
	if err != nil {
		return err
	}

	state.Width = clamp(state.Width, minWindowWidth, maxWindowWidth)
	state.Height = clamp(state.Height, minWindowHeight, maxWindowHeight)

	s.window.SetSize(state.Width, state.Height)
	s.window.SetAlwaysOnTop(state.Pinned)
	s.window.SetResizable(!state.Locked)

	if !s.IsStateVisible(state) {
		s.window.Center()
		return nil
	}

	s.window.SetPosition(state.X, state.Y)
	return nil
}

func (s *WindowService) ScheduleSave() {
	s.timerMu.Lock()
	defer s.timerMu.Unlock()

	if s.saveTimer != nil {
		s.saveTimer.Stop()
	}

	s.saveTimer = time.AfterFunc(500*time.Millisecond, func() {
		if err := s.SaveNow(); err != nil {
			s.app.Logger.Error(
				"save window state failed",
				"error", err,
			)
		}
	})
}

func (s *WindowService) SaveNow() error {
	s.saveMu.Lock()
	defer s.saveMu.Unlock()

	if s.window == nil {
		return nil
	}

	state, err := s.settings.LoadWindowSettings()
	if err != nil {
		return err
	}

	state.X, state.Y = s.window.Position()
	state.Width, state.Height = s.window.Size()

	return s.settings.SaveWindowSettings(state)
}

func (s *WindowService) IsStateVisible(state *model.WindowState) bool {
	const minVisibleSize = 64

	for _, screen := range s.app.Screen.GetAll() {
		area := screen.WorkArea

		visibleWidth := min(
			state.X+state.Width,
			area.X+area.Width,
		) - max(state.X, area.X)

		visibleHeight := min(
			state.Y+state.Height,
			area.Y+area.Height,
		) - max(state.Y, area.Y)

		if visibleWidth >= minVisibleSize &&
			visibleHeight >= minVisibleSize {
			return true
		}
	}

	return false
}

func clamp(value, lower, upper int) int {
	if value < lower {
		return lower
	}
	if value > upper {
		return upper
	}
	return value
}
