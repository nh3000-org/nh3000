// Copyright 2012-2023 The NH3000 Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A Go monitoring client for the NH3000 messaging system (https://newhorizons3000.org).

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nh3000-org/nh3000/config"
)

var idcount int
var MyLogLang string
var MyLogAlias string

type MessageStore struct {
	MSiduuid   string
	MSalias    string
	MShostname string
	MSipadrs   string
	MSmessage  string
	MSnodeuuid string
	MSdate     string
}

// eng esp cmn hin
var MyLogLangs = map[string]string{
	"eng-ls-alias":     "Alias",
	"spa-ls-alias":     "Alias",
	"hin-ls-alias":     "उपनाम",
	"eng-ls-queue":     "Queue",
	"spa-ls-queue":     "Cola",
	"hin-ls-queue":     "कतार",
	"eng-ls-queuepass": "Queue Password",
	"spa-ls-queuepass": "Contraseña de Cola",
	"hin-ls-queuepass": "कतार पासवर्ड",
	"eng-ls-trypass":   "Try Password",
	"spa-ls-trypass":   "Probar Contraseña",
	"hin-ls-trypass":   "पासवर्ड आज़माएं",
	"eng-ls-con":       "Connected",
	"spa-ls-con":       "Conectada",
	"hin-ls-con":       "जुड़े हुए",
	"eng-ls-dis":       "Disconnected",
	"spa-ls-dis":       "Desconectada",
	"hin-ls-dis":       "डिस्कनेक्ट किया गया",
	"eng-ls-err1":      "Error Creating Password Hash 24",
	"spa-ls-err1":      "Error al Crear el Hash de la Contraseña 24",
	"hin-ls-err1":      "पासवर्ड हैश 24 बनाने में त्रुटि",
	"eng-ls-err2":      "Error Loading Password Hash 24",
	"spa-ls-err2":      "Error al Cargar el Hash de la Contraseña 24",
	"hin-ls-err2":      "पासवर्ड हैश 24 लोड करने में त्रुटि",
	"eng-ls-err3":      "Error Invalid Password",
	"spa-ls-err3":      "Error Contraseña no Válida",
	"hin-ls-err3":      "त्रुटि अमान्य पासवर्ड",
	"eng-ls-err4":      "Error URL Incorrect Format",
	"spa-ls-err4":      "URL de Error Formato Incorrecto",
	"hin-ls-err4":      "त्रुटि यूआरएल गलत प्रारूप",
	"eng-ls-err5":      "Error Invalid Queue Password 24",
	"spa-ls-err5":      "Error Contraseña de Cola no Válida 24",
	"hin-ls-err5":      "त्रुटि अमान्य कतार पासवर्ड 24",
	"eng-ls-err6-1":    "Error Queue Password Length is ",
	"spa-ls-err6-1":    "La Longitud de la Contraseña de la Cola de Errores es ",
	"hin-ls-err6-1":    "त्रुटि कतार पासवर्ड की लंबाई है ",
	"eng-ls-err6-2":    " should be length of 24",
	"spa-ls-err6-2":    " Debe Tener una Longitud de 24",
	"hin-ls-err6-2":    " लंबाई 24 होनी चाहिए",
	"eng-ls-err7":      "No NATS connection",
	"spa-ls-err7":      "Sin Conexión NATS",
	"hin-ls-err7":      "कोई NATS कनेक्शन नहीं",
	"eng-ls-erase":     "Security Erase",
	"spa-ls-erase":     "Borrado de seguridad",
	"hin-ls-erase":     "सुरक्षा मिटाएँ",
	"eng-ls-clogon":    "Communications Logon",
	"spa-ls-clogon":    "Inicio de Sesión de Comunicaciones",
	"hin-ls-clogon":    "संचार लॉगऑन",
	"eng-ls-err8":      "No JETSTREAM Connection",
	"spa-ls-err8":      "Sin Conexión JETSTREAM ",
	"hin-ls-err8":      "कोई जेटस्ट्रीम कनेक्शन नहीं",
}

// return translation strings
func GetLogLangs(mystring string) string {
	value, err := MyLogLangs[MyLogLang+"-"+mystring]
	if !err {
		return "xxx"
	}
	return value
}

// send message to nats
func Send(m string) []byte {
	EncMessage := MessageStore{}
	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = "\n" + GetLogLangs("fm-nhn")
	} else {
		EncMessage.MShostname = "\n" + GetLogLangs("fm-hn") + " - " + name
	}
	ifas, err := net.Interfaces()
	if err == nil {
		var as []string
		for _, ifa := range ifas {
			a := ifa.HardwareAddr.String()
			if a != "" {
				as = append(as, a)
			}
		}
		EncMessage.MShostname += "\n" + GetLogLangs("fm-mi")
		for i, s := range as {
			EncMessage.MShostname += "\n- " + strconv.Itoa(i) + " : " + s
		}
		addrs, _ := net.InterfaceAddrs()
		EncMessage.MShostname += "\n" + GetLogLangs("fm-ad")
		for _, addr := range addrs {
			EncMessage.MShostname += "\n- " + addr.String()
		}
	}
	EncMessage.MSalias = MyLogAlias
	idcount++
	EncMessage.MSnodeuuid = "\n" + GetLogLangs("fm-ni") + " - " + strconv.Itoa(idcount)
	iduuid := uuid.New().String()
	EncMessage.MSiduuid = "\n" + GetLogLangs("fm-msg") + " - " + iduuid[0:8]
	EncMessage.MSdate = "\n" + GetLogLangs("fm-on") + " -" + time.Now().Format(time.UnixDate)
	EncMessage.MSmessage = m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	if jsonerr != nil {
		log.Println(GetLogLangs("fm-fm"), jsonerr)
	}
	ejson := config.Encrypt(string(jsonmsg), config.NatsQueuePassword)

	return []byte(ejson)
}

// main loop for receiving pipe
func main() {
	MyLogLang = "eng"
	if strings.HasPrefix(os.Getenv("LANG"), "en") {
		MyLogLang = "eng"
	}
	if strings.HasPrefix(os.Getenv("LANG"), "sp") {
		MyLogLang = "spa"
	}
	logLang := flag.String("loglang", MyLogLang, GetLogLangs("fl-ll"))
	logAlias := flag.String("logalias", "LOGALIAS", GetLogLangs("fl-la"))

	logPattern := flag.String("logpattern", "[ERR]", GetLogLangs("fl-lp"))
	ServerIP := flag.String("serverip", config.NatsServer, GetLogLangs("fl-si"))
	flag.Parse()
	MyLogAlias = *logAlias
	fmt.Println("====================================================== ")
	fmt.Println("EX: tail -f log.file | nhlog ", " -loglang ", *logLang, " -serverip ", *ServerIP, " -logpattern ", *logPattern, " -logalias ", *logAlias)
	fmt.Println("- serverip - NATS nats://xxxxx.yyy:port")
	fmt.Println(" -logalias - make unique for each instance, become DEVICE.device in NATS")
	fmt.Println("====================================================== ")
	
	config.CheckDevice(MyLogAlias)
	r := bufio.NewReader(os.Stdin)
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				time.Sleep(time.Minute)
			}
		}

		if int64(len(buf)) != 0 {
			if strings.Contains(string(buf), *logPattern) {
				config.Send("EVENTS", "events", string(buf), "[logger]" + MyLogAlias)
			}
		}
		if err != nil && err != io.EOF {
			log.Println("log.go ", err)
		}
	}
}
