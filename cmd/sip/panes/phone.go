package panes

import (
	//"strconv"

	"encoding/json"
	"strings"

	"runtime"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"

	//"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget" // Softphone softphone

	"github.com/pion/webrtc/v3"
	"github.com/ringcentral/ringcentral-go"

	"github.com/nh3000-org/nh3000/config"
)

var memoryStats runtime.MemStats

func PhonesScreen(_ fyne.Window) fyne.CanvasObject {
	runtime.GC()
	runtime.ReadMemStats(&memoryStats)
	rc := ringcentral.RestClient{
		ClientID:     "rest",
		ClientSecret: "f0965fb361a90eb4e13a6b7cc805b6a4",
		Server:       "192.168.0.15",
	}
	rc.Authorize(ringcentral.GetTokenRequest{
		GrantType: "password",
		Username:  "Ed Lang",
		Extension: "2000",
		Password:  "12000",
	})
	bytes := rc.Post("/restapi/v1.0/client-info/sip-provision", strings.NewReader(`{"sipInfo":[{"transport":"WSS"}]}`))
	var createSipRegistrationResponse ringcentral.CreateSipRegistrationResponse
	json.Unmarshal(bytes, &createSipRegistrationResponse)
	softphone := config.Softphone{
		CreateSipRegistrationResponse: createSipRegistrationResponse,
	}
	softphone.Register()
	softphone.OnInvite = func(inviteMessage config.SipMessage) {
		softphone.Answer(inviteMessage)
	}
	softphone.OnTrack = func(track *webrtc.TrackRemote) {
		// ...
	}
	errors := widget.NewLabel("...")
	LogDetails := widget.NewLabel("")
	number := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	// custom button

	hangup := widget.NewButton("Hangup", func() {
		runtime.GC()
		runtime.ReadMemStats(&memoryStats)
		config.SIPFyneMainWin.SetTitle("Hanging Up " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
		LogDetails.SetText(LogDetails.Text + "\nHanging Up")
	})
	hangup.Disable()
	dial := widget.NewButton("Dial", func() {
		runtime.GC()
		runtime.ReadMemStats(&memoryStats)
		config.SIPFyneMainWin.SetTitle("Dialing " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")

		LogDetails.SetText(LogDetails.Text + "\nDialing")
		hangup.Enable()
	})
	clear := widget.NewButton("Clear", func() {
		runtime.GC()
		runtime.ReadMemStats(&memoryStats)
		config.SIPFyneMainWin.SetTitle("Cleared " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")

		LogDetails.SetText("Cleared")
		number.SetText("")
		hangup.Enable()
	})
	grid3 := container.New(layout.NewGridLayoutWithColumns(3), dial, hangup, clear)

	pad1 := widget.NewButton("1", func() {
		number.SetText(number.Text + "1")
		dial.Enable()
	})

	pad2 := widget.NewButton("2", func() {
		number.SetText(number.Text + "2")
		dial.Enable()
	})
	pad3 := widget.NewButton("3", func() {
		number.SetText(number.Text + "3")
		dial.Enable()
	})
	pad4 := widget.NewButton("4", func() {
		number.SetText(number.Text + "4")
		dial.Enable()
	})
	pad5 := widget.NewButton("5", func() {
		number.SetText(number.Text + "5")
		dial.Enable()
	})
	pad6 := widget.NewButton("6", func() {
		number.SetText(number.Text + "6")
		dial.Enable()
	})

	pad7 := widget.NewButton("7", func() {
		number.SetText(number.Text + "8")
		dial.Enable()
	})
	pad8 := widget.NewButton("8", func() {
		number.SetText(number.Text + "8")
		dial.Enable()
	})
	pad9 := widget.NewButton("9", func() {
		number.SetText(number.Text + "9")
		dial.Enable()
	})
	pad0 := widget.NewButton("0", func() {
		number.SetText(number.Text + "0")
		dial.Enable()
	})
	padast := widget.NewButton("*", func() {
		number.SetText(number.Text + "*")
		dial.Enable()
	})

	padpound := widget.NewButton("#", func() {
		number.SetText(number.Text + "#")
		dial.Enable()
	})

	LogDetails.SetText("Enter Number: ex 200 for extension or *97 for voicemail")

	var LogDetailsBorder = container.NewBorder(LogDetails, nil, nil, nil, nil)

	DetailsVW := container.NewScroll(LogDetailsBorder)
	DetailsVW.SetMinSize(fyne.NewSize(300, 320))

	phlabel := widget.NewLabel("Sip Phone")
	grid10 := container.New(layout.NewGridLayoutWithColumns(3), pad1, pad2, pad3, pad4, pad5, pad6, pad7, pad8, pad9, padast, pad0, padpound)

	topbox := container.NewVBox(
		phlabel,
		number,
		grid10,
		grid3,
		DetailsVW,
	)
	return container.NewBorder(
		topbox,
		errors,
		nil,
		nil,
		nil,
	)

}
