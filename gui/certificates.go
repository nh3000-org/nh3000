package panes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func certificatesScreen(_ fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")
	if PasswordValid == true {
		calabel := widget.NewLabel(GetLangs("cs-ca"))
		ca := widget.NewMultiLineEntry()
		ca.Resize(fyne.NewSize(320, 240))
		ca.SetText(Caroot)

		cclabel := widget.NewLabel(GetLangs("cs-cc"))
		cc := widget.NewMultiLineEntry()
		cc.SetText(Clientcert)

		cklabel := widget.NewLabel(GetLangs("cs-ck"))
		ck := widget.NewMultiLineEntry()
		ck.SetText(Clientkey)

		ssbutton := widget.NewButton(GetLangs("cs-ss"), func() {
			errors.SetText("...")
			if PasswordValid {
				var iserrors = editEntry("CERTIFICATE", ca.Text)
				if iserrors {
					errors.SetText(GetLangs("cs-err1"))
				}
				iserrors = editEntry("CERTIFICATE", cc.Text)
				if iserrors {
					errors.SetText(GetLangs("cs-err2"))
				}
				iserrors = editEntry("KEY", ck.Text)
				if iserrors {
					errors.SetText(GetLangs("cs-err3"))
				}
				if !iserrors {
					MyJson("SAVE")
				}
			}
		})

		return container.NewCenter(container.NewVBox(
			widget.NewLabelWithStyle(GetLangs("cs-heading"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			calabel,
			ca,
			cclabel,
			cc,
			cklabel,
			cklabel,
			ck,

			ssbutton,
			errors,
			container.NewHBox(
				widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
			),
			widget.NewLabel("_                                                                                             _"),
		))
	}
	errors.SetText("Logon First")
	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle(GetLangs("cs-heading"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		errors,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
		),
		widget.NewLabel(""),
	))

}
