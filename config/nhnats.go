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
	MSiduuid   string
	MSalias    string
	MShostname string
	MSipadrs   string
	MSmessage  string
	MSnodeuuid string
	MSdate     string
}

var shortServerName string
var shortServerName1 string
var memoryStats runtime.MemStats
var NatsMessages = make(map[int]MessageStore)

var pinginterval time.Duration = (30 * time.Minute)

var fyneFilterFound = false

var QuitReceive = make(chan bool)

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
		if iduuid == v.MSiduuid {
			return true
		}
	}
	return false
}

// var TLSConfig = &tls.Config{}

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
	EncMessage.MSdate = "\n" + getLangsNats("ms-on") + time.Now().Format(time.UnixDate)
	EncMessage.MSmessage = m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	if jsonerr != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-err8") + jsonerr.Error())
		}
		log.Println(getLangsNats("ms-err8"), jsonerr.Error())
	}

	sendctx := context.Background()
	sendctx, sendcancel := context.WithTimeout(sendctx, 30*time.Second)

	sendnatsconn, err := nats.Connect(NatsServer, nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(docerts()))
	if err != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error())
		}
		log.Println(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error())
	}
	sendjetstream, err := jetstream.New(sendnatsconn)
	if err != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error())
		}
		log.Println(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error())
	}
	_, errp := sendjetstream.Publish(sendctx, strings.ToLower(NatsQueue)+"."+strings.ToLower(NatsQueue), []byte(Encrypt(string(jsonmsg), NatsQueuePassword)))

	//_, errp := JS.Publish(strings.ToLower(NatsQueue)+".logger", []byte(Encrypt(string(jsonmsg), NatsQueuePassword)))
	if errp != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + errp.Error())
		}
		log.Println("publish ", getLangsNats("ms-snd"), errp)
	}
	sendnatsconn.Drain()
	sendnatsconn.Close()
	sendcancel()
	runtime.GC()
	return false
}

// thread for receiving messages
func Receive() {
	var certpool = docerts()
	msgmaxage, _ := time.ParseDuration(NatsMsgMaxAge)
	for {
		select {
		case <-QuitReceive:
			return
		default:
			ctx := context.TODO()
			ctx, cancel := context.WithTimeout(ctx, 30*time.Hour)

			NC, err := nats.Connect(NatsServer, nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(certpool), nats.PingInterval(pinginterval))
			if err != nil {
				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error())
				}
				log.Println(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error())
			}
			JS, err := jetstream.New(NC)
			if err != nil {
				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error())
				}
				log.Println(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error())
			}

			s, err := JS.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
				Name:      NatsQueue,
				Subjects:  []string{strings.ToLower(NatsQueue) + ".*"},
				Storage:   jetstream.FileStorage,
				Retention: jetstream.LimitsPolicy,
				Discard:   jetstream.DiscardOld,
				MaxAge:    msgmaxage,
			})
			if err != nil {
				log.Fatal(err)
			}
			cons, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
				Durable:       NatsQueueDurable,
				DeliverPolicy: jetstream.DeliverAllPolicy,

				//AckPolicy: jetstream.AckNonePolicy,
			})
			if err != nil {
				log.Fatal("consumer ", err)
			}
			it, err := cons.Messages()

			if err != nil {
				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-carrier") + err.Error())
				}
				log.Println("pull " + err.Error())
			}
			//for {

			msg, err := it.Next()
			if err != nil {
				log.Println("receiving ", err.Error())
			}
			ms := MessageStore{}
			msg.Nak()
			err1 := json.Unmarshal([]byte(string(Decrypt(string(msg.Data()), NatsQueuePassword))), &ms)
			if err1 != nil {
				// send decrypt
				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-mde"))
				}
				log.Println("un marhal")
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
			cn := s.ConsumerNames(ctx)
			errdc := s.DeleteConsumer(ctx, <-cn.Name())
			if errdc != nil {
				log.Println("delete consumer ", errdc)
			}
			errds := JS.DeleteStream(ctx, NatsQueue)
			if errds != nil {
				log.Println("delete stream ", errds)
			}
			cancel()
			ctx.Done()
			it.Stop()
			NC.Drain()
			NC.Close()

			time.Sleep(10 * time.Second)

		}
	}
}

// security erase jetstream data
func Erase() {

	nc, err := nats.Connect(NatsServer, nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(docerts()))
	if err != nil {
		log.Println("Erase Connect", getLangsNats("ms-erac"), err.Error())
	}
	js, err := jetstream.New(nc)
	if err != nil {
		log.Println("Erase Jetstream Make ", getLangsNats("ms-eraj"), err)
	}
	msgmaxage, errmma := time.ParseDuration(NatsMsgMaxAge)
	log.Println("mma ", msgmaxage, errmma)
	cfg := jetstream.StreamConfig{
		Name:      NatsQueue,
		Subjects:  []string{strings.ToLower(NatsQueue) + ".>"},
		Storage:   jetstream.FileStorage,
		Retention: jetstream.LimitsPolicy,
		Discard:   jetstream.DiscardOld,
		MaxAge:    msgmaxage,
	}
	log.Println("max age ", NatsMsgMaxAge)
	log.Println("cfg ", cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	NatsMessages = nil

	err2 := js.DeleteStream(ctx, NatsQueue)
	if err2 != nil {
		log.Println("Erase Jetstream Delete", getLangsNats("ms-dels"), err2)
	}

	js1, err3 := js.CreateStream(ctx, cfg)
	if err3 != nil {
		log.Println("Create Stream ", getLangsNats("ms-adds"), err3)
	}
	fmt.Printf("js1: %v\n", js1)

	Send(getLangsNats("ms-sece"), NatsAlias)
	nc.Close()
}
