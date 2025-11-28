package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nabievarthur/rtf-converter/internal/service"
)

func BuildUI(win fyne.Window) fyne.CanvasObject {
	notifier := NewNotifier()
	rtfConverter := service.NewRtfConverter()

	label := widget.NewLabel("Выберите .rtf файл для конвертации")
	label.Alignment = fyne.TextAlignCenter

	convertBtn := widget.NewButtonWithIcon("Конвертировать", theme.ConfirmIcon(), func() {
		res := rtfConverter.ConvertRtfToExcel(label.Text, "output.xlsx")
		if !res {
			notifier.Show("Ошибка")
			return
		}
		notifier.Show("Файл успешно сохранен")
	})
	convertBtn.Importance = widget.HighImportance
	convertBtn.Hide()

	openBtn := widget.NewButtonWithIcon("Открыть файл", theme.SearchIcon(), func() {
		fileName := OpenFileDialog()
		if fileName == "" {
			notifier.Show("Файл не выбран")
		}
		if fileName != "" {
			convertBtn.Show()
			label.SetText(fileName)
		}
	})

	content := container.NewBorder(nil, notifier.Widget(), nil, nil,
		container.NewVBox(
			label,
			openBtn,
			convertBtn,
		),
	)
	return content
}
