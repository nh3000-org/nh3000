package panes

import (
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/crypto/bcrypt"
)

func passwordScreen(_ fyne.Window) fyne.CanvasObject {

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder(GetLangs("ps-password"))

	passwordc1 := widget.NewPasswordEntry()
	passwordc1.SetPlaceHolder(GetLangs("ps-passwordc1"))
	passwordc1.Disable()

	passwordc2 := widget.NewPasswordEntry()
	passwordc2.SetPlaceHolder(GetLangs("ps-passwordc2"))
	passwordc2.Disable()
	errors := widget.NewLabel("...")
	// try the password
	tpbutton := widget.NewButton(GetLangs("ps-trypassword"), func() {
		var iserrors = false
		Password = password.Text
		pwh, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		Passwordhash = string(pwh)
		if err != nil {
			iserrors = true
			errors.SetText(GetLangs("ps-err1"))
		}
		_, confighasherr := os.Stat(DataStore("config.hash").Path())
		if confighasherr != nil {

			if MyHash("CREATE") {
				errors.SetText(GetLangs("ps-err2"))
			}
		}

		Password = password.Text
		if MyHash("LOAD") {
			errors.SetText(GetLangs("ps-err3"))
		}
		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(Password)); err != nil {

			iserrors = true
			errors.SetText(GetLangs("ps-err4"))
		}
		if !iserrors {
			MyJson("LOAD")
			errors.SetText(GetLangs("ps-err5"))

			password.Disable()
			passwordc1.Enable()
			passwordc2.Enable()

		}
	})

	cpbutton := widget.NewButton(GetLangs("ps-chgpassword"), func() {
		var iserrors = false

		if editEntry("STRING", passwordc1.Text) == true {
			iserrors = true
			errors.SetText(GetLangs("ps-err6"))
		}

		if editEntry("PASSWORD", passwordc1.Text) {
			iserrors = true
			errors.SetText(GetLangs("ps-err7"))
		}
		if passwordc1.Text != passwordc2.Text {
			iserrors = true
			errors.SetText(GetLangs("ps-err8"))
		}
		if !iserrors {
			pwh, err := bcrypt.GenerateFromPassword([]byte(passwordc1.Text), bcrypt.DefaultCost)
			Passwordhash = string(pwh)

			if err != nil {
				log.Fatal(err)
			}

			if MyHash("SAVE") {
				errors.SetText(GetLangs("ps-err9"))

			}

		}
		Password = passwordc1.Text
		if MyHash("LOAD") {
			errors.SetText(GetLangs("ps-err10"))
			iserrors = true
		}
		// Comparing the password with the hash
		err := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(passwordc1.Text))
		if err != nil {
			errors.SetText(GetLangs("ps-err11"))
			iserrors = true
		}
		if !iserrors {
			MyJson("SAVE")

		}

	})

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle(GetLangs("ps-title1"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("config.json", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(GetLangs("ps-title2"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		password,
		tpbutton,
		widget.NewLabelWithStyle(GetLangs("ps-title3"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		passwordc1,
		passwordc2,
		cpbutton,
		errors,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", parseURL("https://github.com/nh3000-org/snats")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))

}
