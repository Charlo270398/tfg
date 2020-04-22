package routes

import (
	"encoding/json"
	"net/http"

	models "../models"
	util "../utils"
)

//POST
func getAdminMenuDataHandler(w http.ResponseWriter, req *http.Request) {
	var userToken util.JSON_Admin_Menu
	json.NewDecoder(req.Body).Decode(&userToken)
	jsonReturn := util.JSON_Admin_Menu{}
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(userToken.UserToken.UserId, userToken.UserToken.Token, models.Rol_administradorC.Id)
	if authorized == true {
		clinica, err := models.GetClinicaByAdmin(userToken.UserToken.UserId)
		if err != nil {
			jsonReturn = util.JSON_Admin_Menu{Error: "Error cargando los datos de la clínica"}
		} else {
			jsonReturn = util.JSON_Admin_Menu{Clinica: clinica, Error: ""}
		}

	} else {
		jsonReturn = util.JSON_Admin_Menu{Error: "No dispones de permisos para realizar esa acción"}
	}
	js, err := json.Marshal(jsonReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
