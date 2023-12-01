package nhpanes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhpref"
)

func CertificatesScreen(_ fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")

	calabel := widget.NewLabel(nhlang.GetLangs("cs-ca"))
	ca := widget.NewMultiLineEntry()
	ca.Resize(fyne.NewSize(320, 240))
	ca.SetText(nhpref.Caroot)

	cclabel := widget.NewLabel(nhlang.GetLangs("cs-cc"))
	cc := widget.NewMultiLineEntry()
	cc.SetText(nhpref.Clientcert)

	cklabel := widget.NewLabel(nhlang.GetLangs("cs-ck"))
	ck := widget.NewMultiLineEntry()
	ck.SetText(nhpref.Clientkey)

	ssbutton := widget.NewButton(nhlang.GetLangs("cs-ss"), func() {
		errors.SetText("...")
		if nhpref.LoggedOn == false {
			errors.SetText(nhlang.GetLangs("cs-lf"))
		}
		if nhpref.LoggedOn {
			var iserrors = false
			var err = nhpref.Edit("CERTIFICATE", ca.Text)
			if !err {
				iserrors = true
				errors.SetText(nhlang.GetLangs("cs-err1"))
			}
			err = nhpref.Edit("CERTIFICATE", cc.Text)
			if !err {
				iserrors = true
				errors.SetText(nhlang.GetLangs("cs-err2"))
			}
			err = nhpref.Edit("KEY", ck.Text)
			if !err {
				iserrors = true
				errors.SetText(nhlang.GetLangs("cs-err3"))
			}
			if !iserrors {
				nhpref.Save()
			}
		}
	})

	topbox := container.NewVBox(
		widget.NewLabelWithStyle(nhlang.GetLangs("cs-heading"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		calabel,
		ca,
		cclabel,
		cc,
		cklabel,
		ck,
		ssbutton,
	)
	if !nhpref.LoggedOn {
		topbox = container.NewVBox(
			widget.NewLabelWithStyle(nhlang.GetLangs("cs-heading"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		)
		errors.SetText(nhlang.GetLangs("cs-lf"))
	}
	return container.NewBorder(
		topbox,
		errors,
		nil,
		nil,
		nil,
	)
}
