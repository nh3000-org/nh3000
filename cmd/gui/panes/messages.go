package panes

import (
	"strings"

	"github.com/nh3000-org/nh3000/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var ackMsgId = ""
var mymessage = ""
var mymessageshort = ""

func MessagesScreen(win fyne.Window) fyne.CanvasObject {

	config.FyneWin = win
	message := widget.NewMultiLineEntry()
	message.SetPlaceHolder(config.GetLangs("ms-mm"))
	message.SetMinRowsVisible(2)

	Filter := widget.NewCheck(config.GetLangs("ms-filter"), func(on bool) { config.FyneFilter = true })

	var Errors = widget.NewLabel("...")
	//var DetailsTop = widget.NewLabel(config.GetLangs("ms-header1") + ".")
	AckButton := widget.NewButtonWithIcon("Ack", theme.MailSendIcon(), func() {
		if !config.LoggedOn {
			Errors.SetText(config.GetLangs("cs-lf"))
			return
		}
		config.SetAck(ackMsgId)

	})
	Details := widget.NewLabel("")
	var DetailsBorder = container.NewBorder(Details, nil, nil, nil, nil)

	DetailsVW := container.NewScroll(DetailsBorder)
	DetailsVW.SetMinSize(fyne.NewSize(300, 240))

	cpybutton := widget.NewButtonWithIcon(config.GetLangs("ms-cpy"), theme.ContentCopyIcon(), func() {
		win.Clipboard().SetContent(Details.Text)
	})
	List := widget.NewList(
		func() int {
			return len(config.NatsMessages)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			mymessage = config.NatsMessages[id].MSmessage
			if len(config.NatsMessages[id].MSmessage) > 100 {
				mymessageshort = strings.ReplaceAll(config.NatsMessages[id].MSmessage, "\n", ".")
				mymessage = mymessageshort[0:100]
			}
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(config.NatsMessages[id].MSalias + " - " + mymessage)
		},
	)
	config.FyneMessageList = List
	List.OnSelected = func(id widget.ListItemID) {
		//config.GetApp().Settings().SetTheme(config.MyTheme{})
		Details.SetText(config.NatsMessages[id].MSmessage + "\n.................." + config.NatsMessages[id].MShostname + config.NatsMessages[id].MSipadrs + config.NatsMessages[id].MSnodeuuid + config.NatsMessages[id].MSiduuid + config.NatsMessages[id].MSdate)
		dlg := fyne.CurrentApp().NewWindow(config.NatsMessages[id].MSalias + config.NatsMessages[id].MSdate)
		DetailsVW := container.NewScroll(DetailsBorder)
		DetailsVW.SetMinSize(fyne.NewSize(300, 240))
		DetailsBottom := container.NewBorder(cpybutton, AckButton, nil, nil, nil)
		dlg.SetContent(container.NewBorder(DetailsVW, DetailsBottom, nil, nil, nil))
		dlg.Show()

		ackMsgId = config.NatsMessages[id].MSiduuid
		List.Unselect(id)
	}
	List.OnUnselected = func(id widget.ListItemID) {
		//Details.SetText(config.GetLangs("ms-header1") + "..")
	}

	List.Resize(fyne.NewSize(500, 5000))
	List.Refresh()

	smbutton := widget.NewButtonWithIcon(config.GetLangs("ms-sm"), theme.MailSendIcon(), func() {
		if !config.LoggedOn {
			Errors.SetText(config.GetLangs("cs-lf"))
		}
		config.Send(message.Text, config.NatsAlias)
		message.SetText("")
	})

	topbox := container.NewBorder(
		nil,
		Filter,
		nil,
		smbutton,
		message,
	)

	bottombox := container.NewBorder(
		nil,
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
		List,
	)

}
