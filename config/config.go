package config

import (
	"errors"
	"log"
	"net/url"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

/*
*  The following fields need to be modified for you production
*  Environment to provide maximum security
*
*  These fields are meant to be distributed at compile time and
*  editable in the gui.
*
 */

var NatsServer = "nats://nats.newhorizons3000.org:4222"
var NatsCaroot = "-----BEGIN CERTIFICATE-----\nMIID7zCCAtegAwIBAgIUaXAPxJvZRRdTq5RWlwxs1XYo+5kwDQYJKoZIhvcNAQEL\nBQAwgYAxCzAJBgNVBAYTAlVTMRAwDgYDVQQIEwdGbG9yaWRhMRIwEAYDVQQHEwlD\ncmVzdHZpZXcxGjAYBgNVBAoTEU5ldyBIb3Jpem9ucyAzMDAwMQwwCgYDVQQLEwNX\nV1cxITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZzAeFw0yMzEyMTkw\nMzA4MDBaFw0yODEyMTcwMzA4MDBaMIGAMQswCQYDVQQGEwJVUzEQMA4GA1UECBMH\nRmxvcmlkYTESMBAGA1UEBxMJQ3Jlc3R2aWV3MRowGAYDVQQKExFOZXcgSG9yaXpv\nbnMgMzAwMDEMMAoGA1UECxMDV1dXMSEwHwYDVQQDExhuYXRzLm5ld2hvcml6b25z\nMzAwMC5vcmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCrVIXA/SxU\n7GeW92UNyiPnQEZgbJIHHQ31AQE2C/vFdpEtv32uoX1SsDl5drWvBrMnd5zrw1tL\nOEPA26tk/ACfQYL0n0HfeutLLu8H9jUWNp8ziX6Qbgd01M+/BixobHQjyDMxulo4\nJU2VK6QBLs9VI6TIihEU2BZhc/XCD9QbWcikAif1JySpz93MjFv3pcQU8ci4vQ0T\nImaGnHesr1qDbX1NuFVuBOPavZ64sQ1RsZtH5CdD+RU772wQWUgkPkwyUn8QBwTS\ne9XV5DNQD5nGEXjKTgjrd9KRf9pmRDnf6gBLi2r6C/l6q2w3ItOOHARdK0mc9CYh\ngY1Nzl59vrWdAgMBAAGjXzBdMA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTAD\nAQH/MB0GA1UdDgQWBBR0qq9ueABC5RDsg/02FZFpBOR1hDAbBgNVHREEFDAShwTA\nqAAFhwTAqFjohwR/AAABMA0GCSqGSIb3DQEBCwUAA4IBAQBfdX0IMya9Dh9dHLJj\nnJZyb96htMWD5nuQQVBAu3ay+8O2GWj5mlsLJXAP2y7p/+3gyvHKTRDdJLux7N79\nHn6AYjmp3PCyZzuL1M/kHhSQxhxqJHGwjGXILt5pLovVkvkl4iukdxWJ5HAPsUGY\nO3QSDDFdoLflsG5VcrtdODm8uyxAjhMPAR2PXKfX8ABI79N7VKcbb98338fifrN8\n9H1r3BXcdsyhpH0gB0ZKJFSpMGWXlfudFEe9mXI9898xbEI2znqlYGhboVsuv5LM\nRESH2zXrkhmZyHqw0RtDROzyZOy5g1LcxbtVMn4w1LI4h3MDuE9B+Vud77A48qtA\ny+5x\n-----END CERTIFICATE-----\n"
var NatsClientcert = "-----BEGIN CERTIFICATE-----\nMIIEMTCCAxmgAwIBAgIUB7+OFX1LQrWtYMl5XIOXsOaLac0wDQYJKoZIhvcNAQEL\nBQAwgYAxCzAJBgNVBAYTAlVTMRAwDgYDVQQIEwdGbG9yaWRhMRIwEAYDVQQHEwlD\ncmVzdHZpZXcxGjAYBgNVBAoTEU5ldyBIb3Jpem9ucyAzMDAwMQwwCgYDVQQLEwNX\nV1cxITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZzAgFw0yMzEyMTkw\nMzA4MDBaGA8yMDUzMTIxMTAzMDgwMFowcjELMAkGA1UEBhMCVVMxEDAOBgNVBAgT\nB0Zsb3JpZGExEjAQBgNVBAcTCUNyZXN0dmlldzEaMBgGA1UEChMRTmV3IEhvcml6\nb25zIDMwMDAxITAfBgNVBAsTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZzCCASIw\nDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAMWARyniHy8r342e3aKSsLDPwVMC\n2mRwuILP2JkXp5FllaFKnu/Z+0mF+iQlSchcC6DOcMQk00Cp/I8cCP865zyxPhqN\n2F2/qVItCU4+PTwe6ZnrfpJgXWwyk1hjS3vVNTT+idI5+pJgFH9YL0lbJ7q1UyPB\n+KP0x/c5T3K2Ec6U4uXhbVt/ePxFmsl1sHw6FE//XrA4EzbqCMEPCTcOfInvFrCJ\ny4/pAqjCxegT/1YDMNEdzmG8vg2tc3jPV+3GIAV3YL5nDE5mprHPEEDJtNQi+E4o\nXXXMobNhrJh9KJ59VbxTF8m5yM3b8fvof97OYhK0KYggplnTH+bhnYU9V5ECAwEA\nAaOBrTCBqjAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAwwCgYIKwYBBQUHAwIwDAYD\nVR0TAQH/BAIwADAdBgNVHQ4EFgQUpffi3LSreerO756B/VnZkyyEVBIwHwYDVR0j\nBBgwFoAUdKqvbngAQuUQ7IP9NhWRaQTkdYQwNQYDVR0RBC4wLIIYbmF0cy5uZXdo\nb3Jpem9uczMwMDAub3JnhwR/AAABhwTAqAAFhwTAqFjoMA0GCSqGSIb3DQEBCwUA\nA4IBAQALlRqqW2HH4flFIgR/nh51gc/Hxv5xivhkzWUHHXRdltECSXknI4yBPchQ\n6Zsy0HZ7XQRlhQSIYd4Bp6eyHbny5t3JA978dHzpGJFCUVQDMY4yHLaCQgFJ+ESn\nwyyDWTRGA3cpEikL0B0ekDfqjWUEMTzmT/gnoSl0vM69nZDLZm1xMx1+EH+bpfFB\nRaVM6gKSAuFJmNYEL2e7JSags+3IHyVHkdo8GDlY//71Z4lxsFxFCF6xF9GDdAr2\niCA4OfydjiBSOz0eLJVgqkk1KGXtMqZXAojX62NrIWnFTW1Vzd46ekOHhq93B3tA\nkjWmHY/KdCZUjQSWss+YXgG4mI8c\n-----END CERTIFICATE-----\n"
var NatsClientkey = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAxYBHKeIfLyvfjZ7dopKwsM/BUwLaZHC4gs/YmRenkWWVoUqe\n79n7SYX6JCVJyFwLoM5wxCTTQKn8jxwI/zrnPLE+Go3YXb+pUi0JTj49PB7pmet+\nkmBdbDKTWGNLe9U1NP6J0jn6kmAUf1gvSVsnurVTI8H4o/TH9zlPcrYRzpTi5eFt\nW394/EWayXWwfDoUT/9esDgTNuoIwQ8JNw58ie8WsInLj+kCqMLF6BP/VgMw0R3O\nYby+Da1zeM9X7cYgBXdgvmcMTmamsc8QQMm01CL4Tihddcyhs2GsmH0onn1VvFMX\nybnIzdvx++h/3s5iErQpiCCmWdMf5uGdhT1XkQIDAQABAoIBAB+Iu9QUJqaBetBB\n7WFnyo5wnY2DhxtCZDN+vDa1cCvm7F00bOwfAeBbY/UhfwZeq/yg+aBXwOMyQQEY\nmNcnsIQgSKo0u7c8Quy8BCBaD6zpwqKw1yTH/iKocJ5MPGEpSbWMbrUCTN/SN3Od\nwO8VfuJw0TWEYw7KpqLyo5zNNUqmczEO438CPGotbkFfzUqkumeUOsGWJFongyZY\na9EwpcTH2TkxuXum9SQVyLy+hSG/AEBp0cQPaRcoNh8sWYk43y5HrkIAqFo7dkMa\n9usAVMz9JCqIH2UNV04cDASFaiDMpYoD2hV2YHlL7/CQ7v5nb6OHT2A9aoSBOAfm\ns+dBzYECgYEA1l8+T9Xux73TCbFO2p7F094xSx4hhBZhaYpvzZoNN7iQdbdUVt2l\n1yHSoRgJUJMZlnKpMoNMLCxo34Lr3ww/TkIE/rrg10pqbqvojIDLCbi103EEB2v9\nWix8MSeOgFCa72T4lg9fDm5T493n4C5dade3LzZczUBF6dgmth3D+nMCgYEA69pa\nlob9n7eNXqDPk9kZUJV1jfLATC8eN4jupEiKfjnxEz9mUewvL/RF8kFhiS1ISC50\nKgM0v+isYBwwX00c7P02L6xCoGT35qOeoutEWVy/tYIHIHsD0jUBBsdnpQVNf58l\n9DDy2hZrpUwrsVHylVHpufBgKOfxgP2Jr3qD0OsCgYEAn4vzTGfkdzSIRMZ58awJ\ngE32Ufny5+PgTDSEUXk+LSJoIbR4SM5eB2dc5BiHljhk6twboUSnBJlo1DEUa8Up\nuIzaOtvLS3BPFl9LjIaulmWqrduHLB7rSJmjNNJD9KwJI/L6MHTwQkVKmmUllmvr\nikLKS5EiMICNiCUfaptsqJECgYEApYaSqzBEUdK1oeMErAPis16hqSTkdtNexqUQ\nrzXGFP6/Rb3qJra3C1XJvVLLjEW+hAIuPsoPPFyklbNS85+gHGc9n0mrXPxfy3ur\nuzWYu4rPdSizrcUIEoBmnwZVpEhLcrUUIwQzfIHdvJ3v0DvuH4PkoD2mjy7xnJDU\nD9bRKk8CgYAqK1lY5waFR0u3eFIPnrV4ATHXYuxcup2DCF+KJ6qwc4nNI6OB/ovU\nttiVZGr1rca42+XdWUQL5ufPFuKymeLbsuVzabbGKi+4RMvL+TIuorYtJRUPF+C7\nA9jlMeckpTZvl0yn5s3lC817N27B+U0M/jGow8sO0NtjBiImuTC5dg==\n-----END RSA PRIVATE KEY-----\n"

// default user
var NatsUser = "natsoperator"
var NatsUserPassword = "hjscr44iod"

var NatsQueuePassword = "987654321098765432109876"
var NatsQueueDurable = "snatsdurable"
var NatsQueue = "MESSAGES"

// default encryption
var KeyAes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}  // must be 16 bytes
var KeyHmac = []byte{36, 45, 53, 21, 87, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05} // must be 16 bytes
const MySecret string = "abd&1*~#^2^#s0^=)^^7%c34"                                   // must be 24 characters
// default gui password
const Password = "123456" // default password shipped with app

var PasswordHash string

var preferedlanguage string
var win fyne.Window
var mainwin fyne.Window
var app fyne.App
var loggedon bool
var receivingmessages bool
var NatsNodeUUID string
var NatsAlias string

// var server string
// var queue string
// var queuepassword string
var caroot string
var clientcert string
var clientkey string
var msgmaxage string
var filter bool

// sip fields
var sipgeneralcallwaiting bool
var sipgeneralhangup3way bool
var sipgeneralincallosd bool
var sipgeneralcallhistory int

var sipsoundringtone string
var sipsoundspeaker string
var sipsoundmicrophone string
var sipsoundvalidate bool
var sipsoundossfragsize int
var sipsoundalsaplayperiodsize int
var sipsoundcaptureperiodsize int

var sipringtoneplay bool
var sipringtone string
var sipringbacktoneplay bool
var sipringbacktone string

var sipaddressbooklookup bool
var sipaddressoverideaddress bool
var sipaddresslookupphoto bool

var sipnetworksipport int
var sipnetworkrtpport int
var sipnetworkmaxudpsize int
var sipnetworkmaxtcpsize int

var siplogmaxsizemb int
var siplogdebug bool
var siplogstun bool
var siplogmemory bool

var sipuseralias string
var sipusername string
var sipdomain string
var siporganization string
var sipuserauthrealm string
var sipuserauthname string
var sipuserpassword string
var sipuserakaop string
var sipuserakaamp string

// var clearmessagedetail bool
var messagelist fyne.Widget

//var messagelistdetails int

//	func SetNatsServer(s string) {
//		server = s
//	}
//
//	func GetNatsServer() string {
//		return server
//	}
func SetMessageList(l fyne.Widget) {
	messagelist = l
}
func GetMessageList() fyne.Widget {
	return messagelist
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

//func SetQueuePassword(a string) {
//	queuepassword = a
//}

//func GetQueuePassword() string {
//	return queuepassword
//}
//func SetQueue(a string) {
//	queue = a
//}

//func GetQueue() string {
//	return queue
//}

//func SetServer(a string) {
//	server = a
//}

//func GetServer() string {
//	return server
//}

//func SetAlias(a string) {
//	alias = a
//}

//func GetAlias() string {
//	return alias
//}

//func SetNodeUUID(n string) {
//	nodeuuid = n
//}

//func GetNodeUUID() string {
//	return nodeuuid
//}

func SetMessageWindow(w fyne.Window) {
	win = w
}

func GetMessageWindow() fyne.Window {
	return win
}
func SetWindow(w fyne.Window) {
	mainwin = w
}

func GetWindow() fyne.Window {
	return mainwin
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
	if action == "FILEEXISTS" {
		_, err := os.Stat(value)
		if errors.Is(err, os.ErrNotExist) {
			return true
		}
		return false
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
		if !valid3 {
			return true
		}

		return false
	}
	if action == "STRING" {
		return len(value) == 0
	}

	if action == "CERTIFICATE" {
		valid := strings.Contains(value, "-----BEGIN CERTIFICATE-----")
		if !valid {
			return false
		}
		valid2 := strings.Contains(value, "-----END CERTIFICATE-----")
		if !valid2 {
			return false
		}
	}
	if action == "KEY" {

		valid := strings.Contains(value, "-----BEGIN RSA PRIVATE KEY-----")
		if !valid {
			return false
		}
		valid2 := strings.Contains(value, "-----END RSA PRIVATE KEY-----")
		if !valid2 {
			return false
		}
	}
	if action == "TRUEFALSE" {
		valid := strings.Contains(value, "True")
		if !valid {
			valid2 := strings.Contains(value, "False")
			if !valid2 {
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
