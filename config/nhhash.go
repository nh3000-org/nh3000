package config

import (
	"log"
	"os"



	"fyne.io/fyne/v2/storage"
	"golang.org/x/crypto/bcrypt"
)

// create and load hash
func LoadHashWithDefault(filename string, password string) (string, bool) {

	nhexists, _ := storage.Exists(DataStore(filename))
	if !nhexists {
		log.Println("err ")
		pwh, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("err ", err)
			return GetLangs("hash-err2"), true
		}
		wrt, errwrite := storage.Writer(DataStore(filename))
		_, err2 := wrt.Write([]byte(pwh))
		if errwrite != nil || err2 != nil {
			log.Println("err ", err, " errwrite ", errwrite)
			return GetLangs("hash-err1"), true
		}
		//Hash = string(pwh)
		return string(pwh), false
	}
	ph, errf := os.ReadFile(DataStore(filename).Path())
	if errf != nil {
		return GetLangs("hash-err3"), true
	}

	return string(ph), false
}

// save hash
func SaveHash(filename string, hash string) (string, bool) {

	errf := storage.Delete(DataStore(filename))
	if errf != nil {
		return GetLangs("hash-err3"), true
	}
	wrt, errwrite := storage.Writer(DataStore(filename))
	_, err2 := wrt.Write([]byte(hash))
	if errwrite != nil || err2 != nil {
		return GetLangs("hash-err2"), true
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
