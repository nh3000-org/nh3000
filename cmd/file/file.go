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
	"flag"
	"fmt"

	"os"

	"strings"

	"github.com/nh3000-org/nh3000/nhcrypt"
)

// var idcount int
var MyFileLang = "eng"

// eng esp cmn
var MyLangs = map[string]string{
	"eng-fl-fl":   "Language to Use eng or esp",
	"spa-fl-fl":   "Idioma a Utilizar eng o esp",
	"eng-fl-fa":   "ENCRYPT or DECRYPT",
	"spa-fl-fa":   "CIFRAR o DESCIFRAR",
	"eng-fl-if":   "Input File",
	"spa-fl-if":   "Fichero de Entrada",
	"eng-fl-of":   "Output File",
	"spa-fl-of":   "Archivo de salida",
	"eng-fl-us":   "Usage:",
	"spa-fl-us":   "Uso:",
	"eng-fl-ro":   "Run Options:",
	"spa-fl-ro":   "Opciones de Ejecución:",
	"eng-fl-err1": "File Does Not Exist",
	"spa-fl-err1": "El Archivo no Existe",
	"eng-fl-err2": "File Already Exists",
	"spa-fl-err2": "El Archivo ya Existe",
	"eng-fl-err3": "Must be ENCRYPT or DECRYPT",
	"spa-fl-err3": "Debe ser CIFRADO o DESCIFRADO",
}

// return translation strings
func GetLangs(mystring string) string {
	value, err := MyLangs[MyFileLang+"-"+mystring]
	if err == false {
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
	fileInput := flag.String("fileinput", "UNKNOWN", GetLangs("fl-if"))
	fileOutput := flag.String("fileoutput", "UNKNOWN", GetLangs("fl-of"))

	flag.Parse()
	fmt.Println(GetLangs("fl-us"))
	fmt.Println("file -filelang FILELANG -fileinput FILEINPUT -fileoutput FILEOUTPUT -fileaction FILEACTION")
	fmt.Println("")
	fmt.Println(GetLangs("fl-ro"))
	fmt.Println("-loglang: ", *fileLang)
	MyFileLang = *fileLang
	fmt.Println("-fileaction: ", *fileAction)
	fmt.Println("-fileinput: ", *fileInput)
	fmt.Println("-fileoutput: ", *fileOutput)

	// edit inputs
	var errors = false
	if *fileAction != "ENCRYPT" && *fileAction != "DECRYPT" {
		errors = true
		fmt.Println(*fileAction + " - " + GetLangs("fl-err3"))
	}
	if _, err := os.Stat(*fileInput); err != nil {
		errors = true
		fmt.Println(*fileInput + " - " + GetLangs("fl-err1"))
	}
	if _, err := os.Stat(*fileOutput); err == nil {
		errors = true
		fmt.Println(*fileOutput + " - " + GetLangs("fl-err2"))
	}
	if errors == false {
		fmt.Println(*fileAction + " " + *fileInput + " > " + *fileOutput)
		if *fileAction == "ENCRYPT" {
			err := nhcrypt.EncryptFile(*fileInput, *fileOutput)
			if err != nil {
				fmt.Println(err)
			}
		}

		if *fileAction == "DECRYPT" {
			err := nhcrypt.DecryptFile(*fileInput, *fileOutput)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
