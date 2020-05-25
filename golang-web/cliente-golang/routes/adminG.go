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
	"strconv"

	util "../utils"
	"github.com/gorilla/mux"
)

//GET
func menuAdminGHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}
	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/adminG/index.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Menú administrador", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getUserListAdminGHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}
	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/list.html", "public/templates/layouts/base.html"),
	)

	page, ok := req.URL.Query()["page"]
	var pageString = "0"

	if ok {
		pageString = page[0]
	}
	//Certificado
	client := GetTLSClient()

	// Request /hello via the created HTTPS client over port 5001 via GET
	response, err := client.Get(SERVER_URL + "/user/adminG/userList?page=" + pageString)
	if err != nil {
		util.PrintErrorLog(err)
	} else {
		//Request al servidor para comprobar usuario/pass
		var serverReq util.UserList_JSON_Pagination
		json.NewDecoder(response.Body).Decode(&serverReq)

		//DESCIFRADO MASTER KEY
		//Recuperamos nuestra clave privada cifrada
		userId, _ := session.Values["userId"].(string)
		userMasterPairkeys := getUserMasterPairKeys(userId)
		userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)

		//Desciframos nuestra clave privada cifrada con AES
		userPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(userMasterPairkeys.PrivateKey))
		userPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(userPrivateKeyString)))

		//DESCIFRADO DE DATOS
		for index, user := range serverReq.UserList {
			//Desciframos la clave AES de los datos cifrados
			claveAESuser := util.RSADecryptOAEP(user.ClaveMaestra, *userPrivateKey)
			claveAESuserByte := util.Base64Decode([]byte(claveAESuser))
			//Desciframos los datos del historial con AES
			serverReq.UserList[index].Nombre, _ = util.AESdecrypt(claveAESuserByte, user.Nombre)
			serverReq.UserList[index].Apellidos, _ = util.AESdecrypt(claveAESuserByte, user.Apellidos)
			serverReq.UserList[index].Email, _ = util.AESdecrypt(claveAESuserByte, user.Email)
		}

		if err := tmp.ExecuteTemplate(w, "base", &util.UserList_Page{Title: "Listado de usuarios", Body: "body", Page: serverReq.Page,
			NextPage: serverReq.NextPage, BeforePage: serverReq.BeforePage, UserList: serverReq.UserList}); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

//form añadir usuario desde admin
func addUserFormGadminHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}
	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}

	//Certificado
	client := GetTLSClient()

	//Cargamos la lista de clinicas
	response, err := client.Get(SERVER_URL + "/clinica/list/query")
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var listaClinicas []util.Clinica_JSON
	json.NewDecoder(response.Body).Decode(&listaClinicas)

	//Cargamos la lista de especialidades
	response, err = client.Get(SERVER_URL + "/especialidad/list/query")
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var listaEspecialidades []util.Especialidad_JSON
	json.NewDecoder(response.Body).Decode(&listaEspecialidades)

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/addUser.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.CitaPage{Title: "Añadir usuario", Body: "body", Clinicas: listaClinicas, Especialidades: listaEspecialidades}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//POST

//post añadir usuario desde admin
func addUserGadminHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}
	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}
	var creds util.User_JSON_AddUsers
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

	//Recuperamos la CLAVE MAESTRA
	userId, _ := session.Values["userId"].(string)
	masterPairKeys := getUserMasterPairKeys(userId)

	//Si es usuario de emergencias o administrador global le damos la CLAVE MAESTRA
	tieneCM := false
	for _, role := range creds.Roles {
		if role == Rol_emergencias.Id || role == Rol_administradorG.Id {
			//Desciframos la clave privada CLAVE MAESTRA cifrada con AES
			userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)
			masterPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(masterPairKeys.PrivateKey))

			//La ciframos con el hash del usuario
			newUserMasterPrivateKeyString, _ := util.AESencrypt(privateKeyHash, masterPrivateKeyString)
			masterPairKeys.PrivateKey = util.Base64Encode([]byte(newUserMasterPrivateKeyString))
			tieneCM = true
		}
	}

	//Si no tiene esos roles borramos la CLAVE MAESTRA PRIVADA
	if !tieneCM {
		masterPairKeys.PrivateKey = nil
	}

	//Pasamos las claves a []byte
	var pairKeys util.PairKeys
	pairKeys.PrivateKey = util.RSAPrivateKeyToBytes(privK)
	pairKeys.PublicKey = util.RSAPublicKeyToBytes(&privK.PublicKey)
	pairKeys.PrivateKey = util.RSAPrivateKeyToBytes(privK)
	//Ciframos clave privada con AES
	privKcifrada, _ := util.AESencrypt(privateKeyHash, string(util.Base64Encode(pairKeys.PrivateKey)))
	pairKeys.PrivateKey = []byte(privKcifrada)

	//Generamos una clave AES aleatoria de 256 bits para cifrar los datos sensibles
	AESkeyDatos := util.AEScreateKey()

	//Ciframos los datos sensibles con la clave
	identificacionCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Identificacion)
	emailCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Email)
	nombreCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Nombre)
	apellidosCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Apellidos)

	//Hacemos HASH del DNI para poder hacer busquedas despues
	sha_256 := sha256.New()
	sha_256.Write([]byte(creds.Identificacion))
	hash := sha_256.Sum(nil)
	identificacionHash := fmt.Sprintf("%x", hash) //Pasamos a hexadecimal el hash

	//Pasamos la clave a base 64
	AESkeyBase64String := string(util.Base64Encode(AESkeyDatos))
	//Ciframos la clave AES usada con nuestra clave pública
	claveAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, privK.PublicKey)
	claveAEScifradaMaestra := util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(masterPairKeys.PublicKey)) //Ciframos con la clave maestra

	nombreDoctor := creds.Nombre + " " + creds.Apellidos
	locJson, err := json.Marshal(util.User_JSON{Identificacion: identificacionCifrado, Nombre: nombreCifrado, Apellidos: apellidosCifrado,
		Email: emailCifrado, Password: loginHash, Roles: creds.Roles, EnfermeroClinica: creds.EnfermeroClinica, MedicoClinica: creds.MedicoClinica,
		AdminClinica: creds.AdminClinica, MedicoEspecialidad: creds.MedicoEspecialidad, UserToken: prepareUserToken(req), PairKeys: pairKeys,
		IdentificacionHash: identificacionHash, NombreDoctor: nombreDoctor, Clave: claveAEScifrada, ClaveMaestra: claveAEScifradaMaestra, MasterPairKeys: masterPairKeys})

	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar usuario
	response, err := client.Post(SERVER_URL+"/user/adminG/userList/add", "application/json", bytes.NewBuffer(locJson))
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

//DELETE

func deleteUserHandler(w http.ResponseWriter, req *http.Request) {

	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}
	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}
	//Cargamos el ID del usuario en la url
	vars := mux.Vars(req)
	userId_int, _ := strconv.Atoi(vars["userId"])
	locJson, err := json.Marshal(util.User_id_JSON{Id: userId_int, UserToken: prepareUserToken(req)})

	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar usuario
	response, err := client.Post(SERVER_URL+"/user/delete", "application/json", bytes.NewBuffer(locJson))
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
