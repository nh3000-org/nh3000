package panes

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/widget"

	"golang.org/x/crypto/bcrypt"

	"github.com/nh3000-org/nh3000/config"

)

func LogonScreen(MyWin fyne.Window) fyne.CanvasObject {

	errors := widget.NewLabel("...")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("Enter Password")

	alias := widget.NewEntry()
	alias.SetPlaceHolder("Alias")
	alias.Disable()
	var aliasShadow = ""

	server := widget.NewEntry()
	server.SetPlaceHolder("URL: sip://xxxxxx:5060")
	server.Disable()
	var serverShadow = ""

	calabel := widget.NewLabel("Certificate Authority CAROOT")
	ca := widget.NewMultiLineEntry()
	ca.Resize(fyne.NewSize(320, 120))
	ca.Disable()
	var caShadow = ""

	cclabel := widget.NewLabel("Client Certificate")
	cc := widget.NewMultiLineEntry()
	cc.Resize(fyne.NewSize(320, 120))
	cc.Disable()
	var ccShadow = ""

	cklabel := widget.NewLabel("Client Key")
	ck := widget.NewMultiLineEntry()
	ck.Resize(fyne.NewSize(320, 120))
	ck.Disable()
	var ckShadow = ""

	TPbutton := widget.NewButtonWithIcon("Try Password", theme.LoginIcon(), func() {
		errors.SetText("...")
		var iserrors = false
		ph, _ := config.LoadHashWithDefault("config.hash", "123456")

		//log.Println("pw ", MyPrefs.Password)
		pwh, err := bcrypt.GenerateFromPassword([]byte(password.Text), bcrypt.DefaultCost)
		config.PasswordHash = string(pwh)
		if err != nil {
			iserrors = true
			errors.SetText("Invalid Password")
		}

		// Comparing the password with the hash
		errpw := bcrypt.CompareHashAndPassword([]byte(ph), []byte(password.Text))
		if errpw != nil {
			iserrors = true
			errors.SetText("Passwords Do Not Match")
		}
		if !iserrors {
			errors.SetText("...")

			var preferedlanguageShadow = config.Decrypt(config.FyneApp.Preferences().StringWithFallback("eng", config.Encrypt(config.PreferedLanguage, config.MySecret)), config.MySecret)
			config.PreferedLanguage = config.Decrypt(preferedlanguageShadow, config.MySecret)

			aliasShadow = config.FyneApp.Preferences().StringWithFallback("SipAlias", config.Encrypt("SipAlias", config.MySecret))
			alias.SetText(config.Decrypt(aliasShadow, config.MySecret))

			serverShadow = config.FyneApp.Preferences().StringWithFallback("SipServerDNS", config.Encrypt("sip://mystic.newhorizons3000.org:5060", config.MySecret))
			server.SetText(config.Decrypt(serverShadow, config.MySecret))

			caShadow = config.FyneApp.Preferences().StringWithFallback("SIPCAROOT", config.Encrypt("-----BEGIN CERTIFICATE-----\nMIID7zCCAtegAwIBAgIUaXAPxJvZRRdTq5RWlwxs1XYo+5kwDQYJKoZIhvcNAQEL\nBQAwgYAxCzAJBgNVBAYTAlVTMRAwDgYDVQQIEwdGbG9yaWRhMRIwEAYDVQQHEwlD\ncmVzdHZpZXcxGjAYBgNVBAoTEU5ldyBIb3Jpem9ucyAzMDAwMQwwCgYDVQQLEwNX\nV1cxITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZzAeFw0yMzEyMTkw\nMzA4MDBaFw0yODEyMTcwMzA4MDBaMIGAMQswCQYDVQQGEwJVUzEQMA4GA1UECBMH\nRmxvcmlkYTESMBAGA1UEBxMJQ3Jlc3R2aWV3MRowGAYDVQQKExFOZXcgSG9yaXpv\nbnMgMzAwMDEMMAoGA1UECxMDV1dXMSEwHwYDVQQDExhuYXRzLm5ld2hvcml6b25z\nMzAwMC5vcmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCrVIXA/SxU\n7GeW92UNyiPnQEZgbJIHHQ31AQE2C/vFdpEtv32uoX1SsDl5drWvBrMnd5zrw1tL\nOEPA26tk/ACfQYL0n0HfeutLLu8H9jUWNp8ziX6Qbgd01M+/BixobHQjyDMxulo4\nJU2VK6QBLs9VI6TIihEU2BZhc/XCD9QbWcikAif1JySpz93MjFv3pcQU8ci4vQ0T\nImaGnHesr1qDbX1NuFVuBOPavZ64sQ1RsZtH5CdD+RU772wQWUgkPkwyUn8QBwTS\ne9XV5DNQD5nGEXjKTgjrd9KRf9pmRDnf6gBLi2r6C/l6q2w3ItOOHARdK0mc9CYh\ngY1Nzl59vrWdAgMBAAGjXzBdMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTAD\nAQH/MB0GA1UdDgQWBBR0qq9ueABC5RDsg/02FZFpBOR1hDAbBgNVHREEFDAShwTA\nqAAFhwTAqFjohwR/AAABMA0GCSqGSIb3DQEBCwUAA4IBAQBfdX0IMya9Dh9dHLJj\nnJZyb96htMWD5nuQQVBAu3ay+8O2GWj5mlsLJXAP2y7p/+3gyvHKTRDdJLux7N79\nHn6AYjmp3PCyZzuL1M/kHhSQxhxqJHGwjGXILt5pLovVkvkl4iukdxWJ5HAPsUGY\nO3QSDDFdoLflsG5VcrtdODm8uyxAjhMPAR2PXKfX8ABI79N7VKcbb98338fifrN8\n9H1r3BXcdsyhpH0gB0ZKJFSpMGWXlfudFEe9mXI9898xbEI2znqlYGhboVsuv5LM\nRESH2zXrkhmZyHqw0RtDROzyZOy5g1LcxbtVMn4w1LI4h3MDuE9B+Vud77A48qtA\ny+5x\n-----END CERTIFICATE-----\n", config.MySecret))
			ca.SetText(config.Decrypt(caShadow, config.MySecret))

			ccShadow = config.FyneApp.Preferences().StringWithFallback("SIPCACLIENT", config.Encrypt("-----BEGIN CERTIFICATE-----\nMIIEMTCCAxmgAwIBAgIUB7+OFX1LQrWtYMl5XIOXsOaLac0wDQYJKoZIhvcNAQEL\nBQAwgYAxCzAJBgNVBAYTAlVTMRAwDgYDVQQIEwdGbG9yaWRhMRIwEAYDVQQHEwlD\ncmVzdHZpZXcxGjAYBgNVBAoTEU5ldyBIb3Jpem9ucyAzMDAwMQwwCgYDVQQLEwNX\nV1cxITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZzAgFw0yMzEyMTkw\nMzA4MDBaGA8yMDUzMTIxMTAzMDgwMFowcjELMAkGA1UEBhMCVVMxEDAOBgNVBAgT\nB0Zsb3JpZGExEjAQBgNVBAcTCUNyZXN0dmlldzEaMBgGA1UEChMRTmV3IEhvcml6\nb25zIDMwMDAxITAfBgNVBAsTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZzCCASIw\nDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMWARyniHy8r342e3aKSsLDPwVMC\n2mRwuILP2JkXp5FllaFKnu/Z+0mF+iQlSchcC6DOcMQk00Cp/I8cCP865zyxPhqN\n2F2/qVItCU4+PTwe6ZnrfpJgXWwyk1hjS3vVNTT+idI5+pJgFH9YL0lbJ7q1UyPB\n+KP0x/c5T3K2Ec6U4uXhbVt/ePxFmsl1sHw6FE//XrA4EzbqCMEPCTcOfInvFrCJ\ny4/pAqjCxegT/1YDMNEdzmG8vg2tc3jPV+3GIAV3YL5nDE5mprHPEEDJtNQi+E4o\nXXXMobNhrJh9KJ59VbxTF8m5yM3b8fvof97OYhK0KYggplnTH+bhnYU9V5ECAwEA\nAaOBrTCBqjAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYD\nVR0TAQH/BAIwADAdBgNVHQ4EFgQUpffi3LSreerO756B/VnZkyyEVBIwHwYDVR0j\nBBgwFoAUdKqvbngAQuUQ7IP9NhWRaQTkdYQwNQYDVR0RBC4wLIIYbmF0cy5uZXdo\nb3Jpem9uczMwMDAub3JnhwR/AAABhwTAqAAFhwTAqFjoMA0GCSqGSIb3DQEBCwUA\nA4IBAQALlRqqW2HH4flFIgR/nh51gc/Hxv5xivhkzWUHHXRdltECSXknI4yBPchQ\n6Zsy0HZ7XQRlhQSIYd4Bp6eyHbny5t3JA978dHzpGJFCUVQDMY4yHLaCQgFJ+ESn\nwyyDWTRGA3cpEikL0B0ekDfqjWUEMTzmT/gnoSl0vM69nZDLZm1xMx1+EH+bpfFB\nRaVM6gKSAuFJmNYEL2e7JSags+3IHyVHkdo8GDlY//71Z4lxsFxFCF6xF9GDdAr2\niCA4OfydjiBSOz0eLJVgqkk1KGXtMqZXAojX62NrIWnFTW1Vzd46ekOHhq93B3tA\nkjWmHY/KdCZUjQSWss+YXgG4mI8c\n-----END CERTIFICATE-----\n", config.MySecret))
			cc.SetText(config.Decrypt(ccShadow, config.MySecret))
			ckShadow = config.FyneApp.Preferences().StringWithFallback("SIPCAKEY", config.Encrypt("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAxYBHKeIfLyvfjZ7dopKwsM/BUwLaZHC4gs/YmRenkWWVoUqe\n79n7SYX6JCVJyFwLoM5wxCTTQKn8jxwI/zrnPLE+Go3YXb+pUi0JTj49PB7pmet+\nkmBdbDKTWGNLe9U1NP6J0jn6kmAUf1gvSVsnurVTI8H4o/TH9zlPcrYRzpTi5eFt\nW394/EWayXWwfDoUT/9esDgTNuoIwQ8JNw58ie8WsInLj+kCqMLF6BP/VgMw0R3O\nYby+Da1zeM9X7cYgBXdgvmcMTmamsc8QQMm01CL4Tihddcyhs2GsmH0onn1VvFMX\nybnIzdvx++h/3s5iErQpiCCmWdMf5uGdhT1XkQIDAQABAoIBAB+Iu9QUJqaBetBB\n7WFnyo5wnY2DhxtCZDN+vDa1cCvm7F00bOwfAeBbY/UhfwZeq/yg+aBXwOMyQQEY\nmNcnsIQgSKo0u7c8Quy8BCBaD6zpwqKw1yTH/iKocJ5MPGEpSbWMbrUCTN/SN3Od\nwO8VfuJw0TWEYw7KpqLyo5zNNUqmczEO438CPGotbkFfzUqkumeUOsGWJFongyZY\na9EwpcTH2TkxuXum9SQVyLy+hSG/AEBp0cQPaRcoNh8sWYk43y5HrkIAqFo7dkMa\n9usAVMz9JCqIH2UNV04cDASFaiDMpYoD2hV2YHlL7/CQ7v5nb6OHT2A9aoSBOAfm\ns+dBzYECgYEA1l8+T9Xux73TCbFO2p7F094xSx4hhBZhaYpvzZoNN7iQdbdUVt2l\n1yHSoRgJUJMZlnKpMoNMLCxo34Lr3ww/TkIE/rrg10pqbqvojIDLCbi103EEB2v9\nWix8MSeOgFCa72T4lg9fDm5T493n4C5dade3LzZczUBF6dgmth3D+nMCgYEA69pa\nlob9n7eNXqDPk9kZUJV1jfLATC8eN4jupEiKfjnxEz9mUewvL/RF8kFhiS1ISC50\nKgM0v+isYBwwX00c7P02L6xCoGT35qOeoutEWVy/tYIHIHsD0jUBBsdnpQVNf58l\n9DDy2hZrpUwrsVHylVHpufBgKOfxgP2Jr3qD0OsCgYEAn4vzTGfkdzSIRMZ58awJ\ngE32Ufny5+PgTDSEUXk+LSJoIbR4SM5eB2dc5BiHljhk6twboUSnBJlo1DEUa8Up\nuIzaOtvLS3BPFl9LjIaulmWqrduHLB7rSJmjNNJD9KwJI/L6MHTwQkVKmmUllmvr\nikLKS5EiMICNiCUfaptsqJECgYEApYaSqzBEUdK1oeMErAPis16hqSTkdtNexqUQ\nrzXGFP6/Rb3qJra3C1XJvVLLjEW+hAIuPsoPPFyklbNS85+gHGc9n0mrXPxfy3ur\nuzWYu4rPdSizrcUIEoBmnwZVpEhLcrUUIwQzfIHdvJ3v0DvuH4PkoD2mjy7xnJDU\nD9bRKk8CgYAqK1lY5waFR0u3eFIPnrV4ATHXYuxcup2DCF+KJ6qwc4nNI6OB/ovU\nttiVZGr1rca42+XdWUQL5ufPFuKymeLbsuVzabbGKi+4RMvL+TIuorYtJRUPF+C7\nA9jlMeckpTZvl0yn5s3lC817N27B+U0M/jGow8sO0NtjBiImuTC5dg==\n-----END RSA PRIVATE KEY-----\n", config.MySecret))

			ck.SetText(config.Decrypt(ckShadow, config.MySecret))
			password.Disable()
			alias.Enable()
			server.Enable()

			ca.Enable()
			cc.Enable()
			ck.Enable()
		}
	})

	SSbutton := widget.NewButtonWithIcon("Logon", theme.LoginIcon(), func() {
		var haserrors = false
		if aliasShadow != alias.Text {
			haserrors = config.Edit("STRING", alias.Text)
			if !haserrors {
				config.Encrypt(alias.Text, config.MySecret)
				config.FyneApp.Preferences().SetString("Alias", config.Encrypt(alias.Text, config.MySecret))
			} else {
				errors.SetText("Invalid Alias")
			}
		}

		if serverShadow != server.Text {
			haserrors = config.Edit("SIP", server.Text)
			if !haserrors {
				config.FyneApp.Preferences().SetString("Server", config.Encrypt(server.Text, config.MySecret))
			} else {
				errors.SetText("Invalid Server")
			}

		}

		if !haserrors {
			config.LoggedOn = true

			config.SIPCaroot = ca.Text
			config.SIPClientCert = cc.Text
			config.SIPClientKey = ck.Text
			config.SIPLoggedOn = true
			password.Disable()
			server.Disable()
			alias.Disable()
			ca.SetText("")
			cc.SetText("")
			ck.SetText("")
			ca.Disable()
			ck.Disable()
			cc.Disable()

			errors.SetText("...")
		}
	})

	if config.SIPLoggedOn {
		TPbutton.Disable()
		TPbutton.Refresh()
		SSbutton.Disable()
		SSbutton.Refresh()

	}
	if !config.SIPLoggedOn {
		password.Enable()
		server.Disable()
		alias.Disable()
		ca.Disable()
		cc.Disable()
		ck.Disable()

	}

	vertbox := container.NewVBox(

		widget.NewLabelWithStyle("Logon", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		password,
		TPbutton,
		alias,
		server,

		calabel,
		ca,
		cclabel,
		cc,
		cklabel,
		ck,
		SSbutton,
		container.NewHBox(
			widget.NewHyperlink("newhorizons3000.org", config.ParseURL("https://newhorizons3000.org/")),
			widget.NewHyperlink("github.com", config.ParseURL("https://github.com/nh3000-org/snats")),
		),
		widget.NewLabel(""),
		//		themes,
		errors,
	)

	return container.NewScroll(
		vertbox,
	)

}
