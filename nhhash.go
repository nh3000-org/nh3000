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
	HashAndSalt([]byte) string
	ComparePasswords(string,[]byte) bool
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
func HashAndSalt(pwd []byte) string {

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

func ComparePasswords(hashedPwd string, plainPwd []byte) bool {
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