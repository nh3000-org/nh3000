// Copyright 2012-2023 The NH3000 Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A Go client for the NH3000 messaging system (https://newhorizons3000.org).

package main

import (
	"log"
	"os"
	"runtime"
	"strconv"

	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"github.com/nh3000-org/nh3000/cmd/sip/panes"
	"github.com/nh3000-org/nh3000/config"

	"fyne.io/fyne/v2/widget"
)

var TopWindow fyne.Window
var memoryStats runtime.MemStats

type Pane struct {
	Title, Intro string
	Icon         fyne.Resource
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

func main() {
	var a = app.NewWithID("org.nh3000.nh3000.SIP")
	config.FyneApp = a
	runtime.GC()
	runtime.ReadMemStats(&memoryStats)
	var w = a.NewWindow("NH3000 SIP" + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
	config.SIPFyneMainWin = w
	config.SIPPreferedLanguage = "eng"
	if strings.HasPrefix(os.Getenv("LANG"), "en") {
		config.SIPPreferedLanguage = "eng"
	}
	if strings.HasPrefix(os.Getenv("LANG"), "sp") {
		config.PreferedLanguage = "spa"
	}
	if strings.HasPrefix(os.Getenv("LANG"), "hn") {
		config.PreferedLanguage = "hin"

	}

	MyLogo, iconerr := fyne.LoadResourceFromPath("Icon.png")
	if iconerr != nil {
		log.Println("Icon.png error ", iconerr.Error())
	}
	config.Selected = config.Dark
	config.FyneApp.Settings().SetTheme(config.MyTheme{})
	config.FyneApp.SetIcon(MyLogo)

	logLifecycle()
	TopWindow = w
	w.SetMaster()

	intro := widget.NewLabel("Intro" + "\n" + "nats.io" + "Intro2")
	intro.Wrapping = fyne.TextWrapWord
	var Panes = map[string]Pane{
		"logon":    {"Logon", "", theme.LoginIcon(), panes.LogonScreen, true},
		"phone":    {"Phone", "", theme.MailSendIcon(), panes.PhonesScreen, true},
		"settings": {"Setings", "", theme.SettingsIcon(), panes.SettingsScreen, true},
		"password": {"Password", "", theme.DocumentIcon(), panes.PasswordScreen, true},
		//"encdec":   {config.GetLangs("es-title"), "", theme.CheckButtonIcon(), panes.EncdecScreen, true},
	}

	config.SIPFyneMainWin.SetContent(container.NewAppTabs(
		container.NewTabItemWithIcon(Panes["logon"].Title, Panes["logon"].Icon, panes.LogonScreen(config.FyneMainWin)),
		container.NewTabItemWithIcon(Panes["phone"].Title, Panes["phone"].Icon, panes.PhonesScreen(config.FyneMainWin)),
		//container.NewTabItemWithIcon(Panes["encdec"].Title, Panes["encdec"].Icon, panes.EncdecScreen(config.FyneMainWin)),
		container.NewTabItemWithIcon(Panes["settings"].Title, Panes["settings"].Icon, panes.SettingsScreen(config.FyneMainWin)),
		container.NewTabItemWithIcon(Panes["password"].Title, Panes["password"].Icon, panes.PasswordScreen(config.FyneMainWin)),
	))

	config.SIPFyneMainWin.Resize(fyne.NewSize(640, 480))
	config.SIPFyneMainWin.ShowAndRun()
}

// handle app close
func logLifecycle() {

	/* 	config.FyneApp.Lifecycle().SetOnStopped(func() {
		if config.LoggedOn {
			config.Send(config.GetLangs("ls-dis"), config.GetAlias())
		}
		if config.GetReceivingMessages() {
			config.QuitReceive <- true
		}
	}) */

}
