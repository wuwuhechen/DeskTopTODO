package main

import (
	"embed"
	"todo/internal/database"
	"todo/internal/repository"

	"log"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	application.RegisterEvent[string]("time")
}

func main() {
	db, err := database.Open()

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	todoRepo := repository.NewTodoRepository(db)
	todoService := NewTodoService(todoRepo)

	app := application.New(application.Options{
		Name:        "DesktopTODO",
		Description: "A demo of using raw HTML & CSS",
		Services: []application.Service{
			application.NewService(todoService),
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
		MinWidth:       280,
		MinHeight:      320,
		Frameless:      true,
		AlwaysOnTop:    false,
		BackgroundType: application.BackgroundTypeTransparent,
		Windows: application.WindowsWindow{
			NonClientRegionSupport: true,
		},
	})

	window.Center()
	window.Show()

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
