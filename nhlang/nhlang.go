package nhlang

import (
	"github.com/nh3000-org/nh3000/nhpref"
)

// var PreferedLanguage string // language string
// eng esp cmn
var MyLangs = map[string]string{
	"eng-mn-intro-1":     "Encrypted Communications Using NATS ",
	"eng-mn-intro-2":     "for Additional Info.",
	"spa-mn-intro-1":     "Comunicaciones Encriptadas Usando NATS ",
	"spa-mn-intro-2":     "para Información Adicional.",
	"hin-mn-intro-1":     "NATS का उपयोग करके एन्क्रिप्टेड संचार ",
	"hin-mn-intro-2":     "अतिरिक्त जानकारी के लिए.",
	"eng-mn-mt":          "Encryption",
	"spa-mn-mt":          "Cifrada",
	"hin-mn-mt":          "कूटलेखन",
	"eng-mn-err1":        "Missing panel: ",
	"spa-mn-err1":        "Falta el panel: ",
	"hin-mn-err1":        "गुम पैनल:",
	"eng-mn-dark":        "Dark",
	"spa-mn-dark":        "Oscura",
	"hin-mn-dark":        "अँधेरा",
	"eng-mn-light":       "Light",
	"spa-mn-light":       "Ligera",
	"hin-mn-light":       "रोशनी",
	"eng-mn-retro":       "Retro",
	"spa-mn-retro":       "Retro",
	"hin-mn-retro":       "रेट्रो",
	"eng-ps-title":       "Password Reset",
	"spa-ps-title":       "Restablecimiento de contraseña",
	"hin-ps-title":       "पासवर्ड रीसेट",
	"eng-ps-password":    "Enter Original Password",
	"spa-ps-password":    "Ingrese la Contraseña Original",
	"hin-ps-password":    "मूल पासवर्ड दर्ज करें",
	"eng-ps-passwordc1":  "Enter New Password",
	"spa-ps-passwordc1":  "Ingrese Nueva Clave",
	"hin-ps-passwordc1":  "नया पासवर्ड दर्ज करें",
	"eng-ps-passwordc2":  "Enter New Password Again",
	"spa-ps-passwordc2":  "Ingrese la Nueva Contraseña Nuevamente",
	"hin-ps-passwordc2":  "नया पासवर्ड दोबारा दर्ज करें",
	"eng-ps-trypassword": "Try Password",
	"spa-ps-trypassword": "Probar Contraseña",
	"hin-ps-trypassword": "पासवर्ड आज़माएं",

	"eng-ps-err1":        "Error Creating Password Hash 1",
	"spa-ps-err1":        "Error al Crear Hash de Contraseña 1",
	"hin-ps-err1":        "पासवर्ड हैश 1 बनाने में त्रुटि",
	"eng-ps-err2":        "Error Creating Password Hash 2",
	"spa-ps-err2":        "Error al Crear Hash de Contraseña 2",
	"hin-ps-err2":        "पासवर्ड हैश 2 बनाने में त्रुटि",
	"eng-ps-err3":        "Error Reading Password Hash",
	"spa-ps-err3":        "Error al Leer el Hash de la Contraseña",
	"hin-ps-err3":        "पासवर्ड हैश पढ़ने में त्रुटि",
	"eng-ps-err4":        "Error Passwords Do Not Match",
	"spa-ps-err4":        "Las Contraseñas de Error no Coinciden",
	"hin-ps-err4":        "त्रुटि पासवर्ड मेल नहीं खाते",
	"eng-ps-err5":        "Password Accepted",
	"spa-ps-err5":        "Contraseña Aceptada",
	"hin-ps-err5":        "पासवर्ड स्वीकृत",
	"eng-ps-chgpassword": "Change Password",
	"spa-ps-chgpassword": "Cambiar la Contraseña",
	"hin-ps-chgpassword": "पासवर्ड बदलें",
	"eng-ps-err6":        "Error Pasword 1 Invalid",
	"spa-ps-err6":        "Error Contraseña 1 Inválida",
	"hin-ps-err6":        "त्रुटि पासवर्ड 1 अमान्य",
	"eng-ps-err7":        "Error Pasword 1 Does Not Meet Requirements",
	"spa-ps-err7":        "Error La contraseña 1 no cumple con los requisitos",
	"hin-ps-err7":        "त्रुटि पासवर्ड 1 आवश्यकताओं को पूरा नहीं करता है",
	"eng-ps-err8":        "Error Password 1 Does Not Match Password 2",
	"spa-ps-err8":        "Error La Contraseña 1 no Coincide con la Contraseña 2",
	"hin-ps-err8":        "त्रुटि पासवर्ड 1 पासवर्ड 2 से मेल नहीं खाता",
	"eng-ps-err9":        "Error Saving Password Hash",
	"spa-ps-err9":        "Error al Guardar el Hash de la Contraseña",
	"hin-ps-err9":        "पासवर्ड हैश सहेजने में त्रुटि",
	"eng-ps-err10":       "Error Reading Password Hash",
	"spa-ps-err10":       "Error al Leer el Hash de la Contraseña",
	"hin-ps-err10":       "पासवर्ड हैश पढ़ने में त्रुटि",
	"eng-ps-err11":       "Error Invalid Password",
	"spa-ps-err11":       "Error Contraseña Inválida",
	"hin-ps-err11":       "त्रुटि अमान्य पासवर्ड",
	"eng-ps-title1":      "Local File Encryption",
	"spa-ps-title1":      "Cifrado de Archivos Locales",
	"hin-ps-title1":      "स्थानीय फ़ाइल एन्क्रिप्शन",
	"eng-ps-title2":      "Enter Password To Reset",
	"spa-ps-title2":      "Ingrese la Contraseña para Restablecer",
	"hin-ps-title2":      "रीसेट करने के लिए पासवर्ड दर्ज करें",
	"eng-ps-title3":      "Enter New Password",
	"spa-ps-title3":      "Ingrese Nueva Clave",
	"hin-ps-title3":      "नया पासवर्ड दर्ज करें",
	"eng-ss-title":       "Settings",
	"spa-ss-title":       "Ajustes",
	"hin-ss-title":       "समायोजन",
	"eng-ss-ss":          "Change Settings",
	"spa-ss-ss":          "Cambiar Ajustes",
	"hin-ss-ss":          "सेटिंग्स परिवर्तित करना",
	"eng-ss-sserr":       "Settings Saved",
	"spa-ss-sserr":       "Ajustes guardados",
	"hin-ss-sserr":       "सेटिंग्स को सहेजा गया",
	"eng-ss-sserr1":      "Logon First",
	"spa-ss-sserr1":      "Iniciar Sesión Primero",
	"hin-ss-sserr1":      "लॉगऑन प्रथम",
	"eng-ss-la":          "Preferred Language",
	"spa-ss-la":          "Idioma Preferido",
	"hin-ss-la":          "पसंदीदा भाषा",
	"eng-ss-pl":          "Minimum Password Length",
	"spa-ss-pl":          "Longitud Mínima de la Contraseña",
	"hin-ss-pl":          "न्यूनतम पासवर्ड लंबाई",
	"eng-ss-ma":          "Message Max Age In Hours",
	"spa-ss-ma":          "Edad Máxima del Mensaje en Horas",
	"hin-ss-ma":          "संदेश अधिकतम आयु घंटों में",
	"eng-ss-mcletter":    "Password Must Contain Letter",
	"spa-ss-mcletter":    "La Contraseña Debe Contener una Letra",
	"hin-ss-mcletter":    "पासवर्ड में पत्र अवश्य होना चाहिए",
	"eng-ss-mcnumber":    "Password Must Contain Number",
	"spa-ss-mcnumber":    "La Contraseña Debe Contener un Número",
	"hin-ss-mcnumber":    "पासवर्ड में नंबर अवश्य होना चाहिए",
	"eng-ss-mcspecial":   "Password Must Contain Special",
	"spa-ss-mcspecial":   "La Contraseña Debe Contener Especial",
	"hin-ss-mcspecial":   "पासवर्ड में विशेष होना चाहिए",
	"eng-ss-heading":     "Change Settings",
	"spa-ss-heading":     "Cambiar Ajustes",
	"hin-ss-heading":     "सेटिंग्स परिवर्तित करना",
	"eng-cs-title":       "Certificates",
	"spa-cs-title":       "Certificados",
	"hin-cs-title":       "प्रमाण पत्र",
	"eng-cs-ca":          "CAROOT Certificate",
	"spa-cs-ca":          "Certificado CAROOT",
	"hin-cs-ca":          "कैरोट प्रमाणपत्र",
	"eng-cs-cc":          "CLIENT Certificate",
	"spa-cs-cc":          "Certificado CLIENTE",
	"hin-cs-cc":          "ग्राहक प्रमाणपत्र",
	"eng-cs-ck":          "CLIENT Key",
	"spa-cs-ck":          "Clave CLIENTE",
	"hin-cs-ck":          "ग्राहक कुंजी",
	"eng-cs-ss":          "Save Certificates",
	"spa-cs-ss":          "Guardar Certificados",
	"hin-cs-ss":          "प्रमाणपत्र सहेजें",
	"eng-cs-err1":        "Error CAROOT is Invalid",
	"spa-cs-err1":        "Error CAROOT no es Válido",
	"hin-cs-err1":        "त्रुटि CAROOT अमान्य है",
	"eng-cs-err2":        "Error CLIENT CERTIFICATE is invalid",
	"spa-cs-err2":        "Error CERTIFICADO DE CLIENTE no es Válido",
	"hin-cs-err2":        "त्रुटि क्लाइंट प्रमाणपत्र अमान्य है",
	"eng-cs-err3":        "Error CLIENT KEY is Invalid",
	"spa-cs-err3":        "Error CLAVE DE CLIENTE no es Cálida",
	"hin-cs-err3":        "त्रुटि क्लाइंट कुंजी अमान्य है",
	"eng-cs-heading":     "Certificate Management",
	"spa-cs-heading":     "Gestión de Certificados",
	"hin-cs-heading":     "प्रमाणपत्र प्रबंधन",
	"eng-cs-lf":          "Logon First",
	"spa-cs-lf":          "Iniciar Sesión Primero",
	"hin-cs-lf":          "लॉगऑन प्रथम",
	"eng-ls-title":       "Logon",
	"spa-ls-title":       "Iniciar sesión",
	"hin-ls-title":       "पर लॉग ऑन करें",
	"eng-ls-password":    "Password For Local Encryption",
	"spa-ls-password":    "Contraseña para el Cifrado Local",
	"hin-ls-password":    "स्थानीय एन्क्रिप्शन के लिए पासवर्ड",
	"eng-ls-alias":       "Alias",
	"spa-ls-alias":       "Alias",
	"hin-ls-alias":       "उपनाम",
	"eng-ls-queue":       "Queue",
	"spa-ls-queue":       "Cola",
	"hin-ls-queue":       "कतार",
	"eng-ls-queuepass":   "Queue Password",
	"spa-ls-queuepass":   "Contraseña de Cola",
	"hin-ls-queuepass":   "कतार पासवर्ड",
	"eng-ls-trypass":     "Try Password",
	"spa-ls-trypass":     "Probar Contraseña",
	"hin-ls-trypass":     "पासवर्ड आज़माएं",
	"eng-ls-con":         "Connected",
	"spa-ls-con":         "Conectada",
	"hin-ls-con":         "जुड़े हुए",
	"eng-ls-dis":         "Disconnected",
	"spa-ls-dis":         "Desconectada",
	"hin-ls-dis":         "डिस्कनेक्ट किया गया",
	"eng-ls-err1":        "Error Creating Password Hash 24",
	"spa-ls-err1":        "Error al Crear el Hash de la Contraseña 24",
	"hin-ls-err1":        "पासवर्ड हैश 24 बनाने में त्रुटि",
	"eng-ls-err2":        "Error Loading Password Hash 24",
	"spa-ls-err2":        "Error al Cargar el Hash de la Contraseña 24",
	"hin-ls-err2":        "पासवर्ड हैश 24 लोड करने में त्रुटि",
	"eng-ls-err3":        "Error Invalid Password",
	"spa-ls-err3":        "Error Contraseña no Válida",
	"hin-ls-err3":        "त्रुटि अमान्य पासवर्ड",
	"eng-ls-err4":        "Error URL Incorrect Format",
	"spa-ls-err4":        "URL de Error Formato Incorrecto",
	"hin-ls-err4":        "त्रुटि यूआरएल गलत प्रारूप",
	"eng-ls-err5":        "Error Invalid Queue Password 24",
	"spa-ls-err5":        "Error Contraseña de Cola no Válida 24",
	"hin-ls-err5":        "त्रुटि अमान्य कतार पासवर्ड 24",
	"eng-ls-err6-1":      "Error Queue Password Length is ",
	"spa-ls-err6-1":      "La Longitud de la Contraseña de la Cola de Errores es ",
	"hin-ls-err6-1":      "त्रुटि कतार पासवर्ड की लंबाई है ",
	"eng-ls-err6-2":      " should be length of 24",
	"spa-ls-err6-2":      " Debe Tener una Longitud de 24",
	"hin-ls-err6-2":      " लंबाई 24 होनी चाहिए",
	"eng-ls-err7":        "No NATS connection",
	"spa-ls-err7":        "Sin Conexión NATS",
	"hin-ls-err7":        "कोई NATS कनेक्शन नहीं",
	"eng-ls-erase":       "Security Erase",
	"spa-ls-erase":       "Borrado de seguridad",
	"hin-ls-erase":       "सुरक्षा मिटाएँ",
	"eng-ls-clogon":      "Communications Logon",
	"spa-ls-clogon":      "Inicio de Sesión de Comunicaciones",
	"hin-ls-clogon":      "संचार लॉगऑन",
	"eng-ls-err8":        "No JETSTREAM Connection",
	"spa-ls-err8":        "Sin Conexión JETSTREAM ",
	"hin-ls-err8":        "कोई जेटस्ट्रीम कनेक्शन नहीं",
	"eng-ms-title":       "Messages",
	"spa-ms-title":       "Mensajes",
	"hin-ms-title":       "संदेशों",
	"eng-ms-mm":          "Enter Message For Encryption",
	"spa-ms-mm":          "Ingrese el Mensaje Para el Cifrado",
	"hin-ms-mm":          "एन्क्रिप्शन के लिए संदेश दर्ज करें",
	"eng-ms-header1":     "Select An Item From The List",
	"spa-ms-header1":     "Seleccione un Elemento de la Lista",
	"hin-ms-header1":     "सूची से एक आइटम का चयन करें",
	"eng-ms-err1":        "NATS No Connection",
	"spa-ms-err1":        "NATS sin Conexión",
	"hin-ms-err1":        "NATS कोई कनेक्शन नहीं",
	"eng-ms-sm":          "Send",
	"spa-ms-sm":          "Enviar",
	"hin-ms-sm":          "भेजना",
	"eng-ms-filter":      "Omit Connected/Disconnected",
	"spa-ms-filter":      "Omitir Conectado/Desconectado",
	"hin-ms-filter":      "कनेक्टेड/डिस्कनेक्टेड को हटा दें",
	"eng-ms-header2":     "NATS Messaging",
	"spa-ms-header2":     "Mensajería NATS",
	"hin-ms-header2":     "NATS मैसेजिंग",
	"eng-ms-cpy":         "Copy To Clipboard",
	"spa-ms-cpy":         "Copiar al Portapapeles",
	"hin-ms-cpy":         "क्लिपबोर्ड पर कॉपी करें",
	"eng-ms-cpyf":        "Copy From Clipboard",
	"spa-ms-cpyf":        "Copiar desde el portapapeles",
	"hin-ms-cpyf":        "क्लिपबोर्ड से कॉपी करें",
	"eng-ms-carrier":     "Carrier",
	"spa-ms-carrier":     "Transportador",
	"वाहक-ms-carrier":    "Carrier",
	"eng-ms-err2":        "NATS No Connection ",
	"spa-ms-err2":        "NATS sin Conexión ",
	"hin-ms-err2":        "NATS कोई कनेक्शन नहीं ",
	"eng-ms-err3":        "Could Not Add Consumer ",
	"spa-ms-err3":        "No se Pudo Agregar el Consumidor ",
	"hin-ms-err3":        "उपभोक्ता नहीं जोड़ा जा सका ",
	"eng-ms-err4":        "Error Pull Subscribe ",
	"spa-ms-err4":        "Error Extraer Suscribirse ",
	"hin-ms-err4":        "त्रुटि खींचो सदस्यता लें ",
	"eng-ms-err5":        "Error Fetch ",
	"spa-ms-err5":        "Recuperación de Errores ",
	"hin-ms-err5":        "लाने में त्रुटि ",
	"eng-ms-err6-1":      "Recieved ",
	"spa-ms-err6-1":      "Recibida ",
	"hin-ms-err6-1":      "प्राप्त ",
	"eng-ms-err6-2":      " Messages ",
	"spa-ms-err6-2":      " Mensajes ",
	"hin-ms-err6-2":      " संदेशों ",
	"eng-ms-err6-3":      " Logs",
	"spa-ms-err6-3":      " Registros",
	"hin-ms-err6-3":      " लॉग्स",
	"eng-ms-err7":        "Please Logon First",
	"spa-ms-err7":        "Por Favor Ingresa Primero",
	"hin-ms-err7":        "कृपया पहले लॉगऑन करें",
	"eng-ms-nhn":         "No Host Name ",
	"spa-ms-nhn":         "Sin Nombre de Host ",
	"hin-ms-nhn":         "कोई होस्ट नाम नहीं ",
	"eng-ms-hn":          "Host ",
	"spa-ms-hn":          "Nombre de Host ",
	"hin-ms-hn":          "मेज़बान ",
	"eng-ms-mi":          "Mac IDS",
	"spa-ms-mi":          "ID de Mac",
	"hin-ms-mi":          "मैक आईडीएस",
	"eng-ms-ad":          "Address",
	"spa-ms-ad":          "Direccion",
	"hin-ms-ad":          "पता",
	"eng-ms-ni":          "Node Id - ",
	"spa-ms-ni":          "ID de Nodo - ",
	"hin-ms-ni":          "नोड आईडी - ",
	"eng-ms-msg":         "Message Id - ",
	"spa-ms-msg":         "ID de Mensaje - ",
	"hin-ms-msg":         "संदेश आईडी - ",
	"eng-ms-on":          "On - ",
	"spa-ms-on":          "En - ",
	"hin-ms-on":          "पर - ",
	"eng-ms-unk":         "Unknown",
	"spa-ms-unk":         "Desconocida",
	"hin-ms-unk":         "अज्ञात",
	"eng-ms-era":         "Erasing",
	"spa-ms-era":         "Borrando",
	"hin-ms-era":         "निकाली जा रही है",
	"eng-ms-erac":        "Erase Connection ",
	"spa-ms-erac":        "Borrar Conexión ",
	"hin-ms-erac":        "कनेक्शन मिटाएँ ",
	"eng-ms-eraj":        "Erase JetStream ",
	"spa-ms-eraj":        "Borrar JetStream ",
	"hin-ms-eraj":        "जेटस्ट्रीम मिटाएँ ",
	"eng-ms-dels":        "Delete JetStream ",
	"spa-ms-dels":        "Eliminar Secuencia ",
	"hin-ms-dels":        "जेटस्ट्रीम हटाएं ",
	"eng-ms-adds":        "Add Stream ",
	"spa-ms-adds":        "Agregar secuencia ",
	"hin-ms-adds":        "स्ट्रीम जोड़ें ",
	"eng-ms-addc":        "Add Consumer ",
	"spa-ms-addc":        "Agregar Consumidora ",
	"hin-ms-addc":        "उपभोक्ता जोड़ें ",
	"eng-ms-sece":        "Security Erase ",
	"spa-ms-sece":        "Borrado de Seguridad ",
	"hin-ms-sece":        "सुरक्षा मिटाएँ ",
	"eng-es-title":       "Enc/Dec",
	"spa-es-title":       "Codificar/Descodificar",
	"hin-es-title":       "एन/दिसंबर",
	"eng-es-pw":          "Enc/Dec",
	"spa-es-pw":          "Codificar/Descodificar",
	"hin-es-pw":          "एन/दिसंबर",
	"eng-es-pass":        "Enter Password to Use For Encryption 24",
	"spa-es-pass":        "Ingrese la contraseña para usar para el cifrado 24",
	"hin-es-pass":        "एन्क्रिप्शन 24 के लिए उपयोग करने के लिए पासवर्ड दर्ज करें",
	"eng-es-mv":          "Enter Value",
	"spa-es-mv":          "Introducir Valor",
	"hin-es-mv":          "मान दर्ज करें",
	"eng-es-mo":          "Output Shows Up Here",
	"spa-es-mo":          "La Salida Aparece Aquí",
	"hin-es-mo":          "आउटपुट यहां दिखता है",
	"eng-es-em":          "Encrypt Message",
	"spa-es-em":          "Cifrar Mensaje",
	"hin-es-em":          "संदेश एन्क्रिप्ट करें",
	"eng-es-dm":          "Decrypt Message",
	"spa-es-dm":          "Descifrar mensaje",
	"hin-es-dm":          "भाषा त्रुटि नहीं मिली",
	"eng-es-err1":        "Error Invalid Password",
	"spa-es-err1":        "Error Contraseña Inválida",
	"hin-es-err1":        "त्रुटि अमान्य पासवर्ड",
	"eng-es-err2-1":      "Error Password Length is ",
	"spa-es-err2-1":      "La Longitud de la Contraseña de Error es ",
	"hin-es-err2-1":      "त्रुटि पासवर्ड की लंबाई है ",
	"eng-es-err2-2":      " Should be Length of 24",
	"spa-es-err2-2":      " Debe Tener una Longitud de 24",
	"hin-es-err2-2":      " लंबाई 24 होनी चाहिए",
	"eng-es-err3":        "Error Input Text",
	"spa-es-err3":        "Texto de Entrada de Error",
	"hin-es-err3":        "त्रुटि इनपुट पाठ",
	"eng-es-err4":        "Cannot Encrypt Input Text",
	"spa-es-err4":        "No se Puede Cifrar el Texto de Entrada",
	"hin-es-err4":        "इनपुट टेक्स्ट को एन्क्रिप्ट नहीं किया जा सकता",
	"eng-es-err5":        "Cannot Decrypt Input Text",
	"spa-es-err5":        "No se Puede Descifrar el Texto Ingresado",
	"hin-es-err5":        "इनपुट टेक्स्ट को डिक्रिप्ट नहीं किया जा सकता",
	"eng-es-head0":       "24 Character Password",
	"spa-es-head0":       "Contraseña de 24 Caracteres",
	"hin-es-head0":       "24 कैरेक्टर का पासवर्ड",
	"eng-es-head1":       "Input",
	"spa-es-head1":       "Aporte",
	"hin-es-head1":       "इनपुट",
	"eng-es-head2":       "Output",
	"spa-es-head2":       "Producción",
	"hin-es-head2":       "उत्पादन",
	"eng-log-nc":         "Path to nNATS Configuration File \nFor TLS Certificate Paths",
	"spa-log-nc":         "Ruta al Archivo de Configuración nNATS \nPara Rutas de Certificados TLS",
	"hin-log-nc":         "नेट्स कॉन्फ़िगरेशन फ़ाइल का पथ \nटीएलएस प्रमाणपत्र पथों के लिए",
	"eng-hash-err1":      "Hash Error on Write",
	"spa-hash-err1":      "Error de Hash al Escribir ",
	"shin-hash-err1":     "लिखते समय हैश त्रुटि ",
	"eng-hash-err2":      "Error Creating Password Hash 24",
	"spa-hash-err2":      "Error al Crear el Hash de la Contraseña 24",
	"hin-hash-err2":      "पासवर्ड हैश 24 बनाने में त्रुटि",
	"eng-hash-err3":      "Hash Error on Read ",
	"spa-hash-err3":      "Error de Hash en Lectura ",
	"hin-hash-err3":      "पढ़ने पर हैश त्रुटि ",
	"eng-lang-err1":      "Language Error not Found ",
	"spa-lang-err1":      "Error de Idioma no Encontrado ",
	"hin-lang-err1":      "भाषा त्रुटि नहीं मिली ",
}

// do translation
func GetLangs(c string) string {
	value, err := MyLangs[nhpref.PreferedLanguage+"-"+c]
	if err == false {
		return "lang-error" + " " + nhpref.PreferedLanguage + "-" + c
	}

	return value
}
