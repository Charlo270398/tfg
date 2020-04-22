package routes

import (
	"encoding/json"
	"net/http"

	util "../utils"
)

//GET

func rolesListHandler(w http.ResponseWriter, req *http.Request) {
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
	}
	// Check user Token
	if !proveToken(req) {
		http.Error(w, "Forbidden", http.StatusForbidden)
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
