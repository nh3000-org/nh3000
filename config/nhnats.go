package config

import (
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
var duration time.Duration = 604800000000
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
func DeleteNatsMsgByUUID(iduuid string) {
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

func Connect() (*nats.Conn, nats.JetStreamContext) {

	NC, err := nats.Connect(NatsServer, nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(docerts()), nats.PingInterval(pinginterval))
	if err != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error())
		}
		log.Println(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error())
	}
	JS, err := NC.JetStream()
	if err != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error() + <-JS.StreamNames())
		}
		log.Println(getLangsNats("ms-snd") + getLangsNats("ms-err7") + err.Error() + <-JS.StreamNames())
	}
	JS.AddStream(&nats.StreamConfig{
		Name:     NatsQueue + NatsNodeUUID,
		Subjects: []string{strings.ToLower(NatsQueue) + ".>"},
	})

	_, err1 := JS.AddConsumer(NatsQueue, &nats.ConsumerConfig{
		Durable:           NatsNodeUUID,
		AckPolicy:         nats.AckExplicitPolicy,
		InactiveThreshold: duration,
		DeliverPolicy:     nats.DeliverAllPolicy,
		ReplayPolicy:      nats.ReplayInstantPolicy,
	})
	if err1 != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-carrier") + err1.Error())
		}
		log.Println("jet stream " + err1.Error())
	}

	return NC, JS
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
	NC, JS := Connect()
	_, errp := JS.Publish(strings.ToLower(NatsQueue)+".logger", []byte(Encrypt(string(jsonmsg), NatsQueuePassword)))
	if errp != nil {
		if FyneMessageWin != nil {
			FyneMessageWin.SetTitle(getLangsNats("ms-snd") + errp.Error())
		}
		log.Println(getLangsNats("ms-snd"), errp)
	}
	NC.Drain()
	NC.Close()
	runtime.GC()
	return false
}

// thread for receiving messages
func Receive() {

	for {
		select {
		case <-QuitReceive:
			return
		default:
			NatsReceivingMessages = true
			NC, JS := Connect()
			sub, errsub := JS.PullSubscribe("", "", nats.BindStream(NatsQueue))
			if errsub != nil {
				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-carrier") + errsub.Error())
				}
				log.Println("pull " + errsub.Error())
			}
			msgs, err := sub.Fetch(100)
			if err != nil {

				if FyneMessageWin != nil {
					FyneMessageWin.SetTitle(getLangsNats("ms-carrier") + err.Error())
				}
				log.Println("fetch" + err.Error())
			}
			if len(msgs) > 0 {
				for i := 0; i < len(msgs); i++ {
					ms := MessageStore{}
					err1 := json.Unmarshal([]byte(string(Decrypt(string(msgs[i].Data), NatsQueuePassword))), &ms)
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
					msgs[i].Ack()
				}
				if FyneMessageWin != nil {
					runtime.GC()
					runtime.ReadMemStats(&memoryStats)
					FyneMessageWin.SetTitle(getLangsNats("ms-err6-1") + strconv.Itoa(len(NatsMessages)) + getLangsNats("ms-err6-2") + " " + strconv.FormatUint(memoryStats.Alloc/1024/1024, 10) + " Mib")
				}
				FyneMessageList.Refresh()
			}
			NC.Drain()
			NC.Close()

			time.Sleep(1 * time.Second)

		}
	}
}

// security erase jetstream data
func Erase() {

	nc, err := nats.Connect(NatsServer, nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(docerts()))
	if err != nil {
		log.Println("Erase Connect", getLangsNats("ms-erac"), err.Error())
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("Erase Jetstream Make ", getLangsNats("ms-eraj"), err)
	}

	NatsMessages = nil
	err1 := js.PurgeStream(NatsQueue)
	if err1 != nil {
		log.Println("Erase Jetstream Purge", getLangsNats("ms-dels"), err1)
	}
	err2 := js.DeleteStream(NatsQueue)
	if err2 != nil {
		log.Println("Erase Jetstream Delete", getLangsNats("ms-dels"), err1)
	}
	msgmaxage, _ := time.ParseDuration(NatsMsgMaxAge)

	js1, err3 := js.AddStream(&nats.StreamConfig{
		Name:     NatsQueue,
		Subjects: []string{strings.ToLower(NatsQueue) + ".>"},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
	})
	if err3 != nil {
		log.Println("Erase Addstream ", getLangsNats("ms-adds"), err3)
	}
	fmt.Printf("js1: %v\n", js1)

	Send(getLangsNats("ms-sece"), NatsAlias)
	nc.Close()
}
