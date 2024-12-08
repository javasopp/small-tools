package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
)

var vmBatFileName = "startVm.bat"

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("只能在windows上使用，GG.")
		return
	}

	a := app.New()
	w := a.NewWindow("sopp-vmware启动bat一键创建")
	vmInstallDirectory := widget.NewEntry()
	vmInstallDirectory.SetPlaceHolder("点击选择vmware安装目录")

	vmSelectButton := widget.NewButton("选择文件", func() {
		path, err := zenity.SelectFile(
			zenity.Title("选择程序路径"),
			zenity.FileFilter{Patterns: []string{"vmrun.exe"}})
		if err == nil {
			if !strings.Contains(path, "vmrun.exe") {
				dialog.ShowInformation("错误", "请选择vmware安装程序目录！", w)
				return
			}
			vmInstallDirectory.SetText(path)
		} else {
			fmt.Println("文件选择被取消或出错：", err)
		}
	})

	vmFileEntries := container.NewVBox()

	addVmFileEntry := func() {
		vmFilePath := widget.NewEntry()
		vmFilePath.SetPlaceHolder("点击选择虚拟机路径")

		vmFileSelectButton := widget.NewButton("选择文件", func() {
			path, err := zenity.SelectFile(
				zenity.Title("选择虚拟机文件路径"),
				zenity.FileFilter{Patterns: []string{"*.vmx"}})
			if err == nil {
				if !strings.Contains(path, ".vmx") {
					dialog.ShowInformation("错误", "请选择虚拟机文件路径！", w)
					return
				}
				vmFilePath.SetText(path)
			} else {
				fmt.Println("文件选择被取消或出错：", err)
			}
		})

		vmFileEntries.Add(container.NewBorder(nil, nil, nil, vmFileSelectButton, vmFilePath))
	}

	addVmFileEntry() // 初始添加一组

	createButton := widget.NewButton("创建脚本", func() {
		vmInDir := vmInstallDirectory.Text
		if vmInDir == "" {
			dialog.ShowInformation("错误", "请选择vmware安装目录。", w)
			return
		}

		var entriesInfo []string
		for _, obj := range vmFileEntries.Objects {
			vmFileEntry := obj.(*fyne.Container).Objects[0].(*widget.Entry)
			vmFile := vmFileEntry.Text

			if vmFile == "" {
				dialog.ShowInformation("错误", "请确保所有虚拟机路径已填写。", w)
				return
			}
			entriesInfo = append(entriesInfo, fmt.Sprintf("\"%s\" start \"%s\" nogui ", vmInDir, vmFile))
		}

		createVmwareBat(entriesInfo, w)
	})

	removeVmFileEntry := func() {
		if len(vmFileEntries.Objects) > 1 {
			vmFileEntries.Remove(vmFileEntries.Objects[len(vmFileEntries.Objects)-1])
		} else {
			dialog.ShowInformation("提示", "至少需要一个虚拟机路径。", w)
		}
	}

	addButton := widget.NewButton("+", func() {
		addVmFileEntry()
	})

	removeButton := widget.NewButton("-", func() {
		removeVmFileEntry()
	})

	buttonContainer := container.NewGridWithColumns(2, addButton, removeButton)

	content := container.NewVBox(
		widget.NewLabel("vmware安装路径:"),
		container.NewBorder(nil, nil, nil, vmSelectButton, vmInstallDirectory),
		widget.NewLabel("虚拟机安装目录:"),
		vmFileEntries,
		//container.NewHBox(addButton, createButton),
		buttonContainer,
		//container.NewHBox(addButton, removeButton),
		createButton,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(700, 300))
	w.CenterOnScreen()
	w.ShowAndRun()
}

func createVmwareBat(entriesInfo []string, w fyne.Window) {
	files, _ := os.Create(vmBatFileName)
	defer func(files *os.File) {
		_ = files.Close()
	}(files)

	for _, entry := range entriesInfo {
		_, writeErr := files.WriteString(entry + "\n")
		if writeErr != nil {
			return
		}
	}
	fmt.Println("写入成功")
	dialog.ShowInformation("成功", "生成成功，请配合计划任务程序使用", w)
}
