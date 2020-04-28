package routes

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
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
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/index.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Menú médico", Body: "body"}); err != nil {
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
		http.Error(w, "Forbidden", http.StatusForbidden)
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
		http.Error(w, "Forbidden", http.StatusForbidden)
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
		http.Error(w, "Forbidden", http.StatusForbidden)
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
		http.Error(w, "Forbidden", http.StatusForbidden)
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
