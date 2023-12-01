package nhpanes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhpref"
)

func SettingsScreen(_ fyne.Window) fyne.CanvasObject {
	errors := widget.NewLabel("...")

	lalabel := widget.NewLabel(nhlang.GetLangs("ss-la"))
	la := widget.NewRadioGroup([]string{"eng", "spa", "hin"}, func(string) {})
	la.Horizontal = true
	la.SetSelected(nhpref.PreferedLanguage)

	pllabel := widget.NewLabel(nhlang.GetLangs("ss-pl"))
	pl := widget.NewRadioGroup([]string{"6", "8", "12"}, func(string) {})
	pl.Horizontal = true
	pl.SetSelected(nhpref.PasswordMinimumSize)

	malabel := widget.NewLabel(nhlang.GetLangs("ss-ma"))
	ma := widget.NewRadioGroup([]string{"12h", "24h", "161h", "8372h"}, func(string) {})
	ma.Horizontal = true
	ma.SetSelected(nhpref.Msgmaxage)

	mcletterlabel := widget.NewLabel(nhlang.GetLangs("ss-mcletter"))
	mcletter := widget.NewRadioGroup([]string{"Yes", "No"}, func(string) {})
	mcletter.Horizontal = true
	mcletter.SetSelected(nhpref.PasswordMustContainLetter)

	mcnumberlabel := widget.NewLabel(nhlang.GetLangs("ss-mcnumber"))
	mcnumber := widget.NewRadioGroup([]string{"Yes", "No"}, func(string) {})
	mcnumber.Horizontal = true
	mcnumber.SetSelected(nhpref.PasswordMustContainNumber)

	mcspeciallabel := widget.NewLabel(nhlang.GetLangs("ss-mcspecial"))
	mcspecial := widget.NewRadioGroup([]string{"Yes", "No"}, func(string) {})
	mcspecial.Horizontal = true
	mcspecial.SetSelected(nhpref.PasswordMustContainSpecial)

	ssbutton := widget.NewButton(nhlang.GetLangs("ss-ss"), func() {
		nhpref.PreferedLanguage = la.Selected
		nhpref.Msgmaxage = ma.Selected
		nhpref.PasswordMustContainNumber = mcnumber.Selected
		nhpref.PasswordMinimumSize = pl.Selected
		nhpref.PasswordMustContainLetter = mcletter.Selected
		nhpref.PasswordMustContainSpecial = mcspecial.Selected
		if nhpref.LoggedOn {
			errors.SetText(nhlang.GetLangs("ss-sserr"))
			nhpref.Save()
		}
		if !nhpref.LoggedOn {
			errors.SetText(nhlang.GetLangs("ss-sserr1"))

		}
	})

	topbox := container.NewVBox(
		widget.NewLabelWithStyle(nhlang.GetLangs("ss-heading"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		lalabel,
		la,
		malabel,
		ma,
		pllabel,
		pl,
		mcletterlabel,
		mcletter,
		mcnumberlabel,
		mcnumber,
		mcspeciallabel,
		mcspecial,
		ssbutton,
	)
	return container.NewBorder(
		topbox,
		errors,
		nil,
		nil,
		nil,
	)
}
