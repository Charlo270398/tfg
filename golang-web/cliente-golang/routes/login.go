package routes

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"unicode"

	util "../utils"
	"github.com/gorilla/sessions"
)

type JSON_Return struct {
	Result string
	Error  string
}

//GET

func loginIndexHandler(w http.ResponseWriter, req *http.Request) {
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/login/index.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Login", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func registerIndexHandler(w http.ResponseWriter, req *http.Request) {
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/login/register.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Registro", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//POST

func loginUserHandler(w http.ResponseWriter, req *http.Request) {
	var creds util.JSON_Credentials_CLIENTE

	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//SHA 512, cogemos la primera mitad
	sha_512 := sha512.New()
	sha_512.Write([]byte(creds.Password))
	hash512 := sha_512.Sum(nil)
	loginHash := make([]byte, len(hash512)-len(hash512)/2)
	privateKeyHash := make([]byte, len(hash512)-len(hash512)/2)

	//Dividimos el hash512 en 2 hashes, uno para login y otro para clave privada
	for index := range loginHash {
		loginHash[index] = hash512[index]
		privateKeyHash[index] = hash512[index+len(hash512)/2]
	}

	//Hacemos HASH del DNI para poder hacer busquedas despues
	sha_256 := sha256.New()
	sha_256.Write([]byte(creds.Identificacion))
	hash := sha_256.Sum(nil)
	identificacionHash := fmt.Sprintf("%x", hash) //Pasamos a hexadecimal el hash

	//USER JSON
	locJson, err := json.Marshal(util.JSON_Credentials_SERVIDOR{Identificacion: identificacionHash, Password: loginHash})

	//Certificado
	client := GetTLSClient()

	// Request /hello via the created HTTPS client over port 5001 via GET
	response, err := client.Post(SERVER_URL+"/login", "application/json", bytes.NewBuffer(locJson))
	if err != nil {
		log.Fatal(err)
	}

	if response != nil {
		//Respuesta del servidor
		var serverRes util.JSON_Login_Return
		json.NewDecoder(response.Body).Decode(&serverRes)
		jsonReturn := JSON_Return{"", ""}
		if serverRes.Error == "" {
			jsonReturn = JSON_Return{"Sesión Iniciada", ""}
			session, _ := store.Get(req, "userSession")
			session.Values["authenticated"] = true
			session.Values["userId"] = serverRes.UserId
			session.Values["userToken"] = serverRes.Token
			session.Values["userName"] = serverRes.Nombre
			session.Values["userSurname"] = serverRes.Apellidos
			session.Values["userEmail"] = serverRes.Email
			session.Values["userDataKey"] = serverRes.Clave
			session.Values["userPublicKey"] = serverRes.PairKeys.PublicKey
			session.Values["userPrivateKeyHash"] = privateKeyHash
			session.Options.MaxAge = 60 * 30
			session.Save(req, w)
		} else {
			jsonReturn = JSON_Return{"", "Usuario y contraseña incorrectos"}
		}
		js, err := json.Marshal(jsonReturn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func registerUserHandler(w http.ResponseWriter, req *http.Request) {
	var creds util.User_JSON_AddUsers
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		util.PrintErrorLog(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Comprobamos si la contraseña cumple los requisitos de seguridad
	rLen := false
	rDigit := false
	rSpecial := false

	//Longitud >= que 8
	rLen = len(creds.Password) >= 8

	//Al menos un dígito
	for _, c := range creds.Password {
		if unicode.IsDigit(c) {
			rDigit = true
			break
		}
	}

	//Un caracter especial
	rSpecial = strings.ContainsAny(creds.Password, "!#$%&()*+,-./:;<=>?@[]^_`{|}~")

	if !(rDigit && rLen && rSpecial) {
		util.PrintLog("La contraseña aportada no cumple los requisitos de seguridad")
		js, _ := json.Marshal(util.JSON_Return{Error: "La contraseña aportada no cumple los requisitos de seguridad"})
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}

	//SHA 512, cogemos la primera mitad
	sha_512 := sha512.New()
	sha_512.Write([]byte(creds.Password))
	hash512 := sha_512.Sum(nil)
	loginHash := make([]byte, len(hash512)-len(hash512)/2)
	privateKeyHash := make([]byte, len(hash512)-len(hash512)/2)

	//Dividimos el hash512 en 2 hashes, uno para login y otro para clave privada
	for index := range loginHash {
		loginHash[index] = hash512[index]
		privateKeyHash[index] = hash512[index+len(hash512)/2]
	}

	//Generamos par de claves RSA
	privK := util.RSAGenerateKeys()
	var masterPairKeys util.PairKeys
	masterPairKeys = getPublicMasterKey() //Obtenemos la CLAVE MAESTRA PUBLICA si existe
	if masterPairKeys.PublicKey == nil {
		//si no existe es que somos el primer usuario, entonces generamos nosotros la clave
		generatedMK := util.RSAGenerateKeys()
		masterPairKeys.PrivateKey = util.RSAPrivateKeyToBytes(generatedMK)
		masterPairKeys.PublicKey = util.RSAPublicKeyToBytes(&generatedMK.PublicKey)
		masterPairKeys.PrivateKey = util.RSAPrivateKeyToBytes(generatedMK)
	}

	//Pasamos las claves a []byte
	var pairKeys util.PairKeys
	pairKeys.PrivateKey = util.RSAPrivateKeyToBytes(privK)
	pairKeys.PublicKey = util.RSAPublicKeyToBytes(&privK.PublicKey)
	pairKeys.PrivateKey = util.RSAPrivateKeyToBytes(privK)

	//Ciframos clave privada con AES
	privKcifrada, _ := util.AESencrypt(privateKeyHash, string(util.Base64Encode(pairKeys.PrivateKey)))
	privKcifradaMaster, _ := util.AESencrypt(privateKeyHash, string(util.Base64Encode(masterPairKeys.PrivateKey)))
	pairKeys.PrivateKey = []byte(privKcifrada)
	masterPairKeys.PrivateKey = []byte(privKcifradaMaster)

	//Generamos una clave AES aleatoria de 256 bits para cifrar los datos sensibles
	AESkeyDatos := util.AEScreateKey()

	//Ciframos los datos sensibles con la clave
	identificacionCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Identificacion)
	nombreCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Nombre)
	apellidosCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Apellidos)
	emailCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Email)
	sexoCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Sexo)
	alergiasCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Alergias)

	//Hacemos HASH del DNI para poder hacer busquedas despues
	sha_256 := sha256.New()
	sha_256.Write([]byte(creds.Identificacion))
	hash := sha_256.Sum(nil)
	identificacionHash := fmt.Sprintf("%x", hash) //Pasamos a hexadecimal el hash

	//Pasamos la clave a base 64
	AESkeyBase64String := string(util.Base64Encode(AESkeyDatos))
	//Ciframos la clave AES usada con nuestra clave pública
	claveAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, privK.PublicKey)
	claveMaestraAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(masterPairKeys.PublicKey))

	locJson, err := json.Marshal(util.User_JSON{Identificacion: identificacionCifrado, IdentificacionHash: identificacionHash, Nombre: nombreCifrado, Apellidos: apellidosCifrado,
		Email: emailCifrado, Password: loginHash, PairKeys: pairKeys, MasterPairKeys: masterPairKeys, Clave: claveAEScifrada, Sexo: sexoCifrado, Alergias: alergiasCifrado, ClaveMaestra: claveMaestraAEScifrada})

	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar usuario
	response, err := client.Post(SERVER_URL+"/register", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		//Respuesta del servidor
		var serverRes util.JSON_Login_Return
		json.NewDecoder(response.Body).Decode(&serverRes)
		jsonReturn := JSON_Return{"", ""}
		if serverRes.Error == "" {
			jsonReturn = JSON_Return{"Sesión Iniciada", ""}
			session, _ := store.Get(req, "userSession")
			session.Values["authenticated"] = true
			session.Values["userId"] = serverRes.UserId
			session.Values["userToken"] = serverRes.Token
			session.Values["userName"] = serverRes.Nombre
			session.Values["userSurname"] = serverRes.Apellidos
			session.Values["userEmail"] = serverRes.Email
			session.Values["userPublicKey"] = pairKeys.PublicKey
			session.Values["userPrivateKeyHash"] = privateKeyHash
			session.Values["userDataKey"] = serverRes.Clave
			session.Options.MaxAge = 60 * 30
			session.Save(req, w)
		} else {
			jsonReturn = JSON_Return{"", "No se ha podido completar el registro"}
		}
		js, err := json.Marshal(jsonReturn)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)
)

func logoutUserHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")

	if session.Values["authenticated"] == true {
		// Revoke users authentication
		session.Values["authenticated"] = false
		session.Save(req, w)
	}
}
