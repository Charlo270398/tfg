package routes

import (
	"encoding/json"
	"net/http"

	models "../models"
	util "../utils"
)

func GetHistorialPaciente(w http.ResponseWriter, req *http.Request) {
	var userToken util.UserToken_JSON
	json.NewDecoder(req.Body).Decode(&userToken)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_paciente.Id)
	authorizedEnfermero, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_enfermero.Id)
	authorizedMedico, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_medico.Id)
	authorizedEmergencias, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_emergencias.Id)
	authorized := authorizedPaciente || authorizedEnfermero || authorizedMedico || authorizedEmergencias
	if authorized == true {
		historialJSON, _ := models.GetHistorialByUserId(userToken.UserId)
		historialJSON.Entradas, _ = models.GetEntradasHistorialByHistorialId(historialJSON.Id)
		historialJSON.Analiticas, _ = models.GetAnaliticasHistorialByHistorialId(historialJSON.Id)
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

func ShareHistorialPaciente(w http.ResponseWriter, req *http.Request) {
	var historialCompartido util.Historial_JSON
	json.NewDecoder(req.Body).Decode(&historialCompartido)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(historialCompartido.UserToken.UserId, historialCompartido.UserToken.Token, models.Rol_paciente.Id)
	authorizedEnfermero, _ := models.GetAuthorizationbyUserId(historialCompartido.UserToken.UserId, historialCompartido.UserToken.Token, models.Rol_enfermero.Id)
	authorizedMedico, _ := models.GetAuthorizationbyUserId(historialCompartido.UserToken.UserId, historialCompartido.UserToken.Token, models.Rol_medico.Id)
	authorizedEmergencias, _ := models.GetAuthorizationbyUserId(historialCompartido.UserToken.UserId, historialCompartido.UserToken.Token, models.Rol_emergencias.Id)
	authorized := authorizedPaciente || authorizedEnfermero || authorizedMedico || authorizedEmergencias
	if authorized == true {
		result, err := models.InsertShareHistorial(historialCompartido)
		var returnJSON util.JSON_Return
		if err != nil {
			returnJSON = util.JSON_Return{Error: err.Error()}
		} else {
			if result == true {
				returnJSON = util.JSON_Return{Result: "OK"}
			} else {
				returnJSON = util.JSON_Return{Error: "Error insertando el historial compartido"}
			}
		}
		js, err := json.Marshal(returnJSON)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return
	} else {

	}
	http.Error(w, "No estas autorizado", http.StatusInternalServerError)
	return
}

func GetHistorialCompartidoPaciente(w http.ResponseWriter, req *http.Request) {
	var userToken util.UserToken_JSON
	json.NewDecoder(req.Body).Decode(&userToken)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_paciente.Id)
	authorizedEnfermero, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_enfermero.Id)
	authorizedMedico, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_medico.Id)
	authorizedEmergencias, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_emergencias.Id)
	authorized := authorizedPaciente || authorizedEnfermero || authorizedMedico || authorizedEmergencias
	if authorized == true {
		historialJSON, _ := models.GetHistorialByUserId(userToken.UserId)
		historialJSON.Entradas, _ = models.GetEntradasHistorialByHistorialId(historialJSON.Id)
		historialJSON.Analiticas, _ = models.GetAnaliticasHistorialByHistorialId(historialJSON.Id)
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

func GetEstadisticasAnaliticas(w http.ResponseWriter, req *http.Request) {
	var userToken util.UserToken_JSON
	json.NewDecoder(req.Body).Decode(&userToken)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedMedico, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_medico.Id)
	authorized := authorizedMedico
	if authorized == true {
		analiticas, _ := models.GetEstadisticasAnaliticas()
		js, err := json.Marshal(analiticas)
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
