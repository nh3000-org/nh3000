package panes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func settingsScreen(_ fyne.Window) fyne.CanvasObject {
	MyJson("LOAD")
	lalabel := widget.NewLabel(GetLangs("ss-la"))
	la := widget.NewRadioGroup([]string{"eng", "spa"}, func(string) {})
	la.Horizontal = true
	la.SetSelected(PreferedLanguage)

	pllabel := widget.NewLabel(GetLangs("ss-pl"))
	pl := widget.NewRadioGroup([]string{"6", "8", "12"}, func(string) {})
	pl.Horizontal = true
	pl.SetSelected(PasswordMinimumSize)

	malabel := widget.NewLabel(GetLangs("ss-ma"))
	ma := widget.NewRadioGroup([]string{"12h", "24h", "161h", "8372h"}, func(string) {})
	ma.Horizontal = true
	ma.SetSelected(Msgmaxage)

	mcletterlabel := widget.NewLabel(GetLangs("ss-mcletter"))
	mcletter := widget.NewRadioGroup([]string{"Yes", "No"}, func(string) {})
	mcletter.Horizontal = true
	mcletter.SetSelected(PasswordMustContainLetter)

	mcnumberlabel := widget.NewLabel(GetLangs("ss-mcnumber"))
	mcnumber := widget.NewRadioGroup([]string{"Yes", "No"}, func(string) {})
	mcnumber.Horizontal = true
	mcnumber.SetSelected(PasswordMustContainNumber)

	mcspeciallabel := widget.NewLabel(GetLangs("ss-mcspecial"))
	mcspecial := widget.NewRadioGroup([]string{"Yes", "No"}, func(string) {})
	mcspecial.Horizontal = true
	mcspecial.SetSelected(PasswordMustContainSpecial)

	ssbutton := widget.NewButton(GetLangs("ss-ss"), func() {
		PreferedLanguage = la.Selected
		Msgmaxage = ma.Selected
		PasswordMustContainNumber = mcnumber.Selected
		PasswordMinimumSize = pl.Selected
		PasswordMustContainLetter = mcletter.Selected
		PasswordMustContainSpecial = mcspecial.Selected
		if PasswordValid {
			MyJson("SAVE")
		}
	})

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle(GetLangs("ss-heading"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
		),
		widget.NewLabel(""),
	))

}
