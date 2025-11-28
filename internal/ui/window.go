package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func Run() {
	a := app.New()
	w := a.NewWindow("RTF-converter")
	w.Resize(fyne.NewSize(400, 200))
	w.CenterOnScreen()
	w.SetContent(BuildUI(w))
	w.ShowAndRun()
}
