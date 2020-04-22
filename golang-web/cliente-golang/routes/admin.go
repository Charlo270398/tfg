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
func menuAdminHandler(w http.ResponseWriter, req *http.Request) {
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
