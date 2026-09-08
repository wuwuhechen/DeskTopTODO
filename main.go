package main

import (
	"embed"
	"todo/internal/database"
	"todo/internal/repository"

	"log"
	"os"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed frontend/public/logo.png
var logo []byte

func init() {
	application.RegisterEvent[string]("time")
	application.RegisterEvent[bool]("autostart-changed")
}

func startHidden() bool {
	for _, arg := range os.Args[1:] {
		if arg == "--hidden" {
			return true
		}
	}
	return false
}

func main() {
	db, err := database.Open()

	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	todoRepository := repository.NewTodoRepository(db)
	settingRepository := repository.NewSettingRepository(db)
	windowRepository := repository.NewWindowRepository(db)

	todoService := NewTodoService(todoRepository)
	settingService := NewSettingService(settingRepository)
	windowService := NewWindowService(windowRepository)

	app := application.New(application.Options{
		Name:        "DesktopTODO",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(todoService),
			application.NewService(settingService),
			application.NewService(windowService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	window := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:          "DesktopTODO",
		Width:          340,
		Height:         460,
		AlwaysOnTop:    true,
		MinWidth:       320,
		MinHeight:      320,
		Frameless:      true,
		BackgroundType: application.BackgroundTypeTransparent,
		Windows: application.WindowsWindow{
			NonClientRegionSupport: true,
		},
	})

	windowService.AttachWindow(app, window)

	app.OnShutdown(func() {
		if err := windowService.SaveNow(); err != nil {
			app.Logger.Error("final window state save failed", "error", err)
		}
	})

	windowService.SetupSystemTray()

	if startHidden() {
		window.Hide()
	} else {
		window.Show()
	}

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			app.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	err = app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
