// Copyright 2012-2023 The NH3000 Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A Go file encryption/decryption client for the NH3000 messaging system (https://newhorizons3000.org).

package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"fmt"
	"io"

	"os"

	"strings"

	"github.com/nh3000-org/nh3000/config"
)

// var idcount int
var MyFileLang = "eng"

// eng esp cmnBAE
var MyFileLangs = map[string]string{
	"eng-fl-fl":   "Language to Use eng or esp or hin",
	"spa-fl-fl":   "Idioma a utilizar eng o esp o hin",
	"hin-fl-fl":   "उपयोग करने के लिए भाषा eng या esp या hin",
	"hin-fl-fa":   "ENCRYPT or DECRYPT",
	"spa-fl-fa":   "CIFRAR o DESCIFRAR",
	"eng-fl-fa":   "एन्क्रिप्ट या डिक्रिप्ट",
	"eng-fl-if":   "Input File",
	"spa-fl-if":   "Fichero de Entrada",
	"hin-fl-if":   "इनपुट फ़ाइल",
	"eng-fl-of":   "Output File",
	"spa-fl-of":   "Archivo de Salida",
	"hin-fl-of":   "आउटपुट फ़ाइल",
	"eng-fl-err1": "File Does Not Exist",
	"spa-fl-err1": "El Archivo no Existe",
	"hin-fl-err1": "फ़ाइल मौजूद नहीं है",
	"eng-fl-err2": "File Already Exists",
	"spa-fl-err2": "El Archivo ya Existe",
	"hin-fl-err2": "फ़ाइल पहले से ही मौजूद है",
	"eng-fl-err3": "Must be ENCRYPT or DECRYPT",
	"spa-fl-err3": "Debe ser CIFRADO o DESCIFRADO",
	"hin-fl-err3": "एन्क्रिप्ट या डिक्रिप्ट होना चाहिए",
}

// return translation strings
func GetLangs(mystring string) string {
	value, err := MyFileLangs[MyFileLang+"-"+mystring]
	if !err {
		return "xxx"
	}
	return value
}

// main loop for receiving pipe
func main() {
	if strings.HasPrefix(os.Getenv("LANG"), "en") {
		MyFileLang = "eng"
	}
	if strings.HasPrefix(os.Getenv("LANG"), "sp") {
		MyFileLang = "spa"
	}

	fileLang := flag.String("filelang", MyFileLang, GetLangs("fl-fl"))
	fileAction := flag.String("fileaction", "UNKNOWN", GetLangs("fl-fa"))
	MyFileLang = *fileLang
	fileInput := flag.String("fileinput", "UNKNOWN", GetLangs("fl-if"))
	fileOutput := flag.String("fileoutput", "UNKNOWN", GetLangs("fl-of"))

	flag.Parse()
	fmt.Println("====================================================== ")
	fmt.Println("file -filelang ", *fileLang, " -fileinput ", *fileInput, " -fileoutput ", *fileOutput, " -fileaction ", *fileAction)
	fmt.Println("====================================================== ")
	// edit inputs
	var errors = false
	if *fileAction != "ENCRYPT" && *fileAction != "DECRYPT" {
		errors = true
		fmt.Println("-fileaction " + " - " + GetLangs("fl-err3"))
	}
	if _, err := os.Stat(*fileInput); err != nil {
		errors = true
		fmt.Println("-fileinput" + " - " + GetLangs("fl-err1"))
	}
	if _, err := os.Stat(*fileOutput); err == nil {
		errors = true
		fmt.Println("-fileoutput" + " - " + GetLangs("fl-err2"))
	}
	if !errors {
		fmt.Println(*fileAction + " " + *fileInput + " > " + *fileOutput)
		if *fileAction == "ENCRYPT" {
			err := EncryptFile(*fileInput, *fileOutput)
			if err != nil {
				fmt.Println(err)
			}
		}
		if *fileAction == "DECRYPT" {
			err := DecryptFile(*fileInput, *fileOutput)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
func EncryptFile(filePathIn, filePathOut string) error {

	infile, err := os.Open(filePathIn)
	if err != nil {
		return err
	}
	defer infile.Close()

	block, err := aes.NewCipher(config.KeyHmac)
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

	block, err := aes.NewCipher(config.KeyHmac)
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
