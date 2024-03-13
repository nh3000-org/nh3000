package panes

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"

	"fyne.io/fyne/v2/widget"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"

	"github.com/nh3000-org/nh3000/config"
)

var nodeuuid = ""

func LogonScreen(MyWin fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder(config.GetLangs("ls-password"))

	alias := widget.NewEntry()
	alias.SetPlaceHolder(config.GetLangs("ls-alias"))
	alias.Disable()
	var aliasShadow = ""

	server := widget.NewEntry()
	server.SetPlaceHolder("URL: nats://xxxxxx:4332")
	server.Disable()
	var serverShadow = ""

	queue := widget.NewEntry()
	queue.SetPlaceHolder(config.GetLangs("ls-queue"))
	queue.Disable()
	var queueShadow = ""

	queuepassword := widget.NewEntry()
	queuepassword.SetPlaceHolder(config.GetLangs("ls-queuepass"))
	queuepassword.Disable()
	var queuepasswordShadow = ""

	calabel := widget.NewLabel(config.GetLangs("cs-ca"))
	ca := widget.NewMultiLineEntry()
	ca.Resize(fyne.NewSize(320, 240))
	ca.Disable()
	var caShadow = ""

	cclabel := widget.NewLabel(config.GetLangs("cs-cc"))
	cc := widget.NewMultiLineEntry()
	cc.Disable()
	var ccShadow = ""

	cklabel := widget.NewLabel(config.GetLangs("cs-ck"))
	ck := widget.NewMultiLineEntry()
	ck.Disable()
	var ckShadow = ""

	TPbutton := widget.NewButton(config.GetLangs("ls-trypass"), func() {
		errors.SetText("...")
		config.SetReceivingMessages(false)
		var iserrors = false
		ph, _ := config.LoadHashWithDefault("config.hash", "123456")

		//log.Println("pw ", MyPrefs.Password)
		pwh, err := bcrypt.GenerateFromPassword([]byte(password.Text), bcrypt.DefaultCost)
		config.PasswordHash = string(pwh)
		if err != nil {
			iserrors = true
			errors.SetText(config.GetLangs("ls-err1"))
		}

		// Comparing the password with the hash
		errpw := bcrypt.CompareHashAndPassword([]byte(ph), []byte(password.Text))
		if errpw != nil {
			iserrors = true
			errors.SetText(config.GetLangs("ls-err3"))
		}
		if !iserrors {
			errors.SetText("...")

			var preferedlanguageShadow = config.Decrypt(config.GetApp().Preferences().StringWithFallback("eng", config.Encrypt(config.GetPreferedLanguage(), config.MySecret)), config.MySecret)
			config.SetPreferedLanguage(config.Decrypt(preferedlanguageShadow, config.MySecret))

			var nodeuuidShadow = config.GetApp().Preferences().StringWithFallback("NodeUUID", config.Encrypt(uuid.New().String(), config.MySecret))
			nodeuuid = config.Decrypt(nodeuuidShadow, config.MySecret)
			config.SetNodeUUID(nodeuuid)
			if nodeuuidShadow != nodeuuid {
				config.GetApp().Preferences().SetString("NodeUUID", config.Encrypt(nodeuuid, config.MySecret))
			}
			aliasShadow = config.GetApp().Preferences().StringWithFallback("Alias", config.Encrypt("MyAlias", config.MySecret))
			alias.SetText(config.Decrypt(aliasShadow, config.MySecret))

			serverShadow = config.GetApp().Preferences().StringWithFallback("Server", config.Encrypt("nats://nats.newhorizons3000.org:4222", config.MySecret))
			server.SetText(config.Decrypt(serverShadow, config.MySecret))

			queueShadow = config.GetApp().Preferences().StringWithFallback("Queue", config.Encrypt("MESSAGES", config.MySecret))
			queue.SetText(config.Decrypt(queueShadow, config.MySecret))

			queuepasswordShadow = config.GetApp().Preferences().StringWithFallback("QueuePassword", config.Encrypt("987654321098765432109876", config.MySecret))
			queuepassword.SetText(config.Decrypt(queuepasswordShadow, config.MySecret))

			config.SetMsgMaxAge(config.Decrypt(config.GetApp().Preferences().StringWithFallback("MsgMaxAge", config.Encrypt("12h", config.MySecret)), config.MySecret))

			config.SetPreferedLanguage(config.Decrypt(config.GetApp().Preferences().StringWithFallback("PreferedLanguage", config.Encrypt("eng", config.MySecret)), config.MySecret))

			caShadow = config.GetApp().Preferences().StringWithFallback("CAROOT", config.Encrypt("-----BEGIN CERTIFICATE-----\nMIID7zCCAtegAwIBAgIUaXAPxJvZRRdTq5RWlwxs1XYo+5kwDQYJKoZIhvcNAQEL\nBQAwgYAxCzAJBgNVBAYTAlVTMRAwDgYDVQQIEwdGbG9yaWRhMRIwEAYDVQQHEwlD\ncmVzdHZpZXcxGjAYBgNVBAoTEU5ldyBIb3Jpem9ucyAzMDAwMQwwCgYDVQQLEwNX\nV1cxITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZzAeFw0yMzEyMTkw\nMzA4MDBaFw0yODEyMTcwMzA4MDBaMIGAMQswCQYDVQQGEwJVUzEQMA4GA1UECBMH\nRmxvcmlkYTESMBAGA1UEBxMJQ3Jlc3R2aWV3MRowGAYDVQQKExFOZXcgSG9yaXpv\nbnMgMzAwMDEMMAoGA1UECxMDV1dXMSEwHwYDVQQDExhuYXRzLm5ld2hvcml6b25z\nMzAwMC5vcmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCrVIXA/SxU\n7GeW92UNyiPnQEZgbJIHHQ31AQE2C/vFdpEtv32uoX1SsDl5drWvBrMnd5zrw1tL\nOEPA26tk/ACfQYL0n0HfeutLLu8H9jUWNp8ziX6Qbgd01M+/BixobHQjyDMxulo4\nJU2VK6QBLs9VI6TIihEU2BZhc/XCD9QbWcikAif1JySpz93MjFv3pcQU8ci4vQ0T\nImaGnHesr1qDbX1NuFVuBOPavZ64sQ1RsZtH5CdD+RU772wQWUgkPkwyUn8QBwTS\ne9XV5DNQD5nGEXjKTgjrd9KRf9pmRDnf6gBLi2r6C/l6q2w3ItOOHARdK0mc9CYh\ngY1Nzl59vrWdAgMBAAGjXzBdMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTAD\nAQH/MB0GA1UdDgQWBBR0qq9ueABC5RDsg/02FZFpBOR1hDAbBgNVHREEFDAShwTA\nqAAFhwTAqFjohwR/AAABMA0GCSqGSIb3DQEBCwUAA4IBAQBfdX0IMya9Dh9dHLJj\nnJZyb96htMWD5nuQQVBAu3ay+8O2GWj5mlsLJXAP2y7p/+3gyvHKTRDdJLux7N79\nHn6AYjmp3PCyZzuL1M/kHhSQxhxqJHGwjGXILt5pLovVkvkl4iukdxWJ5HAPsUGY\nO3QSDDFdoLflsG5VcrtdODm8uyxAjhMPAR2PXKfX8ABI79N7VKcbb98338fifrN8\n9H1r3BXcdsyhpH0gB0ZKJFSpMGWXlfudFEe9mXI9898xbEI2znqlYGhboVsuv5LM\nRESH2zXrkhmZyHqw0RtDROzyZOy5g1LcxbtVMn4w1LI4h3MDuE9B+Vud77A48qtA\ny+5x\n-----END CERTIFICATE-----\n", config.MySecret))
			ca.SetText(config.Decrypt(caShadow, config.MySecret))

			ccShadow = config.GetApp().Preferences().StringWithFallback("CACLIENT", config.Encrypt("-----BEGIN CERTIFICATE-----\nMIIEMTCCAxmgAwIBAgIUB7+OFX1LQrWtYMl5XIOXsOaLac0wDQYJKoZIhvcNAQEL\nBQAwgYAxCzAJBgNVBAYTAlVTMRAwDgYDVQQIEwdGbG9yaWRhMRIwEAYDVQQHEwlD\ncmVzdHZpZXcxGjAYBgNVBAoTEU5ldyBIb3Jpem9ucyAzMDAwMQwwCgYDVQQLEwNX\nV1cxITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZzAgFw0yMzEyMTkw\nMzA4MDBaGA8yMDUzMTIxMTAzMDgwMFowcjELMAkGA1UEBhMCVVMxEDAOBgNVBAgT\nB0Zsb3JpZGExEjAQBgNVBAcTCUNyZXN0dmlldzEaMBgGA1UEChMRTmV3IEhvcml6\nb25zIDMwMDAxITAfBgNVBAsTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZzCCASIw\nDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMWARyniHy8r342e3aKSsLDPwVMC\n2mRwuILP2JkXp5FllaFKnu/Z+0mF+iQlSchcC6DOcMQk00Cp/I8cCP865zyxPhqN\n2F2/qVItCU4+PTwe6ZnrfpJgXWwyk1hjS3vVNTT+idI5+pJgFH9YL0lbJ7q1UyPB\n+KP0x/c5T3K2Ec6U4uXhbVt/ePxFmsl1sHw6FE//XrA4EzbqCMEPCTcOfInvFrCJ\ny4/pAqjCxegT/1YDMNEdzmG8vg2tc3jPV+3GIAV3YL5nDE5mprHPEEDJtNQi+E4o\nXXXMobNhrJh9KJ59VbxTF8m5yM3b8fvof97OYhK0KYggplnTH+bhnYU9V5ECAwEA\nAaOBrTCBqjAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYD\nVR0TAQH/BAIwADAdBgNVHQ4EFgQUpffi3LSreerO756B/VnZkyyEVBIwHwYDVR0j\nBBgwFoAUdKqvbngAQuUQ7IP9NhWRaQTkdYQwNQYDVR0RBC4wLIIYbmF0cy5uZXdo\nb3Jpem9uczMwMDAub3JnhwR/AAABhwTAqAAFhwTAqFjoMA0GCSqGSIb3DQEBCwUA\nA4IBAQALlRqqW2HH4flFIgR/nh51gc/Hxv5xivhkzWUHHXRdltECSXknI4yBPchQ\n6Zsy0HZ7XQRlhQSIYd4Bp6eyHbny5t3JA978dHzpGJFCUVQDMY4yHLaCQgFJ+ESn\nwyyDWTRGA3cpEikL0B0ekDfqjWUEMTzmT/gnoSl0vM69nZDLZm1xMx1+EH+bpfFB\nRaVM6gKSAuFJmNYEL2e7JSags+3IHyVHkdo8GDlY//71Z4lxsFxFCF6xF9GDdAr2\niCA4OfydjiBSOz0eLJVgqkk1KGXtMqZXAojX62NrIWnFTW1Vzd46ekOHhq93B3tA\nkjWmHY/KdCZUjQSWss+YXgG4mI8c\n-----END CERTIFICATE-----\n", config.MySecret))
			cc.SetText(config.Decrypt(ccShadow, config.MySecret))
			ckShadow = config.GetApp().Preferences().StringWithFallback("CAKEY", config.Encrypt("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAxYBHKeIfLyvfjZ7dopKwsM/BUwLaZHC4gs/YmRenkWWVoUqe\n79n7SYX6JCVJyFwLoM5wxCTTQKn8jxwI/zrnPLE+Go3YXb+pUi0JTj49PB7pmet+\nkmBdbDKTWGNLe9U1NP6J0jn6kmAUf1gvSVsnurVTI8H4o/TH9zlPcrYRzpTi5eFt\nW394/EWayXWwfDoUT/9esDgTNuoIwQ8JNw58ie8WsInLj+kCqMLF6BP/VgMw0R3O\nYby+Da1zeM9X7cYgBXdgvmcMTmamsc8QQMm01CL4Tihddcyhs2GsmH0onn1VvFMX\nybnIzdvx++h/3s5iErQpiCCmWdMf5uGdhT1XkQIDAQABAoIBAB+Iu9QUJqaBetBB\n7WFnyo5wnY2DhxtCZDN+vDa1cCvm7F00bOwfAeBbY/UhfwZeq/yg+aBXwOMyQQEY\nmNcnsIQgSKo0u7c8Quy8BCBaD6zpwqKw1yTH/iKocJ5MPGEpSbWMbrUCTN/SN3Od\nwO8VfuJw0TWEYw7KpqLyo5zNNUqmczEO438CPGotbkFfzUqkumeUOsGWJFongyZY\na9EwpcTH2TkxuXum9SQVyLy+hSG/AEBp0cQPaRcoNh8sWYk43y5HrkIAqFo7dkMa\n9usAVMz9JCqIH2UNV04cDASFaiDMpYoD2hV2YHlL7/CQ7v5nb6OHT2A9aoSBOAfm\ns+dBzYECgYEA1l8+T9Xux73TCbFO2p7F094xSx4hhBZhaYpvzZoNN7iQdbdUVt2l\n1yHSoRgJUJMZlnKpMoNMLCxo34Lr3ww/TkIE/rrg10pqbqvojIDLCbi103EEB2v9\nWix8MSeOgFCa72T4lg9fDm5T493n4C5dade3LzZczUBF6dgmth3D+nMCgYEA69pa\nlob9n7eNXqDPk9kZUJV1jfLATC8eN4jupEiKfjnxEz9mUewvL/RF8kFhiS1ISC50\nKgM0v+isYBwwX00c7P02L6xCoGT35qOeoutEWVy/tYIHIHsD0jUBBsdnpQVNf58l\n9DDy2hZrpUwrsVHylVHpufBgKOfxgP2Jr3qD0OsCgYEAn4vzTGfkdzSIRMZ58awJ\ngE32Ufny5+PgTDSEUXk+LSJoIbR4SM5eB2dc5BiHljhk6twboUSnBJlo1DEUa8Up\nuIzaOtvLS3BPFl9LjIaulmWqrduHLB7rSJmjNNJD9KwJI/L6MHTwQkVKmmUllmvr\nikLKS5EiMICNiCUfaptsqJECgYEApYaSqzBEUdK1oeMErAPis16hqSTkdtNexqUQ\nrzXGFP6/Rb3qJra3C1XJvVLLjEW+hAIuPsoPPFyklbNS85+gHGc9n0mrXPxfy3ur\nuzWYu4rPdSizrcUIEoBmnwZVpEhLcrUUIwQzfIHdvJ3v0DvuH4PkoD2mjy7xnJDU\nD9bRKk8CgYAqK1lY5waFR0u3eFIPnrV4ATHXYuxcup2DCF+KJ6qwc4nNI6OB/ovU\nttiVZGr1rca42+XdWUQL5ufPFuKymeLbsuVzabbGKi+4RMvL+TIuorYtJRUPF+C7\nA9jlMeckpTZvl0yn5s3lC817N27B+U0M/jGow8sO0NtjBiImuTC5dg==\n-----END RSA PRIVATE KEY-----\n", config.MySecret))

			ck.SetText(config.Decrypt(ckShadow, config.MySecret))
			password.Disable()
			alias.Enable()
			server.Enable()
			queue.Enable()
			queuepassword.Enable()
			ca.Enable()
			cc.Enable()
			ck.Enable()
		}
	})

	SSbutton := widget.NewButton(config.GetLangs("ls-title"), func() {
		var haserrors = false
		if aliasShadow != alias.Text {
			haserrors = config.Edit("STRING", alias.Text)
			if !haserrors {
				config.Encrypt(alias.Text, config.MySecret)
				config.GetApp().Preferences().SetString("Alias", config.Encrypt(alias.Text, config.MySecret))
			} else {
				errors.SetText(config.GetLangs("ls-err5"))
			}
		}

		if serverShadow != server.Text {
			haserrors = config.Edit("URL", server.Text)
			if !haserrors {
				config.GetApp().Preferences().SetString("Server", config.Encrypt(server.Text, config.MySecret))
			} else {
				errors.SetText(config.GetLangs("ls-err4"))
			}

		}
		if queueShadow != queue.Text {
			haserrors = config.Edit("STRING", queue.Text)
			if !haserrors {
				config.GetApp().Preferences().SetString("Queue", config.Encrypt(queue.Text, config.MySecret))
			} else {
				errors.SetText(config.GetLangs("ls-err0"))
			}

		}

		if queuepasswordShadow != queuepassword.Text {
			haserrors = config.Edit("STRING", queuepassword.Text)
			if !haserrors {
				if len(queuepassword.Text) < 6 {
					haserrors = true
					errors.SetText(config.GetLangs("ls-err6-1") + strconv.Itoa(len(queuepassword.Text)) + "ls-err6-1")
				}
				if !haserrors {
					config.GetApp().Preferences().SetString("QueuePassword", config.Encrypt(queuepassword.Text, config.MySecret))
				}
			} else {
				errors.SetText(config.GetLangs("ls-err5"))
			}

		}

		if !haserrors {
			config.SetLoggedOn()

			config.SetAlias(alias.Text)

			config.SetServer(server.Text)
			config.SetQueue(queue.Text)
			config.SetQueuePassword(queuepassword.Text)
			config.SetCaroot(ca.Text)
			config.SetClientCert(cc.Text)
			config.SetClientKey(ck.Text)
			config.SetLoggedOn()
			password.Disable()
			server.Disable()
			alias.Disable()
			queue.Disable()
			queuepassword.Disable()
			ca.Disable()
			ck.Disable()
			cc.Disable()

			errors.SetText("...")
			config.Send(config.GetLangs("ls-con"), config.GetAlias())
			go config.Receive()
		}
	})

	// security erase
	SEbutton := widget.NewButton(config.GetLangs("ls-erase"), func() {
		if config.GetLoggedOn() {
			config.Erase()
		}
	})
	if config.GetLoggedOn() {
		TPbutton.Disable()
		TPbutton.Refresh()
		SSbutton.Disable()
		SSbutton.Refresh()
		SEbutton.Enable()
		SEbutton.Refresh()
	}
	if !config.GetLoggedOn() {
		password.Enable()
		server.Disable()
		alias.Disable()
		queue.Disable()
		queuepassword.Disable()
		ca.Disable()
		cc.Disable()
		ck.Disable()

	}
	themes := container.NewGridWithColumns(3,
		widget.NewButton(config.GetLangs("mn-dark"), func() {
			config.Selected = config.Dark
			config.GetApp().Settings().SetTheme(config.MyTheme{})

		}),
		widget.NewButton(config.GetLangs("mn-light"), func() {
			config.Selected = config.Light
			config.GetApp().Settings().SetTheme(config.MyTheme{})
		}),
		widget.NewButton(config.GetLangs("mn-retro"), func() {
			config.Selected = config.Retro
			config.GetApp().Settings().SetTheme(config.MyTheme{})
		}),
	)

	vertbox := container.NewVBox(

		widget.NewLabelWithStyle(config.GetLangs("ls-clogon"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		password,
		TPbutton,
		alias,
		server,
		queue,
		queuepassword,
		calabel,
		ca,
		cclabel,
		cc,
		cklabel,
		ck,
		SSbutton,
		SEbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", config.ParseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", config.ParseURL("https://github.com/nh3000-org/snats")),
		),
		widget.NewLabel(""),
		themes,
		errors,
	)
	return container.NewScroll(
		vertbox,
	)

}
