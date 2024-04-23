package panes

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"golang.org/x/crypto/bcrypt"

	"github.com/nh3000-org/nh3000/config"
)

func PasswordScreen(_ fyne.Window) fyne.CanvasObject {

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder(config.GetLangs("ps-password"))

	passwordc1 := widget.NewPasswordEntry()
	passwordc1.SetPlaceHolder(config.GetLangs("ps-passwordc1"))
	passwordc1.Disable()

	passwordc2 := widget.NewPasswordEntry()
	passwordc2.SetPlaceHolder(config.GetLangs("ps-passwordc2"))
	passwordc2.Disable()
	errors := widget.NewLabel("...")
	// try the password
	tpbutton := widget.NewButton(config.GetLangs("ps-trypassword"), func() {
		var iserrors = false

		pwh, err := bcrypt.GenerateFromPassword([]byte(password.Text), bcrypt.DefaultCost)
		config.PasswordHash = string(pwh)
		if err != nil {
			iserrors = true
			errors.SetText(config.GetLangs("ps-err1"))
		}

		myhash, err1 := config.LoadHashWithDefault("config.hash", "123456")
		config.PasswordHash = myhash
		if err1 {
			errors.SetText(config.GetLangs("ps-err2"))
		}

		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(config.PasswordHash), []byte(password.Text)); err != nil {

			iserrors = true
			errors.SetText(config.GetLangs("ps-err4"))
		}
		if !iserrors {

			//errors.SetText(nhlang.GetLangs("ps-err5"))

			password.Disable()
			passwordc1.Enable()
			passwordc2.Enable()

		}
	})

	cpbutton := widget.NewButton(config.GetLangs("ps-chgpassword"), func() {
		var iserrors = false

		if config.Edit("STRING", passwordc1.Text) {
			iserrors = true
			errors.SetText(config.GetLangs("ps-err6"))
		}

		if config.Edit("PASSWORD", passwordc1.Text) {
			iserrors = true
			errors.SetText(config.GetLangs("ps-err7"))
		}
		if passwordc1.Text != passwordc2.Text {
			iserrors = true
			errors.SetText(config.GetLangs("ps-err8"))
		}
		if !iserrors {
			pwh, err := bcrypt.GenerateFromPassword([]byte(passwordc1.Text), bcrypt.DefaultCost)
			config.PasswordHash = string(pwh)

			if err != nil {
				errors.SetText(config.GetLangs("ps-err9") + err.Error())
				log.Fatal(err)
			}

		}

		_, err := config.SaveHash("config.hash", config.PasswordHash)
		if err {
			errors.SetText(config.GetLangs("ps-err10"))
			iserrors = true
		}

	})
	if !config.LoggedOn {
		password.Disable()

		passwordc1.Disable()
		passwordc2.Disable()
		cpbutton.Disable()
	}
	return container.NewVBox(
		widget.NewLabelWithStyle(config.GetLangs("ps-title1"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("config.json", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(config.GetLangs("ps-title2"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		password,
		tpbutton,
		widget.NewLabelWithStyle(config.GetLangs("ps-title3"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		passwordc1,
		passwordc2,
		cpbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", config.ParseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", config.ParseURL("https://github.com/nh3000-org/snats")),
		),
		widget.NewLabel(""),
		errors,
	)

}
