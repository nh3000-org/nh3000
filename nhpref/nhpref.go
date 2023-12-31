package nhpref

import (
	"log"
	//"os"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
	"github.com/nh3000-org/nh3000/nhauth"
	"github.com/nh3000-org/nh3000/nhcrypt"
	"github.com/nh3000-org/nh3000/nhutil"
)

// version
const Version = "snats-beta.3"

//var MyApp fyne.App

var LoggedOn bool = false

//var PasswordValid bool = false

var ErrorMessage = "None"

var Queue string         // server message queue
var Queuepassword string // server message queue password

var Msgmaxage string        // msg age in hours to keep
var PreferedLanguage string // language string
var Password string         // encrypt file password
var Passwordhash string     // hash value of password

var PasswordMinimumSize string        // set minimum password size
var PasswordMustContainNumber string  // password must contain number
var PasswordMustContainLetter string  // password must contain letter
var PasswordMustContainSpecial string // password must contain special character

// Server tab
var Server string // server url

// message
var IdUUID string   // unique message id
var Alias string    // name the queue user
var NodeUUID string // nodeuuid created on logon
var Filter = false
var SearchString = ""

var ReceivingMessages = false
var ClearMessageDetail = true

func DataStore(file string) fyne.URI {
	DataLocation, dlerr := storage.Child(nhutil.GetApp().Storage().RootURI(), file)
	if dlerr != nil {
		log.Println("DataStore error ", dlerr)
	}
	return DataLocation
}

func Load() {
	if nhutil.GetApp() == nil {
		nhutil.SetApp(app.NewWithID("org.nh3000.nh3000"))
	}

	PreferedLanguage = nhutil.GetApp().Preferences().StringWithFallback("PreferedLanguage", "eng")

	xServer, _ := nhcrypt.Encrypt(nhauth.DefaultServer, nhauth.MySecret)
	Server = nhutil.GetApp().Preferences().StringWithFallback("Server", xServer)
	xQueue, _ := nhcrypt.Encrypt("MESSAGES", nhauth.MySecret)
	Queue = nhutil.GetApp().Preferences().StringWithFallback("Queue", xQueue)
	xAlias, _ := nhcrypt.Encrypt("MyAlias", nhauth.MySecret)
	Alias = nhutil.GetApp().Preferences().StringWithFallback("Alias", xAlias)
	xQueuepassword, _ := nhcrypt.Encrypt(nhauth.QueuePassword, nhauth.MySecret)
	Queuepassword = nhutil.GetApp().Preferences().StringWithFallback("Queuepasword", xQueuepassword)

	var xCaroot = strings.ReplaceAll(nhauth.DefaultCaroot, "\n", "<>")
	ycaroot, _ := nhcrypt.Encrypt(xCaroot, nhauth.MySecret)
	nhauth.Caroot = nhutil.GetApp().Preferences().StringWithFallback("Caroot", ycaroot)

	yclientcert, _ := nhcrypt.Encrypt(nhauth.DefaultClientcert, nhauth.MySecret)
	nhauth.Clientcert = nhutil.GetApp().Preferences().StringWithFallback("Clientcert", yclientcert)

	var xClientkey = strings.ReplaceAll(nhauth.DefaultClientkey, "\n", "<>")
	yclientkey, _ := nhcrypt.Encrypt(xClientkey, nhauth.MySecret)
	nhauth.Clientkey = nhutil.GetApp().Preferences().StringWithFallback("Clientkey", yclientkey)

	var ymsgmaxage = []string{"12h", "24h", "161h", "8372h"}
	xmsgmaxage, _ := nhcrypt.Encrypt(strings.Join(ymsgmaxage, ","), nhauth.MySecret)
	Msgmaxage = nhutil.GetApp().Preferences().StringWithFallback("Msgmaxage", xmsgmaxage)

	PasswordMinimumSize = nhutil.GetApp().Preferences().StringWithFallback("PasswordMinimumSize", "12")
	PasswordMustContainNumber = nhutil.GetApp().Preferences().StringWithFallback("PasswordMustContainNumber", "Yes")
	PasswordMustContainLetter = nhutil.GetApp().Preferences().StringWithFallback("PasswordMustContainLetter", "Yes")
	PasswordMustContainSpecial = nhutil.GetApp().Preferences().StringWithFallback("PasswordMustContainSpecial", "Yes")

	// prepare for operationsnhauth.
	yServer, _ := nhcrypt.Decrypt(Server, nhauth.MySecret)
	Server = yServer
	yMsgmaxage, _ := nhcrypt.Decrypt(Msgmaxage, nhauth.MySecret)
	Msgmaxage = yMsgmaxage
	yQueue, _ := nhcrypt.Decrypt(Queue, nhauth.MySecret)
	Queue = yQueue
	yAlias, _ := nhcrypt.Decrypt(Alias, nhauth.MySecret)
	Alias = yAlias
	yQueuepassword, _ := nhcrypt.Decrypt(Queuepassword, nhauth.MySecret)
	Queuepassword = yQueuepassword
	yCaroot, _ := nhcrypt.Decrypt(nhauth.Caroot, nhauth.MySecret)
	nhauth.Caroot = strings.ReplaceAll(yCaroot, "<>", "\n")
	yClientcert, _ := nhcrypt.Decrypt(nhauth.Clientcert, nhauth.MySecret)
	nhauth.Clientcert = strings.ReplaceAll(yClientcert, "<>", "\n")
	yClientkey, _ := nhcrypt.Decrypt(nhauth.Clientkey, nhauth.MySecret)
	nhauth.Clientkey = strings.ReplaceAll(yClientkey, "<>", "\n")
	//log.Println("caroot ", nhauth.Caroot)
	//log.Println("cert ", nhauth.Clientcert)
	//log.Println("key ", nhauth.Clientkey)
}

func Save() {
	xCaroot, _ := nhcrypt.Encrypt(nhauth.Caroot, nhauth.MySecret)
	nhutil.GetApp().Preferences().SetString("Caroot", xCaroot)
	xClientcert, _ := nhcrypt.Encrypt(nhauth.Clientcert, nhauth.MySecret)
	nhutil.GetApp().Preferences().SetString("Clientcert", xClientcert)
	xClientkey, _ := nhcrypt.Encrypt(nhauth.Clientkey, nhauth.MySecret)
	nhutil.GetApp().Preferences().SetString("Clientkey", xClientkey)
	xMsgmaxage, _ := nhcrypt.Encrypt(Msgmaxage, nhauth.MySecret)
	nhutil.GetApp().Preferences().SetString("Msgmaxage", xMsgmaxage)
	xServer, _ := nhcrypt.Encrypt(Server, nhauth.MySecret)
	nhutil.GetApp().Preferences().SetString("Server", xServer)
	xQueue, _ := nhcrypt.Encrypt(Queue, nhauth.MySecret)
	nhutil.GetApp().Preferences().SetString("Queue", xQueue)
	xAlias, _ := nhcrypt.Encrypt(Alias, nhauth.MySecret)
	nhutil.GetApp().Preferences().SetString("Alias", xAlias)
	nhutil.GetApp().Preferences().SetString("PreferedLanguage", PreferedLanguage)
	xQueuepassword, _ := nhcrypt.Encrypt(Queuepassword, nhauth.MySecret)
	nhutil.GetApp().Preferences().SetString("Queuepassword", xQueuepassword)
	nhutil.GetApp().Preferences().SetString("PasswordMinimumSize", PasswordMinimumSize)
	nhutil.GetApp().Preferences().SetString("PasswordMustContainNumber", PasswordMustContainNumber)
	nhutil.GetApp().Preferences().SetString("PasswordMustContainLetter", PasswordMustContainLetter)
	nhutil.GetApp().Preferences().SetString("PasswordMustContainSpecial", PasswordMustContainSpecial)
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
	if action == "URL" {
		valid := strings.Contains(strings.ToLower(value), "nats://")
		if valid == false {
			return true
		}
		valid2 := strings.Contains(value, ".")
		if valid2 == false {
			return true
		}
		valid3 := strings.Contains(value, ":")
		if valid3 == false {
			return true
		}

		return false
	}
	if action == "STRING" {
		if len(value) == 0 {
			return true
		}
		return false
	}

	if action == "PASSWORD" {
		var iserrors = false
		vlen, _ := strconv.Atoi(PasswordMinimumSize)
		if (len(value) <= vlen) == false {
			iserrors = true
		}
		if PasswordMustContainLetter == "Yes" && !iserrors {

			for _, r := range value {
				if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') {
					iserrors = true
					break
				}
			}
		}
		if PasswordMustContainNumber == "Yes" && !iserrors {
			iserrors = true
			for _, r := range value {
				if unicode.IsNumber(r) {
					iserrors = false
					break
				}
			}
		}
		if PasswordMustContainSpecial == "Yes" && !iserrors {
			iserrors = true
			var schars = []string{"|", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "-", "+", "=", "{", "}", "]", "[", "|", ":", ";", ",", ".", "#", "'", "\"", "\\", "%", "?", "\n", "<", "Ø", "ð", ">", "ï", "û"}
			for _, sc := range schars {
				if strings.Contains(value, sc) {
					iserrors = false
					break
				}
			}
		}
		return iserrors
	}
	if action == "CERTIFICATE" {
		valid := strings.Contains(value, "-----BEGIN CERTIFICATE-----")
		if valid == false {
			return false
		}
		valid2 := strings.Contains(value, "-----END CERTIFICATE-----")
		if valid2 == false {
			return false
		}
	}
	if action == "KEY" {
		log.Println("key")
		valid := strings.Contains(value, "-----BEGIN RSA PRIVATE KEY-----")
		if valid == false {
			log.Println("begin")
			return false
		}
		valid2 := strings.Contains(value, "-----END RSA PRIVATE KEY-----")
		if valid2 == false {
			log.Println("end")
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
