package main

import (
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"golang.design/x/clipboard"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("时间转换工具")
	myWindow.CenterOnScreen()
	myWindow.Resize(fyne.NewSize(256, 500))
	// myWindow.SetFixedSize(true)

	// 输入时间戳框
	timestampEntry := widget.NewEntry()
	timestampEntry.SetPlaceHolder("输入时间戳...")

	// 输入时间格式框
	timeFormatEntry := widget.NewEntry()
	timeFormatEntry.SetText("2006-01-02 15:04:05")
	timeFormatEntry.SetPlaceHolder("输入时间格式...")

	// 首次点击修改时间格式, 弹出提示框
	var firstChange bool = true
	timeFormatEntry.OnCursorChanged = func() {
		if firstChange {
			dialog.ShowInformation("", "请不要随意修改时间格式", myWindow)
			firstChange = false
		}
	}

	// 输入时间字符串框
	datetimeEntry := widget.NewEntry()
	datetimeEntry.SetPlaceHolder("输入日期时间...")

	// 转换结果部分
	resultLabel := widget.NewLabel("转换结果将显示在这里")

	// 当前时间戳显示框
	nowTimeStamp := widget.NewEntry()
	nowTimeStamp.SetText(strconv.Itoa(int(time.Now().Unix())))

	// 单位选择
	unitSelect := widget.NewRadioGroup([]string{"秒", "毫秒"}, nil)
	unitSelect.Horizontal = true
	unitSelect.SetSelected("秒")

	// 操作按钮
	copyCurrentTimeBtn := widget.NewButton("复制当前时间戳", func() {})
	convertToTimeBtn := widget.NewButton("时间戳转时间", func() {})
	convertToTimestampBtn := widget.NewButton("时间转时间戳", func() {})

	// 复制当前时间戳功能
	copyCurrentTimeBtn.OnTapped = func() {
		// 写入文本到剪贴板
		clipboard.Write(clipboard.FmtText, []byte(getTimeStampFromTime(unitSelect.Selected, time.Now())))
	}

	// 定时刷新当前时间戳
	go func() {
		for range time.Tick(time.Second) {
			fyne.Do(
				func() {
					nowTimeStamp.SetText(getTimeStampFromTime(unitSelect.Selected, time.Now()))
				})
		}
	}()

	// 时间戳转时间
	convertToTimeBtn.OnTapped = func() {
		tsStr := strings.TrimSpace(timestampEntry.Text)
		if tsStr == "" {
			resultLabel.SetText("请输入时间戳")
			return
		}

		format := timeFormatEntry.Text
		if format == "" {
			format = "2006-01-02 15:04:05"
		}

		ts, err := strconv.ParseInt(tsStr, 10, 64)
		if err != nil {
			resultLabel.SetText("时间戳格式错误: " + err.Error())
			return
		}

		// 如果是毫秒单位，需要调整
		if unitSelect.Selected == "毫秒" && ts > 10000000000 {
			ts = ts / 1000
		}

		t := time.Unix(ts, 0)
		resultLabel.SetText("转换成功")
		datetimeEntry.SetText(t.Format(format))
	}

	// 时间转时间戳
	convertToTimestampBtn.OnTapped = func() {
		dtStr := strings.TrimSpace(datetimeEntry.Text)
		if dtStr == "" {
			resultLabel.SetText("请输入日期时间")
			return
		}

		format := timeFormatEntry.Text
		if format == "" {
			format = "2006-01-02 15:04:05"
		}

		t, err := time.ParseInLocation(format, dtStr, time.Local)
		if err != nil {
			resultLabel.SetText("日期时间格式错误: " + err.Error())
			return
		}

		resultLabel.SetText("转换成功")
		timestampEntry.SetText(getTimeStampFromTime(unitSelect.Selected, t))
	}

	// 布局
	content := container.NewVBox(
		widget.NewLabel("时间戳:"),
		timestampEntry,
		widget.NewLabel("时间格式:"),
		timeFormatEntry,
		widget.NewLabel("日期时间:"),
		datetimeEntry,
		widget.NewLabel("单位:"),
		unitSelect,
		resultLabel,
		container.NewGridWithColumns(2, convertToTimeBtn, convertToTimestampBtn),
		container.NewGridWithColumns(2, widget.NewLabel("当前时间:"), nowTimeStamp),
		copyCurrentTimeBtn,
	)

	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func getTimeStampFromTime(secondFormat string, t time.Time) string {
	var timestamp string
	if secondFormat == "秒" {
		timestamp = strconv.FormatInt(t.Unix(), 10)
	} else {
		timestamp = strconv.FormatInt(t.UnixNano()/1000000, 10)
	}
	return timestamp
}
