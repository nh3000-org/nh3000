package nhnats

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"

	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nh3000-org/nh3000/nhauth"
	"github.com/nh3000-org/nh3000/nhcrypt"
	"github.com/nh3000-org/nh3000/nhlang"

	//"github.com/nh3000-org/nh3000/nhpanes"
	"github.com/nh3000-org/nh3000/nhpref"
	"github.com/nh3000-org/nh3000/nhutil"
)

type MessageStore struct {
	MSiduuid   string
	MSalias    string
	MShostname string
	MSipadrs   string
	MSmessage  string
	MSnodeuuid string
	MSdate     string
}

var NatsMessages []MessageStore
var MyAckMap = make(map[string]bool)
var QuitReceive = make(chan bool)
var TLS tls.Config

func docerts() {
	var done = false
	if !done {
		RootCAs, _ := x509.SystemCertPool()
		if RootCAs == nil {
			RootCAs = x509.NewCertPool()
		}
		ok := RootCAs.AppendCertsFromPEM([]byte(nhauth.Caroot))
		if !ok {
			log.Println("nhnats.go init rootCAs")
		}
		Clientcert, err := tls.X509KeyPair([]byte(nhauth.Clientcert), []byte(nhauth.Clientkey))
		if err != nil {
			log.Println("nhnats.go init Clientcert " + err.Error())
		}
		TLSConfig := &tls.Config{
			RootCAs:            RootCAs,
			Certificates:       []tls.Certificate{Clientcert},
			ServerName:         "nats.newhorizons3000.org",
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		}
		TLS = *TLSConfig.Clone()
		done = true
	}
}

// send message to nats
func Send(m string) bool {
	docerts()
	EncMessage := MessageStore{}
	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = "\n" + nhlang.GetLangs("ms-nhn")
	} else {
		EncMessage.MShostname = "\n" + nhlang.GetLangs("ms-hn") + name
	}

	ifas, err := net.Interfaces()
	if err != nil {
		EncMessage.MShostname += "\n-  " + nhlang.GetLangs("ms-carrier")
	}
	if err == nil {
		var as []string
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				as = append(as, a)
			}
		}
		EncMessage.MShostname += "\n" + nhlang.GetLangs("ms-mi")
		for i, s := range as {
			EncMessage.MShostname += "\n- " + strconv.Itoa(i) + " : " + s
		}
		addrs, _ := net.InterfaceAddrs()
		EncMessage.MShostname += "\n" + nhlang.GetLangs("ms-ad")
		for _, addr := range addrs {
			EncMessage.MShostname += "\n- " + addr.String()
		}
	}
	EncMessage.MSalias = nhpref.Alias
	EncMessage.MSnodeuuid = "\n" + nhlang.GetLangs("ms-ni") + nhpref.NodeUUID
	iduuid := uuid.New().String()
	EncMessage.MSiduuid = "\n" + nhlang.GetLangs("ms-msg") + iduuid
	EncMessage.MSdate = "\n" + nhlang.GetLangs("ms-on") + time.Now().Format(time.UnixDate)
	EncMessage.MSmessage = m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	log.Println("FormatMessage Content", EncMessage.MSmessage)
	if jsonerr != nil {
		log.Println("FormatMessage ", jsonerr)
	}
	ejson, _ := nhcrypt.Encrypt(string(jsonmsg), nhpref.Queuepassword)
	NC, err := nats.Connect(nhpref.Server, nats.UserInfo(nhauth.User, nhauth.UserPassword), nats.Secure(&TLS))
	if err != nil {
		fmt.Println("Send " + nhlang.GetLangs("ls-err7") + err.Error())
	}
	JS, err := NC.JetStream()
	if err != nil {
		fmt.Println("Send " + nhlang.GetLangs("ls-err7") + err.Error() + <-JS.StreamNames())
	}
	_, errp := JS.Publish(strings.ToLower(nhpref.Queue)+"."+nhpref.NodeUUID, []byte(ejson))
	if errp != nil {
		return true
	}
	NC.Drain()

	return false
}

// thread for receiving messages
func Receive() {
	docerts()
	nc, err := nats.Connect(nhpref.Server, nats.UserInfo(nhauth.User, nhauth.UserPassword), nats.Secure(&TLS))
	if err != nil {
		log.Println("Receive ", nhlang.GetLangs("ms-err2"))
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("Receive JetStream ", err)
	}
	js.AddStream(&nats.StreamConfig{
		Name: nhpref.Queue + nhpref.NodeUUID,

		Subjects: []string{strings.ToLower(nhpref.Queue) + ".>"},
	})
	var duration time.Duration = 604800000000
	_, err1 := js.AddConsumer(nhpref.Queue, &nats.ConsumerConfig{
		Durable:           nhpref.NodeUUID,
		AckPolicy:         nats.AckExplicitPolicy,
		InactiveThreshold: duration,
		DeliverPolicy:     nats.DeliverAllPolicy,
		ReplayPolicy:      nats.ReplayInstantPolicy,
	})
	if err1 != nil {
		log.Println("Receive AddConsumer ", nhlang.GetLangs("ms-err3")+err1.Error())
	}
	sub, errsub := js.PullSubscribe("", "", nats.BindStream(nhpref.Queue))
	if errsub != nil {
		log.Println("Receive Pull Subscribe ", nhlang.GetLangs("ms-err4")+errsub.Error())
	}
	nhpref.ReceivingMessages = true
	for {
		select {
		case <-QuitReceive:
			return
		default:

			NatsMessages = nil

			msgs, _ := sub.Fetch(100)
			nhpref.ClearMessageDetail = true
			if len(msgs) > 0 {
				for i := 0; i < len(msgs); i++ {
					handleMessage(msgs[i])
					msgs[i].Nak()
				}
			}
			if nhutil.GetMessageWin() != nil {
				nhutil.GetMessageWin().SetTitle(nhlang.GetLangs("ms-err6-1") + strconv.Itoa(len(msgs)) + nhlang.GetLangs("ms-err6-2"))
			}
			nc.Close()
			time.Sleep(30 * time.Second)
		}
	}
}

// decrypt payload
func handleMessage(m *nats.Msg) string {
	ms := MessageStore{}
	ejson, err := nhcrypt.Decrypt(string(m.Data), nhpref.Queuepassword)
	if err != nil {
		ejson = string(m.Data)
	}
	err1 := json.Unmarshal([]byte(ejson), &ms)
	if err1 != nil {
		ejson = nhlang.GetLangs("ms-unk")
	}
	if nhpref.Filter {
		if strings.Contains(ms.MSmessage, nhlang.GetLangs("ls-con")) || strings.Contains(ms.MSmessage, nhlang.GetLangs("ls-dis")) {
			return ""
		}
	}
	NatsMessages = append(NatsMessages, ms)

	return ms.MSiduuid
}

// security erase jetstream data
func Erase() {
	docerts()
	log.Println(nhlang.GetLangs("ms-era"))

	nc, err := nats.Connect(nhpref.Server, nats.UserInfo(nhauth.User, nhauth.UserPassword), nats.Secure(&TLS))
	if err != nil {
		log.Println("Erase Connect", nhlang.GetLangs("ms-erac"), err.Error())
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("Erase Jetstream Make ", nhlang.GetLangs("ms-eraj"), err)
	}

	NatsMessages = nil
	err1 := js.PurgeStream(nhpref.Queue)
	if err1 != nil {
		log.Println("Erase Jetstream Purge", nhlang.GetLangs("ms-dels"), err1)
	}
	err2 := js.DeleteStream(nhpref.Queue)
	if err2 != nil {
		log.Println("Erase Jetstream Delete", nhlang.GetLangs("ms-dels"), err1)
	}
	msgmaxage, _ := time.ParseDuration(nhpref.Msgmaxage)
	js1, err3 := js.AddStream(&nats.StreamConfig{
		Name:     nhpref.Queue,
		Subjects: []string{strings.ToLower(nhpref.Queue) + ".>"},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
	})
	if err3 != nil {
		log.Println("Erase Addstream ", nhlang.GetLangs("ms-adds"), err3)
	}
	fmt.Printf("js1: %v\n", js1)

	Send(nhlang.GetLangs("ms-sece"))
	nc.Close()
}
