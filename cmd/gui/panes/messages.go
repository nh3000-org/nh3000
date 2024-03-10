package panes

import (
	"strings"

	"github.com/nh3000-org/nh3000/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var Details = widget.NewLabel("")
var ackMsgId = ""

func MessagesScreen(win fyne.Window) fyne.CanvasObject {

	config.SetMessageWindow(win)
	//var Acknode = ""
	var Errors = widget.NewLabel("...")
	var DetailsTop = widget.NewLabel(config.GetLangs("ms-header1"))
	//var DetailsBottom = widget.NewButton(config.GetLangs("ms-header1"))
	DetailButton := widget.NewButton("Ack", func() {
		if !config.GetLoggedOn() {
			Errors.SetText(config.GetLangs("cs-lf"))
			return
		}
		config.SetAck(ackMsgId)

	})
	var DetailsBorder = container.NewBorder(DetailsTop, DetailButton, nil, nil, Details)
	message := widget.NewMultiLineEntry()
	message.SetPlaceHolder(config.GetLangs("ms-mm"))
	message.SetMinRowsVisible(2)

	Filter := widget.NewCheck(config.GetLangs("ms-filter"), func(on bool) { config.SetFilter(on) })

	DetailsVW := container.NewScroll(DetailsBorder)

	DetailsVW.SetMinSize(fyne.NewSize(300, 240))
	DetailsVW.Refresh()
	if config.GetClearMessageDetail() {
		Details.SetText("")
		Details.Refresh()
		DetailsVW.Refresh()
		config.SetClearMessageDetail(false)
	}

	List := widget.NewList(
		func() int {
			return len(config.NatsMessages)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			var short = config.NatsMessages[id].MSmessage
			if len(config.NatsMessages[id].MSmessage) > 20 {
				var short1 = strings.ReplaceAll(config.NatsMessages[id].MSmessage, "\n", ".")
				short = short1[0:12]
			}
			//Acknode = nhnats.NatsMessages[id].MSiduuid
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(config.NatsMessages[id].MSalias + " - " + short)
		},
	)
	config.SetMessageList(List)
	List.OnSelected = func(id widget.ListItemID) {
		var mytext = config.NatsMessages[id].MSmessage + "\n.................." + config.NatsMessages[id].MShostname + config.NatsMessages[id].MSipadrs + config.NatsMessages[id].MSnodeuuid + config.NatsMessages[id].MSiduuid + config.NatsMessages[id].MSdate
		Details.SetText(mytext)
		ackMsgId = config.NatsMessages[id].MSiduuid

	}
	List.OnUnselected = func(id widget.ListItemID) {
		Details.SetText(config.GetLangs("ms-header1"))
	}

	List.Resize(fyne.NewSize(500, 5000))
	List.Refresh()

	smbutton := widget.NewButton(config.GetLangs("ms-sm"), func() {
		if !config.GetLoggedOn() {
			Errors.SetText(config.GetLangs("cs-lf"))

		}
		config.Send(message.Text, config.GetAlias())
		message.SetText("")
	})
	/* 		ackbutton := widget.NewButton(nhlang.GetLangs("ms-ack"), func() {
		SendAck(Acknode)
	}) */
	topbox := container.NewBorder(
		nil,
		Filter,
		nil,
		smbutton,
		message,
	)

	// copy to clipboard messages
	cpybutton := widget.NewButton(config.GetLangs("ms-cpy"), func() {
		win.Clipboard().SetContent(Details.Text)
	})

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
		container.NewHSplit(List, DetailsVW),
	)

}
