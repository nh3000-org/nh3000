package panes

import (
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/goccy/go-json"

	//"github.com/nh30000-org/nats.go"
	"github.com/nats-io/nats.go"
)

var EncMessage MessageStore   // message store
const QueueCheckInterval = 30 // check interval in seconds
var Labeltxt = widget.NewLabel(GetLangs("ms-header1"))
var Errors = widget.NewLabel("...")

func messagesScreen(win fyne.Window) fyne.CanvasObject {
	Errors = widget.NewLabel("...")

	mymessage := widget.NewMultiLineEntry()
	mymessage.SetPlaceHolder(GetLangs("ms-mm"))
	mymessage.SetMinRowsVisible(5)

	icon := widget.NewIcon(nil)
	Labeltxt = widget.NewLabel(GetLangs("ms-header1"))
	label := container.NewHScroll(Labeltxt)
	//label := widget.NewLabel(GetLangs("ms-header1"))
	hbox := container.NewVScroll(label)

	hbox.SetMinSize(fyne.NewSize(300, 240))
	hbox.Refresh()
	List := widget.NewList(
		func() int {
			return len(NatsMessages)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.CheckButtonCheckedIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			var short = NatsMessages[id].MSmessage
			if len(NatsMessages[id].MSmessage) > 12 {
				var short1 = strings.ReplaceAll(NatsMessages[id].MSmessage, "\n", ".")
				short = short1[0:12]
			}

			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(NatsMessages[id].MSalias + " - " + short)
		},
	)

	List.OnSelected = func(id widget.ListItemID) {
		var mytext = NatsMessages[id].MSmessage + "\n.................." + NatsMessages[id].MShostname + NatsMessages[id].MSipadrs + NatsMessages[id].MSnodeuuid + NatsMessages[id].MSiduuid + NatsMessages[id].MSdate
		Labeltxt.SetText(mytext)
		icon.SetResource(theme.DocumentIcon())
	}
	List.OnUnselected = func(id widget.ListItemID) {
		Labeltxt.SetText(GetLangs("ms-header1"))
		icon.SetResource(nil)
	}

	List.Resize(fyne.NewSize(500, 5000))
	List.Refresh()

	if PasswordValid == true {

		smbutton := widget.NewButton(GetLangs("ms-sm"), func() {
			FormatMessage(mymessage.Text)
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
		cpybutton := widget.NewButton(GetLangs("ms-cpy"), func() {
			win.Clipboard().SetContent(Labeltxt.Text)
		})
		go ReceiveMsg()
		if !LoggedOn {
			mymessage.Disable()
			smbutton.Disable()
			//recbutton.Disable()
			ErrorMessage = GetLangs("ms-err7")
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

		widget.NewLabelWithStyle(GetLangs("ms-err7"), fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		nil,
		nil,
		nil,
		container.NewHSplit(List, container.NewCenter(hbox)),
	)
}
func ReceiveMsg() {

	for {
		NatsMessages = nil
		Labeltxt.SetText(GetLangs("ms-header1"))

		nc, err := nats.Connect(Server, nats.RootCAsMem([]byte(Caroot)), nats.ClientCertMem([]byte(Clientcert), []byte(Clientkey)))
		if err != nil {
			Errors.SetText(GetLangs("ms-err2"))

		}

		js, _ := nc.JetStream()
		js.AddStream(&nats.StreamConfig{
			Name: Queue + NodeUUID,

			Subjects: []string{strings.ToLower(Queue) + ".>"},
		})
		var duration time.Duration = 604800000000
		_, err1 := js.AddConsumer(Queue, &nats.ConsumerConfig{
			Durable:           NodeUUID,
			AckPolicy:         nats.AckExplicitPolicy,
			InactiveThreshold: duration,
			DeliverPolicy:     nats.DeliverAllPolicy,
			ReplayPolicy:      nats.ReplayInstantPolicy,
		})
		if err1 != nil {
			Errors.SetText(GetLangs("ms-err3") + err1.Error())
		}
		sub, errsub := js.PullSubscribe("", "", nats.BindStream(Queue))
		if errsub != nil {
			Errors.SetText(GetLangs("ms-err4") + errsub.Error())
		}

		msgs, errfetch := sub.Fetch(100)
		if errfetch != nil {
			Errors.SetText(GetLangs("ms-err5") + errfetch.Error())
			//log.Println("messages.go PullSubscribe Fetch ", errfetch)
		}
		if errfetch != nil {
			Errors.SetText(GetLangs("ms-err5") + errfetch.Error())

		}
		Errors.SetText(GetLangs("ms-err6-1") + strconv.Itoa(len(msgs)) + GetLangs("ms-err6-2"))
		if len(msgs) > 0 {
			for i := 0; i < len(msgs); i++ {
				msgs[i].Nak()
				HandleMessage(msgs[i])
			}
		}

		time.Sleep(time.Minute)
	}
}

func HandleMessage(m *nats.Msg) {
	ms := MessageStore{}
	var inmap = true // unique message id
	ejson, err := Decrypt(string(m.Data), Queuepassword)
	if err != nil {
		ejson = string(m.Data)
	}
	err1 := json.Unmarshal([]byte(ejson), &ms)
	if err1 != nil {
		ejson = "Unknown"
	}

	inmap = NodeMap("MI" + ms.MSiduuid)
	if inmap == false {
		NatsMessages = append(NatsMessages, ms)
	}

}

func NodeMap(node string) bool {
	_, x := MyMap[node]
	return x
}
