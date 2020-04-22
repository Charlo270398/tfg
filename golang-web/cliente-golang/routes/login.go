package routes

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

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
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Register", Body: "body"}); err != nil {
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
	//USER JSON
	locJson, err := json.Marshal(util.JSON_Credentials_SERVIDOR{Email: creds.Email, Password: loginHash})

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
	var creds util.User_JSON
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

	locJson, err := json.Marshal(util.JSON_user_SERVIDOR{Identificacion: creds.Identificacion, Nombre: creds.Nombre, Apellidos: creds.Apellidos,
		Email: creds.Email, Password: loginHash})

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
