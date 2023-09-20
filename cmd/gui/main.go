// MIT License

// Copyright (c) 2023 nh3000-org

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// A Secure client using NATS messaging system (https://newhorizons3000.org).

package main

import (
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/nh3000-org/snats/nhlang"
)

const preferenceCurrentApplication = "logon"

var TopWindow fyne.Window

type Pane struct {
	Title, Intro string
	View         func(w fyne.Window) fyne.CanvasObject
	SupportWeb   bool
}

var Panes = map[string]MyPane{}
var PanesIndex = map[string][]string{}

func main() {

	Panes = map[string]Pane{
		"password":     {GetLangs("ps-title"), "", passwordScreen, true},
		"settings":     {GetLangs("ss-title"), "", settingsScreen, true},
		"certificates": {GetLangs("cs-title"), "", certificatesScreen, true},
		"logon":        {GetLangs("ls-title"), "", logonScreen, true},
		"messages":     {GetLangs("ms-title"), "", messagesScreen, true},
		"encdec":       {GetLangs("es-title"), "", encdecScreen, true},
	}

	// PanesIndex  defines how our panes should be laid out in the index tree
	PanesIndex = map[string][]string{
		"": {"password", "logon", "settings", "certificates", "messages", "encdec"},
	}

	Json("LOAD")



	if strings.HasPrefix(os.Getenv("LANG"), "en") {
		panes.PreferedLanguage = "eng"
	}
	if strings.HasPrefix(os.Getenv("LANG"), "sp") {
		panes.PreferedLanguage = "spa"
	}

	MyLogo, _ := fyne.LoadResourceFromPath("logo.png")
	panes.MyAppDup.SetIcon(MyLogo)
	makeTray(panes.MyAppDup)
	logLifecycle(panes.MyAppDup)

	w := panes.MyAppDup.NewWindow("SNATS BETA.2")
	TopWindow = w
	w.SetMaster()

	content := container.NewMax()
	title := widget.NewLabel("SNATS")

	intro := widget.NewLabel(panes.GetLangs("mn-intro-1") + "\n" + "nats.io" + panes.GetLangs("mn-intro-2"))
	intro.Wrapping = fyne.TextWrapWord
	setPanes := func(t panes.MyPane) {
		if fyne.CurrentDevice().IsMobile() {
			child := panes.MyAppDup.NewWindow(t.Title)
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
		w.SetContent(makeNav(setPanes, false))
	} else {
		split := container.NewHSplit(makeNav(setPanes, true), pane)
		split.Offset = 0.2
		w.SetContent(split)
	}

	w.Resize(fyne.NewSize(640, 460))
	w.ShowAndRun()
}

func logLifecycle(a fyne.App) {

	a.Lifecycle().SetOnStopped(func() {
		panes.FormatMessage("Disconnected")

	})

}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem(panes.GetLangs("mn-mt"), func() {})
		menu := fyne.NewMenu(panes.GetLangs("mn-mt"), h)
		h.Action = func() {
			h.Label = panes.GetLangs("mn-mt")
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}

func unsupportedApplication(t panes.MyPane) bool {
	return !t.SupportWeb && fyne.CurrentDevice().IsBrowser()
}

func makeNav(setTutorial func(panes panes.MyPane), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()
	a.Settings().SetTheme(theme.DarkTheme())
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return panes.MyPanesIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := panes.MyPanesIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := panes.MyPanes[uid]
			if !ok {
				fyne.LogError(panes.GetLangs("mn-err1")+uid, nil)
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
			if t, ok := panes.MyPanes[uid]; ok {
				if unsupportedApplication(t) {
					return
				}
				a.Preferences().SetString(preferenceCurrentApplication, "logon")
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentApplication, "logon")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton(panes.GetLangs("mn-dark"), func() {
			a.Settings().SetTheme(theme.DarkTheme())
		}),
		widget.NewButton(panes.GetLangs("mn-light"), func() {
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
