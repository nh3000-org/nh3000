package nhpanes

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"

	"github.com/google/uuid"
	"github.com/nh3000-org/nh3000/nhhash"
	"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhnats"
	"github.com/nh3000-org/nh3000/nhpref"
	"github.com/nh3000-org/nh3000/nhutil"
	"golang.org/x/crypto/bcrypt"
)

func LogonScreen(MyWin fyne.Window) fyne.CanvasObject {
	errors := widget.NewLabel("...")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder(nhlang.GetLangs("ls-password"))
	password.SetText(nhpref.Password)

	alias := widget.NewEntry()
	alias.SetPlaceHolder(nhlang.GetLangs("ls-alias"))
	alias.Disable()

	server := widget.NewEntry()
	server.SetPlaceHolder("URL: nats://xxxxxx:4332")
	server.Disable()

	queue := widget.NewEntry()
	queue.SetPlaceHolder(nhlang.GetLangs("ls-queue"))
	queue.Disable()

	queuepassword := widget.NewEntry()
	queuepassword.SetPlaceHolder(nhlang.GetLangs("ls-queuepass"))
	queuepassword.Disable()

	TPbutton := widget.NewButton(nhlang.GetLangs("ls-trypass"), func() {
		errors.SetText("...")
		var iserrors = false
		ph, _ := nhhash.LoadWithDefault("config.hash", "123456")
		nhpref.Password = password.Text
		pwh, err := bcrypt.GenerateFromPassword([]byte(nhpref.Password), bcrypt.DefaultCost)
		nhpref.Passwordhash = string(pwh)
		if err != nil {
			iserrors = true
			errors.SetText(nhlang.GetLangs("ls-err1"))
		}

		// Comparing the password with the hash
		errpw := bcrypt.CompareHashAndPassword([]byte(ph), []byte(nhpref.Password))
		if errpw != nil {
			iserrors = true
			errors.SetText(nhlang.GetLangs("ls-err3"))
		}
		if !iserrors {
			errors.SetText("...")
			nhpref.Load()
			alias.SetText(nhpref.Alias)
			server.SetText(nhpref.Server)
			queue.SetText(nhpref.Queue)
			queuepassword.SetText(nhpref.Queuepassword)
			password.Disable()
			server.Enable()
			queue.Enable()
			alias.Enable()
			queuepassword.Enable()
		}
	})

	SSbutton := widget.NewButton(nhlang.GetLangs("ls-title"), func() {
		var iserrors = nhpref.Edit("URL", server.Text)
		if iserrors == true {
			errors.SetText(nhlang.GetLangs("ls-err4"))
		}
		iserrors = nhpref.Edit("STRING", queuepassword.Text)
		if iserrors == true {
			errors.SetText(nhlang.GetLangs("ls-err5"))
			iserrors = true
		}
		if len(queuepassword.Text) != 24 {
			iserrors = true
			errors.SetText(nhlang.GetLangs("ls-err6-1") + strconv.Itoa(len(queuepassword.Text)) + "ls-err6-1")
		}
		if !iserrors {
			nhpref.LoggedOn = true
			nhpref.NodeUUID = uuid.New().String()
			nhpref.Alias = alias.Text
			nhpref.Server = server.Text
			nhpref.Queue = queue.Text
			nhpref.Queuepassword = queuepassword.Text
			password.Disable()
			server.Disable()
			alias.Disable()
			queue.Disable()
			queuepassword.Disable()
			nhpref.Save()
			nhpref.LoggedOn = true
			errors.SetText("...")
			nhnats.Send(nhlang.GetLangs("ls-con"))
			go nhnats.Receive()
		}
	})

	// security erase
	SEbutton := widget.NewButton(nhlang.GetLangs("ls-erase"), func() {
		if nhpref.LoggedOn {
			nhnats.Erase()
		}
	})
	if nhpref.LoggedOn {
		TPbutton.Disable()
		TPbutton.Refresh()
		SSbutton.Disable()
		SSbutton.Refresh()
		SEbutton.Enable()
		SEbutton.Refresh()
	}
	if !nhpref.LoggedOn {
		password.Enable()
		server.Disable()
		alias.Disable()
		queue.Disable()
		queuepassword.Disable()
	}

	topbox := container.NewVBox(
		widget.NewLabelWithStyle(nhlang.GetLangs("ls-clogon"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		password,
		TPbutton,
		alias,
		server,
		queue,
		queuepassword,
		SSbutton,
		SEbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", nhutil.ParseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", nhutil.ParseURL("https://github.com/nh3000-org/snats")),
		),
		widget.NewLabel(""),
		errors,
	)

	return container.NewBorder(
		topbox,
		errors,
		nil,
		nil,
		nil,
	)
}
