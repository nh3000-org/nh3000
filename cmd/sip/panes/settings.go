package panes

import (
	"runtime"
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

	lalabel := widget.NewLabel("Language")
	la := widget.NewRadioGroup([]string{"eng", "spa", "hin"}, func(string) {})
	la.Horizontal = true
	preferredlanguageShadow = config.Decrypt(config.FyneApp.Preferences().StringWithFallback("PreferedLanguage", config.Encrypt("eng", config.MySecret)), config.MySecret)
	la.SetSelected(preferredlanguageShadow)

	preferredthemeShadow = config.FyneApp.Preferences().StringWithFallback("SIPTheme", config.Encrypt("0", config.MySecret))
	config.Selected, _ = strconv.Atoi(config.Decrypt(preferredthemeShadow, config.MySecret))
	themes := container.NewGridWithColumns(3,
		widget.NewButton("Dark", func() {
			config.Selected = config.Dark
			config.FyneApp.Settings().SetTheme(config.MyTheme{})

		}),
		widget.NewButton("Light", func() {
			config.Selected = config.Light
			config.FyneApp.Settings().SetTheme(config.MyTheme{})
		}),
		widget.NewButton("Retro", func() {
			config.Selected = config.Retro
			config.FyneApp.Settings().SetTheme(config.MyTheme{})
		}),
	)
	ssbutton := widget.NewButton("Save Settings", func() {
		x, _ := strconv.Atoi(config.Decrypt(preferredthemeShadow, config.MySecret))
		if x != config.Selected {
			config.FyneApp.Preferences().SetString("SIPTheme", config.Encrypt(strconv.Itoa(config.Selected), config.MySecret))
		}
		if preferredlanguageShadow != la.Selected {
			config.FyneApp.Preferences().SetString("SIPPreferedLanguage", config.Encrypt(la.Selected, config.MySecret))
		}

		if preferredlanguageShadow != config.PreferedLanguage {
			config.FyneApp.Preferences().SetString("PreferedLanguage", config.Encrypt(la.Selected, config.MySecret))
		}
		if config.LoggedOn {
			errors.SetText("Logged In")
		}
		if !config.LoggedOn {
			errors.SetText("Not Logged On")
		}
		runtime.GC()
		runtime.ReadMemStats(&memoryStats)
		config.SIPFyneMainWin.SetTitle("Settings Saved " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")

	})

	topbox := container.NewVBox(
		widget.NewLabelWithStyle("SIP Settings", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		lalabel,
		la,
		ssbutton,
		themes,
	)
	runtime.GC()
	runtime.ReadMemStats(&memoryStats)
	config.SIPFyneMainWin.SetTitle("Saved " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")

	return container.NewBorder(
		topbox,
		errors,
		nil,
		nil,
		nil,
	)
}
