package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/cmd/fyne_demo/data"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func welcome(app *App) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneScene)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(228, 167))
	return container.NewBorder(nil, nil, nil, nil,
		container.NewCenter(container.NewBorder(nil, widget.NewLabelWithStyle("欢迎使用小米云盘Go客户端，本软件仅供学习参考使用！",
			fyne.TextAlignCenter, fyne.TextStyle{Bold: true}), nil, nil, logo)))
}
