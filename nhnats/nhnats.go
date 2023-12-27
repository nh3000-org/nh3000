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
	"github.com/nh3000-org/nh3000/nhauth"
	"github.com/nh3000-org/nh3000/nhcrypt"
	"github.com/nh3000-org/nh3000/nhlang"

	//"github.com/nh3000-org/nh3000/nhpanes"
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
var TLS tls.Config
var MyLogLang = "eng"

// eng esp cmn hin
var MyLangs = map[string]string{
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
}

// return translation strings
func GetLangs(mystring string) string {
	value, err := MyLangs[MyLogLang+"-"+mystring]
	if err == false {
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
		ok := RootCAs.AppendCertsFromPEM([]byte(nhauth.DefaultCaroot))
		if !ok {
			log.Println("nhnats.go init rootCAs")
		}
		Clientcert, err := tls.X509KeyPair([]byte(nhauth.DefaultClientcert), []byte(nhauth.DefaultClientkey))
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

	EncMessage.MSnodeuuid = "\n" + GetLangs("ms-ni") + nhpref.NodeUUID
	iduuid := uuid.New().String()
	EncMessage.MSiduuid = "\n" + GetLangs("ms-msg") + iduuid
	EncMessage.MSdate = "\n" + GetLangs("ms-on") + time.Now().Format(time.UnixDate)
	EncMessage.MSmessage = m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	if jsonerr != nil {
		log.Println("FormatMessage ", jsonerr)
	}
	ejson, _ := nhcrypt.Encrypt(string(jsonmsg), nhauth.QueuePassword)
	NC, err := nats.Connect(nhauth.DefaultServer, nats.UserInfo(nhauth.User, nhauth.UserPassword), nats.Secure(&TLS))
	if err != nil {
		fmt.Println("Send " + GetLangs("ls-err7") + err.Error())
	}
	JS, err := NC.JetStream()
	if err != nil {
		fmt.Println("Send " + GetLangs("ls-err7") + err.Error() + <-JS.StreamNames())
	}
	_, errp := JS.Publish(strings.ToLower(nhauth.Queue)+".logger", []byte(ejson))
	if errp != nil {
		return true
	}
	NC.Drain()

	return false
}

// thread for receiving messages
func Receive() {
	docerts()

	nhpref.ReceivingMessages = true
	for {
		select {
		case <-QuitReceive:
			return
		default:
			NatsMessages = nil
			nc, err := nats.Connect(nhpref.Server, nats.UserInfo(nhauth.User, nhauth.UserPassword), nats.Secure(&TLS))
			if err != nil {
				log.Println("Receive ", GetLangs("ms-err2"))
			}
			js, err := nc.JetStream()
			if err != nil {
				log.Println("Receive JetStream ", err)
			}
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
				log.Println("Receive AddConsumer ", GetLangs("ms-err3")+err1.Error())
			}
			sub, errsub := js.PullSubscribe("", "", nats.BindStream(nhpref.Queue))
			if errsub != nil {
				log.Println("Receive Pull Subscribe ", GetLangs("ms-err4")+errsub.Error())
			}
			msgs, err := sub.Fetch(100)
			if err != nil {
				log.Println("Receive fetch  ", err)
			}
			nhpref.ClearMessageDetail = true
			if len(msgs) > 0 {
				for i := 0; i < len(msgs); i++ {
					handleMessage(msgs[i])
					msgs[i].Nak()
				}
			}
			if nhutil.GetMessageWin() != nil {
				nhutil.GetMessageWin().SetTitle(GetLangs("ms-err6-1") + strconv.Itoa(len(msgs)) + nhlang.GetLangs("ms-err6-2"))
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
		ejson = GetLangs("ms-unk")
	}
	if nhpref.Filter {
		if strings.Contains(ms.MSmessage, GetLangs("ls-con")) || strings.Contains(ms.MSmessage, nhlang.GetLangs("ls-dis")) {
			return ""
		}
	}
	NatsMessages = append(NatsMessages, ms)

	return ms.MSiduuid
}

// security erase jetstream data
func Erase() {
	docerts()
	log.Println(GetLangs("ms-era"))

	nc, err := nats.Connect(nhpref.Server, nats.UserInfo(nhauth.User, nhauth.UserPassword), nats.Secure(&TLS))
	if err != nil {
		log.Println("Erase Connect", GetLangs("ms-erac"), err.Error())
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("Erase Jetstream Make ", GetLangs("ms-eraj"), err)
	}

	NatsMessages = nil
	err1 := js.PurgeStream(nhpref.Queue)
	if err1 != nil {
		log.Println("Erase Jetstream Purge", GetLangs("ms-dels"), err1)
	}
	err2 := js.DeleteStream(nhpref.Queue)
	if err2 != nil {
		log.Println("Erase Jetstream Delete", GetLangs("ms-dels"), err1)
	}
	msgmaxage, _ := time.ParseDuration(nhpref.Msgmaxage)
	js1, err3 := js.AddStream(&nats.StreamConfig{
		Name:     nhpref.Queue,
		Subjects: []string{strings.ToLower(nhpref.Queue) + ".>"},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
	})
	if err3 != nil {
		log.Println("Erase Addstream ", GetLangs("ms-adds"), err3)
	}
	fmt.Printf("js1: %v\n", js1)

	Send(GetLangs("ms-sece"), nhpref.Alias)
	nc.Close()
}
