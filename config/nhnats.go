package config

import (
	//"context"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"

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

var mts string
var shortServerName string
var shortServerName1 string
var memoryStats runtime.MemStats
var NatsMessages = make(map[int]MessageStore)

var pinginterval time.Duration = (30 * time.Minute)

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

// check msg by uuid
func CheckNatsMsgByUUID(iduuid string) bool {
	for _, v := range NatsMessages {
		log.Println("id ", iduuid, "v.id ", v.MSiduuid)
		if iduuid == v.MSiduuid {
			return true
		}
	}
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
func Send(m string, alias string) bool {
	log.Println("in send ", m)
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
	log.Println("sending 2 ", s)
	setMessageToSend(s)
	runtime.GC()
	return false
}
func setMessageToSend(nm string) {
	//	log.Println("settosend ", string(nm))
	MessageToSend = nm
}
func getMessageToSend() string {
	//	log.Println("gettosend ", string(MessageToSend))
	return MessageToSend
}

// thread for receiving messages
func ReceiveJS() {
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
	js, jetstreamerr := jetstream.New(natsconnect)
	if jetstreamerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + jetstreamerr.Error())
		}
		log.Println("jetstreamerror " + jetstreamerr.Error())
	}
	maxage, _ := time.ParseDuration(NatsMsgMaxAge)
	ctx := context.Background()
	//pubctx := context.Background()
	//nats stream add ORDERS --subjects "ORDERS.*" --ack --max-msgs=-1 --max-bytes=-1 --max-age=1y --storage file --retention limits --max-msg-size=-1 --discard=old
	//nats consumer add ORDERS NEW --filter ORDERS.received --ack explicit --pull --deliver all --max-deliver=-1 --sample 100
	//nats consumer add ORDERS DISPATCH --filter ORDERS.processed --ack explicit --pull --deliver all --max-deliver=-1 --sample 100
	//nats consumer add ORDERS MONITOR --filter '' --ack none --target monitor.ORDERS --deliver last --replay instant
	stream, streamerr := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
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

	/*
		natssub, natssuberr := natsconnect.SubscribeSync(strings.ToLower(NatsQueue) + ".>")
		if natssuberr != nil {
			if FyneMessageWin != nil {
				FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + natssuberr.Error())
			}
			log.Println(getLangsNats("ms-snd") + getLangsNats("ms-err7") + natssuberr.Error())
		} */
	//nats consumer add ORDERS MONITOR --filter '' --ack none --target monitor.ORDERS --deliver last --replay instant
	concfg := jetstream.ConsumerConfig{
		Name:          NatsQueueDurable,
		Durable:       NatsQueueDurable,
		DeliverPolicy: jetstream.DeliverAllPolicy,
		//AckPolicy:     jetstream.AckNonePolicy,
		FilterSubject: strings.ToLower(NatsQueue) + ".*",
		ReplayPolicy:  jetstream.ReplayInstantPolicy,
	}
	setupcons, setuperr := stream.CreateOrUpdateConsumer(ctx, concfg)
	//AckPolicy: jetstream.AckNonePolicy,
	if setuperr != nil {
		log.Println("setup error ", setuperr.Error())
	}

	/* 	consume, consumeerr := stream.CreateConsumer(ctx, jetstream.ConsumerConfig{
	   		Durable: NatsQueueDurable,
	   		//AckPolicy: jetstream.AckNonePolicy,
	   	})
	   	if consumeerr != nil {
	   		log.Println("consumeerror " + consumeerr.Error())
	   	} */

	//msgmaxage, _ := time.ParseDuration(NatsMsgMaxAge)
	for {
		select {
		case <-QuitReceive:
			return

		default:
			//log.Println("in ", getMessageToSend())
			mts = getMessageToSend()
			if mts != "" {
				//log.Println("Publish ")
				//sndmsg := &nats.Msg{Subject: strings.ToLower(NatsQueue) + ".>", Data: []byte(getMessageToSend())}
				_, puberr := js.PublishAsync(strings.ToLower(NatsQueue)+".*", []byte(getMessageToSend()))
				if puberr != nil {
					if FyneMessageWin != nil {
						FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + puberr.Error())
					}
					log.Println("Publish Error " + puberr.Error())

					//SendMessage = false
				}
				//log.Println("in error ", pubit.Err())
				setMessageToSend("")
			}

			//NatsMessages = nil
			_, fetcherr := setupcons.Consume(func(msg jetstream.Msg) {
				msg.Ack()
				var ms = MessageStore{}
				err1 := json.Unmarshal([]byte(string(Decrypt(string(msg.Data()), NatsQueuePassword))), &ms)
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
					if !CheckNatsMsgByUUID(ms.MSiduuid) {
						log.Println("adding ms ")
						NatsMessages[len(NatsMessages)] = ms
					}
				}

				if FyneMessageWin != nil {
					runtime.GC()
					runtime.ReadMemStats(&memoryStats)
					FyneMessageWin.SetTitle(getLangsNats("ms-err6-1") + strconv.Itoa(len(NatsMessages)) + getLangsNats("ms-err6-2") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
				}
				//}

				FyneMessageList.Refresh()
				//natssub.Drain()

			})
			/* 			if fetcherr != nil {
				log.Println("natsfetcherr in ", fetcherr)
			} */
			if fetcherr != nil {
				log.Println("natsfetcherr getting ", fetcherr.Error())

			}
		}
		//NatsMessages = nil

		time.Sleep(10 * time.Second)

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
	natssubdurable, natssuberrdurable := natsconnect.SubscribeSync(NatsQueueDurable)
	if natssuberrdurable != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + natssuberrdurable.Error())
		}
		log.Println("Durable " + natssubdurable.Subject + natssuberrdurable.Error())
	}
	natssub, natssuberr := natsconnect.SubscribeSync(strings.ToLower(NatsQueue) + ".*")
	if natssuberr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + natssuberr.Error())
		}
		log.Println(getLangsNats("ms-snd") + getLangsNats("ms-err7") + natssuberr.Error())
	}
	//msgmaxage, _ := time.ParseDuration(NatsMsgMaxAge)
	for {
		select {
		case <-QuitReceive:
			return

		default:
			//log.Println("in ", getMessageToSend())
			if len(getMessageToSend()) != 0 {
				//log.Println("Publish ")
				sndmsg := &nats.Msg{Subject: strings.ToLower(NatsQueue) + ".>", Data: []byte(getMessageToSend())}
				puberr := natsconnect.PublishMsg(sndmsg)
				if puberr != nil {
					if FyneMessageWin != nil {
						FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + puberr.Error())
					}
					log.Println("Publish Error " + puberr.Error())

					//SendMessage = false
				}
				setMessageToSend("")
			}
			msg, natnextmsg := natssub.NextMsg(1 * time.Second)
			//log.Println("natssuberr in ", natssuberr, " msg ", natnextmsg)
			if natnextmsg == nil {
				var ms = MessageStore{}
				err1 := json.Unmarshal([]byte(string(Decrypt(string(msg.Data), NatsQueuePassword))), &ms)
				log.Println("natssuberr ", natssuberr, " msg ", ms)

				msg.Nak()
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
					if !CheckNatsMsgByUUID(ms.MSiduuid) {
						NatsMessages[len(NatsMessages)] = ms
					}
				}

				if FyneMessageWin != nil {
					runtime.GC()
					runtime.ReadMemStats(&memoryStats)
					FyneMessageWin.SetTitle(getLangsNats("ms-err6-1") + strconv.Itoa(len(NatsMessages)) + getLangsNats("ms-err6-2") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
				}
				//}

				FyneMessageList.Refresh()
				//natssub.Drain()
			}
		}
		//NatsMessages = nil

		time.Sleep(10 * time.Second)

	}
}

//}

// security erase jetstream data
func Erase() {

}
