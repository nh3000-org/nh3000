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

// A Go file encryption/decryption client for the NH3000 messaging system (https://newhorizons3000.org).

package main

import (
	"bufio"

	"flag"
	"fmt"
	"io"
	"log"

	"os"

	"strings"
	"time"

	"github.com/nh3000-org/nh3000/nhcrypt"

)


var idcount int
var MyLogLang string




// eng esp cmn
var MyLangs = map[string]string{
	"eng-mn-alias": "Intrusion Detection",
	"spa-mn-alias": "Detección de Intrusos",
	"eng-mn-lc":    "Log Connection",
	"spa-mn-lc":    "Conexión de Registro",
	"eng-fm-nhn":   "No Host Name",
	"spa-fm-nhn":   "Sin Nombre de Host",
	"eng-fm-hn":    "Host Name",
	"spa-fm-hn":    "Nombre de Host",
	"eng-fm-mi":    "Mac Ids",
	"spa-fm-mi":    "Identificadores de Mac",
	"eng-fm-ad":    "Address",
	"spa-fm-ad":    "Direccion",
	"eng-fm-ni":    "Node Id",
	"spa-fm-ni":    "Identificación del Nodo",
	"eng-fm-msg":   "Message Id",
	"spa-fm-msg":   "Identificación del mensaje",
	"eng-fm-on":    "On",
	"spa-fm-on":    "En",
	"eng-fm-fm":    "Format Message",
	"spa-fm-fm":    "Dar Formato al Mensaje",
	"eng-fm-con":   "Connection ",
	"spa-fm-con":   "Conexión ",
	"eng-fm-js":    "Jet Stream ",
	"spa-fm-js":    "Corriente en Chorro ",
}

// return translation strings
func GetLangs(mystring string) string {
	value, err := MyLangs[MyLogLang+"-"+mystring]
	if err == false {
		return "eng"
	}
	return value
}

// main loop for receiving pipe
func main() {
	var lang = "en"
	if strings.HasPrefix(os.Getenv("LANG"), "en") {
		lang = "eng"
	}
	if strings.HasPrefix(os.Getenv("LANG"), "sp") {
		lang = "spa"
	}
	logLang := flag.String("loglang", lang, "NATS Language to Use eng esp")
	logAlias := flag.String("logalias", "Intrusion", "NATS Logging Alias")
	logPattern := flag.String("logpattern", "[ERR]", "Log Pattern to Identify")
	//CA := flag.String("ca", "./ca.pem", "Path to TLS CA Certificate Authority")
	//ClientCert := flag.String("clientcert", "./clientcert.pem", "Path to TLS Client Cert")
	//ClientKey := flag.String("clientkey", "./clientkey.pem", "Path to TLS Client Key")
	ServerIP := flag.String("serverip", "nats://127.0.0.1:4222", "Server IP or DNS Name")
	flag.Parse()
	fmt.Println("Usage:")
	fmt.Println("tail -f log.file | log -serverip nats://?.?.?.? -logpattern ??? -logalias ????")
	fmt.Println("LOGPATTERN will be logged to the nats server using the alias")
	fmt.Println("fully encrypted")
	fmt.Println("")
	fmt.Println("Run Options:")
	fmt.Println("-loglang: ", *logLang)
	MyLogLang = *logLang


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
				clientcert, err := tls.X509KeyPair([]byte(LogClientcert), []byte(LogClientkey))
				if err != nil {
					log.Println("nhnats.go clientcert " + err.Error())
				}
				rootCAs, _ := x509.SystemCertPool()
				if rootCAs == nil {
					rootCAs = x509.NewCertPool()
				}
				ok := rootCAs.AppendCertsFromPEM([]byte(LogCaroot))
				if !ok {
					log.Println("nhnats.go rootCAs")
				}
				tlsConfig := &tls.Config{
					RootCAs:      rootCAs,
					Certificates: []tls.Certificate{clientcert},
				}
				nc, err := nats.Connect(*ServerIP, nats.Secure(tlsConfig))
				if err != nil {
					log.Println(GetLangs("mn-con"), err.Error())
				}
				js, errjs := nc.JetStream()
				if errjs != nil {
					log.Println(GetLangs("mn-js"), errjs.Error())
				}
				_, jserr := js.Publish("messages.log", []byte(Send(string(buf))))
				if jserr != nil {
					log.Println(GetLangs("mn-js"), jserr)
				}
			}
		}
		if err != nil && err != io.EOF {
			log.Println("log.go ", err)
		}
	}
}
