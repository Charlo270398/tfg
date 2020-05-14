package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "../models"
	util "../utils"
)

//GET

func getUsersPaginationHandler(w http.ResponseWriter, req *http.Request) {
	page, ok := req.URL.Query()["page"]
	var usersListReturn util.UserList_JSON_Pagination
	var usersList []util.User_JSON
	if !ok || len(page[0]) < 1 {
		usersList = models.GetUsersPagination(0) //Devolvemos primera pagina
	} else {
		pageNumber, err := strconv.Atoi(page[0])
		usersListReturn.Page = pageNumber
		usersListReturn.BeforePage = pageNumber - 1
		usersListReturn.NextPage = pageNumber + 1
		if err != nil {
			usersList = models.GetUsersPagination(0) //Devolvemos primera pagina
		} else {
			usersList = models.GetUsersPagination(pageNumber)
		}
	}
	usersListReturn.UserList = usersList
	js, err := json.Marshal(usersListReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func getUserHandler(w http.ResponseWriter, req *http.Request) {
	userIdURL, ok := req.URL.Query()["userId"]
	var usersReturn util.User_JSON
	var user util.User_JSON
	if !ok || len(userIdURL[0]) < 1 {
		http.Error(w, "No hay parámetros", http.StatusInternalServerError)
		return
	} else {
		userId, err := strconv.Atoi(userIdURL[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			user, err = models.GetUserById(userId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}
	usersReturn.Nombre = user.Nombre
	usersReturn.Apellidos = user.Apellidos
	usersReturn.Identificacion = user.Identificacion
	usersReturn.Email = user.Email
	js, err := json.Marshal(usersReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getUserRolesHandler(w http.ResponseWriter, req *http.Request) {
	var userRolesListJSON util.Roles_List_json
	roles, err := models.GetRolesList() //Aqui usar función para un ID de usuario
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

//PAIRKEYS

func getUserPairKeysHandler(w http.ResponseWriter, req *http.Request) {
	userIdURL, ok := req.URL.Query()["userId"]
	var userReturn util.User_JSON
	if !ok || len(userIdURL[0]) < 1 {
		http.Error(w, "No hay parámetros", http.StatusInternalServerError)
		return
	} else {
		pairKeys, err := models.GetUserPairKeys(userIdURL[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userId, err := strconv.Atoi(userIdURL[0])
		userReturn.Id = userId
		userReturn.PairKeys = pairKeys
	}
	js, err := json.Marshal(userReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getUserMasterPairKeysHandler(w http.ResponseWriter, req *http.Request) {
	userIdURL, ok := req.URL.Query()["userId"]
	var userReturn util.User_JSON
	if !ok || len(userIdURL[0]) < 1 {
		http.Error(w, "No hay parámetros", http.StatusInternalServerError)
		return
	} else {
		masterPairKeys, err := models.GetUserMasterPairKeys(userIdURL[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userId, err := strconv.Atoi(userIdURL[0])
		userReturn.Id = userId
		userReturn.MasterPairKeys = masterPairKeys
	}
	js, err := json.Marshal(userReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getPublicMasterKeyHandler(w http.ResponseWriter, req *http.Request) {
	masterPairKeys, err := models.GetPublicMasterKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var user util.User_JSON
	user.MasterPairKeys = masterPairKeys
	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getUserByHistorialIdHandler(w http.ResponseWriter, req *http.Request) {
	masterPairKeys, err := models.GetPublicMasterKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var user util.User_JSON
	user.MasterPairKeys = masterPairKeys
	js, err := json.Marshal(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getUserPairKeysByHistorialIdHandler(w http.ResponseWriter, req *http.Request) {
	historialIdURL, ok := req.URL.Query()["historialId"]
	var userReturn util.User_JSON
	if !ok || len(historialIdURL[0]) < 1 {
		http.Error(w, "No hay parámetros", http.StatusInternalServerError)
		return
	} else {
		pairKeys, err := models.GetUserPairKeysByHistorialId(historialIdURL[0])
		pairKeys.PrivateKey = nil
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		userReturn.PairKeys = pairKeys
	}
	js, err := json.Marshal(userReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
