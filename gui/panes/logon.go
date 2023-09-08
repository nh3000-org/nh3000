package panes

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func logonScreen(MyWin fyne.Window) fyne.CanvasObject {
	errors := widget.NewLabel("...")

	MyJson("LOAD")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder(GetLangs("ls-password"))
	password.SetText(Password)

	alias := widget.NewEntry()
	alias.SetPlaceHolder(GetLangs("ls-alias"))
	alias.SetText(Alias)
	alias.Disable()

	server := widget.NewEntry()
	server.SetPlaceHolder("URL: nats://xxxxxx:4332")
	server.Disable()

	queue := widget.NewEntry()
	queue.SetPlaceHolder(GetLangs("ls-queue"))
	queue.Disable()

	queuepassword := widget.NewEntry()
	queuepassword.SetPlaceHolder(GetLangs("ls-queuepass"))
	queuepassword.Disable()

	tpbutton := widget.NewButton(GetLangs("ls-trypass"), func() {
		errors.SetText("...")
		var iserrors = false
		Password = password.Text
		pwh, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
		Passwordhash = string(pwh)
		if err != nil {
			iserrors = true
			log.Println(GetLangs("ls-err1"))
			errors.SetText(GetLangs("ls-err1"))
		}
		confighasherr, _ := storage.Exists(DataStore("config.hash"))
		//log.Println("hash logon ", confighasherr, " at ", DataStore("config.hash"))
		if confighasherr == false {
			if MyHash("CREATE") {
				log.Println(GetLangs("ls-err1"))
				errors.SetText(GetLangs("ls-err1"))
			}
		}

		if MyHash("LOAD") {
			log.Println(GetLangs("ls-err2"))
			errors.SetText(GetLangs("ls-err2"))
		}
		// Comparing the password with the hash
		errpw := bcrypt.CompareHashAndPassword([]byte(Passwordhash), []byte(Password))

		if errpw != nil {
			iserrors = true
			log.Println(GetLangs("ls-err3"))
			errors.SetText(GetLangs("ls-err3"))
		}

		if !iserrors {
			errors.SetText("...")
			PasswordValid = true
			MyJson("LOAD")
			alias.SetText(Alias)
			server.SetText(Server)
			queue.SetText(Queue)
			queuepassword.SetText(Queuepassword)
			password.Disable()
			server.Enable()
			queue.Enable()
			alias.Enable()
			queuepassword.Enable()
		}
	})

	SSbutton := widget.NewButton(GetLangs("ls-title"), func() {

		var iserrors = editEntry("URL", server.Text)
		if iserrors == true {
			log.Println(GetLangs("ls-err4"))
			errors.SetText(GetLangs("ls-err4"))
		}
		iserrors = editEntry("STRING", queuepassword.Text)
		if iserrors == true {
			errors.SetText(GetLangs("ls-err5"))
			iserrors = true
		}
		if len(queuepassword.Text) != 24 {
			iserrors = true
			errors.SetText(GetLangs("ls-err6-1") + strconv.Itoa(len(queuepassword.Text)) + "ls-err6-1")
		}

		if !iserrors && PasswordValid {
			NodeUUID = uuid.New().String()
			Alias = alias.Text
			Server = server.Text
			Queue = queue.Text
			Queuepassword = queuepassword.Text
			password.Disable()
			server.Disable()

			alias.Disable()
			queue.Disable()
			queuepassword.Disable()

			MyJson("SAVE")

			LoggedOn = true

			errors.SetText("...")

			FormatMessage("Connected")

		}

	})

	// security erase
	SEbutton := widget.NewButton(GetLangs("ls-erase"), func() {
		if PasswordValid {
			NATSErase()
		}

	})

	if !PasswordValid {
		password.Enable()
		server.Disable()
		alias.Disable()
		queue.Disable()
		queuepassword.Disable()
	}

	return container.NewCenter(container.NewVBox(
		widget.NewLabelWithStyle(GetLangs("ls-clogon"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		password,
		tpbutton,
		alias,
		server,
		queue,
		queuepassword,
		SSbutton,
		SEbutton,
		errors,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", parseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", parseURL("https://github.com/nh3000-org/snats")),
		),
		widget.NewLabel(""), // balance the header on the tutorial screen we leave blank on this content
	))
}
