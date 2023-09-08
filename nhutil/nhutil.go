package nhutil

import (
	"log"
	"net/url"
	"strconv"
	"strings"
	"unicode"

	"fyne.io/fyne/v2"
	"golang.org/x/crypto/bcrypt"
)

type util interface {
	ParseURL(string) *url.URL

	SetApp(fyne.App)
	GetApp() fyne.App
	Edit(string, string) bool
}

var App fyne.App

func ParseURL(urlStr string) *url.URL {
	link, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return link
}

func SetApp(a fyne.App) {
	FyneApp = a
}
func GetApp() fyne.App {
	return FyneApp
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
func hashAndSalt(pwd []byte) string {

	// Use GenerateFromPassword to hash & salt pwd
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
