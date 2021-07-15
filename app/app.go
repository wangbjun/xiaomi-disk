package app

import (
	"bytes"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"time"
	"xiaomi_cloud/api"
	"xiaomi_cloud/theme"
)

const previousSelect = "previous_select"

type App struct {
	Window fyne.Window
	Center *fyne.Container
	User   *api.User
	Api    *api.Api
	Menu   map[string]Menu
}

type Menu struct {
	Title string
	View  func(app *App) fyne.CanvasObject
}

func New() *App {
	newApp := app.NewWithID("xiaomi_disk")
	newApp.Settings().SetTheme(&theme.MyTheme{})
	user := api.NewUser()
	return &App{
		Window: newApp.NewWindow("MiCloud"),
		User:   user,
		Api:    api.NewApi(user),
		Menu: map[string]Menu{
			"welcome": {"欢迎", welcome},
			"disk":    {"云盘", disk},
		},
	}
}

func (app *App) ShowWindow() {
	content := container.NewMax()
	menus := func(t Menu, app *App) {
		content.Objects = []fyne.CanvasObject{t.View(app)}
		content.Refresh()
	}
	split := container.NewBorder(nil, nil, app.makeNav(menus), nil,
		container.NewBorder(nil, nil, widget.NewSeparator(), nil, content))
	app.Window.SetContent(split)
	app.Window.Resize(fyne.NewSize(800, 600))
	if !app.User.IsLogin {
		app.ShowLoginPop()
	}
	app.Window.ShowAndRun()
}

func (app *App) ShowLoginPop() {
	image, err := app.User.GetQrImage()
	if err != nil {
		app.Alert("获取二维码失败")
		return
	}
	qrImage := canvas.NewImageFromReader(bytes.NewReader(image), "download.png")
	qrImage.SetMinSize(fyne.NewSize(240, 240))
	card := widget.NewCard("扫码登录", "", qrImage)
	popUp := widget.NewModalPopUp(card, app.Window.Canvas())
	go func() {
		ticker := time.NewTicker(time.Second * 1)
		for {
			<-ticker.C
			if app.User.IsLogin {
				popUp.Hide()
			}
		}
	}()
	popUp.Show()
}

func (app *App) makeNav(setTree func(menu Menu, app *App)) fyne.CanvasObject {
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return []string{"welcome", "disk"}
		},
		IsBranch: func(uid string) bool {
			if uid == "" {
				return true
			}
			return false
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := app.Menu[uid]
			if !ok {
				fyne.LogError("Missing menu panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := app.Menu[uid]; ok {
				setTree(t, app)
			}
		},
	}
	tree.Select("welcome")
	return container.NewBorder(nil, nil, nil, nil, tree)
}

func (app *App) Alert(msg string) {
	dialog.ShowInformation("提示", msg, app.Window)
	return
}
