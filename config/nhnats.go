package config

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"

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
)

type MessageStore struct {
	MSiduuid    string // message id
	MSalias     string // alias
	MShostname  string // hostname
	MSipadrs    string // ip address
	MSmacid     string // macids
	MSmessage   string // message payload
	MSnodeuuid  string // unique id
	MSdate      string // message date
	MSsubject   string // message subject
	MSos        string // device os
	MSsequence  uint64 // consumer sequence for secure delete
	MSelementid int    // order in array
}
type Natsjs struct {
	Nc  *nats.Conn
	Js  jetstream.Stream
	Con jetstream.Consumer
	Ctx context.Context
}

func NewNatsJSmessages() (*Natsjs, error) {
	var n = new(Natsjs)

	var certpool = docerts()
	//var lastseq uint64
	ctxmessage, _ := context.WithTimeout(context.Background(), 10*time.Second)
	n.Ctx = ctxmessage
	//defer cancelmessage()
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
	n.Nc = natsconnect
	jetst, jetstreamerr := jetstream.New(natsconnect)
	if jetstreamerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + jetstreamerr.Error())
		}
		log.Println(" Recieve MESSAGE jetstream " + jetstreamerr.Error())
	}
	jsstream, err := jetst.CreateOrUpdateStream(ctxmessage, jetstream.StreamConfig{
		Name:     "MESSAGES",
		Subjects: []string{"messages.*"},
	})

	if err != nil {
		log.Fatal(err)
	}
	n.Js = jsstream
	/* 	dcerror := jsstream.DeleteConsumer(ctxmessage, "messages"+NatsAlias)
	   	if dcerror != nil {
	   		log.Println("dc ", dcerror)
	   	} */
	cons, conserr := jsstream.CreateOrUpdateConsumer(ctxmessage, jetstream.ConsumerConfig{
		Durable:   "messages" + NatsAlias,
		AckPolicy: jetstream.AckNonePolicy,
		//DeliverPolicy: jetstream.DeliverByStartSequencePolicy,
	})
	if conserr != nil {
		log.Fatal(err)
	}
	n.Con = cons
	//NatsCONSUMER = cons

	return n, nil

}

func NewNatsJSdevices() (*Natsjs, error) {
	var d = new(Natsjs)
	var certpool = docerts()
	//var lastseq uint64
	ctxdevice, _ := context.WithTimeout(context.Background(), 10*time.Second)
	d.Ctx = ctxdevice
	//defer canceldevice()
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
		log.Println("Recieve DEVICES connect" + getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
	}
	d.Nc = natsconnect
	jetst, jetstreamerr := jetstream.New(natsconnect)
	if jetstreamerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + jetstreamerr.Error())
		}
		log.Println(" Recieve DEVICE jetstream " + jetstreamerr.Error())
	}
	jsstream, err := jetst.CreateOrUpdateStream(ctxdevice, jetstream.StreamConfig{
		Name:     "DEVICES",
		Subjects: []string{"devices.*"},
	})

	if err != nil {
		log.Fatal("devices ", err)
	}
	d.Js = jsstream
	/* 	dcerror := jsstream.DeleteConsumer(ctxmessage, "messages"+NatsAlias)
	   	if dcerror != nil {
	   		log.Println("dc ", dcerror)
	   	} */
	cons, conserr := jsstream.CreateOrUpdateConsumer(ctxdevice, jetstream.ConsumerConfig{
		Durable:   "devices" + NatsAlias,
		AckPolicy: jetstream.AckNonePolicy,
		//DeliverPolicy: jetstream.DeliverByStartSequencePolicy,
	})
	if conserr != nil {
		log.Fatal("devices ", err)
	}
	d.Con = cons
	return d, nil

}

var ms = MessageStore{}
var devicefound = false
var messageloop = false
var shortServerName string
var shortServerName1 string
var memoryStats runtime.MemStats
var NatsMessages = make(map[int]MessageStore)
var NatsMessagesIndex = make(map[string]bool)
var fyneFilterFound = false
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
func docerts() *tls.Config {
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
	EncMessage.MSsubject = queue
	EncMessage.MSos = runtime.GOOS
	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = getLangsNats("ms-nhn")
	} else {
		EncMessage.MShostname = name
	}

	ifas, err := net.Interfaces()
	if err != nil {
		EncMessage.MShostname += "-  " + getLangsNats("ms-carrier")
	}
	if err == nil {
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				EncMessage.MSmacid += a + ", "
			}
		}

		addrs, _ := net.InterfaceAddrs()
		for _, addr := range addrs {
			EncMessage.MSipadrs += addr.String() + ", "
		}
	}
	EncMessage.MSalias = alias
	EncMessage.MSnodeuuid = NatsNodeUUID
	msiduuid := uuid.New().String()
	EncMessage.MSiduuid = msiduuid
	EncMessage.MSdate = time.Now().Format(time.UnixDate)
	EncMessage.MSmessage = m
	jsonmsg, jsonerr := json.Marshal(EncMessage)

	if jsonerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-err8") + jsonerr.Error())
		}
		log.Println(getLangsNats("ms-err8"), jsonerr.Error())
	}
	SendMessage(queue, subject, Encrypt(string(jsonmsg), NatsQueuePassword))
	runtime.GC()
	return false
}

func SendMessage(queue string, subject string, m string) {
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
	ctx, ctxcancel := context.WithCancel(context.Background())
	_, puberr := sendjs.Publish(ctx, subject, []byte(m))
	if puberr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + puberr.Error())
		}
		log.Println("Publish Error "+puberr.Error(), " queue ", queue, " subject ", subject)
	}
	ctxcancel()
}
func DeleteConsumer(queue, subject string) {
	var certpool = docerts()
	ctxmessage, cancelmessage := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelmessage()
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
		log.Println(queue + " " + "Recieve connect" + getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
	}
	js, jetstreamerr := jetstream.New(natsconnect)
	if jetstreamerr != nil {
		log.Println(queue + " " + " Recieve MESSAGE jetstream " + jetstreamerr.Error())
	}
	jsstream, err := js.CreateOrUpdateStream(ctxmessage, jetstream.StreamConfig{
		Name:     queue,
		Subjects: []string{subject},
	})

	if err != nil {
		log.Fatal(err)
	}
	dcerror := jsstream.DeleteConsumer(ctxmessage, subject+NatsAlias)
	if dcerror != nil {
		log.Println(queue+"."+subject, " DeleteConsumer ", dcerror)
	}
}

// thread for receiving messages
func ReceiveMESSAGE(a *Natsjs) {

	for {

		select {

		default:
			for {

				msg, errsub := a.Con.Next()

				if errsub != nil {
					time.Sleep(20 * time.Second)
				}
				if errsub == nil {
					meta, _ := msg.Metadata()
					//lastseq = meta.Sequence.Consumer
					//log.Println("Stream seq " + strconv.FormatUint(meta.Sequence.Stream, 10))
					//log.Println("Consumer seq " + strconv.FormatUint(meta.Sequence.Consumer, 10))
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
							ms.MSsequence = meta.Sequence.Stream
							ms.MSelementid = len(NatsMessages)
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

			}

		}

	}

}

// secure delete messages
func DeleteNatsMessage(queue string, seq uint64) {
	a, aerr := NewNatsJSmessages()
	//fmt.Printf("%+v\n", a)
	if aerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + aerr.Error())
		}
		log.Println(" Delete Message  jetstream " + aerr.Error())
	}
	//fmt.Fprintln(" Delete Message  jetstream %v " ,a)
	errdelete := a.Js.SecureDeleteMsg(a.Ctx, seq)

	if errdelete != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + errdelete.Error())
		}
		log.Println(" Delete Message  jetstream " + errdelete.Error())

	}
}

func CheckDEVICE(a *Natsjs, alias string) {
	b, _ := NewNatsJSdevices()
	consdevice, errdevice := b.Js.CreateOrUpdateConsumer(a.Ctx, jetstream.ConsumerConfig{
		Durable:       "devices" + alias,
		AckPolicy:     jetstream.AckNonePolicy,
		FilterSubject: "devices." + alias,
		//DeliverPolicy: jetstream.DeliverByStartSequencePolicy,
	})
	if errdevice != nil {
		log.Println("DEVICE CheckDEVICE consumer", errdevice)
	}
	messageloop = true
	for messageloop {
		//_, errsub := cons.Consume(func(msg jetstream.Msg) {
		msgdevice, errsubdevice := consdevice.Next()
		//log.Println("Errsub ", errsub)
		log.Println("DEVICE Receive ")
		if errsubdevice == nil {
			runtime.GC()
			runtime.ReadMemStats(&memoryStats)

			msgdevice.Nak()
			ms = MessageStore{}
			err1 := json.Unmarshal([]byte(string(Decrypt(string(msgdevice.Data()), NatsQueuePassword))), &ms)
			if err1 != nil {
				log.Println("DEVICE ReceiveMESSAGE Un Marhal", err1)
			}
			if ms.MSalias == alias {
				devicefound = true
				messageloop = false
			}

		}
		if errsubdevice != nil {
			messageloop = false
			continue
		}

	}
	if !devicefound {
		Send("DEVICES", "devices."+alias, "Add", alias)
	}

	dcerror := b.Js.DeleteConsumer(b.Ctx, "devices"+alias)
	if dcerror != nil {
		log.Println("DEVICE Delete Consumer ", dcerror)
	}

}

// }
func SetupDetails(queue string, age string) {
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
func NatsSetup() {

	SetupDetails("MESSAGES", "2048000h")
	SetupDetails("EVENTS", "2048000h")
	SetupDetails("COMMANDS", "2048000h")
	SetupDetails("DEVICES", "20480000h")
	SetupDetails("AUTHORIZATIONS", "20480000h")
	NatsMessages = nil
}
