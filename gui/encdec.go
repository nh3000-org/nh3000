package panes

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func encdecScreen(win fyne.Window) fyne.CanvasObject {
	errors := widget.NewLabel("...")

	password := widget.NewEntry()
	password.SetPlaceHolder(GetLangs("es-pass"))

	myinputtext := widget.NewMultiLineEntry()
	myinputtext.SetPlaceHolder(GetLangs("es-mv"))
	myinputtext.SetMinRowsVisible(6)

	myinputtext.SetText(win.Clipboard().Content())
	myoutputtext := widget.NewMultiLineEntry()
	myoutputtext.SetPlaceHolder(GetLangs("es-mo"))
	myoutputtext.SetMinRowsVisible(6)
	var iserrors = false
	errors.SetText("...")
	encbutton := widget.NewButton(GetLangs("es-em"), func() {
		iserrors = editEntry("STRING", password.Text)
		if iserrors == true {
			errors.SetText(GetLangs("es-err1"))
			iserrors = true
		}
		if len(password.Text) != 24 {
			iserrors = true
			errors.SetText(GetLangs("es-err2-1") + strconv.Itoa(len(password.Text)) + GetLangs("es-err2-2"))
		}
		iserrors = editEntry("STRING", myinputtext.Text)
		if iserrors == true {
			errors.SetText(GetLangs("es-err3"))
			iserrors = true
		}
		if iserrors == false {
			t, err := Encrypt(myinputtext.Text, password.Text)
			if err != nil {
				errors.SetText(GetLangs("es-err4"))
			} else {
				myoutputtext.SetText(string(t))
				win.Clipboard().SetContent(t)
				errors.SetText("...")
			}
		}
	})

	decbutton := widget.NewButton("Decrypt Message", func() {
		iserrors = editEntry("STRING", password.Text)
		if iserrors == true {
			errors.SetText(GetLangs("es-err1"))
			iserrors = true
		}
		if len(password.Text) != 24 {
			iserrors = true
			errors.SetText(GetLangs("es-err2-1") + strconv.Itoa(len(password.Text)) + GetLangs("es-err2-2"))
		}
		iserrors = editEntry("STRING", myinputtext.Text)
		if iserrors == true {
			errors.SetText(GetLangs("es-err3"))
			iserrors = true
		}

		if iserrors == false {
			t, err := Decrypt(myinputtext.Text, password.Text)
			if err != nil {
				errors.SetText(GetLangs("es-err5"))
			} else {
				myoutputtext.SetText(t)
				win.Clipboard().SetContent(t)
				errors.SetText("...")
			}
		}

	})

	if iserrors == true {
		//encbutton.Disable()
		//decbutton.Disable()
	}
	keybox := container.NewBorder(
		password,
		nil,
		nil,
		nil,
		nil,
	)
	inputbox := container.NewBorder(
		widget.NewLabelWithStyle(GetLangs("es-head1"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		myinputtext,
		nil,
		nil,
		nil,
	)
	outputbox := container.NewBorder(
		widget.NewLabelWithStyle(GetLangs("es-head2"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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

	return container.NewBorder(
		c3box,
		errors,
		nil,
		nil,
		nil,
	)

}
