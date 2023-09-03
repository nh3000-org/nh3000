package nhhash

import (
	"fyne.io/fyne/v2/storage"
	"golang.org/x/crypto/bcrypt"
	"nhlang"
)

var Hash string

type hash interface {
	LoadWithDefault(string) // file name
	Save(string, string)    // file name, hash
}

// provide file name and password to hash
func LoadWithDefault(filename string, password string) (string, bool) {
	lwderr, _ := storage.Exists(DataStore(filename))
	if lwderr == true {
		wrt, errwrite := storage.Writer(DataStore(filename))
		_, err2 := wrt.Write([]byte(Hash))
		pwh, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return GetLangs("hash-err2"), true
		}
		if errwrite != nil || err2 != nil {
			return GetLangs("hash-err1"), true
		}
		Hash = string(pwh)
		return nil, false
	}
	ph, errf := os.ReadFile(DataStore(filename).Path())
	if errf != nil {
		return GetLangs("hash-err3"), true
	}
	Hash = string(ph)
	return nil, false
}
func Save(filename string, hash string) (string, bool) {
	errf := storage.Delete(DataStore(filename))

	if errf != nil {
		return GetLangs("hash-err3"), true
	}
	wrt, errwrite := storage.Writer(DataStore(filename))
	_, err2 := wrt.Write([]byte(hash))
	if errwrite != nil || err2 != nil {
		return GetLangs("hash-err2"), true
	}
	Hash = hash
	return nil, false
}
