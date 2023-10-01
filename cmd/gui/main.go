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
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhnats"
	"github.com/nh3000-org/nh3000/nhpanes"
	"github.com/nh3000-org/nh3000/nhpref"
	"github.com/nh3000-org/nh3000/nhutil"
)

const preferenceCurrentApplication = "logon"

var TopWindow fyne.Window

type Pane struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var Panes = map[string]Pane{}
var PanesIndex = map[string][]string{}

func main() {
	nhpref.Load()
	nhpref.Save()
	if strings.HasPrefix(os.Getenv("LANG"), "en") {
		nhpref.PreferedLanguage = "eng"
	}
	if strings.HasPrefix(os.Getenv("LANG"), "sp") {
		nhpref.PreferedLanguage = "spa"
	}

	Panes = map[string]Pane{
		"password":     {nhlang.GetLangs("ps-title"), "", nhpanes.PasswordScreen, true},
		"settings":     {nhlang.GetLangs("ss-title"), "", nhpanes.SettingsScreen, true},
		"certificates": {nhlang.GetLangs("cs-title"), "", nhpanes.CertificatesScreen, true},
		"logon":        {nhlang.GetLangs("ls-title"), "", nhpanes.LogonScreen, true},
		"messages":     {nhlang.GetLangs("ms-title"), "", nhpanes.MessagesScreen, true},
		"encdec":       {nhlang.GetLangs("es-title"), "", nhpanes.EncdecScreen, true},
	}

	// PanesIndex  defines how our panes should be laid out in the index tree
	PanesIndex = map[string][]string{
		"": {"password", "logon", "settings", "certificates", "messages", "encdec"},
	}

	//if strings.HasPrefix(os.Getenv("LANG"), "en") {

	MyLogo, _ := fyne.LoadResourceFromPath("logo.png")

	w := nhutil.GetApp().NewWindow("SNATS BETA.3")
	nhutil.GetApp().SetIcon(MyLogo)
	makeTray(nhutil.GetApp())
	logLifecycle(nhutil.GetApp())
	TopWindow = w
	w.SetMaster()

	content := container.NewStack()
	title := widget.NewLabel("SNATS")

	intro := widget.NewLabel(nhlang.GetLangs("mn-intro-1") + "\n" + "nats.io" + nhlang.GetLangs("mn-intro-2"))
	intro.Wrapping = fyne.TextWrapWord
	SetPanes := func(t Pane) {
		if fyne.CurrentDevice().IsMobile() {
			child := nhutil.GetApp().NewWindow(t.Title)
			TopWindow = child
			child.SetContent(t.View(TopWindow))
			child.Show()
			child.SetOnClosed(func() {
				TopWindow = w
			})
			return
		}

		title.SetText(t.Title)
		intro.SetText(t.Intro)

		content.Objects = []fyne.CanvasObject{t.View(w)}
		content.Refresh()
	}

	pane := container.NewBorder(
		container.NewVBox(title, widget.NewSeparator(), intro), nil, nil, nil, content)
	if fyne.CurrentDevice().IsMobile() {
		w.SetContent(makeNav(SetPanes, false))
	} else {
		split := container.NewHSplit(makeNav(SetPanes, true), pane)
		split.Offset = 0.2
		w.SetContent(split)
	}
	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}

func logLifecycle(a fyne.App) {
	a.Lifecycle().SetOnStopped(func() {
		if nhpref.LoggedOn {
			nhnats.Send(nhlang.GetLangs("ls-dis"))
			nhpref.CertsOffFS()
		}
	})

}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem(nhlang.GetLangs("mn-mt"), func() {})
		menu := fyne.NewMenu(nhlang.GetLangs("mn-mt"), h)
		h.Action = func() {
			h.Label = nhlang.GetLangs("mn-mt")
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}

func unsupportedApplication(t Pane) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}

func makeNav(setGui func(panes Pane), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()
	a.Settings().SetTheme(theme.DarkTheme())
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return PanesIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := PanesIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := Panes[uid]
			if !ok {
				fyne.LogError(nhlang.GetLangs("mn-err1")+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
			if unsupportedApplication(t) {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{Italic: true}
			} else {
				obj.(*widget.Label).TextStyle = fyne.TextStyle{}
			}
		},
		OnSelected: func(uid string) {
			if t, ok := Panes[uid]; ok {
				if unsupportedApplication(t) {
					return
				}
				a.Preferences().SetString(preferenceCurrentApplication, "logon")
				setGui(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentApplication, "logon")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton(nhlang.GetLangs("mn-dark"), func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton(nhlang.GetLangs("mn-light"), func() {
			a.Settings().SetTheme(theme.LightTheme())
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

func shortcutFocused(s fyne.Shortcut, w fyne.Window) {
	switch sh := s.(type) {
	case *fyne.ShortcutCopy:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutCut:
		sh.Clipboard = w.Clipboard()
	case *fyne.ShortcutPaste:
		sh.Clipboard = w.Clipboard()
	}
	if focused, ok := w.Canvas().Focused().(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}
