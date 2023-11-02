package nhpanes

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/nh3000-org/nh3000/nhhash"
	"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhpref"
	"github.com/nh3000-org/nh3000/nhutil"
	"golang.org/x/crypto/bcrypt"
)

func PasswordScreen(_ fyne.Window) fyne.CanvasObject {

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder(nhlang.GetLangs("ps-password"))

	passwordc1 := widget.NewPasswordEntry()
	passwordc1.SetPlaceHolder(nhlang.GetLangs("ps-passwordc1"))
	passwordc1.Disable()

	passwordc2 := widget.NewPasswordEntry()
	passwordc2.SetPlaceHolder(nhlang.GetLangs("ps-passwordc2"))
	passwordc2.Disable()
	errors := widget.NewLabel("...")
	// try the password
	tpbutton := widget.NewButton(nhlang.GetLangs("ps-trypassword"), func() {
		var iserrors = false
		nhpref.Password = password.Text
		pwh, err := bcrypt.GenerateFromPassword([]byte(nhpref.Password), bcrypt.DefaultCost)
		nhpref.Passwordhash = string(pwh)
		if err != nil {
			iserrors = true
			errors.SetText(nhlang.GetLangs("ps-err1"))
		}

		myhash, err1 := nhhash.LoadWithDefault("config.hash", "123456")
		nhpref.Passwordhash = myhash
		if err1 {
			errors.SetText(nhlang.GetLangs("ps-err2"))
		}

		nhpref.Password = password.Text

		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(nhpref.Passwordhash), []byte(password.Text)); err != nil {

			iserrors = true
			errors.SetText(nhlang.GetLangs("ps-err4"))
		}
		if !iserrors {
			nhpref.Load()

			//errors.SetText(nhlang.GetLangs("ps-err5"))

			password.Disable()
			passwordc1.Enable()
			passwordc2.Enable()

		}
	})

	cpbutton := widget.NewButton(nhlang.GetLangs("ps-chgpassword"), func() {
		var iserrors = false

		if nhpref.Edit("STRING", passwordc1.Text) == true {
			iserrors = true
			errors.SetText(nhlang.GetLangs("ps-err6"))
		}

		if nhpref.Edit("PASSWORD", passwordc1.Text) {
			iserrors = true
			errors.SetText(nhlang.GetLangs("ps-err7"))
		}
		if passwordc1.Text != passwordc2.Text {
			iserrors = true
			errors.SetText(nhlang.GetLangs("ps-err8"))
		}
		if !iserrors {
			pwh, err := bcrypt.GenerateFromPassword([]byte(passwordc1.Text), bcrypt.DefaultCost)
			nhpref.Passwordhash = string(pwh)

			if err != nil {
				errors.SetText(nhlang.GetLangs("ps-err9") + err.Error())
				log.Fatal(err)
			}
			if !iserrors {
				nhpref.Save()

			}

		}
		nhpref.Password = passwordc1.Text
		_, err := nhhash.Save("config.hash", nhpref.Passwordhash)
		if err {
			errors.SetText(nhlang.GetLangs("ps-err10"))
			iserrors = true
		}

	})

	return container.NewVBox(
		widget.NewLabelWithStyle(nhlang.GetLangs("ps-title1"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("config.json", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle(nhlang.GetLangs("ps-title2"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		password,
		tpbutton,
		widget.NewLabelWithStyle(nhlang.GetLangs("ps-title3"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		passwordc1,
		passwordc2,
		cpbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", nhutil.ParseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", nhutil.ParseURL("https://github.com/nh3000-org/snats")),
		),
		errors,

		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	)

}
