package nhpref

import (
	"log"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/storage"
	"github.com/nh3000-org/nh3000/nhcrypt"
	"github.com/nh3000-org/nh3000/nhutil"
)

/*
*  The following fields need to be modified for you production
*  Environment to provide maximum security
*
*  These fields are meant to be distributed at compile time and
*  editable in the gui.
*
 */
var MyBytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05} // must be 16 bytes
const MySecret string = "abd&1*~#^2^#s0^=)^^7%c34"                                   // must be 24 characters
const MyDurable string = "snatsdurable"
const PasswordDefault = "123456" // default password shipped with app

// tls
var Caroot = ""
var Clientcert = ""
var Clientkey = ""

// version
const Version = "snats-beta.1"

var MyApp fyne.App

var LoggedOn bool = false
var PasswordValid bool = false

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

var IdUUID string   // unique message id
var Alias string    // name the queue user
var NodeUUID string // nodeuuid created on logon

//type json interface {
//	Load()
//	Save()
//	DataStore(string) fyne.URI
//}

func DataStore(file string) fyne.URI {
	DataLocation, dlerr := storage.Child(MyApp.Storage().RootURI(), file)
	if dlerr != nil {
		log.Println("DataStore error ", dlerr)
	}
	return DataLocation
}

func Load() {
	if nhutil.GetApp() == nil {
		MyApp = app.NewWithID("org.nh3000.snats")
		nhutil.SetApp(app.NewWithID("org.nh3000.snats"))
	}

	PreferedLanguage = MyApp.Preferences().StringWithFallback("PreferedLanguage", "eng")

	xServer, _ := nhcrypt.Encrypt("nats://nats.newhorizons3000.org:4222", MySecret)
	Server = MyApp.Preferences().StringWithFallback("Server", xServer)
	xQueue, _ := nhcrypt.Encrypt("MESSAGES", MySecret)
	Queue = MyApp.Preferences().StringWithFallback("Queue", xQueue)
	xAlias, _ := nhcrypt.Encrypt("MyAlias", MySecret)
	Alias = MyApp.Preferences().StringWithFallback("Alias", xAlias)
	xQueuepassword, _ := nhcrypt.Encrypt("123456789012345678901234", MySecret)
	Queuepassword = MyApp.Preferences().StringWithFallback("Queuepasword", xQueuepassword)

	var xCaroot = strings.ReplaceAll("-----BEGIN CERTIFICATE-----\nMIICFDCCAbugAwIBAgIUDkHxHO1DwrlkTzUimG5PoiswB6swCgYIKoZIzj0EAwIw\nZjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UE\nChMDU0VDMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMz\nMDAwLm9yZzAgFw0yMzAzMzExNzI5MDBaGA8yMDUzMDMyMzE3MjkwMFowZjELMAkG\nA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UEChMDU0VD\nMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9y\nZzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABHXwMUfMXiJix3tuzFymcA+3RkeY\nZE7urUzVgaqkv/Oef3jhqhtf1XzK/qVYGxWWmpvADGB252PG1Mp7Z5wmzqyjRTBD\nMA4GA1UdDwEB/wQEAwIBBjASBgNVHRMBAf8ECDAGAQH/AgEBMB0GA1UdDgQWBBQm\nFA5caanuqxGFOf9DtZkVYv5dCzAKBggqhkjOPQQDAgNHADBEAiB3BheNP4XdBZ27\nxVBQ7ztMJqK7wDi1V3LuMy5jmXr7rQIgHCse0oaiAwcl4VwF00aSshlV+T/da0Tx\n1ANkaM+rie4=\n-----END CERTIFICATE-----\n", "\n", "<>")
	ycaroot, _ := nhcrypt.Encrypt(xCaroot, MySecret)
	Caroot = MyApp.Preferences().StringWithFallback("Caroot", ycaroot)

	var xClientcert = strings.ReplaceAll("-----BEGIN CERTIFICATE-----\nMIIDUzCCAvigAwIBAgIUUyhlJt8mp1XApRbSkdrUS55LGV8wCgYIKoZIzj0EAwIw\nZjELMAkGA1UEBhMCVVMxCzAJBgNVBAgTAkZMMQswCQYDVQQHEwJDVjEMMAoGA1UE\nChMDU0VDMQwwCgYDVQQLEwNuaDExITAfBgNVBAMTGG5hdHMubmV3aG9yaXpvbnMz\nMDAwLm9yZzAeFw0yMzAzMzExNzI5MDBaFw0yODAzMjkxNzI5MDBaMHIxCzAJBgNV\nBAYTAlVTMRAwDgYDVQQIEwdGbG9yaWRhMRIwEAYDVQQHEwlDcmVzdHZpZXcxGjAY\nBgNVBAoTEU5ldyBIb3Jpem9ucyAzMDAwMSEwHwYDVQQLExhuYXRzLm5ld2hvcml6\nb25zMzAwMC5vcmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDFttVH\nQ131JYwazAQMm0XAQvRvTjTjOY3aei1++mmQ+NQ9mrOFk6HlZFoKqsy6+HPXsB9x\nQbWlYvUOuqBgb9xFQZoL8jiKskLLrXoIxUAlIBTlyf76r4SV+ZpxJYoGzXNTedaU\n0EMTyAiUQ6nBbFMXiehN5q8VzxtTESk7QguGdAUYXYsCmYBvQtBXoFYO5CHyhPqu\nOZh7PxRAruYypEWVFBA+29+pwVeaRHzpfd/gKLY4j2paInFn7RidYUTqRH97BjdR\nSZpOJH6fD7bI4L09pnFtII5pAARSX1DntS0nWIWhYYI9use9Hi/B2DRQLcDSy1G4\n0t1z4cdyjXxbFENTAgMBAAGjgawwgakwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQM\nMAoGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFAzgPVB2/sfT7R0U\ne3iXRSvUkfoQMB8GA1UdIwQYMBaAFCYUDlxpqe6rEYU5/0O1mRVi/l0LMDQGA1Ud\nEQQtMCuCGG5hdHMubmV3aG9yaXpvbnMzMDAwLm9yZ4IJMTI3LDAsMCwxhwTAqABn\nMAoGCCqGSM49BAMCA0kAMEYCIQCDlUH2j69mJ4MeKvI8noOmvLHfvP4qMy5nFW2F\nPT5UxgIhAL6pHFyEbANtSkcVJqxTyKE4GTXcHc4DB43Z1F7VxSJj\n-----END CERTIFICATE-----\n", "\n", "<>")
	yclientcert, _ := nhcrypt.Encrypt(xClientcert, MySecret)
	Clientcert = MyApp.Preferences().StringWithFallback("Clientcert", yclientcert)

	var xClientkey = strings.ReplaceAll("-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAxbbVR0Nd9SWMGswEDJtFwEL0b0404zmN2notfvppkPjUPZqz\nhZOh5WRaCqrMuvhz17AfcUG1pWL1DrqgYG/cRUGaC/I4irJCy616CMVAJSAU5cn+\n+q+ElfmacSWKBs1zU3nWlNBDE8gIlEOpwWxTF4noTeavFc8bUxEpO0ILhnQFGF2L\nApmAb0LQV6BWDuQh8oT6rjmYez8UQK7mMqRFlRQQPtvfqcFXmkR86X3f4Ci2OI9q\nWiJxZ+0YnWFE6kR/ewY3UUmaTiR+nw+2yOC9PaZxbSCOaQAEUl9Q57UtJ1iFoWGC\nPbrHvR4vwdg0UC3A0stRuNLdc+HHco18WxRDUwIDAQABAoIBACe0XMZP4Al//c/P\n0qxZbjt69q13jiVnhHYwfPx3+0UywySP8adMi4GOkop73Ftb05+n7diHspvA8KeB\nkP1s2VZLI01s2i/4NnPCpbQnMIeEFs5Cr2LWZpDbrEk2ma5eCd/kotQFssLBM//a\nSrfeMh2TA0TJo7WEft9Cnf4ZeEkKnycplfvwTyv286iFZCYo2dv66BfTej6kkVCo\nAi+ZVCe2zSqRYyr0u4/j/kE3b3eSkCnY2IVcqlP7epuEGVOZyxeFLwM5ljbWL816\npA6WIJgQo2EQ1N7L531neg5WjXQ/UwTQoXP1jvuuVtKtOBFqm1IshEyFk3WpsfpD\nr16OTdECgYEA6FB6NYxYtnWPaIYAOqP7GtMKoJujH8MtZy6J33LkxI7nPkMkn8Mv\nva32tvjU4Bu1FVNp9k5guC+b+8ixXK0URj25IOhDs6K57tck22W9WiTZlmnkCO01\nJOavrelWbvYt5xNWIdnPualoPfGB0iJKXsKY/bpH4eVfhWwpNPI5sMkCgYEA2d9G\nEPuWN6gUjZ+JfdS+0WHK1yGD7thXs7MPUlhGqDzBryh2dkywyo8U8+tMLuDok1RZ\njnT3PYkLQEpzoV0qBkpFFShL6ubaGmDz1UZsozl0YcIg4diZeuPHnIAeXOFrhgYf\n825163LmT3jYHCROFEMLtTYyIQP0EznE+qFT3TsCgYEApgtvbfqkJbWdDL5KR5+R\nCLky7VyQmVEtkIRI8zbxoDPrwCrJcI9X/iDrKBhuPshPA7EdGXkn1D3jJXFqo6zp\nwtK3EXgxe6Ghd766jz4Guvl/s+x3mpHA3GEtzAXtS14VrQW7GHLP8AnPggauHX14\n3oYER8XvPtxtC7YlNbyz01ECgYAe2b7SKM3ck7BVXYHaj4V1oKNYUyaba4b/qxtA\nTb+zkubaJqCfn7xo8lnFMExZVv+X3RnRUj6wN/ef4ur8rnSE739Yv5wAZy/7DD96\ns74uXrRcI2EEmechv59ESeACxuiy0as0jS+lZ1+1YSc41Os5c0T1I/d1NVoaXtPF\nqZJ2gQKBgBp/XavdULBPzC7B8tblySzmL01qJZV7MSSVo2/1vJ7gPM0nQPZdTDog\nTfA5QKSX9vFTGC9CZHSJ+fabYDDd6+3UNYUKINfr+kwu9C2cysbiPaM3H27WR5mW\n5LhStAfwuRRYBDsG2ndjraxcBrrPdtkbS0dpeQUDJxvkMIuLHnhQ\n-----END RSA PRIVATE KEY-----\n", "\n", "<>")
	yclientkey, _ := nhcrypt.Encrypt(xClientkey, MySecret)
	Clientkey = MyApp.Preferences().StringWithFallback("Clientkey", yclientkey)

	var ymsgmaxage = []string{"12h", "24h", "161h", "8372h"}
	xmsgmaxage, _ := nhcrypt.Encrypt(strings.Join(ymsgmaxage, ","), MySecret)
	Msgmaxage = MyApp.Preferences().StringWithFallback("Msgmaxage", xmsgmaxage)

	PasswordMinimumSize = MyApp.Preferences().StringWithFallback("PasswordMinimumSize", "6")
	PasswordMustContainNumber = MyApp.Preferences().StringWithFallback("PasswordMustContainNumber", "Yes")
	PasswordMustContainLetter = MyApp.Preferences().StringWithFallback("PasswordMustContainLetter", "Yes")
	PasswordMustContainSpecial = MyApp.Preferences().StringWithFallback("PasswordMustContainSpecial", "Yes")

	// prepare for operations
	yServer, _ := nhcrypt.Decrypt(Server, MySecret)
	Server = yServer
	yMsgmaxage, _ := nhcrypt.Decrypt(Msgmaxage, MySecret)
	Msgmaxage = yMsgmaxage
	yQueue, _ := nhcrypt.Decrypt(Queue, MySecret)
	Queue = yQueue
	yAlias, _ := nhcrypt.Decrypt(Alias, MySecret)
	Alias = yAlias
	yQueuepassword, _ := nhcrypt.Decrypt(Queuepassword, MySecret)
	Queuepassword = yQueuepassword
	yCaroot, _ := nhcrypt.Decrypt(Caroot, MySecret)
	Caroot = strings.ReplaceAll(yCaroot, "<>", "\n")
	yClientcert, _ := nhcrypt.Decrypt(Clientcert, MySecret)
	Clientcert = strings.ReplaceAll(yClientcert, "<>", "\n")
	yClientkey, _ := nhcrypt.Decrypt(Clientkey, MySecret)
	Clientkey = strings.ReplaceAll(yClientkey, "<>", "\n")

}

func Save() {
	xCaroot, _ := nhcrypt.Encrypt(Caroot, MySecret)
	MyApp.Preferences().SetString("Caroot", xCaroot)
	xClientcert, _ := nhcrypt.Encrypt(Clientcert, MySecret)
	MyApp.Preferences().SetString("Clientcert", xClientcert)
	xClientkey, _ := nhcrypt.Encrypt(Clientkey, MySecret)
	MyApp.Preferences().SetString("Clientkey", xClientkey)

	xMsgmaxage, _ := nhcrypt.Encrypt(Msgmaxage, MySecret)
	MyApp.Preferences().SetString("Msgmaxage", xMsgmaxage)

	xServer, _ := nhcrypt.Encrypt(Server, MySecret)
	MyApp.Preferences().SetString("Server", xServer)
	xQueue, _ := nhcrypt.Encrypt(Queue, MySecret)
	MyApp.Preferences().SetString("Queue", xQueue)
	xAlias, _ := nhcrypt.Encrypt(Alias, MySecret)
	MyApp.Preferences().SetString("Alias", xAlias)
	MyApp.Preferences().SetString("PreferedLanguage", PreferedLanguage)
	xQueuepassword, _ := nhcrypt.Encrypt(Queuepassword, MySecret)
	MyApp.Preferences().SetString("Queuepassword", xQueuepassword)
	MyApp.Preferences().SetString("PasswordMinimumSize", PasswordMinimumSize)
	MyApp.Preferences().SetString("PasswordMustContainNumber", PasswordMustContainNumber)
	MyApp.Preferences().SetString("PasswordMustContainLetter", PasswordMustContainLetter)
	MyApp.Preferences().SetString("PasswordMustContainSpecial", PasswordMustContainSpecial)

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
		valid := strings.Contains(value, "-----BEGIN RSA PRIVATE KEY-----")
		if valid == false {
			return false
		}
		valid2 := strings.Contains(value, "-----END RSA PRIVATE KEY-----")
		if valid2 == false {
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
