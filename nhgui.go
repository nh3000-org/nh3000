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
	"strings"

	"github.com/nh3000-org/nh3000/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/widget"
)

var TopWindow fyne.Window

type Pane struct {
	Title, Intro string
	Icon         fyne.Resource
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var Panes = map[string]Pane{}
var PanesIndex = map[string][]string{}

func main() {

	var a = (app.NewWithID("org.nh3000.nh3000"))
	config.SetApp(a)

	if strings.HasPrefix(os.Getenv("LANG"), "en") {
		config.SetPreferedLanguage("eng")
	}
	if strings.HasPrefix(os.Getenv("LANG"), "sp") {
		config.SetPreferedLanguage("spa")

	}
	if strings.HasPrefix(os.Getenv("LANG"), "hn") {
		config.SetPreferedLanguage("hin")

	}

	Panes = map[string]Pane{
		"logon":    {config.GetLangs("ls-title"), "", theme.LoginIcon(), LogonScreen, true},
		"messages": {config.GetLangs("ms-title"), "", theme.MailSendIcon(), MessagesScreen, true},
		"settings": {config.GetLangs("ss-title"), "", theme.SettingsIcon(), SettingsScreen, true},
		"password": {config.GetLangs("ps-title"), "", theme.DocumentIcon(), PasswordScreen, true},
		"encdec":   {config.GetLangs("es-title"), "", theme.CheckButtonIcon(), EncdecScreen, true},
	}

	MyLogo, iconerr := fyne.LoadResourceFromPath("logo.png")
	if iconerr != nil {
		log.Println("logo error ", iconerr.Error())
	}

	var w = config.GetApp().NewWindow("NH3000")
	config.GetApp().SetIcon(MyLogo)
	Selected = Dark

	config.GetApp().Settings().SetTheme(MyTheme{})

	config.GetApp().SetIcon(MyLogo)

	logLifecycle()
	TopWindow = w
	w.SetMaster()

	intro := widget.NewLabel(config.GetLangs("mn-intro-1") + "\n" + "nats.io" + config.GetLangs("mn-intro-2"))
	intro.Wrapping = fyne.TextWrapWord

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon(Panes["logon"].Title, Panes["logon"].Icon, LogonScreen(w)),
		container.NewTabItemWithIcon(Panes["messages"].Title, Panes["messages"].Icon, MessagesScreen(w)),
		container.NewTabItemWithIcon(Panes["encdec"].Title, Panes["encdec"].Icon, EncdecScreen(w)),
		container.NewTabItemWithIcon(Panes["settings"].Title, Panes["settings"].Icon, SettingsScreen(w)),
		container.NewTabItemWithIcon(Panes["password"].Title, Panes["password"].Icon, PasswordScreen(w)),
	)

	w.SetContent(tabs)
	w.Resize(fyne.NewSize(640, 480))
	w.ShowAndRun()
}

// handle app close
func logLifecycle() {
	config.GetApp().Lifecycle().SetOnStopped(func() {
		if config.GetLoggedOn() {
			Send(config.GetLangs("ls-dis"), config.GetAlias())
		}
		if config.GetReceivingMessages() {
			QuitReceive <- true
		}
	})

}
