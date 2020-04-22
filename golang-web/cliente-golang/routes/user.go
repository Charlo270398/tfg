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

func menuUserHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	userId, _ := session.Values["userId"].(string)
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

	var rolesList []int
	//Certificado
	client := GetTLSClient()

	//Request al servidor para recibir lista de roles
	response, err := client.Get(SERVER_URL + "/rol/list/user?userId=" + userId)
	if err != nil {
		util.PrintErrorLog(err)
	}
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&rolesList)
		if rolesList == nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/index.html", "public/templates/layouts/base.html"),
	)
	//Si solo hay un rol redirigimos a la pagina de ese rol
	if len(rolesList) == 1 {
		switch rolesList[0] {
		case 1: //Paciente
			http.Redirect(w, req, "/user/patient", http.StatusSeeOther)
		case 2: //Enfermero
			http.Redirect(w, req, "/user/nurse", http.StatusSeeOther)
		case 3: //Medico
			http.Redirect(w, req, "/user/doctor", http.StatusSeeOther)
		case 4: //AdminC
			http.Redirect(w, req, "/user/admin", http.StatusSeeOther)
		case 5: //AdminG
			http.Redirect(w, req, "/user/adminG", http.StatusSeeOther)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	//Si no cargamos la pagina de roles
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Usuario", Body: "body", UserRoles: rolesList}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func menuEditUserFormHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/user/menu/edit.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Cambiar datos", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func menuEditUserHandler(w http.ResponseWriter, req *http.Request) {
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
	var responseJSON JSON_Return
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	locJson, err := json.Marshal(util.JSON_user_SERVIDOR{Identificacion: creds.Identificacion, Nombre: creds.Nombre, Apellidos: creds.Apellidos,
		Email: creds.Email})

	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar usuario
	response, err := client.Post(SERVER_URL+"/user/menu/edit", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&responseJSON)
		js, err := json.Marshal(responseJSON)
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

func historialUserHandler(w http.ResponseWriter, req *http.Request) {
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/paciente/historial/index.html", "public/templates/layouts/menuUsuario.html", "public/templates/layouts/base.html"),
	)
	// Check user Token
	if !proveToken(req) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Historial", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
