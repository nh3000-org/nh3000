package nhhash

import (
	"log"
	"os"

	"fyne.io/fyne/v2/storage"
	"golang.org/x/crypto/bcrypt"

	//"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhlang"
	"github.com/nh3000-org/nh3000/nhpref"
)

// create and load hash
func LoadWithDefault(filename string, password string) (string, bool) {
	nhexists, _ := storage.Exists(nhpref.DataStore(filename))
	if !nhexists {
		log.Println("err ")
		pwh, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("err ", err)
			return nhlang.GetLangs("hash-err2"), true
		}
		wrt, errwrite := storage.Writer(nhpref.DataStore(filename))
		_, err2 := wrt.Write([]byte(pwh))
		if errwrite != nil || err2 != nil {
			log.Println("err ", err, " errwrite ", errwrite)
			return nhlang.GetLangs("hash-err1"), true
		}
		//Hash = string(pwh)
		return string(pwh), false
	}
	ph, errf := os.ReadFile(nhpref.DataStore(filename).Path())
	if errf != nil {
		return nhlang.GetLangs("hash-err3"), true
	}

	return string(ph), false
}

// save hash
func Save(filename string, hash string) (string, bool) {
	errf := storage.Delete(nhpref.DataStore(filename))
	if errf != nil {
		return nhlang.GetLangs("hash-err3"), true
	}
	wrt, errwrite := storage.Writer(nhpref.DataStore(filename))
	_, err2 := wrt.Write([]byte(hash))
	if errwrite != nil || err2 != nil {
		return nhlang.GetLangs("hash-err2"), true
	}

	return hash, false
}

// hash and salt
func HashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

// validate password
func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
