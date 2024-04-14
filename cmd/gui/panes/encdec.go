package panes

import (
	"strconv"
	"strings"

	"github.com/nh3000-org/nh3000/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var fileOpened = ""

func EncdecScreen(win fyne.Window) fyne.CanvasObject {
	var Details = widget.NewLabel("")
	errors := widget.NewLabel("...")

	password := widget.NewEntry()
	password.SetText(config.MySecret)

	myinputtext := widget.NewMultiLineEntry()
	myinputtext.SetPlaceHolder(config.GetLangs("es-mv"))
	myinputtext.SetMinRowsVisible(3)

	myoutputtext := widget.NewMultiLineEntry()
	myoutputtext.SetPlaceHolder(config.GetLangs("es-mo"))
	myoutputtext.SetMinRowsVisible(3)

	encbutton := widget.NewButtonWithIcon(config.GetLangs("es-em"), theme.MediaFastForwardIcon(), func() {
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
	encbuttonfile := widget.NewButtonWithIcon(config.GetLangs("ms-encf"), theme.MediaFastForwardIcon(), func() {
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
			iserrors = config.Edit("FILEEXISTS", fileOpened)
		}
		if !iserrors {
			config.EncryptFile(fileOpened, fileOpened+".nh3000")
		}
	})
	// copy from clipboard
	cpyFrombutton := widget.NewButtonWithIcon(config.GetLangs("ms-cpyf"), theme.ContentCopyIcon(), func() {
		myinputtext.SetText(win.Clipboard().Content())
	})

	// copy to clipboard
	cpyTobutton := widget.NewButtonWithIcon(config.GetLangs("ms-cpy"), theme.ContentPasteIcon(), func() {
		win.Clipboard().SetContent(Details.Text)
	})

	decbutton := widget.NewButtonWithIcon(config.GetLangs("es-dm"), theme.MediaFastRewindIcon(), func() {
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
	decbuttonfile := widget.NewButtonWithIcon(config.GetLangs("ms-decf"), theme.MediaFastRewindIcon(), func() {
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
			iserrors = config.Edit("FILEEXISTS", fileOpened)
		}
		if !iserrors {
			var d = strings.Replace(fileOpened, ".nh3000", "", 1)
			config.DecryptFile(fileOpened, d)
			errors.SetText(d)
		}

	})
	fileOpened = ""
	openFile := widget.NewButton(config.GetLangs("ms-self"), func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				return
			}

			fileOpened = reader.URI().Path()
			errors.SetText(fileOpened)
		}, win)
		fd.Show()
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
		myinputtext,
		nil,
		nil,
		nil,
	)
	outputbox := container.NewBorder(
		widget.NewLabelWithStyle(config.GetLangs("es-head2"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		myoutputtext,
		nil,
		nil,
		nil,
	)
	buttonbox := container.NewGridWithColumns(4,
		cpyFrombutton,
		encbutton,
		decbutton,
		cpyTobutton,
	)
	buttonboxfile := container.NewGridWithColumns(2,

		encbuttonfile,
		decbuttonfile,
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
		outputbox,
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
		buttonbox,
		nil,
		nil,
		nil,
	)
	c4box := container.NewBorder(
		c3box,
		widget.NewLabelWithStyle(config.GetLangs("ms-file"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
	)
	c5box := container.NewBorder(
		c4box,
		openFile,
		nil,
		nil,
		nil,
	)

	c6box := container.NewBorder(
		c5box,
		buttonboxfile,
		nil,
		nil,
		nil,
	)
	return container.NewBorder(
		c6box,
		errors,
		nil,
		nil,
		nil,
	)

}
