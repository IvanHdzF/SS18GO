package main

// import fyne
import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// New app
	a := app.New()
	// New Window & title
	w := a.NewWindow("MFS")
	//Resize main/parent window
	w.Resize(fyne.NewSize(600, 300))
	// 1st Widget Entry
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter your name")
	entry.Resize(fyne.NewSize(250, 30)) // my widget size
	entry.Move(fyne.NewPos(40, 50))     // position of widget
	entry_email := widget.NewEntry()
	entry_email.SetPlaceHolder("Enter your email")
	entry_email.Resize(fyne.NewSize(250, 30)) // my widget size
	entry_email.Move(fyne.NewPos(40, 100))    // position of widget
	entry_address := widget.NewEntry()
	entry_address.SetPlaceHolder("Enter your name")
	entry_address.Resize(fyne.NewSize(250, 30)) // my widget size
	entry_address.Move(fyne.NewPos(40, 150))    // position of widget
	// button
	btn_submit := widget.NewButton("Submit", func() { OptnRun(false) })
	btn_submit.Resize(fyne.NewSize(150, 30)) // my widget size
	btn_submit.Move(fyne.NewPos(40, 200))    // position of widget
	w.SetContent(
		container.NewWithoutLayout(
			entry,
			entry_email,
			entry_address,
			btn_submit,
		),
	)
	//show and run
	w.ShowAndRun()
}
