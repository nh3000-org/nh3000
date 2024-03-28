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

var NatsMessages []MessageStore
var ackMap = make(map[string]bool)
var QuitReceive = make(chan bool)
var TLS tls.Config
var MyLogLang = "eng"

// eng esp cmn hin
var MyLangsNats = map[string]string{
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
	"eng-ms-con":      "Connected",
	"spa-ms-con":      "Conectada",
	"hin-ms-con":      "जुड़े हुए",
	"eng-ms-dis":      "Disconnected",
	"spa-ms-dis":      "Desconectada",
	"hin-ms-dis":      "डिस्कनेक्ट किया गया",
}

// return translation strings
func GetLangsNats(mystring string) string {
	value, err := MyLangs[MyLogLang+"-"+mystring]
	if !err {
		return "xxx"
	}
	return value
}

func docerts() {

	var done = false
	if !done {
		RootCAs, _ := x509.SystemCertPool()
		if RootCAs == nil {
			RootCAs = x509.NewCertPool()
		}
		ok := RootCAs.AppendCertsFromPEM([]byte(GetCaroot()))
		if !ok {
			log.Println("nhnats.go init rootCAs")
		}
		Clientcert, err := tls.X509KeyPair([]byte(GetClientCert()), []byte(GetClientKey()))
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
func Send(m string, alias string) bool {

	docerts()
	EncMessage := MessageStore{}
	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = "\n" + GetLangs("ms-nhn")
	} else {
		EncMessage.MShostname = "\n" + GetLangs("ms-hn") + name
	}

	ifas, err := net.Interfaces()
	if err != nil {
		EncMessage.MShostname += "\n-  " + GetLangs("ms-carrier")
	}
	if err == nil {
		var as []string
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				as = append(as, a)
			}
		}
		EncMessage.MShostname += "\n" + GetLangs("ms-mi")
		for i, s := range as {
			EncMessage.MShostname += "\n- " + strconv.Itoa(i) + " : " + s
		}
		addrs, _ := net.InterfaceAddrs()
		EncMessage.MShostname += "\n" + GetLangs("ms-ad")
		for _, addr := range addrs {
			EncMessage.MShostname += "\n- " + addr.String()
		}
	}
	EncMessage.MSalias = alias
	EncMessage.MSnodeuuid = "\n" + GetLangs("ms-ni") + GetNodeUUID()
	msiduuid := uuid.New().String()
	EncMessage.MSiduuid = "\n" + GetLangs("ms-msg") + msiduuid
	EncMessage.MSdate = "\n" + GetLangs("ms-on") + time.Now().Format(time.UnixDate)
	EncMessage.MSmessage = m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	if jsonerr != nil {
		log.Println("FormatMessage ", jsonerr)
	}
	//log.Println("jsonmsg ", string(jsonmsg))
	ejson := Encrypt(string(jsonmsg), GetQueuePassword())
	//log.Println("ejson ", string(ejson))
	NC, err := nats.Connect(GetServer(), nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(&TLS))
	if err != nil {
		fmt.Println("Send " + GetLangs("ms-err7") + err.Error())
	}
	JS, err := NC.JetStream()
	if err != nil {
		fmt.Println("Send " + GetLangs("ms-err7") + err.Error() + <-JS.StreamNames())
	}
	_, errp := JS.Publish(strings.ToLower(GetQueue())+".logger", []byte(ejson))
	if errp != nil {
		return true
	}
	NC.Drain()

	return false
}

// thread for receiving messages
func Receive() {

	docerts()

	SetReceivingMessages(true)
	for {
		select {
		case <-QuitReceive:
			return
		default:
			NatsMessages = nil
			nc, err := nats.Connect(GetServer(), nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(&TLS))
			if err != nil {
				if GetMessageWindow() != nil {
					GetMessageWindow().SetTitle(GetLangs("ms-carrier") + err.Error())
				}

			}
			js, err := nc.JetStream()
			if err != nil {
				if GetMessageWindow() != nil {
					GetMessageWindow().SetTitle(GetLangs("ms-carrier") + err.Error())
				}
			}
			js.AddStream(&nats.StreamConfig{
				Name:     GetQueue() + GetNodeUUID(),
				Subjects: []string{strings.ToLower(GetQueue()) + ".>"},
			})
			var duration time.Duration = 604800000000
			_, err1 := js.AddConsumer(GetQueue(), &nats.ConsumerConfig{
				Durable:           GetNodeUUID(),
				AckPolicy:         nats.AckExplicitPolicy,
				InactiveThreshold: duration,
				DeliverPolicy:     nats.DeliverAllPolicy,
				ReplayPolicy:      nats.ReplayInstantPolicy,
			})
			if err1 != nil {
				//log.Println(err1.Error())
				if GetMessageWindow() != nil {
					GetMessageWindow().SetTitle(GetLangs("ms-carrier") + err1.Error())
				}
			}
			sub, errsub := js.PullSubscribe("", "", nats.BindStream(GetQueue()))
			if errsub != nil {
				//log.Println(errsub.Error())
				if GetMessageWindow() != nil {
					GetMessageWindow().SetTitle(GetLangs("ms-carrier") + errsub.Error())
				}
			}
			msgs, err := sub.Fetch(100)
			if err != nil {
				//log.Println(err.Error())
				if GetMessageWindow() != nil {

					GetMessageWindow().SetTitle(GetLangs("ms-carrier") + err.Error())
				}
			}
			SetClearMessageDetail(true)
			var acked = 0
			if len(msgs) > 0 {
				for i := 0; i < len(msgs); i++ {
					if !handleMessage(msgs[i]) {

						acked++
					}
					msgs[i].Ack()
				}

			}
			//var shadowackMap = ackMap
			if len(ackMap) > 0 {
				for k, v := range ackMap {
					if !v {
						delete(ackMap, k)
					}
				}
			}
			//shadowackMap = nil
			if GetMessageWindow() != nil {
				var m runtime.MemStats
				runtime.ReadMemStats(&m)
				GetMessageWindow().SetTitle(GetLangs("ms-err6-1") + strconv.Itoa(len(msgs)) + GetLangs("ms-err6-2") + " - " + strconv.Itoa(acked) + " Acked" + " " + strconv.FormatUint(m.Alloc/1024/1024, 10) + " Mib")
				GetMessageList().Refresh()
			}
			nc.Close()
			time.Sleep(30 * time.Second)
		}
	}
}

// decrypt payload
func handleMessage(m *nats.Msg) bool {

	ms := MessageStore{}

	var ejson = string(Decrypt(string(m.Data), GetQueuePassword()))

	err1 := json.Unmarshal([]byte(ejson), &ms)
	if err1 != nil {
		log.Println("NATS Receive ", ejson[0])
	}
	if GetFilter() {
		if strings.Contains(ms.MSmessage, GetLangs("ms-con")) {
			return false
		}
		if strings.Contains(ms.MSmessage, GetLangs("ms-dis")) {
			return false
		}
	}
	if !ackMap[ms.MSiduuid] {
		ackMap[ms.MSiduuid] = false
		NatsMessages = append(NatsMessages, ms)
		return true
	}
	//ackMap[ms.MSiduuid] = false
	return false
}
func SetAck(a string) {
	ackMap[a] = true
}

// security erase jetstream data
func Erase() {

	docerts()

	nc, err := nats.Connect(GetServer(), nats.UserInfo(NatsUser, NatsUserPassword), nats.Secure(&TLS))
	if err != nil {
		log.Println("Erase Connect", GetLangs("ms-erac"), err.Error())
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("Erase Jetstream Make ", GetLangs("ms-eraj"), err)
	}

	NatsMessages = nil
	err1 := js.PurgeStream(GetQueue())
	if err1 != nil {
		log.Println("Erase Jetstream Purge", GetLangs("ms-dels"), err1)
	}
	err2 := js.DeleteStream(GetQueue())
	if err2 != nil {
		log.Println("Erase Jetstream Delete", GetLangs("ms-dels"), err1)
	}
	msgmaxage, _ := time.ParseDuration(GetMsgMaxAge())

	js1, err3 := js.AddStream(&nats.StreamConfig{
		Name:     GetQueue(),
		Subjects: []string{strings.ToLower(GetQueue()) + ".>"},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
	})
	if err3 != nil {
		log.Println("Erase Addstream ", GetLangs("ms-adds"), err3)
	}
	fmt.Printf("js1: %v\n", js1)

	Send(GetLangs("ms-sece"), GetAlias())
	nc.Close()
}
