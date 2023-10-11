package nhutil

import (
	"net/url"

	"fyne.io/fyne/v2"
)

type util interface {
	ParseURL(string) *url.URL

	SetApp(fyne.App)
	GetApp() fyne.App
	Edit(string, string) bool
}

var App fyne.App
var Win fyne.Window

func ParseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func SetApp(a fyne.App) {
	App = a
}
func GetApp() fyne.App {
	return App
}
func SetMessageWindow(a fyne.Window) {
	Win = a
}
func GetMessageWin() fyne.Window {
	return Win
}