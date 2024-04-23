package panes

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/nh3000-org/nh3000/config"
)

var preferredlanguageShadow string
var msgmaxageShadow string
var preferredthemeShadow string

func SettingsScreen(_ fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")

	lalabel := widget.NewLabel(config.GetLangs("ss-la"))
	la := widget.NewRadioGroup([]string{"eng", "spa", "hin"}, func(string) {})
	la.Horizontal = true
	preferredlanguageShadow = config.Decrypt(config.FyneApp.Preferences().StringWithFallback("PreferedLanguage", config.Encrypt("eng", config.MySecret)), config.MySecret)
	la.SetSelected(preferredlanguageShadow)
	malabel := widget.NewLabel(config.GetLangs("ss-ma"))
	ma := widget.NewRadioGroup([]string{"1h", "12h", "24h", "161h", "8372h"}, func(string) {})
	ma.Horizontal = true
	msgmaxageShadow = config.FyneApp.Preferences().StringWithFallback("MsgMaxAge", config.Encrypt("12h", config.MySecret))
	ma.SetSelected(config.Decrypt(msgmaxageShadow, config.MySecret))

	preferredthemeShadow = config.FyneApp.Preferences().StringWithFallback("Theme", config.Encrypt("0", config.MySecret))
	config.Selected, _ = strconv.Atoi(config.Decrypt(preferredthemeShadow, config.MySecret))
	themes := container.NewGridWithColumns(3,
		widget.NewButton(config.GetLangs("mn-dark"), func() {
			config.Selected = config.Dark
			config.FyneApp.Settings().SetTheme(config.MyTheme{})

		}),
		widget.NewButton(config.GetLangs("mn-light"), func() {
			config.Selected = config.Light
			config.FyneApp.Settings().SetTheme(config.MyTheme{})
		}),
		widget.NewButton(config.GetLangs("mn-retro"), func() {
			config.Selected = config.Retro
			config.FyneApp.Settings().SetTheme(config.MyTheme{})
		}),
	)
	ssbutton := widget.NewButton(config.GetLangs("ss-ss"), func() {
		x, _ := strconv.Atoi(config.Decrypt(preferredthemeShadow, config.MySecret))
		if x != config.Selected {
			config.FyneApp.Preferences().SetString("Theme", config.Encrypt(strconv.Itoa(config.Selected), config.MySecret))
		}
		if preferredlanguageShadow != la.Selected {
			config.FyneApp.Preferences().SetString("PreferedLanguage", config.Encrypt(la.Selected, config.MySecret))
		}
		if msgmaxageShadow != ma.Selected {
			config.FyneApp.Preferences().SetString("MsgMaxAge", config.Encrypt(ma.Selected, config.MySecret))
		}

		if preferredlanguageShadow != config.PreferedLanguage {
			config.FyneApp.Preferences().SetString("PreferedLanguage", config.Encrypt(la.Selected, config.MySecret))
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
		themes,
	)
	return container.NewBorder(
		topbox,
		errors,
		nil,
		nil,
		nil,
	)
}