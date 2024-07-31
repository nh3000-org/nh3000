package panes

import (
	"log"
	"runtime"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"golang.org/x/crypto/bcrypt"

	"github.com/nh3000-org/nh3000/config"
)

func PasswordScreen(_ fyne.Window) fyne.CanvasObject {

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Original Password")

	passwordc1 := widget.NewPasswordEntry()
	passwordc1.SetPlaceHolder("New Password")
	passwordc1.Disable()

	passwordc2 := widget.NewPasswordEntry()
	passwordc2.SetPlaceHolder("New Password Again")
	passwordc2.Disable()
	errors := widget.NewLabel("...")
	runtime.GC()
	runtime.ReadMemStats(&memoryStats)
	config.SIPFyneMainWin.SetTitle("Password " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")

	// try the password
	tpbutton := widget.NewButton("Try Password", func() {
		var iserrors = false

		pwh, err := bcrypt.GenerateFromPassword([]byte(password.Text), bcrypt.DefaultCost)
		config.PasswordHash = string(pwh)
		if err != nil {
			iserrors = true
			errors.SetText("Invalid Password")
		}

		myhash, err1 := config.LoadHashWithDefault("config.hash", "123456")
		config.PasswordHash = myhash
		if err1 {
			errors.SetText("Password Invalid")
		}

		// Comparing the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(config.PasswordHash), []byte(password.Text)); err != nil {

			iserrors = true
			errors.SetText("Passwords Mismatch")
		}
		if !iserrors {

			//errors.SetText(nhlang.GetLangs("ps-err5"))

			password.Disable()
			passwordc1.Enable()
			passwordc2.Enable()

		}
	})

	cpbutton := widget.NewButton("Change Password", func() {
		var iserrors = false

		if config.Edit("STRING", passwordc1.Text) {
			iserrors = true
			errors.SetText("Password Blank Error")
		}

		if config.Edit("PASSWORD", passwordc1.Text) {
			iserrors = true
			errors.SetText("Password Error")
		}
		if passwordc1.Text != passwordc2.Text {
			iserrors = true
			errors.SetText("Password Mismatch")
		}
		if !iserrors {
			pwh, err := bcrypt.GenerateFromPassword([]byte(passwordc1.Text), bcrypt.DefaultCost)
			config.PasswordHash = string(pwh)

			if err != nil {
				errors.SetText("Password Hash " + err.Error())
				log.Fatal(err)
			}

		}

		_, err := config.SaveHash("config.hash", config.PasswordHash)
		if err {
			errors.SetText("Could Not Save to File System")
			iserrors = true
		}

	})
	if !config.SIPLoggedOn {
		password.Disable()

		passwordc1.Disable()
		passwordc2.Disable()
		cpbutton.Disable()
	}
	runtime.GC()
	runtime.ReadMemStats(&memoryStats)
	config.SIPFyneMainWin.SetTitle("Saved " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")

	return container.NewVBox(
		widget.NewLabelWithStyle("Password Reset", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("config.json", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewLabelWithStyle("Change Local Encryption Password", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),

		password,
		tpbutton,
		widget.NewLabelWithStyle("Enter Passwords To Change", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
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
