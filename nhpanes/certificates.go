package panes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhpref"
	"github.com/nh3000-org/nh3000/nhutil"
)

func certificatesScreen(_ fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")
	if nhpref.PasswordValid == false {
		errors.SetText(nhlang.GetLangs("cs-lf"))
	}
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
		if nhpref.PasswordValid {
			var iserrors = nhpref.Edit("CERTIFICATE", ca.Text)
			if iserrors {
				errors.SetText(nhlang.GetLangs("cs-err1"))
			}
			iserrors = nhpref.Edit("CERTIFICATE", cc.Text)
			if iserrors {
				errors.SetText(nhlang.GetLangs("cs-err2"))
			}
			iserrors = nhpref.Edit("KEY", ck.Text)
			if iserrors {
				errors.SetText(nhlang.GetLangs("cs-err3"))
			}
			if !iserrors {
				nhpref.Save()
			}
		}
	})

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle(nhlang.GetLangs("cs-heading"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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
			widget.NewHyperlink("newhorizons3000.org", nhutil.ParseURL("https://newhorizons3000.org/")),
			widget.NewLabel("_                                                                                             _"),
		)))

}
