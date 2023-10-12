package nhutil

import (
	"net/url"

	"fyne.io/fyne/v2"
)

var App fyne.App
var Win fyne.Window

// parse a url
func ParseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

// set app id
func SetApp(a fyne.App) {
	App = a
}

// return app id
func GetApp() fyne.App {
	return App
}

// set message window
func SetMessageWindow(a fyne.Window) {
	Win = a
}
func GetMessageWin() fyne.Window {
	return Win
}
