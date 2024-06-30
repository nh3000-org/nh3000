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
	MSiduuid   string
	MSalias    string
	MShostname string
	MSipadrs   string
	MSmessage  string
	MSnodeuuid string
	MSdate     string
}

var ms = MessageStore{}

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

/* // check msg by uuid
func CheckNatsMsgByUUID(iduuid string) bool {
	for _, v := range NatsMessages {
		//log.Println("id ", iduuid, "v.id ", v.MSiduuid)
		if iduuid == v.MSiduuid {
			return true
		}
	}
	return false
} */

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
func Send(m string, alias string) bool {

	EncMessage := MessageStore{}
	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = "\n" + getLangsNats("ms-nhn")
	} else {
		EncMessage.MShostname = "\n" + getLangsNats("ms-hn") + name
	}

	ifas, err := net.Interfaces()
	if err != nil {
		EncMessage.MShostname += "\n-  " + getLangsNats("ms-carrier")
	}
	if err == nil {
		var as []string
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				as = append(as, a)
			}
		}
		EncMessage.MShostname += "\n" + getLangsNats("ms-mi")
		for i, s := range as {
			EncMessage.MShostname += "\n- " + strconv.Itoa(i) + " : " + s
		}
		addrs, _ := net.InterfaceAddrs()
		EncMessage.MShostname += "\n" + getLangsNats("ms-ad")
		for _, addr := range addrs {
			EncMessage.MShostname += "\n- " + addr.String()
		}
	}
	EncMessage.MSalias = alias
	EncMessage.MSnodeuuid = "\n" + getLangsNats("ms-ni") + NatsNodeUUID
	msiduuid := uuid.New().String()
	EncMessage.MSiduuid = "\n" + getLangsNats("ms-msg") + msiduuid
	//log.Println("in send ", m, " ", msiduuid)
	EncMessage.MSdate = "\n" + getLangsNats("ms-on") + time.Now().Format(time.UnixDate)
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
	Sendjs(s)
	runtime.GC()
	return false
}

func Sendjs(m string) {
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
	maxage, _ := time.ParseDuration(NatsMsgMaxAge)
	ctx, ctxcancel := context.WithCancel(context.Background())
	//pubctx := context.Background()
	//nats stream add ORDERS --subjects "ORDERS.*" --ack --max-msgs=-1 --max-bytes=-1 --max-age=1y --storage file --retention limits --max-msg-size=-1 --discard=old
	//nats consumer add ORDERS NEW --filter ORDERS.received --ack explicit --pull --deliver all --max-deliver=-1 --sample 100
	//nats consumer add ORDERS DISPATCH --filter ORDERS.processed --ack explicit --pull --deliver all --max-deliver=-1 --sample 100
	//nats consumer add ORDERS MONITOR --filter '' --ack none --target monitor.ORDERS --deliver last --replay instant
	_, streamerr := sendjs.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Description: "ClientStream",
		Name:        NatsQueue,
		Subjects:    []string{strings.ToLower(NatsQueue + ".*")},
		MaxAge:      maxage,
		MaxMsgs:     -1,
		MaxBytes:    -1,
		MaxMsgSize:  -1,
		Retention:   jetstream.LimitsPolicy,
		Discard:     jetstream.DiscardOld,
		Storage:     jetstream.FileStorage,
		//Replicas:   5,
	})
	if streamerr != nil {
		log.Println("streamerror " + streamerr.Error())
	}

	if streamerr != nil {
		log.Println("streamerror " + streamerr.Error())
	}
	//log.Println("in ", getMessageToSend())

	//log.Println("Publish ")
	//sndmsg := &nats.Msg{Subject: strings.ToLower(NatsQueue) + ".>", Data: []byte(getMessageToSend())}
	_, puberr := sendjs.Publish(ctx, strings.ToLower(NatsQueue)+".client", []byte(m))
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
func ReceiveJS() {
	var certpool = docerts()

	for {
		select {
		case <-QuitReceive:
			return

		default:
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
				log.Println("Connect " + getLangsNats("ms-snd") + " " + getLangsNats("ms-err7") + connecterr.Error())
			}
			js, jetstreamerr := natsconnect.JetStream()
			if jetstreamerr != nil {
				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + jetstreamerr.Error())
				}
				log.Println("jetstreamerror " + jetstreamerr.Error())
			}
			maxage, _ := time.ParseDuration(NatsMsgMaxAge)
			js.AddStream(&nats.StreamConfig{
				Name:     NatsQueue,
				Subjects: []string{strings.ToLower(NatsQueue) + ".*"},
			})
			_, err1 := js.AddConsumer(NatsQueue, &nats.ConsumerConfig{
				Durable:           NatsNodeUUID,
				AckPolicy:         nats.AckExplicitPolicy,
				InactiveThreshold: maxage,
				DeliverPolicy:     nats.DeliverAllPolicy,
				ReplayPolicy:      nats.ReplayInstantPolicy,
			})
			if err1 != nil {
				//log.Println(err1.Error())
				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-carrier") + " " + err1.Error())
					log.Println("Consumer " + getLangsNats("ms-carrier") + " " + err1.Error())
				}
			}
			sub, errsub := js.PullSubscribe("", "", nats.BindStream(NatsQueue))
			if errsub != nil {
				//log.Println(errsub.Error())
				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-carrier") + errsub.Error())
					log.Println("Pull " + getLangsNats("ms-carrier") + " " + errsub.Error())
				}
			}
			msgs, err := sub.Fetch(100)
			if err != nil {
				//log.Println(err.Error())
				if FyneMessageWin != nil {
					runtime.GC()
					runtime.ReadMemStats(&memoryStats)
					FyneMessageWin.SetTitle(getLangsNats("ms-nnm") + err.Error() + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
					//yulog.Println("Fetch " + GetLangs("ms-carrier") + " " + err.Error())
				}
			}
			if len(msgs) > 0 {
				for i := 0; i < len(msgs); i++ {
					msgs[i].Nak()
					ms = MessageStore{}
					err1 := json.Unmarshal([]byte(string(Decrypt(string(msgs[i].Data), NatsQueuePassword))), &ms)
					//log.Println("natssuberr ", fetcherr, " msg ", ms)

					//err1 := json.Unmarshal([]byte(string(Decrypt(string(msg.Data), NatsQueuePassword))), &ms)
					if err1 != nil {
						// send decrypt
						if FyneMessageWin != nil {
							FyneMessageWin.SetTitle(getLangsNats("ms-mde"))
						}
						log.Println("un marhal", err1)
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
							//log.Println("adding ms ", ms.MSiduuid)
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
					//}

					//FyneMessageList.Refresh()
					//natssub.Drain()
				}
				FyneMessageList.Refresh()
			}
			/* 			if fetcherr != nil {
				log.Println("natsfetcherr in ", fetcherr)
			} */

			//sub.Drain()
			dltconserr := js.DeleteConsumer(NatsQueue, NatsNodeUUID)
			if dltconserr != nil {
				log.Println("dltconserr ", dltconserr)
			}
			natsconnect.Close()
			//log.Println("big ", bigerr)
		}

		//NatsMessages = nil

		time.Sleep(20 * time.Second)

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
			FyneMessageWin.SetTitle(GetLangs("ms-carrier") + err.Error())
		}
	}
	js.AddStream(&nats.StreamConfig{
		Name:     NatsQueue + NatsNodeUUID,
		Subjects: []string{strings.ToLower(NatsQueue) + ".>"},
	})

	msgmaxage, _ := time.ParseDuration(NatsMsgMaxAge)
	_, err1 := js.AddConsumer(NatsQueue, &nats.ConsumerConfig{
		Durable:           NatsNodeUUID,
		AckPolicy:         nats.AckExplicitPolicy,
		InactiveThreshold: msgmaxage,
		DeliverPolicy:     nats.DeliverAllPolicy,
		ReplayPolicy:      nats.ReplayInstantPolicy,
	})

	if err1 != nil {
		//log.Println(err1.Error())
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(GetLangs("ms-carrier") + err1.Error())
		}
	}

	sub, errsub := js.PullSubscribe("", "", nats.BindStream(NatsQueue))
	if errsub != nil {
		//log.Println(errsub.Error())
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(GetLangs("ms-carrier") + errsub.Error())
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
					FyneMessageWin.SetTitle(GetLangs("ms-carrier") + err.Error())
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
						log.Println("receive un marhal", err1)
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

//}

// security erase jetstream data
func Erase() {
	//docerts()
	//log.Println(config.GetLangs("ms-era"))

	nc, err := nats.Connect(NatsServer, nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(docerts()))
	if err != nil {
		log.Println("Erase Connect", GetLangs("ms-erac"), err.Error())
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("Erase Jetstream Make ", GetLangs("ms-eraj"), err)
	}

	NatsMessages = nil
	err1 := js.PurgeStream(NatsQueue)
	if err1 != nil {
		log.Println("Erase Jetstream Purge", GetLangs("ms-dels"), err1)
	}
	err2 := js.DeleteStream(NatsQueue)
	if err2 != nil {
		log.Println("Erase Jetstream Delete", GetLangs("ms-dels"), err1)
	}
	msgmaxage, _ := time.ParseDuration(NatsMsgMaxAge)
	js1, err3 := js.AddStream(&nats.StreamConfig{
		Name:     NatsQueue,
		Subjects: []string{strings.ToLower(NatsQueue) + ".>"},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
	})
	if err3 != nil {
		log.Println("Erase Addstream ", GetLangs("ms-adds"), err3)
	}
	fmt.Printf("js1: %v\n", js1)

	Send(GetLangs("ms-sece"), NatsAlias)
	nc.Close()
}
