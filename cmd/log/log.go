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
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
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
	"github.com/nats-io/nats.go"
)

var MyBytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05} // must be 16 bytes
var idcount int
var MyLogLang string
var MyLogAlias string
var LogCaroot = "-----BEGIN CERTIFICATE-----\nMIICFDCCAbugAwIBAgIUDkHxHO1DwrlkTzUimG5PoiswB6swCgYIKoZIzj0EAwIw\nZjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UE\nChMDU0VDMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMz\nMDAwLm9yZzAgFw0yMzAzMzExNzI5MDBaGA8yMDUzMDMyMzE3MjkwMFowZjELMAkG\nA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UEChMDU0VD\nMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9y\nZzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHXwMUfMXiJix3tuzFymcA+3RkeY\nZE7urUzVgaqkv/Oef3jhqhtf1XzK/qVYGxWWmpvADGB252PG1Mp7Z5wmzqyjRTBD\nMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEBMB0GA1UdDgQWBBQm\nFA5caanuqxGFOf9DtZkVYv5dCzAKBggqhkjOPQQDAgNHADBEAiB3BheNP4XdBZ27\nxVBQ7ztMJqK7wDi1V3LuMy5jmXr7rQIgHCse0oaiAwcl4VwF00aSshlV+T/da0Tx\n1ANkaM+rie4=\n-----END CERTIFICATE-----\n"
var LogClientcert = "-----BEGIN CERTIFICATE-----\nMIIDUzCCAvigAwIBAgIUUyhlJt8mp1XApRbSkdrUS55LGV8wCgYIKoZIzj0EAwIw\nZjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UE\nChMDU0VDMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMz\nMDAwLm9yZzAeFw0yMzAzMzExNzI5MDBaFw0yODAzMjkxNzI5MDBaMHIxCzAJBgNV\nBAYTAlVTMRAwDgYDVQQIEwdGbG9yaWRhMRIwEAYDVQQHEwlDcmVzdHZpZXcxGjAY\nBgNVBAoTEU5ldyBIb3Jpem9ucyAzMDAwMSEwHwYDVQQLExhuYXRzLm5ld2hvcml6\nb25zMzAwMC5vcmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDFttVH\nQ131JYwazAQMm0XAQvRvTjTjOY3aei1++mmQ+NQ9mrOFk6HlZFoKqsy6+HPXsB9x\nQbWlYvUOuqBgb9xFQZoL8jiKskLLrXoIxUAlIBTlyf76r4SV+ZpxJYoGzXNTedaU\n0EMTyAiUQ6nBbFMXiehN5q8VzxtTESk7QguGdAUYXYsCmYBvQtBXoFYO5CHyhPqu\nOZh7PxRAruYypEWVFBA+29+pwVeaRHzpfd/gKLY4j2paInFn7RidYUTqRH97BjdR\nSZpOJH6fD7bI4L09pnFtII5pAARSX1DntS0nWIWhYYI9use9Hi/B2DRQLcDSy1G4\n0t1z4cdyjXxbFENTAgMBAAGjgawwgakwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQM\nMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFAzgPVB2/sfT7R0U\ne3iXRSvUkfoQMB8GA1UdIwQYMBaAFCYUDlxpqe6rEYU5/0O1mRVi/l0LMDQGA1Ud\nEQQtMCuCGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZ4IJMTI3LDAsMCwxhwTAqABn\nMAoGCCqGSM49BAMCA0kAMEYCIQCDlUH2j69mJ4MeKvI8noOmvLHfvP4qMy5nFW2F\nPT5UxgIhAL6pHFyEbANtSkcVJqxTyKE4GTXcHc4DB43Z1F7VxSJj\n-----END CERTIFICATE-----\n"
var LogClientkey = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAxbbVR0Nd9SWMGswEDJtFwEL0b0404zmN2notfvppkPjUPZqz\nhZOh5WRaCqrMuvhz17AfcUG1pWL1DrqgYG/cRUGaC/I4irJCy616CMVAJSAU5cn+\n+q+ElfmacSWKBs1zU3nWlNBDE8gIlEOpwWxTF4noTeavFc8bUxEpO0ILhnQFGF2L\nApmAb0LQV6BWDuQh8oT6rjmYez8UQK7mMqRFlRQQPtvfqcFXmkR86X3f4Ci2OI9q\nWiJxZ+0YnWFE6kR/ewY3UUmaTiR+nw+2yOC9PaZxbSCOaQAEUl9Q57UtJ1iFoWGC\nPbrHvR4vwdg0UC3A0stRuNLdc+HHco18WxRDUwIDAQABAoIBACe0XMZP4Al//c/P\n0qxZbjt69q13jiVnhHYwfPx3+0UywySP8adMi4GOkop73Ftb05+n7diHspvA8KeB\nkP1s2VZLI01s2i/4NnPCpbQnMIeEFs5Cr2LWZpDbrEk2ma5eCd/kotQFssLBM//a\nSrfeMh2TA0TJo7WEft9Cnf4ZeEkKnycplfvwTyv286iFZCYo2dv66BfTej6kkVCo\nAi+ZVCe2zSqRYyr0u4/j/kE3b3eSkCnY2IVcqlP7epuEGVOZyxeFLwM5ljbWL816\npA6WIJgQo2EQ1N7L531neg5WjXQ/UwTQoXP1jvuuVtKtOBFqm1IshEyFk3WpsfpD\nr16OTdECgYEA6FB6NYxYtnWPaIYAOqP7GtMKoJujH8MtZy6J33LkxI7nPkMkn8Mv\nva32tvjU4Bu1FVNp9k5guC+b+8ixXK0URj25IOhDs6K57tck22W9WiTZlmnkCO01\nJOavrelWbvYt5xNWIdnPualoPfGB0iJKXsKY/bpH4eVfhWwpNPI5sMkCgYEA2d9G\nEPuWN6gUjZ+JfdS+0WHK1yGD7thXs7MPUlhGqDzBryh2dkywyo8U8+tMLuDok1RZ\njnT3PYkLQEpzoV0qBkpFFShL6ubaGmDz1UZsozl0YcIg4diZeuPHnIAeXOFrhgYf\n825163LmT3jYHCROFEMLtTYyIQP0EznE+qFT3TsCgYEApgtvbfqkJbWdDL5KR5+R\nCLky7VyQmVEtkIRI8zbxoDPrwCrJcI9X/iDrKBhuPshPA7EdGXkn1D3jJXFqo6zp\nwtK3EXgxe6Ghd766jz4Guvl/s+x3mpHA3GEtzAXtS14VrQW7GHLP8AnPggauHX14\n3oYER8XvPtxtC7YlNbyz01ECgYAe2b7SKM3ck7BVXYHaj4V1oKNYUyaba4b/qxtA\nTb+zkubaJqCfn7xo8lnFMExZVv+X3RnRUj6wN/ef4ur8rnSE739Yv5wAZy/7DD96\ns74uXrRcI2EEmechv59ESeACxuiy0as0jS+lZ1+1YSc41Os5c0T1I/d1NVoaXtPF\nqZJ2gQKBgBp/XavdULBPzC7B8tblySzmL01qJZV7MSSVo2/1vJ7gPM0nQPZdTDog\nTfA5QKSX9vFTGC9CZHSJ+fabYDDd6+3UNYUKINfr+kwu9C2cysbiPaM3H27WR5mW\n5LhStAfwuRRYBDsG2ndjraxcBrrPdtkbS0dpeQUDJxvkMIuLHnhQ\n-----END RSA PRIVATE KEY-----\n"
var LogQueuePassword = "123456789012345678901234"

type MessageStore struct {
	MSiduuid   string
	MSalias    string
	MShostname string
	MSipadrs   string
	MSmessage  string
	MSnodeuuid string
	MSdate     string
}

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

// encode string to base64
func Encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// encrypt string
func Encrypt(text string, MySecret string) (string, error) {
	block, err := aes.NewCipher([]byte(MySecret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, MyBytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)

	return Encode(cipherText), nil
}

// send message to nats
func Send(m string) []byte {
	EncMessage := MessageStore{}
	name, err := os.Hostname()
	if err != nil {
		EncMessage.MShostname = "\n" + GetLangs("fm-nhn")
	} else {
		EncMessage.MShostname = "\n" + GetLangs("fm-hn") + " - " + name
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
		EncMessage.MShostname += "\n" + GetLangs("fm-mi")
		for i, s := range as {
			EncMessage.MShostname += "\n- " + strconv.Itoa(i) + " : " + s
		}
		addrs, _ := net.InterfaceAddrs()
		EncMessage.MShostname += "\n" + GetLangs("fm-ad")
		for _, addr := range addrs {
			EncMessage.MShostname += "\n- " + addr.String()
		}
	}
	EncMessage.MSalias = MyLogAlias
	idcount++
	EncMessage.MSnodeuuid = "\n" + GetLangs("fm-ni") + " - " + strconv.Itoa(idcount)
	iduuid := uuid.New().String()
	EncMessage.MSiduuid = "\n" + GetLangs("fm-msg") + " - " + iduuid[0:8]
	EncMessage.MSdate = "\n" + GetLangs("fm-on") + " -" + time.Now().Format(time.UnixDate)
	EncMessage.MSmessage = m
	jsonmsg, jsonerr := json.Marshal(EncMessage)
	if jsonerr != nil {
		log.Println(GetLangs("fm-fm"), jsonerr)
	}
	ejson, _ := Encrypt(string(jsonmsg), LogQueuePassword)

	return []byte(ejson)
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
	CA := flag.String("ca", "./ca.pem", "Path to TLS CA Certificate Authority")
	ClientCert := flag.String("clientcert", "./clientcert.pem", "Path to TLS Client Cert")
	ClientKey := flag.String("clientkey", "./clientkey.pem", "Path to TLS Client Key")
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
	fmt.Println("logAlias", *logAlias)
	MyLogAlias = *logAlias
	fmt.Println("serverip", *ServerIP)
	fmt.Println("-logpattern: ", *logPattern)
	fmt.Println("-ca: ", *CA)
	fmt.Println("-clientkey: ", *ClientKey)
	fmt.Println("-clientcert", *ClientCert)

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
