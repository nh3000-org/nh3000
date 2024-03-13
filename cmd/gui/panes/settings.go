package panes

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/nh3000-org/nh3000/config"
)

var preferredlanguageShadow string
var msgmaxageShadow string

func SettingsScreen(_ fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")

	lalabel := widget.NewLabel(config.GetLangs("ss-la"))
	la := widget.NewRadioGroup([]string{"eng", "spa", "hin"}, func(string) {})
	la.Horizontal = true
	preferredlanguageShadow = config.Decrypt(config.GetApp().Preferences().StringWithFallback("PreferedLanguage", config.Encrypt("eng", config.MySecret)), config.MySecret)
	la.SetSelected(preferredlanguageShadow)
	malabel := widget.NewLabel(config.GetLangs("ss-ma"))
	ma := widget.NewRadioGroup([]string{"1h", "12h", "24h", "161h", "8372h"}, func(string) {})
	ma.Horizontal = true
	msgmaxageShadow = config.GetApp().Preferences().StringWithFallback("MsgMaxAge", config.Encrypt("12h", config.MySecret))
	ma.SetSelected(config.Decrypt(msgmaxageShadow, config.MySecret))

	ssbutton := widget.NewButton(config.GetLangs("ss-ss"), func() {

		if preferredlanguageShadow != la.Selected {
			config.GetApp().Preferences().SetString("PreferedLanguage", config.Encrypt(la.Selected, config.MySecret))
		}
		if msgmaxageShadow != ma.Selected {
			config.GetApp().Preferences().SetString("MsgMaxAge", config.Encrypt(ma.Selected, config.MySecret))
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
