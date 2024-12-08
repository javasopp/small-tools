package main

import (
	"bytes"
	"fmt"
	"fyne.io/fyne/v2"
	"os/exec"
	"runtime"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
)

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("只能在windows上使用，GG.")
		return
	}

	a := app.New()
	w := a.NewWindow("sopp-计划任务创建器")

	taskNameEntry := widget.NewEntry()
	taskNameEntry.SetPlaceHolder("输入计划任务名称")

	programPath := widget.NewEntry()
	programPath.SetPlaceHolder("点击选择程序路径")

	selectButton := widget.NewButton("选择文件", func() {
		path, err := zenity.SelectFile(zenity.Title("选择程序路径"))
		if err == nil {
			programPath.SetText(path)
		} else {
			fmt.Println("文件选择被取消或出错：", err)
		}
	})

	createTaskFunc := func() {
		taskName := taskNameEntry.Text
		path := programPath.Text

		if taskName == "" {
			dialog.ShowInformation("错误", "请输入计划任务名称。", w)
			return
		}

		if path == "" {
			dialog.ShowInformation("错误", "请选择程序路径。", w)
			return
		}

		if taskExists(taskName) {
			dialog.ShowInformation("信息", "计划任务已存在。", w)
			return
		}

		createTask(taskName, path, w)
	}

	taskNameEntry.OnSubmitted = func(text string) {
		if taskExists(taskNameEntry.Text) {
			dialog.ShowInformation("信息", "计划任务已存在。", w)
		} else {
			dialog.ShowInformation("信息", "计划任务不存在，可以创建新任务。", w)
		}
	}

	programPath.OnSubmitted = func(text string) {
		createTaskFunc()
	}

	createButton := widget.NewButton("创建任务", func() {
		createTaskFunc()
	})

	content := container.NewVBox(
		widget.NewLabel("计划任务名称:"),
		taskNameEntry,
		widget.NewLabel("程序路径:"),
		container.NewBorder(nil, nil, nil, selectButton, programPath),
		createButton,
	)

	w.SetContent(content)
	w.Resize(fyne.NewSize(700, 300)) // 调整窗口大小
	w.CenterOnScreen()               // 窗口居中显示
	w.ShowAndRun()
}

func taskExists(taskName string) bool {
	cmd := exec.Command("schtasks", "/query", "/tn", taskName)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	// 如果命令执行成功，则任务存在；如果失败，任务不存在
	return err == nil
}

func createTask(taskName, path string, w fyne.Window) {
	args := []string{
		"/create",
		"/tn", taskName,
		"/tr", path,
		"/sc", "onstart",
		"/ru", "SYSTEM",
		"/f", // 强制覆盖已有任务
	}

	cmd := exec.Command("schtasks", args...)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()

	// 将输出转换为 UTF-8，假定系统返回为 GBK
	output, _ := ioutil.ReadAll(transform.NewReader(&out, simplifiedchinese.GBK.NewDecoder()))

	if err != nil {
		dialog.ShowError(fmt.Errorf("创建任务时出错: %s\n%s", err, output), w)
	} else {
		dialog.ShowInformation("成功", fmt.Sprintf("任务创建成功:\n%s", output), w)
	}
}
