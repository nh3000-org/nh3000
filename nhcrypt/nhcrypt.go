package nhcrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
)

var keyAes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}  // must be 16 bytes
var keyHmac = []byte{36, 45, 53, 21, 87, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05} // must be 16 bytes
const BUFFER_SIZE int = 4096
const IV_SIZE int = 16

func EncryptFile(filePathIn, filePathOut string) error {
	inFile, err := os.Open(filePathIn)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(filePathOut)
	if err != nil {
		return err
	}
	defer outFile.Close()

	iv := make([]byte, IV_SIZE)
	_, err = rand.Read(iv)
	if err != nil {
		return err
	}

	aes, err := aes.NewCipher(keyAes)
	if err != nil {
		return err
	}

	ctr := cipher.NewCTR(aes, iv)
	hmac := hmac.New(sha256.New, keyHmac)

	buf := make([]byte, BUFFER_SIZE)
	for {
		n, err := inFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		outBuf := make([]byte, n)
		ctr.XORKeyStream(outBuf, buf[:n])
		hmac.Write(outBuf)
		outFile.Write(outBuf)

		if err == io.EOF {
			break
		}
	}

	outFile.Write(iv)
	hmac.Write(iv)
	outFile.Write(hmac.Sum(nil))

	return nil
}
func DecryptFile(filePathIn, filePathOut string) error {
	inFile, err := os.Open(filePathIn)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(filePathOut)
	if err != nil {
		return err
	}
	defer outFile.Close()

	iv := make([]byte, IV_SIZE)
	_, err = rand.Read(iv)
	if err != nil {
		return err
	}

	aes, err := aes.NewCipher(keyAes)
	if err != nil {
		return err
	}

	ctr := cipher.NewCTR(aes, iv)
	hmac := hmac.New(sha256.New, keyHmac)

	buf := make([]byte, BUFFER_SIZE)
	for {
		n, err := inFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		outBuf := make([]byte, n)
		ctr.XORKeyStream(outBuf, buf[:n])
		hmac.Write(outBuf)
		outFile.Write(outBuf)

		if err == io.EOF {
			break
		}
	}

	outFile.Write(iv)
	hmac.Write(iv)
	outFile.Write(hmac.Sum(nil))

	return nil
}

// encrypt string
func Encrypt(text string, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, keyAes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

// decrypt string
func Decrypt(text string, secret string) (string, error) {
	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", err
	}
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, keyAes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}

// encode to base64
func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// decode base64
func decode(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		data = []byte(s)
	}
	return data
}
