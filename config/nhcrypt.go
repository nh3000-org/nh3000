package config

import (
	"crypto/aes"
	"crypto/cipher"
	"log"

	"crypto/rand"
	"encoding/base64"
	"io"

	"os"
)
/*
todo: use circl from cloudfare to encrypt data
post quantum  crystals-kyber, crystals-dilithium, sphincs=, FALCON
*/
const BUFFER_SIZE int = 4096
const IV_SIZE int = 16

func EncryptFile(filePathIn, filePathOut string) error {

	infile, err := os.Open(filePathIn)
	if err != nil {
		return err
	}
	defer infile.Close()

	block, err := aes.NewCipher(KeyHmac)
	if err != nil {
		return err
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	iv := make([]byte, block.BlockSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	outfile, err := os.OpenFile(filePathOut, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer outfile.Close()
	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}
	// Append the IV
	outfile.Write(iv)
	return nil
}
func DecryptFile(filePathIn, filePathOut string) error {

	infile, err := os.Open(filePathIn)
	if err != nil {
		return err
	}
	defer infile.Close()

	block, err := aes.NewCipher(KeyHmac)
	if err != nil {
		return err
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of repeat.
	fi, err := infile.Stat()
	if err != nil {
		return err
	}

	iv := make([]byte, block.BlockSize())
	msgLen := fi.Size() - int64(len(iv))
	_, err = infile.ReadAt(iv, msgLen)
	if err != nil {
		return err
	}

	outfile, err := os.OpenFile(filePathOut, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer outfile.Close()

	// The buffer size must be multiple of 16 bytes
	buf := make([]byte, 1024)
	stream := cipher.NewCTR(block, iv)
	for {
		n, err := infile.Read(buf)
		if n > 0 {
			// The last bytes are the IV, don't belong the original message
			if n > int(msgLen) {
				n = int(msgLen)
			}
			msgLen -= int64(n)
			stream.XORKeyStream(buf, buf[:n])
			// Write into file
			outfile.Write(buf[:n])
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}
	}
	return nil
}

// encrypt string
func Encrypt(text string, secret string) string {

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		log.Println("enNotable stories and conversation starterscrypt 1", err)
		return ""
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, KeyAes)
	cipherText := make([]byte, len(plainText))
	//log.Println("cyphertext ", plainText)
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText)
}

// decrypt string
func Decrypt(text string, secret string) string {

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return ""
	}
	cipherText := decode(text)
	cfb := cipher.NewCFBDecrypter(block, KeyAes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText)
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
