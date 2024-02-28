package config

import (
	"log"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

var preferedlanguage string
var win fyne.Window
var app fyne.App
var loggedon bool
var receivingmessages bool
var nodeuuid string
var alias string
var server string
var queue string
var queuepassword string
var caroot string
var clientcert string
var clientkey string
var msgmaxage string
var filter bool
var clearmessagedetail bool

func SetClearMessageDetail(a bool) {
	clearmessagedetail = a
}
func GetClearMessageDetail() bool {
	return clearmessagedetail
}

func SetFilter(a bool) {
	filter = a
}
func GetFilter() bool {
	return filter
}

func DataStore(file string) fyne.URI {
	DataLocation, dlerr := storage.Child(GetApp().Storage().RootURI(), file)
	if dlerr != nil {
		log.Println("DataStore error ", dlerr)
	}
	return DataLocation
}
func SetMsgMaxAge(a string) {
	msgmaxage = a
}
func GetMsgMaxAge() string {
	return msgmaxage
}

func SetReceivingMessages(b bool) {
	receivingmessages = b
}
func GetReceivingMessages() bool {
	return receivingmessages
}
func SetPreferedLanguage(a string) {
	preferedlanguage = a
}
func GetPreferedLanguage() string {
	return preferedlanguage
}

func SetClientKey(a string) {
	clientkey = a
}
func GetClientKey() string {
	return clientkey
}
func SetClientCert(a string) {
	clientcert = a
}
func GetClientCert() string {
	return clientcert
}
func SetCaroot(a string) {
	caroot = a
}
func GetCaroot() string {
	return caroot
}
func SetQueuePassword(a string) {
	queuepassword = a
}

func GetQueuePassword() string {
	return queuepassword
}
func SetQueue(a string) {
	queue = a
}

func GetQueue() string {
	return queue
}
func SetServer(a string) {
	server = a
}

func GetServer() string {
	return server
}

func SetAlias(a string) {
	alias = a
}

func GetAlias() string {
	return alias
}

func SetNodeUUID(n string) {
	nodeuuid = n
}

func GetNodeUUID() string {
	return nodeuuid
}

func SetMessageWindow(w fyne.Window) {
	win = w
}

func GetMessageWindow() fyne.Window {
	return win
}

func SetApp(a fyne.App) {
	app = a
}

func GetApp() fyne.App {
	return app
}

func SetLoggedOn() {
	loggedon = true
}

func GetLoggedOn() bool {
	return loggedon
}
func Edit(action string, value string) bool {
	if action == "cvtbool" {
		if value == "True" {
			return true
		}
		if value == "False" {
			return false
		}

	}
	if action == "QUEUEPASSWORD" {
		if len(value) == 0 {
			return true
		}
		if len(value) != 24 {
			return true
		}
		return false
	}
	if action == "URL" {
		valid := strings.Contains(strings.ToLower(value), "nats://")
		if !valid {
			return true
		}
		valid2 := strings.Contains(value, ".")
		if !valid2 {
			return true
		}
		valid3 := strings.Contains(value, ":")
		if !valid3  {
			return true
		}

		return false
	}
	if action == "STRING" {
		return len(value) == 0 
	}

	if action == "CERTIFICATE" {
		valid := strings.Contains(value, "-----BEGIN CERTIFICATE-----")
		if !valid  {
			return false
		}
		valid2 := strings.Contains(value, "-----END CERTIFICATE-----")
		if !valid2  {
			return false
		}
	}
	if action == "KEY" {

		valid := strings.Contains(value, "-----BEGIN RSA PRIVATE KEY-----")
		if !valid  {
			return false
		}
		valid2 := strings.Contains(value, "-----END RSA PRIVATE KEY-----")
		if !valid2  {
			return false
		}
	}
	if action == "TRUEFALSE" {
		valid := strings.Contains(value, "True")
		if valid == false {
			valid2 := strings.Contains(value, "False")
			if valid2 == false {
				return false
			}
		}
	}
	return true
}
func ParseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}
