package routes

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	util "../utils"
)

//GET
func addEntradaHistorialFormMedicoHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/historial/addEntrada.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", Page{Title: "Añadir entrada historial", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func addAnaliticaHistorialFormMedicoHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/historial/addAnalitica.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", Page{Title: "Añadir analitica historial", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//POST
func addEntradaHistorialConsultaMedicoHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	}

	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}

	//Recuperamos datos del form
	var entradaHistorial util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entradaHistorial)

	//Certificado
	client := GetTLSClient()

	//Recuperamos clave publica del paciente
	var user util.User_JSON
	response, err := client.Get(SERVER_URL + "/user/pairkeys?userId=" + strconv.Itoa(entradaHistorial.PacienteId))
	if response != nil {
		json.NewDecoder(response.Body).Decode(&user)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Generamos una clave AES aleatoria de 256 bits para cifrar los datos sensibles
	AESkeyDatos := util.AEScreateKey()

	//Ciframos los datos sensibles con la clave
	entradaHistorial.JuicioDiagnostico, _ = util.AESencrypt(AESkeyDatos, entradaHistorial.JuicioDiagnostico)
	entradaHistorial.MotivoConsulta, _ = util.AESencrypt(AESkeyDatos, entradaHistorial.MotivoConsulta)

	//Pasamos la clave a base 64
	AESkeyBase64String := string(util.Base64Encode(AESkeyDatos))
	//Ciframos la clave AES usada con nuestra clave pública
	claveAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(user.PairKeys.PublicKey))

	//Preparamos los datos para enviar
	entradaHistorial.UserToken = prepareUserToken(req)
	entradaHistorial.Clave = claveAEScifrada
	locJson, err := json.Marshal(entradaHistorial)

	//Request al servidor para añadir entrada
	response, err = client.Post(SERVER_URL+"/user/doctor/citas/addEntrada", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		err := json.NewDecoder(response.Body).Decode(&result)
		js, err := json.Marshal(result)
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
