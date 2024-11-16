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

var NATSSTREAM = "NATS"

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
	NatsConnect *nats.Conn
	Js          jetstream.Stream
	Jetstream   jetstream.JetStream
	//Con         jetstream.Consumer
	Ctx    context.Context
	Ctxcan context.CancelFunc
}

func NewNatsJS() (*Natsjs, error) {
	var d = new(Natsjs)
	ctxdevice, ctxcan := context.WithTimeout(context.Background(), 2048*time.Hour)
	d.Ctxcan = ctxcan
	d.Ctx = ctxdevice
	natsopts := nats.Options{
		//Name:           "OPTS-" + alias,
		Url:            NatsServer,
		Verbose:        true,
		TLSConfig:      docerts(),
		AllowReconnect: true,
		MaxReconnect:   -1,
		ReconnectWait:  2,
		PingInterval:   20 * time.Second,
		Timeout:        20480 * time.Hour,
		User:           NatsUser,
		Password:       NatsUserPassword,
	}
	natsconnect, connecterr := natsopts.Connect()
	if connecterr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
		}
		log.Println("NewNatsJS  connect" + getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
	}
	d.NatsConnect = natsconnect

	jetstream, jetstreamerr := jetstream.New(natsconnect)
	if jetstreamerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + jetstreamerr.Error())
		}
		log.Println("NewNatsJS jetstreamnew ", getLangsNats("ms-eraj"), jetstreamerr)
	}
	d.Jetstream = jetstream
	js, jserr := jetstream.Stream(ctxdevice, "NATS")
	if jserr != nil {
		log.Println("NewNatsJS test ", getLangsNats("ms-eraj"), jserr)

	}
	d.Js = js
	return d, nil
}

var ms = MessageStore{}

// var devicefound = false
var messageloopauth = true
var shortServerName string
var shortServerName1 string
var memoryStats runtime.MemStats
var NatsMessages = make(map[int]MessageStore)
var NatsMessagesIndex = make(map[string]bool)
var NatsMessagesDevice = make(map[int]MessageStore)
var NatsMessagesIndexDevice = make(map[string]bool)
var fyneFilterFound = false
var MessageToSend string
var myNatsLang = "eng"

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
	"eng-ms-eraj":     "Erase JetStream ",
	"spa-ms-eraj":     "Borrar JetStream ",
	"hin-ms-eraj":     "जेटस्ट्रीम मिटाएँ ",

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
func Send(subject, m, alias string) bool {

	EncMessage := MessageStore{}
	EncMessage.MSsubject = subject
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

	snd, _ := NewNatsJS()
	snd.Jetstream.Publish(snd.Ctx, subject, []byte(Encrypt(string(jsonmsg), NatsQueuePassword)))

	snd.Ctxcan()

	//SendMessage(subject, Encrypt(string(jsonmsg), NatsQueuePassword), alias)
	runtime.GC()
	return false
}
func SetupNATS() {
	// queue = NATS
	// subjects = TYPE.alias
	// TYPE = devices, events, auth, messages
	natsoptsmissing := nats.Options{
		Name:           "SN-" + NatsAlias,
		Url:            NatsServer,
		Verbose:        true,
		TLSConfig:      docerts(),
		AllowReconnect: true,
		MaxReconnect:   -1,
		ReconnectWait:  2,
		PingInterval:   3 * time.Second,
		Timeout:        5 * time.Second,
		User:           NatsUser,
		Password:       NatsUserPassword,
	}
	natsconnectmissing, connecterrmissing := natsoptsmissing.Connect()
	if connecterrmissing != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterrmissing.Error())
		}
		log.Println("SetupNATS connect" + getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterrmissing.Error())
	}

	jsmissingctx, jsmissingerr := natsconnectmissing.JetStream()
	if jsmissingerr != nil {
		log.Println("SetupNATS JetStream ", getLangsNats("ms-eraj"), jsmissingerr)

	}
	_, streammissing := jsmissingctx.StreamInfo("NATS")
	if streammissing != nil {
		_, createerr := jsmissingctx.AddStream(&nats.StreamConfig{
			Name:      "NATS",
			Subjects:  []string{"messages.*", "events.*", "authorizations.*", "devices.*"},
			Storage:   nats.FileStorage,
			MaxAge:    204800 * time.Hour,
			FirstSeq:  1,
			Retention: nats.LimitsPolicy,
		})
		if createerr != nil {
			log.Println("SetupNATS streammissing ", getLangsNats("ms-eraj"), streammissing)
		}
	}
}

// thread for receiving messages
var startseqdev uint64
var startseqmsg uint64

func ReceiveMESSAGE() {
	//log.Println("RECIEVEMESSAGE")
	NatsReceivingMessages = true
	startseqmsg = 1

	a, aerr := NewNatsJS()
	if aerr != nil {
		log.Println("ReceiveMessage loop", aerr)
	}

	for {

		consumer, conserr := a.Js.CreateConsumer(a.Ctx, jetstream.ConsumerConfig{
			Name: "RcvMsg-" + NatsAlias,
			//Durable:           subject + alias,
			AckPolicy:         jetstream.AckExplicitPolicy,
			DeliverPolicy:     jetstream.DeliverByStartSequencePolicy,
			InactiveThreshold: 5 * time.Second,
			FilterSubject:     "messages.*",
			ReplayPolicy:      jetstream.ReplayInstantPolicy,
			OptStartSeq:       startseqmsg,
		})
		if conserr != nil {
			log.Panicln("MESSAGE Consumer", conserr)
		}
		msg, errsub := consumer.Next()
		if MsgCancel {
			a.Js.DeleteConsumer(a.Ctx, "RcvMsg-"+NatsAlias)
			a.Ctxcan()
			runtime.GC()
			return
		}
		if errsub == nil {
			meta, _ := msg.Metadata()
			//lastseq = meta.Sequence.Consumer
			//log.Println("RecieveMESSAGE seq " + strconv.FormatUint(meta.Sequence.Stream, 10))
			//log.Println("Consumer seq " + strconv.FormatUint(meta.Sequence.Consumer, 10))
			startseqmsg = meta.Sequence.Stream + 1
			if FyneMessageWin != nil {
				runtime.GC()
				runtime.ReadMemStats(&memoryStats)
				FyneMessageWin.SetTitle("RecieveMESSAGE Received " + getLangsNats("ms-nnm") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
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
				log.Println("ReceiveMESSAGE Un Marhal", err1)
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
					FyneMessageList.Refresh()
				}
			}
			if FyneDeviceWin != nil {
				runtime.GC()
				runtime.ReadMemStats(&memoryStats)
				FyneMessageWin.SetTitle(getLangsNats("ms-err6-1") + strconv.Itoa(len(NatsMessages)) + getLangsNats("ms-err6-2") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
			}
			FyneMessageList.Refresh()
		}

		if errsub == nil {
			a.Js.DeleteConsumer(a.Ctx, "RcvMsg-"+NatsAlias)
			runtime.GC()
		}
	}
}

// thread for receiving messages
func ReceiveDEVICE(alias string) {
	log.Println("CHECKDEVICE")
	devchk, _ := NewNatsJS()
	startseqdev = 1
	consumedevice, conserr := devchk.Js.CreateOrUpdateConsumer(devchk.Ctx, jetstream.ConsumerConfig{
		Name: NatsAlias + "-" + NatsNodeUUID,
		//Durable:           subject + alias,
		AckPolicy:         jetstream.AckExplicitPolicy,
		DeliverPolicy:     jetstream.DeliverByStartSequencePolicy,
		InactiveThreshold: 5 * time.Second,
		//FilterSubject:     "devices.*",
		FilterSubjects: []string{"devices.*", "authorizations.*"},
		OptStartSeq:    startseqdev,
	})
	if conserr != nil {
		log.Panicln("CheckDEVICE Consumer", conserr)
	}
	//	for {
	msgdevice, errsubdevice := consumedevice.Next()

	if errsubdevice == nil {
		runtime.GC()
		runtime.ReadMemStats(&memoryStats)

		msgdevice.Nak()
		ms = MessageStore{}
		err1 := json.Unmarshal([]byte(string(Decrypt(string(msgdevice.Data()), NatsQueuePassword))), &ms)
		if err1 != nil {
			log.Println("CheckDEVICE Un Marhal", err1)
		}
		if ms.MSalias == alias {
			devicefound = true
		}
	}
	if errsubdevice != nil {
		log.Println("CheckDEVICE exiting", errsubdevice)
		Send("devices."+alias, "Add", alias)
	}
	devchk.Ctxcan()

	log.Println("RECIEVEDEVICE")
	startseqdev = 1
	rcvdev, rcvdeverr := NewNatsJS()
	if rcvdeverr != nil {
		log.Println("ReceiveDevice aerr", rcvdeverr)
	}

	for {
		rdconsumer, rdconserr := rcvdev.Js.CreateConsumer(rcvdev.Ctx, jetstream.ConsumerConfig{
			Name: "RcvDEVICE-" + alias,
			//Durable:           subject + alias,
			AckPolicy:         jetstream.AckExplicitPolicy,
			DeliverPolicy:     jetstream.DeliverByStartSequencePolicy,
			InactiveThreshold: 1 * time.Second,
			//FilterSubject:     "devices.>",
			FilterSubjects: []string{"devices.*", "authorizations.*"},
			ReplayPolicy:   jetstream.ReplayInstantPolicy,
			OptStartSeq:    startseqdev,
		})
		if rdconserr != nil {
			log.Panicln("ReceiveDEVICE Consumer", rdconserr)
		}
		msgdev, errsubdev := rdconsumer.Next()
		if MsgCancel {
			dcerror := rcvdev.Jetstream.DeleteConsumer(rcvdev.Ctx, "DEVICES", "RcvDEVICE-"+alias)
			if dcerror != nil {
				log.Println("RecieveDEVICE Consumer not found:", dcerror)
			}
			rcvdev.Ctxcan()
			return
		}
		if errsubdev != nil {
			//log.Println("ReceiveDEVICE errsub", errsubdev)
			delerr := rcvdev.Js.DeleteConsumer(rcvdev.Ctx, "RcvDEVICE-"+alias)
			if delerr != nil {
				log.Println("ReceiveDEVICE delerr", delerr)
			}
			runtime.GC()
		}
		if errsubdev == nil {
			meta, merr := msgdev.Metadata()
			if merr != nil {
				log.Println("RecieveDEVICE meta ", merr)
			}
			//lastseq = meta.Sequence.Consumer
			log.Println("RecieveDEVICE seq " + strconv.FormatUint(meta.Sequence.Stream, 10))
			//log.Println("Consumer seq " + strconv.FormatUint(meta.Sequence.Consumer, 10))
			startseqdev = meta.Sequence.Stream + 1
			if FyneDeviceWin != nil {
				runtime.GC()
				runtime.ReadMemStats(&memoryStats)
				FyneDeviceWin.SetTitle(" Received " + getLangsNats("ms-nnm") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
				//yulog.Println("Fetch " + GetLangs("ms-carrier") + " " + err.Error())
			}
			msgdev.Nak()
			ms = MessageStore{}
			err1 := json.Unmarshal([]byte(string(Decrypt(string(msgdev.Data()), NatsQueuePassword))), &ms)
			if err1 != nil {
				// send decrypt
				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-mde"))
				}
				log.Println("ReceiveDEVICE Un Marhal", err1)
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
					ms.MSelementid = len(NatsMessagesDevice)
					NatsMessagesDevice[len(NatsMessagesDevice)] = ms

					NatsMessagesIndexDevice[ms.MSiduuid] = true
					FyneDeviceList.Refresh()
				}
			}

			if FyneDeviceWin != nil {
				runtime.GC()
				runtime.ReadMemStats(&memoryStats)
				FyneDeviceWin.SetTitle(getLangsNats("ms-err6-1") + strconv.Itoa(len(NatsMessages)) + getLangsNats("ms-err6-2") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
			}
			delerr := rcvdev.Js.DeleteConsumer(rcvdev.Ctx, "RcvDEVICE-"+alias)
			if delerr != nil {
				log.Println("ReceiveDEVICE delerr", delerr)
			}
			runtime.GC()
			FyneDeviceList.Refresh()

		}

	}

}

// secure delete messages
func DeleteNatsMessage(seq uint64) {
	a, aerr := NewNatsJS()
	//fmt.Printf("%+v\n", a)
	if aerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + aerr.Error())
		}
		log.Println("DeleteNatsMessage " + aerr.Error())
	}
	//fmt.Fprintln(" Delete Message  jetstream %v " ,a)
	errdelete := a.Js.SecureDeleteMsg(a.Ctx, seq)

	if errdelete != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + errdelete.Error())
		}
		log.Println("DeleteNatsMessage " + errdelete.Error())

	}
	a.Ctxcan()
}

var devicefound = false

func CheckDEVICE(alias string) bool {
	if devicefound {
		return true
	}
	log.Println("CHECKDEVICE")
	runtime.GC()
	runtime.ReadMemStats(&memoryStats)

	log.Println("DEVICE Memory Start: " + strconv.FormatUint(memoryStats.Alloc/1024, 10) + " K")
	cd, cderr := NewNatsJS()
	log.Println("CHECKDEVICE", cderr)
	consumedevice, conserr := cd.Js.CreateOrUpdateConsumer(cd.Ctx, jetstream.ConsumerConfig{
		Name: "CHECKDEVICE-" + alias,
		//Durable:           subject + alias,
		AckPolicy: jetstream.AckExplicitPolicy,
		//DeliverPolicy:     jetstream.DeliverLastPerSubjectPolicy,
		InactiveThreshold: 2 * time.Second,
		FilterSubject:     "devices." + alias,
		//OptStartSeq:       start,
	})
	if conserr != nil {
		log.Panicln("CheckDEVICE Consumer", conserr)
	}
	//messageloopdevice = true
	//for messageloopdevice {
	log.Println("CheckDEVICE filter", "devices."+alias)
	msgdevice, errsubdevice := consumedevice.Next()
	if errsubdevice == nil {
		msgdevice.Nak()
		ms = MessageStore{}
		err1 := json.Unmarshal([]byte(string(Decrypt(string(msgdevice.Data()), NatsQueuePassword))), &ms)
		if err1 != nil {
			log.Println("nhnats.go Receive Un Marhal", err1)
		}
		log.Println("nhnats.go Receive ", ms.MSalias, alias)
		if ms.MSalias == alias {
			devicefound = true
		}

	}

	if errsubdevice != nil {
		log.Println("CheckDEVICE exiting", errsubdevice)
		//messageloopdevice = false
		Send("devices."+alias, "Add", alias)

	}

	/* 	if !devicefound {

		Send(NatsUser, NatsUserPassword, "DEVICES", "devices."+alias, "Add", alias)
	} */

	cd.Ctxcan()
	runtime.GC()
	runtime.ReadMemStats(&memoryStats)

	log.Println("DEVICE Memory End: " + strconv.FormatUint(memoryStats.Alloc/1024, 10) + " K")

	return devicefound
}

var deviceauthorized = false
var deviceloopauth = true

func CheckAUTHORIZATIONS(alias string) bool {
	log.Println("CHECKAUTHORIZATIONS", alias)

	if deviceauthorized {
		return true
	}

	runtime.GC()
	runtime.ReadMemStats(&memoryStats)

	log.Println("AUTH Memory Start: " + strconv.FormatUint(memoryStats.Alloc/1024, 10) + " K")

	cd, cderr := NewNatsJS()
	if cderr != nil {
		log.Println("CheckDEVICE test ", getLangsNats("ms-err2"), cderr)
	}
	consumeauth, conserr := cd.Js.CreateOrUpdateConsumer(cd.Ctx, jetstream.ConsumerConfig{
		Name: "CheckAUTH" + "-" + alias,
		//Durable:           subject + alias,
		AckPolicy:         jetstream.AckExplicitPolicy,
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		InactiveThreshold: 5 * time.Second,
		FilterSubject:     "authorizations." + alias,
		//OptStartSeq:       start,
	})
	if conserr != nil {
		log.Panicln("CheckAUTHORIZATIONS consumer", conserr)
	}
	for deviceloopauth {

		runtime.GC()
		runtime.ReadMemStats(&memoryStats)

		log.Println("AUTHORIZATIONS Memory Start: " + strconv.FormatUint(memoryStats.Alloc/1024, 10) + " K")

		msgauthorizations, errsubauthorizations := consumeauth.Next()
		log.Println("CheckAUTHORIZATIONS next", "authorizations."+alias, " error ", errsubauthorizations)
		if errsubauthorizations == nil {
			msgauthorizations.Nak()
			ms = MessageStore{}
			err1 := json.Unmarshal([]byte(string(Decrypt(string(msgauthorizations.Data()), NatsQueuePassword))), &ms)
			if err1 != nil {
				log.Println("nhnats.go AUTHORIZATIONS Receive Un Marhal", err1)
			}
			log.Println("CheckAUTHORIZATIONS ", ms.MSalias, alias)
			deviceauthorized = true
			deviceloopauth = false

		}
		if errsubauthorizations != nil {
			log.Println("CheckAUTHORIZATIONS error ", errsubauthorizations, alias)
			deviceloopauth = true

			cd.Js.DeleteConsumer(cd.Ctx, "authorizations"+alias)

			time.Sleep(120 * time.Second)
		}

	}
	cd.Ctxcan()
	return deviceauthorized
}

// }
func SetupDetails(queue string, age string) {

	log.Println("nhnats.go Erase Connect", queue, " ", age)
	nc, connecterr := nats.Connect(NatsServer, nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(docerts()))
	if connecterr != nil {
		log.Println("nhnats.go Erase Connect", getLangsNats("ms-erac"), connecterr.Error())
	}
	js, jserr := nc.JetStream()
	if jserr != nil {
		log.Println("nhnats.go Erase Jetstream Make ", getLangsNats("ms-eraj"), jserr)
	}

	jspurge := js.PurgeStream(queue)
	if jspurge != nil {
		log.Println("nhnats.go Erase Jetstream Purge "+queue, getLangsNats("ms-dels"), jspurge)
	}
	jsdelete := js.DeleteStream(queue)
	if jsdelete != nil {
		log.Println("nhnats.go Erase Jetstream Delete "+queue, getLangsNats("ms-dels"), jsdelete)
	}

	msgmaxage, ageerr := time.ParseDuration("24h")
	if ageerr != nil {
		log.Println("nhnats.go Erase Jetstream parse ", ageerr)
	}

	queuestr, queueerr := js.AddStream(&nats.StreamConfig{
		Name:     queue,
		Subjects: []string{strings.ToLower(queue)},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
		FirstSeq: 1,
	})
	if queueerr != nil {
		log.Println("nhnats.go ", queue+" Addstream ", getLangsNats("ms-adds"), queueerr)
	}
	fmt.Printf(queue+": %v\n", queuestr)
	//Send(queue, strings.ToLower(queue), getLangsNats("ms-sece"), NatsAlias+":" +NatsNodeUUID+" created subject: " + queue)
	nc.Close()
}

// security erase jetstream data
func NatsSetup() {

	/*
		 	SetupDetails("MESSAGES", "24h")
			SetupDetails("EVENTS", "24h")
			SetupDetails("COMMANDS", "24h")
			SetupDetails("DEVICES", "24h")
			SetupDetails("AUTHORIZATIONS", "24h")
			NatsMessages = nil
	*/
}
