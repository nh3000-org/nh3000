package nhnats

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type nats interface {
	Send(string) bool
	Recieve() bool
	Erase() bool
}

type MessageStore struct {
	MSiduuid   string
	MSalias    string
	MShostname string
	MSipadrs   string
	MSmessage  string
	MSnodeuuid string
	MSdate     string
}

func Send(m string) bool {
	EncMessage := MessageStore{}

	//ID , err := exec.Command("uuidgen").Output()

	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = "\nNo Host Name"
		//strings.Replace(EncMessage, "#HOSTNAME", "No Host Name", -1)

	} else {
		EncMessage.MShostname = "\nHost - " + name
		//strings.Replace(EncMessage, "#HOSTNAME", name, -1)
	}
	ifas, err := net.Interfaces()
	if err == nil {

		var as []string
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				as = append(as, a)
			}
		}
		EncMessage.MShostname += "\nMac ids"
		for i, s := range as {
			EncMessage.MShostname += "\n- " + strconv.Itoa(i) + " : " + s
		}
		addrs, _ := net.InterfaceAddrs()
		EncMessage.MShostname += "\nAddress"
		for _, addr := range addrs {
			EncMessage.MShostname += "\n- " + addr.String()
		}

	}

	EncMessage.MSalias = Alias

	EncMessage.MSnodeuuid = "\nNode Id - " + NodeUUID[0:8]
	iduuid := uuid.New().String()
	EncMessage.MSiduuid = "\nMessage Id - " + iduuid[0:8]
	EncMessage.MSdate = "\nOn -" + time.Now().Format(time.UnixDate)
	//EncMessage.MSdate = "\nOn -"
	EncMessage.MSmessage = m
	//EncMessage += m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	if jsonerr != nil {
		log.Println("FormatMessage ", jsonerr)
	}
	ejson, _ := Encrypt(string(jsonmsg), Queuepassword)
	//return []byte(ejson)

	var errflag = false
	nc, err := nats.Connect(Server, nats.RootCAsMem([]byte(Caroot)), nats.ClientCertMem([]byte(Clientcert), []byte(Clientkey)))
	if err != nil {
		fmt.Println(GetLangs("ls-err7") + err.Error())
		errflag = true
	}
	js, err := nc.JetStream()
	if err != nil {
		fmt.Println(GetLangs("ls-err7") + err.Error())
		errflag = true
	}
	if errflag == false {
		_, errp := js.Publish(strings.ToLower(Queue)+"."+NodeUUID, []byte(ejson))
		if errp != nil {
			errflag = true
		}
	}
	return errflag
}
func Receive() {

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
				handleMessage(msgs[i])
			}
		}

		time.Sleep(time.Minute)
	}
}

func handleMessage(m *nats.Msg) {
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

	inmap = nodeMap("MI" + ms.MSiduuid)
	if inmap == false {
		NatsMessages = append(NatsMessages, ms)
	}

}

func nodeMap(node string) bool {
	_, x := MyMap[node]
	return x
}

func Erase() {
	log.Println("Erasing  ")
	//msgmaxage, _ := time.ParseDuration("148h")
	msgmaxage, _ := time.ParseDuration(Msgmaxage)
	nc, err := nats.Connect(Server, nats.RootCAsMem([]byte(Caroot)), nats.ClientCertMem([]byte(Clientcert), []byte(Clientkey)))
	if err != nil {
		log.Println("NatsErase Connection ", err.Error())
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("NatsErase JetStream ", err)
	}

	NatsMessages = nil

	err1 := js.DeleteStream(Queue)
	if err != nil {
		log.Println("NatsErase DeleteStream ", err1)
	}

	js1, err1 := js.AddStream(&nats.StreamConfig{
		Name:     Queue,
		Subjects: []string{strings.ToLower(Queue) + ".>"},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
	})

	if err1 != nil {
		log.Println("NatsErase AddStream ", err1)
	}
	fmt.Printf("js1: %v\n", js1)

	ac, err1 := js.AddConsumer(Queue, &nats.ConsumerConfig{
		Durable:       MyDurable,
		AckPolicy:     nats.AckExplicitPolicy,
		DeliverPolicy: nats.DeliverAllPolicy,
		ReplayPolicy:  nats.ReplayInstantPolicy,
	})
	if err1 != nil {
		log.Println("NatsErase AddConsumer ", err1, " ", ac)
	}

	FormatMessage("Security Erase")

	nc.Close()

}
