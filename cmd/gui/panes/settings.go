package panes

import (
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/nh3000-org/nh3000/config"
)

var preferredlanguageShadow string
var msgmaxageShadow string
var preferredthemeShadow string
var filterShadow string

func SettingsScreen(_ fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")

	lalabel := widget.NewLabel(config.GetLangs("ss-la"))
	la := widget.NewRadioGroup([]string{"eng", "spa", "hin"}, func(string) {})
	la.Horizontal = true

	//preferredlanguageShadow = config.Decrypt(config.FyneApp.Preferences().StringWithFallback("PreferedLanguage", config.Encrypt("eng", config.MySecret)), config.MySecret)
	//la.SetSelected(preferredlanguageShadow)

	var swf = config.FyneApp.Preferences().StringWithFallback("PreferredLanguage", config.Encrypt("eng", config.MySecret))
	preferredlanguageShadow = config.Decrypt(swf, config.MySecret)
	config.PreferedLanguage = preferredlanguageShadow
	la.SetSelected(preferredlanguageShadow)

	malabel := widget.NewLabel(config.GetLangs("ss-ma"))
	ma := widget.NewRadioGroup([]string{"1h", "12h", "24h", "161h", "8372h"}, func(string) {})
	ma.Horizontal = true
	msgmaxageShadow = config.FyneApp.Preferences().StringWithFallback("NatsMsgMaxAge", config.Encrypt("12h", config.MySecret))
	ma.SetSelected(config.Decrypt(msgmaxageShadow, config.MySecret))

	fllabel := widget.NewLabel(config.GetLangs("ms-filter"))
	filter := widget.NewRadioGroup([]string{"True", "False"}, func(string) {})
	filterShadow = config.FyneApp.Preferences().StringWithFallback("NatsMsgFilter", config.Encrypt("False", config.MySecret))
	filter.SetSelected(config.Decrypt(filterShadow, config.MySecret))
	config.FyneFilter = false
	if strings.Contains(filter.Selected, "True") {
		config.FyneFilter = true
	}
	preferredthemeShadow = config.FyneApp.Preferences().StringWithFallback("FyneTheme", config.Encrypt("0", config.MySecret))
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
			config.FyneApp.Preferences().SetString("FyneTheme", config.Encrypt(strconv.Itoa(config.Selected), config.MySecret))
		}
		if preferredlanguageShadow != la.Selected {
			config.FyneApp.Preferences().SetString("PreferedLanguage", config.Encrypt(la.Selected, config.MySecret))
		}
		if msgmaxageShadow != ma.Selected {
			config.FyneApp.Preferences().SetString("NatsMsgMaxAge", config.Encrypt(ma.Selected, config.MySecret))
		}
		log.Println("settings ", ma.Selected)
		if filterShadow != filter.Selected {
			config.FyneApp.Preferences().SetString("FyneFilter", config.Encrypt(filter.Selected, config.MySecret))
		}

		if config.LoggedOn {
			errors.SetText(config.GetLangs("ss-sserr"))
		}
		if !config.LoggedOn {
			errors.SetText(config.GetLangs("ss-sserr1"))
		}
	})

	topbox := container.NewVBox(
		widget.NewLabelWithStyle(config.GetLangs("ss-heading"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		lalabel,
		la,
		malabel,
		ma,
		fllabel,
		filter,
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
