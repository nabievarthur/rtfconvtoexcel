package ui

import "github.com/sqweek/dialog"

func OpenFileDialog() string {
	fileName, err := dialog.File().
		Filter("rtf файлы", "rtf").
		Title("Выберите файл .rtf").
		Load()
	if err != nil {
		return ""
	}
	return fileName
}
