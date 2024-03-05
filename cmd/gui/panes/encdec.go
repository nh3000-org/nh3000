package panes

import (
	"strconv"

	"github.com/nh3000-org/nh3000/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func EncdecScreen(win fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")

	password := widget.NewEntry()
	password.SetText(config.MySecret)

	myinputtext := widget.NewMultiLineEntry()
	myinputtext.SetPlaceHolder(config.GetLangs("es-mv"))
	myinputtext.SetMinRowsVisible(6)

	//myinputtext.SetText(win.Clipboard().Content())
	myoutputtext := widget.NewMultiLineEntry()
	myoutputtext.SetPlaceHolder(config.GetLangs("es-mo"))
	myoutputtext.SetMinRowsVisible(6)

	encbutton := widget.NewButton(config.GetLangs("es-em"), func() {
		var iserrors = false

		iserrors = config.Edit("STRING", password.Text)
		if iserrors {
			errors.SetText(config.GetLangs("es-err1"))
			iserrors = true
		}
		if !iserrors {
			if len(password.Text) != 24 {
				iserrors = true
				errors.SetText(config.GetLangs("es-err2-1") + strconv.Itoa(len(password.Text)) + config.GetLangs("es-err2-2"))
			}
		}
		if !iserrors {
			iserrors = config.Edit("STRING", myinputtext.Text)
			if iserrors {
				errors.SetText(config.GetLangs("es-err3"))
			}
		}
		if !iserrors {
			t := config.Encrypt(myinputtext.Text, password.Text)

			myoutputtext.SetText(string(t))
			//win.Clipboard().SetContent(t)
			errors.SetText("...")

		}
	})
	// copy from clipboard
	cpyFrombutton := widget.NewButton(config.GetLangs("ms-cpyf"), func() {
		myinputtext.SetText(win.Clipboard().Content())
	})

	// copy to clipboard
	cpyTobutton := widget.NewButton(config.GetLangs("ms-cpy"), func() {
		win.Clipboard().SetContent(Details.Text)
	})

	decbutton := widget.NewButton(config.GetLangs("es-dm"), func() {
		var iserrors = false
		iserrors = config.Edit("STRING", password.Text)
		if iserrors {
			errors.SetText(config.GetLangs("es-err1"))
			iserrors = true
		}
		if !iserrors {
			if len(password.Text) != 24 {
				iserrors = true
				errors.SetText(config.GetLangs("es-err2-1") + strconv.Itoa(len(password.Text)) + config.GetLangs("es-err2-2"))
			}
		}
		if !iserrors {
			iserrors = config.Edit("STRING", myinputtext.Text)
			if iserrors {
				errors.SetText(config.GetLangs("es-err3"))
			}
		}
		if !iserrors {
			t := config.Decrypt(myinputtext.Text, password.Text)

			myoutputtext.SetText(t)
			win.Clipboard().SetContent(t)
			errors.SetText("...")

		}

	})

	keybox := container.NewBorder(
		widget.NewLabelWithStyle(config.GetLangs("es-head0"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		password,
		nil,
		nil,
		nil,
	)
	inputbox := container.NewBorder(
		widget.NewLabelWithStyle(config.GetLangs("es-head1"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		cpyFrombutton,
		nil,
		nil,
		myinputtext,
	)
	outputbox := container.NewBorder(
		widget.NewLabelWithStyle(config.GetLangs("es-head2"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		myoutputtext,
		nil,
		nil,
		nil,
	)
	buttonbox := container.NewBorder(
		nil,
		nil,
		nil,
		encbutton,
		decbutton,
	)
	c0box := container.NewBorder(
		keybox,
		nil,
		nil,
		nil,
		nil,
	)
	c1box := container.NewBorder(
		inputbox,
		buttonbox,
		nil,
		nil,
		nil,
	)
	c2box := container.NewBorder(
		c0box,
		c1box,
		nil,
		nil,
		nil,
	)
	c3box := container.NewBorder(
		c2box,
		outputbox,
		nil,
		nil,
		nil,
	)
	c4box := container.NewBorder(
		c3box,
		cpyTobutton,
		nil,
		nil,
		nil,
	)
	return container.NewBorder(
		c4box,
		errors,
		nil,
		nil,
		nil,
	)

}
