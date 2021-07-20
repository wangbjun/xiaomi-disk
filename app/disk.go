package app

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/dustin/go-humanize"
	"time"
	"unicode/utf8"
	"xiaomi_disk/api"
	"xiaomi_disk/utils"
)

func disk(app *App) fyne.CanvasObject {
	var diskNavs = []string{"全部文件"}
	folders, err := app.Api.GetFoldersById("0")
	if err != nil {
		if err != api.UnauthorizedError {
			app.Alert("获取目录列表失败")
		} else {
			err = app.User.SecondLoginCheck()
			if err != nil {
				app.Alert("二次安全验证失败")
			}
			folders, err = app.Api.GetFoldersById("0")
		}
	}
	navs := widget.NewLabelWithStyle("全部文件", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	list := widget.NewTable(
		func() (int, int) { return len(folders), 4 },
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			data := folders[id.Row]
			modifyTime := time.Unix(data.ModifyTime/1000, data.ModifyTime%1000)
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%d", id.Row+1))
			case 1:
				if utf8.RuneCountInString(data.Name) >= 25 {
					runes := []rune(data.Name)
					label.SetText(string(runes[:25]) + "...")
				} else {
					label.SetText(data.Name)
				}
			case 2:
				label.SetText(humanize.Bytes(uint64(data.Size)))
			case 3:
				label.SetText(modifyTime.Format(utils.YmdHis))
			}
		})
	list.SetColumnWidth(0, 30)
	list.SetColumnWidth(1, 350)
	list.SetColumnWidth(2, 80)
	list.OnSelected = func(id widget.TableCellID) {
		selected := folders[id.Row]
		if selected.Type == api.TypeFolder {
			diskNavs = append(diskNavs, selected.Name)
			folders, err = app.Api.GetFoldersById(selected.Id)
			if err != nil {
				app.Alert(err.Error())
			}
			navs.SetText(getDiskNavs(diskNavs))
			list.Refresh()
			list.Unselect(widget.TableCellID{})
		}
	}
	return container.NewBorder(container.NewVBox(navs, widget.NewSeparator()), nil, nil, nil, list)
}

func getDiskNavs(navs []string) string {
	var labels = ""
	labels += navs[0]
	for _, v := range navs[1:] {
		labels += " > "
		labels += v
	}
	return labels
}
