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
	"github.com/nh3000-org/nh3000/nhcrypt"
	"github.com/nh3000-org/nh3000/nhlang"
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

// send message to nats
func Send(m string) bool {
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
	if jsonerr != nil {
		log.Println("FormatMessage ", jsonerr)
	}
	ejson, _ := nhcrypt.Encrypt(string(jsonmsg), nhpref.Queuepassword)
	clientcert, err := tls.X509KeyPair([]byte(nhpref.Clientcert), []byte(nhpref.Clientkey))
	if err != nil {
		log.Println("nhnats.go clientcert " + err.Error())
	}
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	ok := rootCAs.AppendCertsFromPEM([]byte(nhpref.Caroot))
	if !ok {
		log.Println("nhnats.go rootCAs")
	}
	tlsConfig := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{clientcert},
	}
	var errflag = false
	nc, err := nats.Connect(nhpref.Server, nats.Secure(tlsConfig))
	if err != nil {
		fmt.Println(nhlang.GetLangs("ls-err7") + err.Error())
		errflag = true
	}
	js, err := nc.JetStream()
	if err != nil {
		fmt.Println(nhlang.GetLangs("ls-err7") + err.Error())
		errflag = true
	}
	if errflag == false {
		_, errp := js.Publish(strings.ToLower(nhpref.Queue)+"."+nhpref.NodeUUID, []byte(ejson))
		if errp != nil {
			errflag = true
		}
	}
	nc.Close()
	return errflag
}

// thread for receiving messages
func Receive() {
	nhpref.ReceivingMessages = true
	for {
		select {
		case <-QuitReceive:
			return
		default:
			//clientcert, err := tls.LoadX509KeyPair(nhpref.DataStore("cert.pem").Path(), nhpref.DataStore("key.pem").Path())
			clientcert, err := tls.X509KeyPair([]byte(nhpref.Clientcert), []byte(nhpref.Clientkey))
			if err != nil {
				log.Println("nhnats.go clientcert " + err.Error())
			}
			rootCAs, _ := x509.SystemCertPool()
			if rootCAs == nil {
				rootCAs = x509.NewCertPool()
			}
			ok := rootCAs.AppendCertsFromPEM([]byte(nhpref.Caroot))
			if !ok {
				log.Println("nhnats.go rootCAs")
			}
			NatsMessages = nil
			tlsConfig := &tls.Config{
				RootCAs:      rootCAs,
				Certificates: []tls.Certificate{clientcert},
			}
			nc, err := nats.Connect(nhpref.Server, nats.Secure(tlsConfig))
			if err != nil {
				log.Println(nhlang.GetLangs("ms-err2"))
			}
			js, _ := nc.JetStream()
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
				log.Println(nhlang.GetLangs("ms-err3") + err1.Error())
			}
			sub, errsub := js.PullSubscribe("", "", nats.BindStream(nhpref.Queue))
			if errsub != nil {
				log.Println(nhlang.GetLangs("ms-err4") + errsub.Error())
			}
			msgs, _ := sub.Fetch(100)
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
	log.Println(nhlang.GetLangs("ms-era"))
	//msgmaxage, _ := time.ParseDuration("148h")
	clientcert, err := tls.X509KeyPair([]byte(nhpref.Clientcert), []byte(nhpref.Clientkey))
	if err != nil {
		log.Println("nhnats.go clientcert " + err.Error())
	}
	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	ok := rootCAs.AppendCertsFromPEM([]byte(nhpref.Caroot))
	if !ok {
		log.Println("nhnats.go rootCAs")
	}
	tlsConfig := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{clientcert},
		//ClientAuth:   tls.RequireAndVerifyClientCert,
	}
	msgmaxage, _ := time.ParseDuration(nhpref.Msgmaxage)
	nc, err := nats.Connect(nhpref.Server, nats.Secure(tlsConfig))
	if err != nil {
		log.Println(nhlang.GetLangs("ms-erac"), err.Error())
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println(nhlang.GetLangs("ms-eraj"), err)
	}
	NatsMessages = nil
	err1 := js.DeleteStream(nhpref.Queue)
	if err != nil {
		log.Println(nhlang.GetLangs("ms-dels"), err1)
	}
	js1, err1 := js.AddStream(&nats.StreamConfig{
		Name:     nhpref.Queue,
		Subjects: []string{strings.ToLower(nhpref.Queue) + ".>"},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
	})
	if err1 != nil {
		log.Println(nhlang.GetLangs("ms-adds"), err1)
	}
	fmt.Printf("js1: %v\n", js1)

	ac, err1 := js.AddConsumer(nhpref.Queue, &nats.ConsumerConfig{
		Durable:       nhpref.MyDurable,
		AckPolicy:     nats.AckExplicitPolicy,
		DeliverPolicy: nats.DeliverAllPolicy,
		ReplayPolicy:  nats.ReplayInstantPolicy,
	})
	if err1 != nil {
		log.Println(nhlang.GetLangs("ms-addc"), err1, " ", ac)
	}
	Send(nhlang.GetLangs("ms-sece"))
	nc.Close()
}
