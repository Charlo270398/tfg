package routes

import (
	"encoding/json"
	"net/http"

	util "../utils"
)

var Rol_paciente util.Rol
var Rol_enfermero util.Rol
var Rol_medico util.Rol
var Rol_administradorC util.Rol
var Rol_administradorG util.Rol
var Rol_emergencias util.Rol

func LoadRoles() {
	//Definimos los roles basicos
	Rol_paciente = util.Rol{Id: 1, Nombre: "paciente", Descripcion: "Paciente"}
	Rol_enfermero = util.Rol{Id: 2, Nombre: "enfermero", Descripcion: "Enfermero"}
	Rol_medico = util.Rol{Id: 3, Nombre: "medico", Descripcion: "Medico"}
	Rol_administradorC = util.Rol{Id: 4, Nombre: "administradorC", Descripcion: "Administrador clinica"}
	Rol_administradorG = util.Rol{Id: 5, Nombre: "administradorG", Descripcion: "Administrador global"}
	Rol_emergencias = util.Rol{Id: 6, Nombre: "emergencias", Descripcion: "Emergencias"}
}

//GET

func rolesListHandler(w http.ResponseWriter, req *http.Request) {
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

	var rolesJSON util.Roles_List_json
	//Certificado
	client := GetTLSClient()

	//Request al servidor para recibir lista de roles
	response, err := client.Get(SERVER_URL + "/rol/list")
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&rolesJSON)
		js, err := json.Marshal(rolesJSON)
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

func rolesListByUserHandler(w http.ResponseWriter, req *http.Request) {

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

	userIdparam, ok := req.URL.Query()["userId"]
	var userId = "-1"

	if ok {
		userId = userIdparam[0]
	}

	var rolesList []int
	//Certificado
	client := GetTLSClient()

	//Request al servidor para recibir lista de roles
	response, err := client.Get(SERVER_URL + "/rol/list/user?userId=" + userId)
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&rolesList)
		js, err := json.Marshal(rolesList)
		if rolesList == nil {
			jsondat := &util.JSON_Return{Result: "", Error: "Usuario sin roles"}
			js, err = json.Marshal(jsondat)
		}
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
