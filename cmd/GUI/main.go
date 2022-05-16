package main

// import fyne
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// New app
	a := app.New()
	// New Window & title
	w := a.NewWindow("Multi-File Simulator and Optimizer")
	//Resize main/parent window
	w.Resize(fyne.NewSize(600, 300))

	//check
	replace := false
	check_replace := widget.NewCheck("Replace configs", func(b bool) {
		if b {
			replace = true
			// refresh

		} else {
			replace = false

		}
	})
	check_replace.SetChecked(true)
	// button
	btn_run := widget.NewButton("Just run", func() { OptnRunFunc(false, replace) })

	// button
	btn_optrun := widget.NewButton("Opt n' run", func() { OptnRunFunc(true, replace) })

	row1 := container.New(layout.NewVBoxLayout(), btn_run, btn_optrun, check_replace)

	content := container.New(layout.NewHBoxLayout(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), layout.NewSpacer(), row1)

	w.SetContent(content)
	//show and run
	w.ShowAndRun()
}
