package gui

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/goccy/go-json"

	"github.com/nh3000-org/nh3000/nhcrypt"
	"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhpref"

	"github.com/nats-io/nats.go"
	"github.com/nh3000-org/nh3000/nhnats"
)

var EncMessage nhnats.MessageStore // message store
const QueueCheckInterval = 30      // check interval in seconds
var Labeltxt = widget.NewLabel(nhlang.GetLangs("ms-header1"))
var Errors = widget.NewLabel("...")

func messagesScreen(win fyne.Window) fyne.CanvasObject {
	Errors = widget.NewLabel("...")

	mymessage := widget.NewMultiLineEntry()
	mymessage.SetPlaceHolder(nhlang.GetLangs("ms-mm"))
	mymessage.SetMinRowsVisible(5)

	icon := widget.NewIcon(nil)
	Labeltxt = widget.NewLabel(nhlang.GetLangs("ms-header1"))
	label := container.NewHScroll(Labeltxt)
	//label := widget.NewLabel(GetLangs("ms-header1"))
	hbox := container.NewVScroll(label)

	hbox.SetMinSize(fyne.NewSize(300, 240))
	hbox.Refresh()
	List := widget.NewList(
		func() int {
			return len(nhnats.NatsMessages)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.CheckButtonCheckedIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			var short = nhnats.NatsMessages[id].MSmessage
			if len(nhnats.NatsMessages[id].MSmessage) > 12 {
				var short1 = strings.ReplaceAll(nhnats.NatsMessages[id].MSmessage, "\n", ".")
				short = short1[0:12]
			}

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(nhnats.NatsMessages[id].MSalias + " - " + short)
		},
	)

	List.OnSelected = func(id widget.ListItemID) {
		var mytext = nhnats.NatsMessages[id].MSmessage + "\n.................." + nhpref.NatsMessages[id].MShostname + nhpref.NatsMessages[id].MSipadrs + nhpref.NatsMessages[id].MSnodeuuid + nhpref.NatsMessages[id].MSiduuid + nhpref.NatsMessages[id].MSdate
		Labeltxt.SetText(mytext)
		icon.SetResource(theme.DocumentIcon())
	}
	List.OnUnselected = func(id widget.ListItemID) {
		Labeltxt.SetText(nhlang.GetLangs("ms-header1"))
		icon.SetResource(nil)
	}

	List.Resize(fyne.NewSize(500, 5000))
	List.Refresh()

	if nhpref.PasswordValid == true {

		smbutton := widget.NewButton(nhlang.GetLangs("ms-sm"), func() {
			nhnats.Send(mymessage.Text)
		})

		topbox := container.NewBorder(
			//widget.NewLabelWithStyle(GetLangs("ms-header2"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
			mymessage,
			nil,
			nil,
			nil,
			smbutton,
		)

		// copy to clipboard messages
		cpybutton := widget.NewButton(nhlang.GetLangs("ms-cpy"), func() {
			win.Clipboard().SetContent(Labeltxt.Text)
		})

		if !nhpref.LoggedOn {
			mymessage.Disable()
			smbutton.Disable()
			//recbutton.Disable()
			nhpref.ErrorMessage = nhlang.GetLangs("ms-err7")
		}
		bottombox := container.NewBorder(
			cpybutton,
			Errors,
			nil,
			nil,
			nil,
		)
		return container.NewBorder(
			topbox,
			bottombox,
			nil,
			nil,
			container.NewHSplit(List, container.NewCenter(hbox)),
		)

	}
	return container.NewBorder(

		widget.NewLabelWithStyle(nhlang.GetLangs("ms-err7"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		container.NewHSplit(List, container.NewCenter(hbox)),
	)
}

func HandleMessage(m *nats.Msg) {
	ms := nhnats.MessageStore{}
	var inmap = true // unique message id
	ejson, err := nhcrypt.Decrypt(string(m.Data), nhpref.Queuepassword)
	if err != nil {
		ejson = string(m.Data)
	}
	err1 := json.Unmarshal([]byte(ejson), &ms)
	if err1 != nil {
		ejson = "Unknown"
	}

	inmap = NodeMap("MI" + ms.MSiduuid)
	if inmap == false {
		nhnats.NatsMessages = append(nhnats.NatsMessages, ms)
	}

}

func NodeMap(node string) bool {
	_, x := nhpref.MyMap[node]
	return x
}
