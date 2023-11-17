package nhpanes

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhnats"
	"github.com/nh3000-org/nh3000/nhpref"
	"github.com/nh3000-org/nh3000/nhutil"
)

var Details = widget.NewLabel(nhlang.GetLangs("ms-header1"))

func MessagesScreen(win fyne.Window) fyne.CanvasObject {
	nhutil.SetMessageWindow(win)
	//var Acknode = ""
	var Errors = widget.NewLabel("...")
	var Details = widget.NewLabel(nhlang.GetLangs("ms-header1"))
	message := widget.NewMultiLineEntry()
	message.SetPlaceHolder(nhlang.GetLangs("ms-mm"))
	message.SetMinRowsVisible(2)

	Filter := widget.NewCheck(nhlang.GetLangs("ms-filter"), func(on bool) { nhpref.Filter = on })

	DetailsHS := container.NewHScroll(Details)
	DetailsHS.Refresh()
	DetailsVS := container.NewVScroll(DetailsHS)
	DetailsVS.SetMinSize(fyne.NewSize(300, 240))
	DetailsVS.Refresh()
	if nhpref.ClearMessageDetail {
		Details.SetText("")
		Details.Refresh()
		DetailsHS.Refresh()
		DetailsVS.Refresh()
		nhpref.ClearMessageDetail = false
	}

	List := widget.NewList(
		func() int {
			return len(nhnats.NatsMessages)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			var short = nhnats.NatsMessages[id].MSmessage
			if len(nhnats.NatsMessages[id].MSmessage) > 12 {
				var short1 = strings.ReplaceAll(nhnats.NatsMessages[id].MSmessage, "\n", ".")
				short = short1[0:12]
			}
			//Acknode = nhnats.NatsMessages[id].MSiduuid
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(nhnats.NatsMessages[id].MSalias + " - " + short)
		},
	)

	List.OnSelected = func(id widget.ListItemID) {
		var mytext = nhnats.NatsMessages[id].MSmessage + "\n.................." + nhnats.NatsMessages[id].MShostname + nhnats.NatsMessages[id].MSipadrs + nhnats.NatsMessages[id].MSnodeuuid + nhnats.NatsMessages[id].MSiduuid + nhnats.NatsMessages[id].MSdate
		Details.SetText(mytext)

	}
	List.OnUnselected = func(id widget.ListItemID) {
		Details.SetText(nhlang.GetLangs("ms-header1"))
	}

	List.Resize(fyne.NewSize(500, 5000))
	List.Refresh()

	if nhpref.LoggedOn == true {

		smbutton := widget.NewButton(nhlang.GetLangs("ms-sm"), func() {
			nhnats.Send(message.Text)
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
		cpybutton := widget.NewButton(nhlang.GetLangs("ms-cpy"), func() {
			win.Clipboard().SetContent(Details.Text)
		})

		if !nhpref.LoggedOn {
			message.Disable()
			smbutton.Disable()
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
			container.NewHSplit(List, DetailsHS),
		)

	}
	return container.NewBorder(

		widget.NewLabelWithStyle(nhlang.GetLangs("ms-err7"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		container.NewHSplit(List, container.NewCenter(DetailsHS)),
	)
}
