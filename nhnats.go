package main

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

	"github.com/nh3000-org/nh3000/config"

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
var MyAckMap = make(map[string]bool)
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
	"eng-ms-con":      "Connected",
	"spa-ms-con":      "Conectada",
	"hin-ms-con":      "जुड़े हुए",
	"eng-ms-dis":      "Disconnected",
	"spa-ms-dis":      "Desconectada",
	"hin-ms-dis":      "डिस्कनेक्ट किया गया",
}

// return translation strings
func GetLangsNats(mystring string) string {
	value, err := config.MyLangs[MyLogLang+"-"+mystring]
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
		ok := RootCAs.AppendCertsFromPEM([]byte(config.GetCaroot()))
		if !ok {
			log.Println("nhnats.go init rootCAs")
		}
		Clientcert, err := tls.X509KeyPair([]byte(config.GetClientCert()), []byte(config.GetClientKey()))
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
		EncMessage.MShostname = "\n" + config.GetLangs("ms-nhn")
	} else {
		EncMessage.MShostname = "\n" + config.GetLangs("ms-hn") + name
	}

	ifas, err := net.Interfaces()
	if err != nil {
		EncMessage.MShostname += "\n-  " + config.GetLangs("ms-carrier")
	}
	if err == nil {
		var as []string
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				as = append(as, a)
			}
		}
		EncMessage.MShostname += "\n" + config.GetLangs("ms-mi")
		for i, s := range as {
			EncMessage.MShostname += "\n- " + strconv.Itoa(i) + " : " + s
		}
		addrs, _ := net.InterfaceAddrs()
		EncMessage.MShostname += "\n" + config.GetLangs("ms-ad")
		for _, addr := range addrs {
			EncMessage.MShostname += "\n- " + addr.String()
		}
	}
	EncMessage.MSalias = alias
	EncMessage.MSnodeuuid = "\n" + config.GetLangs("ms-ni") + config.GetNodeUUID()
	iduuid := uuid.New().String()
	EncMessage.MSiduuid = "\n" + config.GetLangs("ms-msg") + iduuid
	EncMessage.MSdate = "\n" + config.GetLangs("ms-on") + time.Now().Format(time.UnixDate)
	EncMessage.MSmessage = m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	if jsonerr != nil {
		log.Println("FormatMessage ", jsonerr)
	}
	//log.Println("jsonmsg ", string(jsonmsg))
	ejson := config.Encrypt(string(jsonmsg), config.GetQueuePassword())
	//log.Println("ejson ", string(ejson))
	NC, err := nats.Connect(config.GetServer(), nats.UserInfo(config.NatsUser, config.NatsUserPassword), nats.Secure(&TLS))
	if err != nil {
		fmt.Println("Send " + config.GetLangs("ls-err7") + err.Error())
	}
	JS, err := NC.JetStream()
	if err != nil {
		fmt.Println("Send " + config.GetLangs("ls-err7") + err.Error() + <-JS.StreamNames())
	}
	_, errp := JS.Publish(strings.ToLower(config.GetQueue())+".logger", []byte(ejson))
	if errp != nil {
		return true
	}
	NC.Drain()

	return false
}

// thread for receiving messages
func Receive() {

	docerts()

	config.SetReceivingMessages(true)
	for {
		select {
		case <-QuitReceive:
			return
		default:
			NatsMessages = nil
			nc, err := nats.Connect(config.GetServer(), nats.UserInfo(config.NatsUser, config.NatsUserPassword), nats.Secure(&TLS))
			if err != nil {
				if config.GetMessageWindow() != nil {
					config.GetMessageWindow().SetTitle(config.GetLangs("ms-carrier") + err.Error())
				}

			}
			js, err := nc.JetStream()
			if err != nil {
				if config.GetMessageWindow() != nil {
					config.GetMessageWindow().SetTitle(config.GetLangs("ms-carrier") + err.Error())
				}
			}
			js.AddStream(&nats.StreamConfig{
				Name:     config.GetQueue() + config.GetNodeUUID(),
				Subjects: []string{strings.ToLower(config.GetQueue()) + ".>"},
			})
			var duration time.Duration = 604800000000
			_, err1 := js.AddConsumer(config.GetQueue(), &nats.ConsumerConfig{
				Durable:           config.GetNodeUUID(),
				AckPolicy:         nats.AckExplicitPolicy,
				InactiveThreshold: duration,
				DeliverPolicy:     nats.DeliverAllPolicy,
				ReplayPolicy:      nats.ReplayInstantPolicy,
			})
			if err1 != nil {
				//log.Println(err1.Error())
				if config.GetMessageWindow() != nil {
					config.GetMessageWindow().SetTitle(config.GetLangs("ms-carrier") + err1.Error())
				}
			}
			sub, errsub := js.PullSubscribe("", "", nats.BindStream(config.GetQueue()))
			if errsub != nil {
				//log.Println(errsub.Error())
				if config.GetMessageWindow() != nil {
					config.GetMessageWindow().SetTitle(config.GetLangs("ms-carrier") + errsub.Error())
				}
			}
			msgs, err := sub.Fetch(100)
			if err != nil {
				//log.Println(err.Error())
				if config.GetMessageWindow() != nil {
					config.GetMessageWindow().SetTitle(config.GetLangs("ms-carrier") + err.Error())
				}
			}
			config.SetClearMessageDetail(true)
			if len(msgs) > 0 {
				for i := 0; i < len(msgs); i++ {
					handleMessage(msgs[i])
					msgs[i].Nak()
				}
			}
			if config.GetMessageWindow() != nil {
				config.GetMessageWindow().SetTitle(config.GetLangs("ms-err6-1") + strconv.Itoa(len(msgs)) + config.GetLangs("ms-err6-2"))
			}
			nc.Close()
			time.Sleep(30 * time.Second)
		}
	}
}

// decrypt payload
func handleMessage(m *nats.Msg) string {

	ms := MessageStore{}
	//ejson := config.Decrypt(string(m.Data), config.GetQueuePassword())
	//log.Println("m.data ", m.Data)
	var ejson = string(config.Decrypt(string(m.Data), config.GetQueuePassword()))

	err1 := json.Unmarshal([]byte(ejson), &ms)
	if err1 != nil {
		//ejson = config.GetLangs("ms-unk")
		//log.Println("ejson ", ejson)
	}
	if config.GetFilter() {
		if strings.Contains(ms.MSmessage, config.GetLangs("ms-con")) {
			return ""
		}
		if strings.Contains(ms.MSmessage, config.GetLangs("ms-dis")) {
			return ""
		}
	}
	NatsMessages = append(NatsMessages, ms)

	return ms.MSiduuid
}

// security erase jetstream data
func Erase() {

	docerts()
	//log.Println(config.GetLangs("ms-era"))

	nc, err := nats.Connect(config.GetServer(), nats.UserInfo(config.NatsUser, config.NatsUserPassword), nats.Secure(&TLS))
	if err != nil {
		log.Println("Erase Connect", config.GetLangs("ms-erac"), err.Error())
	}
	js, err := nc.JetStream()
	if err != nil {
		log.Println("Erase Jetstream Make ", config.GetLangs("ms-eraj"), err)
	}

	NatsMessages = nil
	err1 := js.PurgeStream(config.GetQueue())
	if err1 != nil {
		log.Println("Erase Jetstream Purge", config.GetLangs("ms-dels"), err1)
	}
	err2 := js.DeleteStream(config.GetQueue())
	if err2 != nil {
		log.Println("Erase Jetstream Delete", config.GetLangs("ms-dels"), err1)
	}
	msgmaxage, _ := time.ParseDuration(config.GetMsgMaxAge())
	js1, err3 := js.AddStream(&nats.StreamConfig{
		Name:     config.GetQueue(),
		Subjects: []string{strings.ToLower(config.GetQueue()) + ".>"},
		Storage:  nats.FileStorage,
		MaxAge:   msgmaxage,
	})
	if err3 != nil {
		log.Println("Erase Addstream ", config.GetLangs("ms-adds"), err3)
	}
	fmt.Printf("js1: %v\n", js1)

	Send(config.GetLangs("ms-sece"), config.GetAlias())
	nc.Close()
}
