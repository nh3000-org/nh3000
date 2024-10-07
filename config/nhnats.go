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
	NatsConnect *nats.Conn
	Js          jetstream.Stream
	Jetstream   jetstream.JetStream
	//Con         jetstream.Consumer
	Ctx    context.Context
	Ctxcan context.CancelFunc
}

func NewNatsJS(queue, subject, alias string, start uint64) (*Natsjs, error) {
	var d = new(Natsjs)
	log.Println("NewNasJS q ", queue, " sub ", subject, " alias ", alias)
	var certpool = docerts()
	//var lastseq uint64
	ctxdevice, ctxcan := context.WithTimeout(context.Background(), 2048*time.Hour)
	d.Ctxcan = ctxcan
	d.Ctx = ctxdevice
	//defer canceldevice()
	natsopts := nats.Options{
		Name:           "OPTS-" + NatsAlias,
		Url:            NatsServer,
		Verbose:        true,
		TLSConfig:      certpool,
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
	jsctx, ncerr := natsconnect.JetStream()
	if ncerr != nil {
		log.Println("NewNatsJS jsctx ", getLangsNats("ms-eraj"), ncerr)

	}
	_, streammissing := jsctx.StreamInfo(queue)
	if streammissing != nil {
		_, createerr := jsctx.AddStream(&nats.StreamConfig{
			Name:     queue,
			Subjects: []string{strings.ToLower(queue)},
			Storage:  nats.FileStorage,
			MaxAge:   204800 * time.Hour,
			FirstSeq: 1,
		})
		if createerr != nil {
			log.Println("NewNatsJS streammissing ", getLangsNats("ms-eraj"), streammissing)
		}
	}
	jetstream, jetstreamerr := jetstream.New(natsconnect)
	if jetstreamerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + jetstreamerr.Error())
		}
		log.Println("NewNatsJS jetstreamnew ", getLangsNats("ms-eraj"), jetstreamerr)
	}
	d.Jetstream = jetstream
	js, jserr := jetstream.Stream(ctxdevice, queue)
	if jserr != nil {
		log.Println("NewNatsJS test ", getLangsNats("ms-eraj"), jserr)

	}
	d.Js = js
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
func Send(user, password, queue, subject, m, alias string) bool {

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
	SendMessage(user, password, queue, subject, Encrypt(string(jsonmsg), NatsQueuePassword))
	runtime.GC()
	return false
}
func CheckQueue(user, password, queue string) {
	var certpool = docerts()

	//defer canceldevice()
	natsoptsmissing := nats.Options{
		Name:           "OPTS-" + NatsAlias,
		Url:            NatsServer,
		Verbose:        true,
		TLSConfig:      certpool,
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
		log.Println("NewNatsJS  connect" + getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterrmissing.Error())
	}

	jsmissingctx, jsmissingerr := natsconnectmissing.JetStream()
	if jsmissingerr != nil {
		log.Println("NewNatsJS jsctx ", getLangsNats("ms-eraj"), jsmissingerr)

	}
	_, streammissing := jsmissingctx.StreamInfo(queue)
	if streammissing != nil {
		_, createerr := jsmissingctx.AddStream(&nats.StreamConfig{
			Name:      queue,
			Subjects:  []string{strings.ToLower(queue) + ".*"},
			Storage:   nats.FileStorage,
			MaxAge:    204800 * time.Hour,
			FirstSeq:  1,
			Retention: nats.LimitsPolicy,
		})
		if createerr != nil {
			log.Println("NewNatsJS streammissing ", getLangsNats("ms-eraj"), streammissing)
		}
	}
}
func SendMessage(user, password, queue, subject, m string) {

	log.Println("SendMessage", queue, subject, m)
	certpool := docerts()
	//var lastseq uint64
	_, ctxsendcancel := context.WithTimeout(context.Background(), 30*time.Second)

	//defer canceldevice()
	natsopts := nats.Options{
		Name:           "OPTS-" + NatsAlias,
		Url:            NatsServer,
		Verbose:        true,
		TLSConfig:      certpool,
		AllowReconnect: true,
		MaxReconnect:   -1,
		ReconnectWait:  2,
		PingInterval:   20 * time.Second,
		Timeout:        1 * time.Hour,
		User:           NatsUser,
		Password:       NatsUserPassword,
	}
	//natsconnect, connecterr := nats.Connect(Na)
	natsconnect, connecterr := natsopts.Connect()
	if connecterr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
		}
		log.Println("NewNatsJS  connect" + getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
	}
	// Use the JetStream context to produce and consumer messages
	// that have been persisted.
	js, err := natsconnect.JetStream(nats.PublishAsyncMaxPending(2560))
	if err != nil {
		log.Fatal(err)
	}

	js.AddStream(&nats.StreamConfig{
		Name:     queue,
		Subjects: []string{subject},
	})

	js.Publish(subject, []byte(m))

	ctxsendcancel()
	//ctxcancel()
}

// thread for receiving messages
var startseq uint64

func ReceiveMESSAGE() {
	log.Println("RECIEVEMESSAGE")
	NatsReceivingMessages = true
	startseq = 1

	a, aerr := NewNatsJS("MESSAGES", "messages", "RcvMsg-"+NatsAlias, startseq)
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
			OptStartSeq:       startseq,
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
			startseq = meta.Sequence.Stream + 1
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
			//createstream.DeleteConsumer(ctx, "MESSAGESCONSUMER")
			FyneMessageList.Refresh()

		}

		if errsub == nil {
			//log.Println("ReceiveMESSAGE errsub", errsub)

			a.Js.DeleteConsumer(a.Ctx, "RcvMsg-"+NatsAlias)
			runtime.GC()

			//time.Sleep(5 * time.Second)
		}
	}
}

// thread for receiving messages
func ReceiveDEVICE(alias string) {
	log.Println("RECIEVEDEVICE")
	startseq = 1
	a, aerr := NewNatsJS("DEVICES", "devices", "RecDevice-"+alias, 1)
	if aerr != nil {
		log.Println("ReceiveDevice ", aerr)
	}

	//if aerr != nil {
	//	log.Println("ReceiveDEVICE err ", aerr)
	//}time.Sleep(20 * time.Second)
	for {
		consumer, conserr := a.Js.CreateConsumer(a.Ctx, jetstream.ConsumerConfig{
			Name: "RecDevice-" + alias,
			//Durable:           subject + alias,
			AckPolicy:         jetstream.AckExplicitPolicy,
			DeliverPolicy:     jetstream.DeliverByStartSequencePolicy,
			InactiveThreshold: 5 * time.Second,
			ReplayPolicy:      jetstream.ReplayInstantPolicy,
			FilterSubject:     "devices.*",
			OptStartSeq:       startseq,
		})
		if conserr != nil {
			log.Panicln("ReceiveDEVICE Consumer", conserr)
		}
		msg, errsub := consumer.Next()
		if MsgCancel {
			dcerror := a.Jetstream.DeleteConsumer(a.Ctx, "DEVICES", "RecDevice-"+alias)
			if dcerror != nil {
				log.Println("RecieveDEVICE Consumer not found:", dcerror)
			}
			a.Ctxcan()
			return
		}

		if errsub == nil {
			meta, _ := msg.Metadata()
			//lastseq = meta.Sequence.Consumer
			log.Println("RecieveDEVICE seq " + strconv.FormatUint(meta.Sequence.Stream, 10))
			//log.Println("Consumer seq " + strconv.FormatUint(meta.Sequence.Consumer, 10))
			startseq = meta.Sequence.Stream + 1
			if FyneMessageWin != nil {
				runtime.GC()
				runtime.ReadMemStats(&memoryStats)
				FyneMessageWin.SetTitle("RecieveDEVICE Received " + getLangsNats("ms-nnm") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
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
					//FyneMessageList.Refresh()
				}
			}

			if FyneDeviceWin != nil {
				runtime.GC()
				runtime.ReadMemStats(&memoryStats)
				FyneMessageWin.SetTitle(getLangsNats("ms-err6-1") + strconv.Itoa(len(NatsMessages)) + getLangsNats("ms-err6-2") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
			}
			//createstream.DeleteConsumer(ctx, "MESSAGESCONSUMER")
			FyneDeviceList.Refresh()

		}

		if errsub != nil {
			a.Js.DeleteConsumer(a.Ctx, "RecDevice-"+NatsAlias)
			runtime.GC()
		}

	}

}

// secure delete messages
func DeleteNatsMessage(queue, subject string, seq uint64) {
	a, aerr := NewNatsJS(queue, subject, "DelMsg"+NatsAlias, 1)
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

func CheckDEVICE(alias string) bool {
	devicefound = false
	log.Println("CHECKDEVICE")

	//a, _ := NewNatsJS("DEVICES", "devices", "ChkDev-"+alias, 1)
	//	var d = new(Natsjs)
	//log.Println("NewNasJS CheckDevice", queue, " sub ", subject, " alias ", alias)
	var certpool = docerts()
	//var lastseq uint64
	devctx, devctxcan := context.WithTimeout(context.Background(), 10*time.Second)
	//defer devctxcan()

	natsopts := nats.Options{
		Name:           "OPTS-" + NatsAlias,
		Url:            NatsServer,
		Verbose:        true,
		TLSConfig:      certpool,
		AllowReconnect: false,
		MaxReconnect:   -1,
		ReconnectWait:  2,
		PingInterval:   5 * time.Second,
		Timeout:        1 * time.Hour,
		User:           NatsUser,
		Password:       NatsUserPassword,
	}
	natsconnect, connecterr := natsopts.Connect()
	if connecterr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
		}
		log.Println("CheckDEVICE natsconnect" + getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
	}
	// check if stream exists, if not create it
	devcheckstream, devcheckstreamerr := natsconnect.JetStream()
	if devcheckstreamerr != nil {
		log.Println("CheckDEVICE devjsctxerr ", getLangsNats("ms-eraj"), devcheckstream)

	}
	_, devicestreammissing := devcheckstream.StreamInfo("DEVICES")
	if devicestreammissing != nil {
		_, createerr := devcheckstream.AddStream(&nats.StreamConfig{
			Name:     "DEVICES",
			Subjects: []string{"devices"},
			Storage:  nats.FileStorage,
			MaxAge:   204800 * time.Hour,
			FirstSeq: 1,
		})
		if createerr != nil {
			log.Println("CheckDEVICE stream create", getLangsNats("ms-eraj"), createerr)
		}
	}

	// receive the device info
	devstream, jetstreamerr := jetstream.New(natsconnect)
	if jetstreamerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + jetstreamerr.Error())
		}
		log.Println("CheckDEVICE jetstreamerr ", getLangsNats("ms-eraj"), jetstreamerr)
	}

	devjs, devjserr := devstream.Stream(devctx, "DEVICES")
	if devjserr != nil {
		log.Println("CheckDEVICE test ", getLangsNats("ms-eraj"), devjserr)

	}

	consumedevice, conserr := devjs.CreateOrUpdateConsumer(devctx, jetstream.ConsumerConfig{
		Name: NatsAlias + "-" + NatsNodeUUID,
		//Durable:           subject + alias,
		AckPolicy:         jetstream.AckNonePolicy,
		DeliverPolicy:     jetstream.DeliverAllPolicy,
		InactiveThreshold: 5 * time.Second,
		//OptStartSeq:       start,
	})
	if conserr != nil {
		log.Panicln("CheckDEVICE Consumer", conserr)
	}
	messageloop = true
	for messageloop {

		msgdevice, errsubdevice := consumedevice.Next()

		if errsubdevice == nil {
			runtime.GC()
			runtime.ReadMemStats(&memoryStats)

			msgdevice.Nak()
			ms = MessageStore{}
			err1 := json.Unmarshal([]byte(string(Decrypt(string(msgdevice.Data()), NatsQueuePassword))), &ms)
			if err1 != nil {
				log.Println("nhnats.go Receive Un Marhal", err1)
			}
			if ms.MSalias == alias {
				devicefound = true
				messageloop = false
			}

		}
		if errsubdevice != nil {
			log.Println("CheckDEVICE exiting", errsubdevice)
			messageloop = false
			continue
		}

	}
	if !devicefound {
		Send(NatsUser, NatsUserPassword, "DEVICES", "devices."+alias, "Add", alias)
	}
	/* 	//dcerror := jsstream.DeleteConsumer(ctxmessage, subject+NatsAlias)
	   	dcerror := devctx.DeleteConsumer(a.Ctx, NatsAlias+"-"+NatsNodeUUID)
	   	if dcerror != nil {
	   		log.Println("nhnats.go CheckDevice Consumer not found:", dcerror)
	   	} */
	devctxcan()
	return devicefound
}

var deviceauthorized bool

func CheckAUTHORIZATIONS(a *Natsjs, alias string) bool {
	log.Println("CHECKAUTHORIZATIONS")
	time.Sleep(20 * time.Second)
	deviceauthorized = false
	b, _ := NewNatsJS("AUTHORIZATIONS", "authorizations", "ChkAuth-"+alias, 1)
	consauthorizations, errdevice := b.Js.CreateOrUpdateConsumer(a.Ctx, jetstream.ConsumerConfig{
		Name:          "authorizations" + alias,
		AckPolicy:     jetstream.AckNonePolicy,
		FilterSubject: "authorizations.*" + alias,
		//DeliverPolicy: jetstream.DeliverByStartSequencePolicy,
	})
	if errdevice != nil {
		log.Println("nhnats.go AUTHORIZATIONS CheckAUTHORIZATIONS consumer", errdevice)
	}
	messageloop = true
	for messageloop {

		msgauthorizations, errsubauthorizations := consauthorizations.Next()

		if errsubauthorizations == nil {
			runtime.GC()
			runtime.ReadMemStats(&memoryStats)

			msgauthorizations.Nak()
			ms = MessageStore{}
			err1 := json.Unmarshal([]byte(string(Decrypt(string(msgauthorizations.Data()), NatsQueuePassword))), &ms)
			if err1 != nil {
				log.Println("nhnats.go AUTHORIZATIONS Receive Un Marhal", err1)
			}
			if ms.MSalias == alias {
				deviceauthorized = true
				messageloop = false
			}

		}
		if errsubauthorizations != nil {
			messageloop = false
		}

	}
	//	if !deviceauthorized {
	//		Send(NatsUserDevices, NatsUserDevicesPassword, "DEVICES", "devices."+alias, "Add", alias)
	//	}

	dcerror := b.Js.DeleteConsumer(b.Ctx, "authorizations"+alias)
	if dcerror != nil {
		log.Println("nhnats.go CheckAUTH Consumer not found: ", dcerror)
	}
	a.Ctxcan()
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
