package nhpanes

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nh3000-org/nh3000/nhauth"
	"github.com/nh3000-org/nh3000/nhcrypt"
	"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhpref"
)

func EncdecScreen(win fyne.Window) fyne.CanvasObject {
	errors := widget.NewLabel("...")

	password := widget.NewEntry()
	password.SetText(nhauth.MySecret)

	myinputtext := widget.NewMultiLineEntry()
	myinputtext.SetPlaceHolder(nhlang.GetLangs("es-mv"))
	myinputtext.SetMinRowsVisible(6)

	myinputtext.SetText(win.Clipboard().Content())
	myoutputtext := widget.NewMultiLineEntry()
	myoutputtext.SetPlaceHolder(nhlang.GetLangs("es-mo"))
	myoutputtext.SetMinRowsVisible(6)

	encbutton := widget.NewButton(nhlang.GetLangs("es-em"), func() {
		var iserrors = false
		iserrors = nhpref.Edit("STRING", password.Text)
		if iserrors {
			errors.SetText(nhlang.GetLangs("es-err1"))
			iserrors = true
		}
		if !iserrors {
			if len(password.Text) != 24 {
				iserrors = true
				errors.SetText(nhlang.GetLangs("es-err2-1") + strconv.Itoa(len(password.Text)) + nhlang.GetLangs("es-err2-2"))
			}
		}
		if !iserrors {
			iserrors = nhpref.Edit("STRING", myinputtext.Text)
			if iserrors {
				errors.SetText(nhlang.GetLangs("es-err3"))
			}
		}
		if !iserrors {
			t, err := nhcrypt.Encrypt(myinputtext.Text, password.Text)
			if err != nil {
				errors.SetText(nhlang.GetLangs("es-err4"))
			} else {
				myoutputtext.SetText(string(t))
				//win.Clipboard().SetContent(t)
				errors.SetText("...")
			}
		}
	})
	// copy to clipboard messages
	cpyFrombutton := widget.NewButton(nhlang.GetLangs("ms-cpyf"), func() {
		myinputtext.SetText(win.Clipboard().Content())
	})

	// copy to clipboard messages
	cpyTobutton := widget.NewButton(nhlang.GetLangs("ms-cpy"), func() {
		win.Clipboard().SetContent(Details.Text)
	})

	decbutton := widget.NewButton(nhlang.GetLangs("es-dm"), func() {
		var iserrors = false
		iserrors = nhpref.Edit("STRING", password.Text)
		if iserrors {
			errors.SetText(nhlang.GetLangs("es-err1"))
			iserrors = true
		}
		if !iserrors {
			if len(password.Text) != 24 {
				iserrors = true
				errors.SetText(nhlang.GetLangs("es-err2-1") + strconv.Itoa(len(password.Text)) + nhlang.GetLangs("es-err2-2"))
			}
		}
		if !iserrors {
			iserrors = nhpref.Edit("STRING", myinputtext.Text)
			if iserrors == true {
				errors.SetText(nhlang.GetLangs("es-err3"))
			}
		}
		if !iserrors {
			t, err := nhcrypt.Decrypt(myinputtext.Text, password.Text)
			if err != nil {
				errors.SetText(nhlang.GetLangs("es-err5"))
			} else {
				myoutputtext.SetText(t)
				win.Clipboard().SetContent(t)
				errors.SetText("...")
			}
		}

	})
	if !nhpref.LoggedOn {
		encbutton.Disable()
		decbutton.Disable()
	}

	keybox := container.NewBorder(
		widget.NewLabelWithStyle(nhlang.GetLangs("es-head0"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		password,
		nil,
		nil,
		nil,
	)
	inputbox := container.NewBorder(
		widget.NewLabelWithStyle(nhlang.GetLangs("es-head1"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		cpyFrombutton,
		nil,
		nil,
		myinputtext,
	)
	outputbox := container.NewBorder(
		widget.NewLabelWithStyle(nhlang.GetLangs("es-head2"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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
