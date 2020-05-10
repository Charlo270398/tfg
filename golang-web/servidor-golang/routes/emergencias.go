package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "../models"
	util "../utils"
)

func GetHistorialEmergencias(w http.ResponseWriter, req *http.Request) {
	var user util.User_JSON
	json.NewDecoder(req.Body).Decode(&user)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(user.UserToken.UserId, user.UserToken.Token, models.Rol_emergencias.Id)
	if authorized == true {
		userId, err := models.CheckUserDniHash(user.Identificacion)
		if userId == -1 || err != nil {
			http.Error(w, "Error buscando el usuario", http.StatusInternalServerError)
			return
		}
		userIdString := strconv.Itoa(userId)
		userData, _ := models.GetUserById(userId)
		historialJSON, _ := models.GetHistorialByUserId(userIdString)
		historialJSON.NombrePaciente = userData.Nombre
		historialJSON.ApellidosPaciente = userData.Apellidos
		historialJSON.Entradas, _ = models.GetEntradasHistorialByHistorialId(historialJSON.Id)
		js, err := json.Marshal(historialJSON)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	}
	http.Error(w, "No estas autorizado", http.StatusInternalServerError)
	return
}
