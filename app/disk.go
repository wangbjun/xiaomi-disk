package app

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"xiaomi_cloud/api"
)

func disk(app *App) fyne.CanvasObject {
	var diskNavs = []string{"全部文件"}
	folders, err := app.Api.GetFoldersById("0")
	if err != nil {
		if err == api.UnauthorizedError { // 手机安全验证码
			err = app.User.GetPhoneCode()
			if err != nil {
				if err == api.NoPhoneCodeError {
					folders, err = app.Api.GetFoldersById("0")
					if err != nil {
						app.Alert(err.Error())
					}
					goto ShowList
				}
			} else {
				code := widget.NewEntry()
				code.Validator = validation.NewRegexp(`^[0-9]+$`, "Only Numbers")
				codeItem := widget.NewFormItem("验证码", code)
				codeItem.HintText = "请输入短信验证码"
				items := []*widget.FormItem{
					codeItem,
				}
				dialog.ShowForm("安全验证", "确定", "取消", items, func(b bool) {
					fmt.Printf("%v\n", code.Text)
				}, app.Window)
			}
		} else {
			app.Alert("获取目录列表失败")
		}
	}

ShowList:
	navs := widget.NewLabel("全部文件")
	list := widget.NewList(
		func() int {
			return len(folders)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			data := folders[i]
			o.(*widget.Label).SetText(data.Name)
		})

	list.OnSelected = func(i widget.ListItemID) {
		selected := folders[i]
		if selected.Type == api.TypeFolder {
			diskNavs = append(diskNavs, selected.Name)
			folders, err = app.Api.GetFoldersById(selected.Id)
			if err != nil {
				app.Alert(err.Error())
			}
			navs.SetText(getDiskNavs(diskNavs))
			list.Refresh()
		}
	}
	return container.NewBorder(container.NewVBox(navs), nil, nil, nil, list)
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
