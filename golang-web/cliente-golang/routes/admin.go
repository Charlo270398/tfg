package routes

import (
	"bytes"
	"crypto/sha512"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	util "../utils"
)

//GET
func menuAdminHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	// Check user Token
	if !proveToken(req) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	userToken := prepareUserToken(req)
	locJson, err := json.Marshal(util.JSON_Admin_Menu{UserToken: userToken})
	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar usuario
	response, err := client.Post(SERVER_URL+"/user/admin", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var JSON_Admin_Menu util.JSON_Admin_Menu
		err := json.NewDecoder(response.Body).Decode(&JSON_Admin_Menu)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		} else {
			if JSON_Admin_Menu.Error != "" {
				util.PrintLog(JSON_Admin_Menu.Error)
				http.Error(w, JSON_Admin_Menu.Error, http.StatusInternalServerError)
				return
			} else {
				session.Values["clinicaId"] = JSON_Admin_Menu.Clinica.Id
				session.Save(req, w)
				var tmp = template.Must(
					template.New("").ParseFiles("public/templates/user/admin/index.html", "public/templates/layouts/base.html"),
				)
				if err := tmp.ExecuteTemplate(w, "base", &util.PageMenuAdmin{Title: "Menú administrador clínica", Body: "body", Clinica: JSON_Admin_Menu.Clinica}); err != nil {
					log.Printf("Error executing template: %v", err)
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func adminAddMedicoFormHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	// Check user Token
	if !proveToken(req) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	clinica := util.Clinica{Id: session.Values["clinicaId"].(int)}

	//Certificado
	client := GetTLSClient()

	//Cargamos la lista de especialidades
	response, err := client.Get(SERVER_URL + "/especialidad/list/query")
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var listaEspecialidades []util.Especialidad_JSON
	json.NewDecoder(response.Body).Decode(&listaEspecialidades)

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/admin/addMedico.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.PageAdminAddMedico{Title: "Menú administrador clínica", Body: "body", Clinica: clinica, Especialidades: listaEspecialidades}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func adminAddEnfermeroFormHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	// Check user Token
	if !proveToken(req) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	clinica := util.Clinica{Id: session.Values["clinicaId"].(int)}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/admin/addEnfermero.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.PageAdminAddMedico{Title: "Menú administrador clínica", Body: "body", Clinica: clinica}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//POST

//post añadir enfermero desde admin
func addEnfermeroAdminHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	// Check user Token
	if !proveToken(req) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var creds util.User_JSON
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		util.PrintErrorLog(err)
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

	//Generamos par de claves RSA
	privK := util.RSAGenerateKeys()
	//Pasamos las claves a []byte
	var pairKeys util.PairKeys
	pairKeys.PrivateKey = util.RSAPrivateKeyToBytes(privK)
	pairKeys.PublicKey = util.RSAPublicKeyToBytes(&privK.PublicKey)
	pairKeys.PrivateKey = util.RSAPrivateKeyToBytes(privK)
	//Ciframos clave privada con AES
	privKcifrada, _ := util.AESencrypt(privateKeyHash, string(pairKeys.PrivateKey))
	pairKeys.PrivateKey = []byte(privKcifrada)

	locJson, err := json.Marshal(util.User_JSON{Identificacion: creds.Identificacion, Nombre: creds.Nombre, Apellidos: creds.Apellidos,
		Email: creds.Email, Password: loginHash, Roles: []int{Rol_enfermero.Id}, EnfermeroClinica: creds.EnfermeroClinica,
		UserToken: prepareUserToken(req), PairKeys: pairKeys})

	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar usuario
	response, err := client.Post(SERVER_URL+"/user/admin/nurse/add", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var responseJSON JSON_Return
		err := json.NewDecoder(response.Body).Decode(&responseJSON)
		js, err := json.Marshal(responseJSON)
		if err != nil {
			util.PrintErrorLog(err)
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

//post añadir medico desde admin
func addMedicoAdminHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}
	// Check user Token
	if !proveToken(req) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var creds util.User_JSON
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		util.PrintErrorLog(err)
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

	//Generamos par de claves RSA
	privK := util.RSAGenerateKeys()
	//Pasamos las claves a []byte
	var pairKeys util.PairKeys
	pairKeys.PrivateKey = util.RSAPrivateKeyToBytes(privK)
	pairKeys.PublicKey = util.RSAPublicKeyToBytes(&privK.PublicKey)
	pairKeys.PrivateKey = util.RSAPrivateKeyToBytes(privK)
	//Ciframos clave privada con AES
	privKcifrada, _ := util.AESencrypt(privateKeyHash, string(pairKeys.PrivateKey))
	pairKeys.PrivateKey = []byte(privKcifrada)

	locJson, err := json.Marshal(util.User_JSON{Identificacion: creds.Identificacion, Nombre: creds.Nombre, Apellidos: creds.Apellidos,
		Email: creds.Email, Password: loginHash, Roles: []int{Rol_medico.Id}, MedicoClinica: creds.MedicoClinica,
		MedicoEspecialidad: creds.MedicoEspecialidad, UserToken: prepareUserToken(req), PairKeys: pairKeys})

	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar usuario
	response, err := client.Post(SERVER_URL+"/user/admin/doctor/add", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var responseJSON JSON_Return
		err := json.NewDecoder(response.Body).Decode(&responseJSON)
		js, err := json.Marshal(responseJSON)
		if err != nil {
			util.PrintErrorLog(err)
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
