package main

import (
	"embed"
	"os"
	"path/filepath"
	"strings"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

type DashboardMsg struct {
	Path    string `json:"path"`
	State   int    `json:"state"`
	Total   int    `json:"total"`
	Success int    `json:"success"`
	Fail    int    `json:"fail"`
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	AppMenu := menu.NewMenu()
	FileMenu := AppMenu.AddSubmenu("File")

	FileMenu.AddText("文件转换", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		files, err := runtime.OpenMultipleFilesDialog(app.ctx, runtime.OpenDialogOptions{
			DefaultDirectory:     "",
			Title:                "Select HEIC Files",
			CanCreateDirectories: true,
			ShowHiddenFiles:      true,
			Filters: []runtime.FileFilter{
				{
					DisplayName: "Images (*.HEIC;*.heic)",
					Pattern:     "*.HEIC;*.heic",
				},
			},
		})
		if err != nil {
			panic("crashed!")
		}
		if len(files) > 0 {
			total := len(files)
			success := 0
			fail := 0
			path := filepath.Dir(files[0])

			runtime.EventsEmit(app.ctx, "ProcessStart", DashboardMsg{
				State:   1,
				Path:    path,
				Total:   total,
				Success: 0,
				Fail:    0,
			})
			for i := 0; i < len(files); i++ {
				err := app.Heic2jpg(files[i])
				if err != "" {
					fail++
				} else {
					success++
				}
				runtime.EventsEmit(app.ctx, "ProcessUpdate", DashboardMsg{
					State:   1,
					Path:    path,
					Total:   total,
					Success: success,
					Fail:    fail,
				})
			}
			runtime.EventsEmit(app.ctx, "ProcessDone", DashboardMsg{
				State:   2,
				Path:    path,
				Total:   total,
				Success: success,
				Fail:    fail,
			})
		}
	})

	FileMenu.AddText("文件夹转换", keys.CmdOrCtrl("d"), func(_ *menu.CallbackData) {
		dir, err := runtime.OpenDirectoryDialog(app.ctx, runtime.OpenDialogOptions{
			DefaultDirectory:     "",
			Title:                "Select HEIC Files",
			CanCreateDirectories: true,
			ShowHiddenFiles:      true,
			Filters: []runtime.FileFilter{
				{
					DisplayName: "Images (*.HEIC;*.heic)",
					Pattern:     "*.HEIC;*.heic",
				},
			},
		})
		if err != nil {
			panic("crashed!")
		}
		if dir != "" {
			total := 0
			success := 0
			fail := 0
			filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				// 检查文件名是否以.heic结尾
				if strings.HasSuffix(info.Name(), ".HEIC") || strings.HasSuffix(info.Name(), ".heic") {
					total++
				}
				return nil
			})
			if total == 0 {
				return
			}

			runtime.EventsEmit(app.ctx, "ProcessStart", DashboardMsg{
				State:   1,
				Path:    dir,
				Total:   total,
				Success: 0,
				Fail:    0,
			})

			filepath.Walk(dir, func(filePath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				// 检查文件名是否以.heic结尾
				if strings.HasSuffix(info.Name(), ".HEIC") || strings.HasSuffix(info.Name(), ".heic") {
					aerr := app.Heic2jpg(filePath)
					if aerr != "" {
						fail++
					} else {
						success++
					}
					runtime.EventsEmit(app.ctx, "ProcessUpdate", DashboardMsg{
						State:   1,
						Path:    dir,
						Total:   total,
						Success: success,
						Fail:    fail,
					})

				}
				return nil
			})
			runtime.EventsEmit(app.ctx, "ProcessDone", DashboardMsg{
				State:   2,
				Path:    dir,
				Total:   total,
				Success: success,
				Fail:    fail,
			})
		}
	})

	FileMenu.AddSeparator()
	FileMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(app.ctx)
	})
	// FileMenu.AddSeparator()
	// FileMenu.AddText("Setting", keys.CmdOrCtrl("s"), func(_ *menu.CallbackData) {
	// 	dir, err := runtime.OpenDirectoryDialog(app.ctx, runtime.OpenDialogOptions{
	// 		DefaultDirectory:     "",
	// 		Title:                "",
	// 		CanCreateDirectories: true,
	// 		ShowHiddenFiles:      true,
	// 	})
	// 	if err != nil {
	// 		panic("奔溃了")
	// 	}
	// 	if dir != "" {
	// 		print(dir)
	// 		runtime.EventsEmit(app.ctx, "SettingSavePath", dir)
	// 	}
	// })

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "HEIC转换器",
		Width:  550,
		Height: 300,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Menu:             AppMenu,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
