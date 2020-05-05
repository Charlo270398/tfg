package routes

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	util "../utils"
)

//GET
func menuMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	locJson, err := json.Marshal(prepareUserToken(req))

	//Certificado
	client := GetTLSClient()
	var cita util.CitaJSON

	//Request al servidor para obtener citas futuras
	response, err := client.Post(SERVER_URL+"/user/doctor/citas/actual", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&cita)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/index.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.PageMenuMedico{Title: "Menú médico", Body: "body", CitaActual: cita.Id}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func solicitarHistorialMedicoFormHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/user/medico/historial/solicitar.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Solicitar historial", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func solicitarHistorialMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	//Preparamos datos request
	var user util.User
	json.NewDecoder(req.Body).Decode(&user)
	var solicitarHistorial util.SolicitarHistorial_JSON
	solicitarHistorial.UserDNI = user.Identificacion
	solicitarHistorial.UserToken = prepareUserToken(req)
	locJson, err := json.Marshal(solicitarHistorial)

	//Certificado
	client := GetTLSClient()

	//Request al servidor para recibir lista de roles
	response, err := client.Post(SERVER_URL+"/user/doctor/historial/solicitar", "application/json", bytes.NewBuffer(locJson))
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

//Listar medicos dada una especialidad de la clinica
func getMedicoDiasDisponiblesHandler(w http.ResponseWriter, req *http.Request) {
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
	medicoId, _ := req.URL.Query()["doctorId"]
	//Certificado
	client := GetTLSClient()

	response, err := client.Get(SERVER_URL + "/user/doctor/disponible/dia?doctorId=" + medicoId[0])
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		var citasListJSON []util.Cita
		err := json.NewDecoder(response.Body).Decode(&citasListJSON)
		js, err := json.Marshal(citasListJSON)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//Listar medicos dada una especialidad de la clinica
func getMedicoHorasDiaDisponiblesHandler(w http.ResponseWriter, req *http.Request) {
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
	medicoId, _ := req.URL.Query()["doctorId"]
	dia, _ := req.URL.Query()["dia"]
	//Certificado
	client := GetTLSClient()
	diaTratado := strings.Replace(dia[0], " ", "", -1)
	response, err := client.Get(SERVER_URL + "/user/doctor/disponible/hora?doctorId=" + medicoId[0] + "&dia=" + diaTratado)
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		var citasListJSON []util.Cita
		err := json.NewDecoder(response.Body).Decode(&citasListJSON)
		js, err := json.Marshal(citasListJSON)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func medicoCitaListHandler(w http.ResponseWriter, req *http.Request) {

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

	locJson, err := json.Marshal(prepareUserToken(req))

	//Certificado
	client := GetTLSClient()
	var citasList []util.CitaJSON

	//Request al servidor para obtener citas futuras
	response, err := client.Post(SERVER_URL+"/user/doctor/citas/list", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&citasList)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/citas/list.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.CitaListPage{Title: "Citas pendientes", Body: "body", Citas: citasList}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//GET
func getCitaFormMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	//Preparamos Id cita
	citaId, _ := req.URL.Query()["citaId"]
	userId_int, _ := strconv.Atoi(citaId[0])

	//Certificado
	client := GetTLSClient()
	var cita util.CitaJSON
	cita.UserToken = prepareUserToken(req)
	cita.Id = userId_int
	locJson, err := json.Marshal(cita)

	//Request al servidor para obtener la cita
	response, err := client.Post(SERVER_URL+"/user/doctor/citas", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&cita)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/citas/index.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", util.ConsultaPage{Title: "Pasar consulta", Body: "body", Cita: cita}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getListHistorialMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	locJson, err := json.Marshal(prepareUserToken(req))

	//Certificado
	client := GetTLSClient()
	var historialList []util.Historial_JSON

	//Request al servidor para obtener historiales compartidos
	response, err := client.Post(SERVER_URL+"/user/doctor/historial/list", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&historialList)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Recuperamos nuestra clave privada cifrada
	userId, _ := session.Values["userId"].(string)
	userPairkeys := getUserPairKeys(userId)
	userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)

	//Desciframos nuestra clave privada cifrada con AES
	userPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(userPairkeys.PrivateKey))
	userPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(userPrivateKeyString)))

	//DESCIFRADO DE DATOS
	for index, historial := range historialList {
		//Desciframos la clave AES de los datos cifrados
		claveAEShistorial := util.RSADecryptOAEP(historial.Clave, *userPrivateKey)
		claveAEShistorialByte := util.Base64Decode([]byte(claveAEShistorial))
		//Desciframos los datos del historial con AES
		historialList[index].NombrePaciente, _ = util.AESdecrypt(claveAEShistorialByte, historial.NombrePaciente)
		historialList[index].Sexo, _ = util.AESdecrypt(claveAEShistorialByte, historial.Sexo)
	}
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/historial/list.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.HistorialListPage{Title: "Historiales compartidos", Body: "body", Historiales: historialList}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
