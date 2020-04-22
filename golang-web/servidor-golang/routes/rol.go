package routes

import (
	"encoding/json"
	"net/http"

	models "../models"
	util "../utils"
)

//GET

func getRolesListHandler(w http.ResponseWriter, req *http.Request) {
	var userRolesListJSON util.Roles_List_json
	roles, err := models.GetRolesList()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userRolesListJSON.Roles = roles
	js, err := json.Marshal(userRolesListJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getRolesByUserHandler(w http.ResponseWriter, req *http.Request) {
	userId, ok := req.URL.Query()["userId"]

	if !ok || len(userId[0]) < 1 {
		//Si no hay userId
		http.Error(w, "No hay userId", http.StatusInternalServerError)
		return
	} else {
		//Si hay userId
		roles, err := models.GetRolesbyUserId(userId[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		js, err := json.Marshal(roles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}
