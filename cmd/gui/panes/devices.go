package panes

import (
	//"log"
	"bytes"
	//"log"
	"strings"

	"github.com/nh3000-org/nh3000/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

var mydevice = ""
var mydevicehort = ""

// var mymessagedlg []byte
var selectedmsdevice int
var selectedseqdevice uint64
var selecteduuiddevice string
var selectedalias string
var ErrorsDevice = widget.NewLabel("...")

func DevicesScreen(win fyne.Window) fyne.CanvasObject {

	//message := widget.NewMultiLineEntry()
	//message.SetPlaceHolder(config.GetLangs("ms-mm"))
	//message.SetMinRowsVisible(2)

	Details := widget.NewLabel("")
	var DetailsBorder = container.NewBorder(Details, nil, nil, nil, nil)

	DetailsVW := container.NewScroll(DetailsBorder)
	DetailsVW.SetMinSize(fyne.NewSize(300, 240))

	cpybutton := widget.NewButtonWithIcon(config.GetLangs("ms-cpy"), theme.ContentCopyIcon(), func() {
		win.Clipboard().SetContent(Details.Text)
	})
	delbutton := widget.NewButtonWithIcon(config.GetLangs("ms-del"), theme.ContentCopyIcon(), func() {
		config.DeleteNatsMessage(selectedseqdevice)
		delete(config.NatsMessagesDevice, selectedms)
		delete(config.NatsMessagesIndexDevice, selecteduuid)
	})
	authbutton := widget.NewButtonWithIcon(config.GetLangs("dv-auth"), theme.ContentCopyIcon(), func() {

		config.Send("authorizations."+selectedalias, "AUTHORIZED", config.NatsAlias)
		config.DeleteNatsMessage(selectedseqdevice)
		delete(config.NatsMessagesDevice, selectedms)
		delete(config.NatsMessagesIndexDevice, selecteduuid)
	})
	//a, aerr := config.NewNatsJS("DEVICES", "devices"+config.NatsAlias, config.NatsAlias)

	config.FyneDeviceWin = win
	List := widget.NewList(
		func() int {
			return len(config.NatsMessagesDevice)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {

			mymessage = config.NatsMessagesDevice[id].MSmessage
			if len(config.NatsMessagesDevice[id].MSmessage) > 100 {
				mymessageshort = strings.ReplaceAll(config.NatsMessagesDevice[id].MSmessage, "\n", ".")
				mymessage = mymessageshort[0:100]
			}
			item.(*fyne.Container).Objects[0].(*widget.Label).SetText(config.NatsMessagesDevice[id].MSsubject + " " + config.NatsMessagesDevice[id].MSalias + " - " + mymessage)
		},
	)
	config.FyneDeviceList = List
	List.OnSelected = func(id widget.ListItemID) {
		selectedmsdevice = id
		selectedseqdevice = config.NatsMessagesDevice[id].MSsequence
		selecteduuiddevice = config.NatsMessagesDevice[id].MSnodeuuid
		selectedalias = config.NatsMessagesDevice[id].MSalias
		Details.SetText(config.NatsMessagesDevice[id].MSmessage + "\n." + "\n" + config.NatsMessagesDevice[id].MSsubject + " " + config.NatsMessagesDevice[id].MSos + " " + config.NatsMessagesDevice[id].MShostname + " on " + config.NatsMessagesDevice[id].MSdate + "\n" + config.NatsMessagesDevice[id].MSipadrs + "\n" + config.NatsMessagesDevice[id].MSmacid + "\nNode ID: " + config.NatsMessagesDevice[id].MSnodeuuid + "\nMsg ID:" + config.NatsMessagesDevice[id].MSiduuid)
		dlg := fyne.CurrentApp().NewWindow(msg2dlgdev(id))
		DetailsVW := container.NewScroll(DetailsBorder)
		DetailsVW.SetMinSize(fyne.NewSize(300, 240))
		DetailsBottom := container.NewBorder(cpybutton, delbutton, nil, nil, nil)
		dlg.SetContent(container.NewBorder(DetailsVW, DetailsBottom, nil, nil, authbutton))
		dlg.Show()
		List.Unselect(id)
	}

	return container.NewBorder(
		nil,
		Errors,
		nil,
		nil,
		List,
	)

}
func msg2dlgdev(id int) string {
	//config.NatsMessages[id].MSmessage + "\n.................." + "\n" + config.NatsMessages[id].MSsubject + "\n" + config.NatsMessages[id].MSos + "\n" + config.NatsMessages[id].MShostname + "\n" + config.NatsMessages[id].MSipadrs + "\n" + config.NatsMessages[id].MSnodeuuid + "\n" + config.NatsMessages[id].MSiduuid + "\n" + config.NatsMessages[id].MSdate
	buf := &bytes.Buffer{}
	buf.WriteString(config.NatsMessagesDevice[id].MSmessage)
	buf.WriteString("\n[")

	buf.WriteString(config.NatsMessagesDevice[id].MSsubject)
	buf.WriteString("] from ")
	buf.WriteString(config.NatsMessagesDevice[id].MSalias)
	buf.WriteString(" on ")
	buf.WriteString(config.NatsMessagesDevice[id].MSos)
	buf.WriteString(" host ")
	buf.WriteString(config.NatsMessagesDevice[id].MShostname)
	buf.WriteString("\n[net] ")
	buf.WriteString(config.NatsMessagesDevice[id].MSipadrs)
	buf.WriteString("\n[mac] ")
	buf.WriteString(config.NatsMessagesDevice[id].MSmacid)
	buf.WriteString("\n[msgid] ")
	buf.WriteString(config.NatsMessagesDevice[id].MSiduuid)
	buf.WriteString("\n[node] ")
	buf.WriteString(config.NatsMessagesDevice[id].MSnodeuuid)
	buf.WriteString("\n[on] ")
	buf.WriteString(config.NatsMessagesDevice[id].MSdate)
	return buf.String()
}
