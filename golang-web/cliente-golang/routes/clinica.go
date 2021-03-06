package routes

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	util "../utils"
)

//form añadir especialidad desde admin
func addClinicaFormGadminHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/clinica/addClinica.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Añadir clinica", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//form añadir especialidad desde admin
func addClinicaEspecialidadFormGadminHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/clinica/addEspecialidad.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Añadir clinica", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//listar clinicas
func getClinicaListGadminHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/clinica/list.html", "public/templates/layouts/base.html"),
	)

	page, ok := req.URL.Query()["page"]
	var pageString = "0"

	if ok {
		pageString = page[0]
	}
	//Certificado
	client := GetTLSClient()

	// Request /hello via the created HTTPS client over port 5001 via GET
	response, err := client.Get(SERVER_URL + "/clinica/list?page=" + pageString)
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		//Request al servidor para comprobar usuario/pass
		var serverReq util.Clinica_JSON_Pagination
		json.NewDecoder(response.Body).Decode(&serverReq)
		if err := tmp.ExecuteTemplate(w, "base", &util.ClinicaList_Page{Title: "Listado de clinicas", Body: "body", Page: serverReq.Page,
			NextPage: serverReq.NextPage, BeforePage: serverReq.BeforePage, ClinicaList: serverReq.ClinicaList}); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

//POST

//añadir clinica desde admin
func addClinicaGadminHandler(w http.ResponseWriter, req *http.Request) {
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
	var clinica util.Clinica_JSON
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&clinica)
	if err != nil {
		util.PrintErrorLog(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	locJson, err := json.Marshal(util.Clinica_JSON{Nombre: clinica.Nombre, Direccion: clinica.Direccion, Telefono: clinica.Telefono, UserToken: prepareUserToken(req)})

	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar clinica
	response, err := client.Post(SERVER_URL+"/clinica/add", "application/json", bytes.NewBuffer(locJson))
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

//añadir especialidad a una clinica
func addClinicaEspecialidadGadminHandler(w http.ResponseWriter, req *http.Request) {
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
	var clinica util.Clinica_JSON
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&clinica)
	if err != nil {
		util.PrintErrorLog(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	locJson, err := json.Marshal(util.Clinica_JSON{Nombre: clinica.Nombre, Direccion: clinica.Direccion, Telefono: clinica.Telefono, UserToken: prepareUserToken(req)})

	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar clinica
	response, err := client.Post(SERVER_URL+"/clinica/add", "application/json", bytes.NewBuffer(locJson))
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
