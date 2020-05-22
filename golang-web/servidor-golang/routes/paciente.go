package routes

import (
	"encoding/json"
	"net/http"

	models "../models"
	util "../utils"
)

//POST
func PacienteInsertCita(w http.ResponseWriter, req *http.Request) {
	var cita util.CitaJSON
	json.NewDecoder(req.Body).Decode(&cita)
	jsonReturn := util.JSON_Return{}

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(cita.UserToken.UserId, cita.UserToken.Token, models.Rol_paciente.Id)
	if authorized == true {
		//Insertamos la cita
		result, _ := models.InsertCita(cita)
		if result == true {
			jsonReturn = util.JSON_Return{Result: "Cita reservada correctamente"}
		} else {
			jsonReturn = util.JSON_Return{Error: "Error reservando cita"}
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

func PacienteGetCitasFuturasList(w http.ResponseWriter, req *http.Request) {
	var userToken util.UserToken_JSON
	json.NewDecoder(req.Body).Decode(&userToken)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_paciente.Id)
	if authorized == true {
		jsonReturn, _ := models.GetCitasFuturasPaciente(userToken.UserId)
		js, err := json.Marshal(jsonReturn)
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

func PacienteGetEntradaHandler(w http.ResponseWriter, req *http.Request) {
	var entrada util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entrada)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(entrada.UserToken.UserId, entrada.UserToken.Token, models.Rol_paciente.Id)
	if authorized == true {
		entradaJSON, _ := models.GetEntradaById(entrada.Id)
		js, err := json.Marshal(entradaJSON)
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

func PacienteGetAnaliticaHandler(w http.ResponseWriter, req *http.Request) {
	var analitica util.AnaliticaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&analitica)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(analitica.UserToken.UserId, analitica.UserToken.Token, models.Rol_paciente.Id)
	if authorized == true {
		analiticaJSON, _ := models.GetAnaliticaById(analitica.Id)
		js, err := json.Marshal(analiticaJSON)
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
