package panes

import (
	//"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/nh3000-org/nh3000/config"
)

//var preferredlanguageShadow string
//var msgmaxageShadow string
//var preferredthemeShadow string

func PhonesScreen(_ fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")
	pad1 := widget.NewButton("1", func() {

	})

	pad2 := widget.NewButton("2", func() {

	})
	pad3 := widget.NewButton("3", func() {

	})
	pad4 := widget.NewButton("4", func() {

	})
	pad5 := widget.NewButton("5", func() {

	})
	pad6 := widget.NewButton("6", func() {

	})

	pad7 := widget.NewButton("7", func() {

	})
	pad8 := widget.NewButton("8", func() {

	})
	pad9 := widget.NewButton("9", func() {

	})
	pad0 := widget.NewButton("0", func() {

	})
	padast := widget.NewButton("*", func() {

	})
	padpound := widget.NewButton("#", func() {

	})
	LogDetails := widget.NewLabel("")
	LogDetails.SetText("Enter Number: ex 200 for extension or *97 for voicemail")
	dial := widget.NewButtonWithIcon("", theme.MailSendIcon(), func() {
		LogDetails.SetText(LogDetails.Text + "\nDialing")
	})

	var LogDetailsBorder = container.NewBorder(LogDetails, nil, nil, nil, nil)

	DetailsVW := container.NewScroll(LogDetailsBorder)
	DetailsVW.SetMinSize(fyne.NewSize(300, 320))

	phlabel := widget.NewLabel("Sip Phone")
	grid10 := container.New(layout.NewGridLayoutWithColumns(3), pad1, pad2, pad3, pad4, pad5, pad6, pad7, pad8, pad9, padast, pad0, padpound)

	topbox := container.NewVBox(
		widget.NewLabelWithStyle(config.GetLangs("ss-heading"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		phlabel,
		grid10,
		dial,
		DetailsVW,
	)
	return container.NewBorder(
		topbox,
		errors,
		nil,
		nil,
		nil,
	)

}
