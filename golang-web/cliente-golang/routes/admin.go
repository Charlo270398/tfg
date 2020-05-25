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

	util "../utils"
)

//GET
func menuAdminHandler(w http.ResponseWriter, req *http.Request) {
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
	masterPairKeys := getPublicMasterKey() //Obtenemos la CLAVE MAESTRA PUBLICA si existe

	//Pasamos las claves a []byte
	var pairKeys util.PairKeys
	pairKeys.PrivateKey = util.RSAPrivateKeyToBytes(privK)
	pairKeys.PublicKey = util.RSAPublicKeyToBytes(&privK.PublicKey)
	pairKeys.PrivateKey = util.RSAPrivateKeyToBytes(privK)

	//Generamos una clave AES aleatoria de 256 bits para cifrar los datos sensibles
	AESkeyDatos := util.AEScreateKey()

	//Ciframos los datos sensibles con la clave
	identificacionCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Identificacion)
	nombreCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Nombre)
	apellidosCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Apellidos)
	emailCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Email)

	//Pasamos la clave a base 64
	AESkeyBase64String := string(util.Base64Encode(AESkeyDatos))

	//Ciframos clave privada con AES
	privKcifrada, _ := util.AESencrypt(privateKeyHash, string(util.Base64Encode(pairKeys.PrivateKey)))
	pairKeys.PrivateKey = []byte(privKcifrada)

	//Ciframos la clave AES usada con nuestra clave pública
	claveAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, privK.PublicKey)
	claveMaestraAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(masterPairKeys.PublicKey))

	locJson, err := json.Marshal(util.User_JSON{Identificacion: identificacionCifrado, Nombre: nombreCifrado, Apellidos: apellidosCifrado,
		Email: emailCifrado, Password: loginHash, Roles: []int{Rol_enfermero.Id}, EnfermeroClinica: creds.EnfermeroClinica,
		Clave: claveAEScifrada, ClaveMaestra: claveMaestraAEScifrada, UserToken: prepareUserToken(req), PairKeys: pairKeys})

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
	masterPairKeys := getPublicMasterKey() //Obtenemos la CLAVE MAESTRA PUBLICA si existe

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
	nombreCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Nombre)
	apellidosCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Apellidos)
	emailCifrado, _ := util.AESencrypt(AESkeyDatos, creds.Email)

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

	nombreDoctor := creds.Nombre + " " + creds.Apellidos
	locJson, err := json.Marshal(util.User_JSON{Identificacion: identificacionCifrado, Nombre: nombreCifrado, Apellidos: apellidosCifrado,
		Email: emailCifrado, Password: loginHash, Roles: []int{Rol_medico.Id}, MedicoClinica: creds.MedicoClinica,
		MedicoEspecialidad: creds.MedicoEspecialidad, UserToken: prepareUserToken(req), PairKeys: pairKeys,
		IdentificacionHash: identificacionHash, Clave: claveAEScifrada, ClaveMaestra: claveMaestraAEScifrada, NombreDoctor: nombreDoctor})

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
