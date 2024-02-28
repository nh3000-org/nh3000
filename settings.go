package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/nh3000-org/nh3000/config"
)

func SettingsScreen(_ fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")

	lalabel := widget.NewLabel(config.GetLangs("ss-la"))
	la := widget.NewRadioGroup([]string{"eng", "spa", "hin"}, func(string) {})
	la.Horizontal = true
	la.SetSelected(config.GetPreferedLanguage())

	malabel := widget.NewLabel(config.GetLangs("ss-ma"))
	ma := widget.NewRadioGroup([]string{"12h", "24h", "161h", "8372h"}, func(string) {})
	ma.Horizontal = true
	ma.SetSelected(config.GetMsgMaxAge())

	ssbutton := widget.NewButton(config.GetLangs("ss-ss"), func() {

		if config.GetPreferedLanguage() != la.Selected {
			config.SetPreferedLanguage(la.Selected)
		}
		if config.GetMsgMaxAge() != ma.Selected {
			config.SetMsgMaxAge(ma.Selected)
		}
		if config.GetLoggedOn() {
			errors.SetText(config.GetLangs("ss-sserr"))

		}
		if !config.GetLoggedOn() {
			errors.SetText(config.GetLangs("ss-sserr1"))

		}
	})

	topbox := container.NewVBox(
		widget.NewLabelWithStyle(config.GetLangs("ss-heading"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		lalabel,
		la,
		malabel,
		ma,
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
