package routes

import (
	"encoding/json"
	"net/http"

	models "../models"
	util "../utils"
)

//POST
func MedicoSolicitarHistorialHandler(w http.ResponseWriter, req *http.Request) {
	var solicitarHistorial util.SolicitarHistorial_JSON
	json.NewDecoder(req.Body).Decode(&solicitarHistorial)
	jsonReturn := util.JSON_Return{}
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(solicitarHistorial.UserToken.UserId, solicitarHistorial.UserToken.Token, models.Rol_medico.Id)
	if authorized == true {
		checkDNI, _ := models.CheckUserDniHash(solicitarHistorial.UserToken.UserId, solicitarHistorial.UserDNI)
		if checkDNI != -1 {
			jsonReturn = util.JSON_Return{Result: "OK"}
		} else {
			jsonReturn = util.JSON_Return{Error: "El documento de identificación no existe"}
		}
	} else {
		jsonReturn = util.JSON_Return{Error: "No dispones de permisos para realizar esa acción"}
	}
	js, err := json.Marshal(jsonReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
