package config

import (
	//"context"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"

	//	"fmt"
	"runtime"

	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	//"github.com/nats-io/nats.go/jetstream"
)

type MessageStore struct {
	MSiduuid   string // message id
	MSalias    string // alias
	MShostname string // hostname
	MSipadrs   string // ip address
	MSmessage  string // message payload
	MSnodeuuid string // unique id
	MSdate     string // message date
	MSsubject  string // message subject
	MSos       string // device os
}

var ms = MessageStore{}
var devicefound = false
var shortServerName string
var shortServerName1 string
var memoryStats runtime.MemStats
var NatsMessages = make(map[int]MessageStore)
var NatsMessagesIndex = make(map[string]bool)

//var pinginterval time.Duration = (30 * time.Minute)

var fyneFilterFound = false

var QuitReceive = make(chan bool)

var MessageToSend string

var myNatsLang = "eng"

// subjects - natsoperator: messages peer: 4 hours
//            natsevents: events: 96 hours
//            natscommands: commands: 96 hours
//            natsdevices: devices: forever
// MESSAGES.events - results from log, monitoring (display payload)
// MESSAGES.commands - commands sent to agent (device spacifig payload command and control)
// MESSAGES.devices - device info by node id (display payload)
//
// subjects define what is to be done with the payload

// eng esp cmn hin
var myLangsNats = map[string]string{
	"eng-fl-ll":       "NATS Language to Use eng or esp",
	"eng-ms-err2":     "NATS No Connection ",
	"spa-ms-err2":     "NATS sin Conexión ",
	"hin-ms-err2":     "NATS कोई कनेक्शन नहीं ",
	"eng-ms-carrier":  "Carrier",
	"spa-ms-carrier":  "Transportador",
	"वाहक-ms-carrier": "Carrier",
	"eng-ms-nhn":      "No Host Name ",
	"spa-ms-nhn":      "Sin Nombre de Host ",
	"hin-ms-nhn":      "कोई होस्ट नाम नहीं ",
	"eng-ms-hn":       "Host ",
	"spa-ms-hn":       "Nombre de Host ",
	"hin-ms-hn":       "मेज़बान ",
	"eng-ms-mi":       "Mac IDS",
	"spa-ms-mi":       "ID de Mac",
	"hin-ms-mi":       "मैक आईडीएस",
	"eng-ms-ad":       "Address",
	"spa-ms-ad":       "Direccion",
	"hin-ms-ad":       "पता",
	"eng-ms-ni":       "Node Id - ",
	"spa-ms-ni":       "ID de Nodo - ",
	"hin-ms-ni":       "नोड आईडी - ",
	"eng-ms-msg":      "Message Id - ",
	"spa-ms-msg":      "ID de Mensaje - ",
	"hin-ms-msg":      "संदेश आईडी - ",
	"eng-ms-on":       "On - ",
	"spa-ms-on":       "En - ",
	"hin-ms-on":       "पर - ",
	"eng-ms-err6-1":   "Recieved ",
	"spa-ms-err6-1":   "Recibida ",
	"hin-ms-err6-1":   "प्राप्त ",
	"eng-ms-err6-2":   " Messages ",
	"spa-ms-err6-2":   " Mensajes ",
	"hin-ms-err6-2":   " संदेशों ",
	"eng-ms-err6-3":   " Logs",
	"spa-ms-err6-3":   " Registros",
	"hin-ms-err6-3":   " लॉग्स",
	"eng-ms-err7":     " NATS Server Missing",
	"spa-ms-err7":     " Falta el servidor NATS",
	"hin-ms-err7":     " NATS सर्वर गायब है",

	"eng-ms-err8": " JSON Marshall",
	"spa-ms-err8": " Mariscal JSON",
	"hin-ms-err8": " JSON मार्शल",

	"eng-ms-con": "Connected",
	"spa-ms-con": "Conectada",
	"hin-ms-con": "जुड़े हुए",
	"eng-ms-dis": "Disconnected",
	"spa-ms-dis": "Desconectada",
	"hin-ms-dis": "डिस्कनेक्ट किया गया",

	"eng-ms-snd": "Send ",
	"spa-ms-snd": "Enviar ",
	"hin-ms-snd": "भेजना ",

	"eng-ms-mde": "Message Decode Error ",
	"spa-ms-mde": "Error de Decodificación de Mensaje ",
	"hin-ms-mde": "संदेश डिकोड त्रुटि ",

	"eng-ms-root": "nhnats.go docerts() rootCAs Error ",
	"spa-ms-root": "Error de CA Raíz de nhnats.go docerts() ",
	"hin-ms-root": "nhnats.go docerts() rootCAs त्रुटि ",

	"eng-ms-client": "nhnats.go docerts() client cert Error",
	"spa-ms-client": "Error de Certificado de Cliente de nhnats.go docerts()",
	"hin-ms-client": "nhnats.go docerts() क्लाइंट प्रमाणपत्र त्रुटि",

	"eng-ms-sece": "Security Erase ",
	"spa-ms-sece": "Borrado de Seguridad ",
	"hin-ms-sece": "सुरक्षा मिटाएँ ",
	"eng-ms-nnm":  "No New Messages On Server ",
	"spa-ms-nnm":  "No hay Mensajes Nuevos en el Servidor ",
	"hin-ms-nnm":  "सर्वर पर कोई नया संदेश नहीं ",
}

// return translation strings
func getLangsNats(mystring string) string {
	value, err := myLangsNats[myNatsLang+"-"+mystring]
	if !err {
		return myNatsLang + "-" + mystring
	}
	return value
}

// delete msg by uuid
func DEPRECATEDDeleteNatsMsgByUUID(iduuid string) {
	for i, v := range NatsMessages {
		if iduuid == v.MSiduuid {
			delete(NatsMessages, i)
		}
	}
}

/*
// check msg by uuid

	func CheckNatsMsgByUUID(iduuid string) bool {
		for _, v := range NatsMessages {
			//log.Println("id ", iduuid, "v.id ", v.MSiduuid)
			if iduuid == v.MSiduuid {
				return true
			}
		}
		return false
	}
*/
func CheckDevice(alias string) bool {

	var certpool = docerts()
	natsopts := nats.Options{
		Name:           NatsAlias,
		Url:            NatsServer,
		Verbose:        true,
		TLSConfig:      certpool,
		AllowReconnect: true,
		MaxReconnect:   100,
		Timeout:        30 * time.Second,
		User:           NatsUser,
		Password:       NatsUserPassword,
	}
	natsconnect, connecterr := natsopts.Connect()
	if connecterr != nil {
		log.Println("Connect" + getLangsNats("ms-snd") + getLangsNats("ms-err7") + connecterr.Error())
	}

	js, jserr := natsconnect.JetStream(nats.PublishAsyncMaxPending(256))
	if jserr != nil {
		log.Println("Connect" + getLangsNats("ms-snd") + getLangsNats("ms-err7") + connecterr.Error())
	}
	js.AddStream(&nats.StreamConfig{
		Name:     "DEVICES" + NatsNodeUUID,
		Subjects: []string{"devices"},
	})
	log.Println("Subscribe ", "devices*")
	sub, errsub := js.Subscribe("devices*", func(msg *nats.Msg) {
		log.Println("Looking for ", "devices."+alias)
		meta, metaerr := msg.Metadata()
		log.Println("in", meta.Stream)
		if metaerr != nil {
			log.Println("CheckDevice Found ", metaerr.Error())

		} else {
			log.Println("CheckDevice NotFound ", meta.NumDelivered)
			Send("DEVICES", "devices"+alias, "Device Added", alias)
		}

	}, nats.OrderedConsumer())
	if errsub != nil {
		log.Println(sub, errsub.Error())
		log.Println("CheckDevice errsub NotFound ")
		Send("DEVICES", "devices."+alias, "Device Added", alias)

	}
	sub.Drain()
	return false

}
func docerts() *tls.Config {

	//if !tlsDone {
	RootCAs, _ := x509.SystemCertPool()
	if RootCAs == nil {
		RootCAs = x509.NewCertPool()
	}
	ok := RootCAs.AppendCertsFromPEM([]byte(NatsCaroot))
	if !ok {
		log.Println(getLangsNats("ms-root"))
	}
	Clientcert, err := tls.X509KeyPair([]byte(NatsClientcert), []byte(NatsClientkey))
	if err != nil {
		log.Println(getLangsNats("ms-client") + err.Error())
	}
	shortServerName = strings.ReplaceAll(NatsServer, "nats://", "")
	shortServerName1 = strings.ReplaceAll(shortServerName, ":4222", "")
	TLSConfig := &tls.Config{
		RootCAs:            RootCAs,
		Certificates:       []tls.Certificate{Clientcert},
		ServerName:         shortServerName1,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}

	return TLSConfig
}

// send message to nats
func Send(queue string, subject string, m string, alias string) bool {

	EncMessage := MessageStore{}
	EncMessage.MSsubject = "[" + subject + "]"
	EncMessage.MSos = runtime.GOOS
	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = getLangsNats("ms-nhn")
	} else {
		EncMessage.MShostname = getLangsNats("ms-hn") + name
	}

	ifas, err := net.Interfaces()
	if err != nil {
		EncMessage.MShostname += "-  " + getLangsNats("ms-carrier")
	}
	if err == nil {
		var as []string
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				as = append(as, a)
			}
		}
		EncMessage.MShostname += getLangsNats("ms-mi")
		for i, s := range as {
			EncMessage.MShostname += "- " + strconv.Itoa(i) + " : " + s
		}
		addrs, _ := net.InterfaceAddrs()
		EncMessage.MShostname += getLangsNats("ms-ad")
		for _, addr := range addrs {
			EncMessage.MShostname += "- " + addr.String()
		}
	}
	EncMessage.MSalias = alias
	EncMessage.MSnodeuuid = getLangsNats("ms-ni") + NatsNodeUUID
	msiduuid := uuid.New().String()
	EncMessage.MSiduuid = getLangsNats("ms-msg") + msiduuid
	//log.Println("in send ", m, " ", msiduuid)
	EncMessage.MSdate = getLangsNats("ms-on") + time.Now().Format(time.UnixDate)
	EncMessage.MSmessage = m
	jsonmsg, jsonerr := json.Marshal(EncMessage)

	if jsonerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-err8") + jsonerr.Error())
		}
		log.Println(getLangsNats("ms-err8"), jsonerr.Error())
	}
	var s = Encrypt(string(jsonmsg), NatsQueuePassword)
	//log.Println("sending 2 ", s)
	Sendjs(queue, subject, s)
	runtime.GC()
	return false
}

func Sendjs(queue string, subject string, m string) {
	var certpool = docerts()
	sendnatsopts := nats.Options{
		Name:           NatsAlias,
		Url:            NatsServer,
		Verbose:        true,
		TLSConfig:      certpool,
		AllowReconnect: true,
		MaxReconnect:   1000,
		Timeout:        2048 * time.Hour,
		User:           NatsUser,
		Password:       NatsUserPassword,
	}
	sendnatsconnect, sendconnecterr := sendnatsopts.Connect()
	if sendconnecterr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + sendconnecterr.Error())
		}
		log.Println("Connect" + getLangsNats("ms-snd") + getLangsNats("ms-err7") + sendconnecterr.Error())
	}
	sendjs, sendjetstreamerr := jetstream.New(sendnatsconnect)
	if sendjetstreamerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + sendjetstreamerr.Error())
		}
		log.Println("jetstreamerror " + sendjetstreamerr.Error())
	}
	//maxage, _ := time.ParseDuration(NatsMsgMaxAge)
	ctx, ctxcancel := context.WithCancel(context.Background())
	log.Println("sendjs publish ", subject)
	_, puberr := sendjs.Publish(ctx, subject, []byte(m))
	if puberr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + puberr.Error())
		}
		log.Println("Publish Error " + puberr.Error())

		//SendMessage = false
	}
	//log.Println("in error ", pubit.Err())
	ctxcancel()
}

// thread for receiving messages
func ReceiveMESSAGE() {
	var certpool = docerts()
	//var lastseq uint64
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	natsopts := nats.Options{
		Name:           NatsAlias,
		Url:            NatsServer,
		Verbose:        true,
		TLSConfig:      certpool,
		AllowReconnect: true,
		MaxReconnect:   1000,
		Timeout:        2048 * time.Hour,
		User:           NatsUser,
		Password:       NatsUserPassword,
	}
	natsconnect, connecterr := natsopts.Connect()
	if connecterr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
		}
		log.Println("Recieve MESSAGE connect" + getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
	}
	js, jetstreamerr := jetstream.New(natsconnect)
	if jetstreamerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + jetstreamerr.Error())
		}
		log.Println(" Recieve MESSAGE jetstream " + jetstreamerr.Error())
	}
	jsstream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name: "MESSAGES",

		Subjects: []string{"messages"},
	})

	if err != nil {
		log.Fatal(err)

		//maxage, _ := time.ParseDuration(NatsMsgMaxAge)

	}
	cons, err := jsstream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "messages",
		AckPolicy: jetstream.AckNonePolicy,
		//DeliverPolicy: jetstream.DeliverByStartSequencePolicy,
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-QuitReceive:
			return

		default:
			for {
				//_, errsub := cons.Consume(func(msg jetstream.Msg) {
				msg, errsub := cons.Next()
				log.Println("Errsub ", errsub)
				if errsub != nil {
					time.Sleep(20 * time.Second)
				}
				if errsub == nil {
					meta, _ := msg.Metadata()
					//lastseq = meta.Sequence.Consumer
					log.Println("Stream seq " + strconv.FormatUint(meta.Sequence.Stream, 10))
					log.Println("Consumer seq " + strconv.FormatUint(meta.Sequence.Consumer, 10))
					if FyneMessageWin != nil {
						runtime.GC()
						runtime.ReadMemStats(&memoryStats)
						FyneMessageWin.SetTitle("ReceiveJS " + getLangsNats("ms-nnm") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
						//yulog.Println("Fetch " + GetLangs("ms-carrier") + " " + err.Error())
					}
					msg.Nak()
					ms = MessageStore{}
					err1 := json.Unmarshal([]byte(string(Decrypt(string(msg.Data()), NatsQueuePassword))), &ms)
					if err1 != nil {
						// send decrypt
						if FyneMessageWin != nil {
							FyneMessageWin.SetTitle(getLangsNats("ms-mde"))
						}
						log.Println("ReceiveJS Un Marhal", err1)
					}
					fyneFilterFound = false
					if FyneFilter {
						if strings.Contains(ms.MSmessage, getLangsNats("ms-con")) {
							fyneFilterFound = true
						}
						if strings.Contains(ms.MSmessage, getLangsNats("ms-dis")) {
							fyneFilterFound = true
						}
					}
					if !fyneFilterFound {
						//if !CheckNatsMsgByUUID(ms.MSiduuid) {
						//log.Println("check ", ms.MSiduuid, " ", NatsMessagesIndex[ms.MSiduuid])
						if !NatsMessagesIndex[ms.MSiduuid] {
							//log.Println("adding , nats.OrderedConsumer()ms ", ms.MSiduuid)
							NatsMessages[len(NatsMessages)] = ms
							NatsMessagesIndex[ms.MSiduuid] = true
							//FyneMessageList.Refresh()
						}
					}

					if FyneMessageWin != nil {
						runtime.GC()
						runtime.ReadMemStats(&memoryStats)
						FyneMessageWin.SetTitle(getLangsNats("ms-err6-1") + strconv.Itoa(len(NatsMessages)) + getLangsNats("ms-err6-2") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
					}
					//createstream.DeleteConsumer(ctx, "MESSAGESCONSUMER")
					FyneMessageList.Refresh()

				}

				//sub.Drain()
				//natsconnect.Close()

			}

		}

	}
}

// thread for receiving messages
func ReceiveNats() {
	var certpool = docerts()
	natsopts := nats.Options{
		Name:           NatsAlias,
		Url:            NatsServer,
		Verbose:        true,
		TLSConfig:      certpool,
		AllowReconnect: true,
		MaxReconnect:   100,
		Timeout:        30 * time.Second,
		User:           NatsUser,
		Password:       NatsUserPassword,
	}
	natsconnect, connecterr := natsopts.Connect()
	if connecterr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + connecterr.Error())
		}
		log.Println("Connect" + getLangsNats("ms-snd") + getLangsNats("ms-err7") + connecterr.Error())
	}

	js, err := natsconnect.JetStream()
	if err != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-carrier") + err.Error())
		}
	}
	js.AddStream(&nats.StreamConfig{
		Name:     "MESSAGES" + NatsNodeUUID,
		Subjects: []string{strings.ToLower("MESSAGES") + ".>"},
	})

	msgmaxage, _ := time.ParseDuration(NatsMsgMaxAge)
	_, err1 := js.AddConsumer("MESSAGES", &nats.ConsumerConfig{
		Durable:           NatsNodeUUID,
		AckPolicy:         nats.AckExplicitPolicy,
		InactiveThreshold: msgmaxage,
		DeliverPolicy:     nats.DeliverAllPolicy,
		ReplayPolicy:      nats.ReplayInstantPolicy,
	})

	if err1 != nil {
		//log.Println(err1.Error())
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-carrier") + err1.Error())
		}
	}

	sub, errsub := js.PullSubscribe("", "", nats.BindStream("MESSAGES"))
	if errsub != nil {
		//log.Println(errsub.Error())
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-carrier") + errsub.Error())
		}
	}

	for {
		select {
		case <-QuitReceive:
			return

		default:
			//log.Println("in ", getMessageToSend())
			msgs, err := sub.Fetch(100)
			if err != nil {
				//log.Println(err.Error())
				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-carrier") + err.Error())
				}
			}

			if len(msgs) > 0 {
				for i := 0; i < len(msgs); i++ {
					var ms = MessageStore{}
					err1 := json.Unmarshal([]byte(string(Decrypt(string(msgs[i].Data), NatsQueuePassword))), &ms)
					//log.Println("natssuberr ", natssuberr, " msg ", ms)

					//err1 := json.Unmarshal([]byte(string(Decrypt(string(msg.Data), NatsQueuePassword))), &ms)
					if err1 != nil {
						// send decrypt
						if FyneMessageWin != nil {
							FyneMessageWin.SetTitle(getLangsNats("ms-mde"))
						}
						log.Println("ReceiveNats Unmarshal", err1)
					}
					fyneFilterFound = false
					if FyneFilter {
						if strings.Contains(ms.MSmessage, getLangsNats("ms-con")) {
							fyneFilterFound = true
						}
						if strings.Contains(ms.MSmessage, getLangsNats("ms-dis")) {
							fyneFilterFound = true
						}
					}
					if !fyneFilterFound {
						//if !CheckNatsMsgByUUID(ms.MSiduuid) {
						NatsMessages[len(NatsMessages)] = ms
						//}
					}

					if FyneMessageWin != nil {
						runtime.GC()
						runtime.ReadMemStats(&memoryStats)
						FyneMessageWin.SetTitle(getLangsNats("ms-err6-1") + strconv.Itoa(len(NatsMessages)) + getLangsNats("ms-err6-2") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
					}

					msgs[i].Nak()
					FyneMessageList.Refresh()
				}
			}

			//natssub.Drain()

		}
		//NatsMessages = nil

		time.Sleep(10 * time.Second)

	}
}

// }
func EraseDetails(queue string, age string) {
	//var NatsQueues = []string{"MESSAGES","EVENTS", "COMMANDS", "DEVICES"}
	log.Println("Erase Connect", queue, " ", age)
	nc, connecterr := nats.Connect(NatsServer, nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(docerts()))
	if connecterr != nil {
		log.Println("Erase Connect", getLangsNats("ms-erac"), connecterr.Error())
	}
	js, jserr := nc.JetStream()
	if jserr != nil {
		log.Println("Erase Jetstream Make ", getLangsNats("ms-eraj"), jserr)
	}

	NatsMessages = nil
	jspurge := js.PurgeStream(queue)
	if jspurge != nil {
		log.Println("Erase Jetstream Purge "+queue, getLangsNats("ms-dels"), jspurge)
	}
	jsdelete := js.DeleteStream(queue)
	if jsdelete != nil {
		log.Println("Erase Jetstream Delete "+queue, getLangsNats("ms-dels"), jsdelete)
	}
	msgmaxage, _ := time.ParseDuration(age)
	queuestr, queueerr := js.AddStream(&nats.StreamConfig{
		Name:     queue,
		Subjects: []string{strings.ToLower(queue)},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
	})
	if queueerr != nil {
		log.Println(queue+" Addstream ", getLangsNats("ms-adds"), queueerr)
	}
	fmt.Printf(queue+": %v\n", queuestr)
	//Send(queue, strings.ToLower(queue), getLangsNats("ms-sece"), NatsAlias+":" +NatsNodeUUID+" created subject: " + queue)
	nc.Close()
}

// security erase jetstream data
func Erase() {

	EraseDetails("MESSAGES", "4h")
	EraseDetails("EVENTS", "96h")
	EraseDetails("COMMANDS", "204800h")
	EraseDetails("DEVICES", "2048000h")

}
