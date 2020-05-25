package routes

import (
	"encoding/json"
	"net/http"

	util "../utils"
)

//GET

func tagsListHandler(w http.ResponseWriter, req *http.Request) {
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

	var tagsJSON []util.Tag_JSON
	//Certificado
	client := GetTLSClient()

	//Request al servidor para recibir lista de roles
	response, err := client.Get(SERVER_URL + "/tag/list")
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&tagsJSON)
		js, err := json.Marshal(tagsJSON)
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
