package routes

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	util "../utils"
)

//GET

//form añadir especialidad desde admin
func addEspecialidadFormGadminHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/especialidad/addEspecialidad.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Añadir especialidad", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//listar especialidades
func getEspecialidadListGadminHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/especialidad/list.html", "public/templates/layouts/base.html"),
	)
	page, ok := req.URL.Query()["page"]
	var pageString = "0"

	if ok {
		pageString = page[0]
	}
	//Certificado
	client := GetTLSClient()

	// Request /hello via the created HTTPS client over port 5001 via GET
	response, err := client.Get(SERVER_URL + "/especialidad/list?page=" + pageString)
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		//Request al servidor para comprobar usuario/pass
		var serverReq util.Especialidad_JSON_Pagination
		json.NewDecoder(response.Body).Decode(&serverReq)
		if err := tmp.ExecuteTemplate(w, "base", &util.EspecialidadList_Page{Title: "Listado de especialidades", Body: "body", Page: serverReq.Page,
			NextPage: serverReq.NextPage, BeforePage: serverReq.BeforePage, EspecialidadList: serverReq.EspecialidadList}); err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

//POST

//añadir especialidad desde admin
func addEspecialidadGadminHandler(w http.ResponseWriter, req *http.Request) {
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
	var especialidad util.Especialidad_JSON
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&especialidad)
	if err != nil {
		util.PrintErrorLog(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	locJson, err := json.Marshal(util.Especialidad_JSON{Nombre: especialidad.Nombre, UserToken: prepareUserToken(req)})

	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar usuario
	response, err := client.Post(SERVER_URL+"/especialidad/add", "application/json", bytes.NewBuffer(locJson))
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
