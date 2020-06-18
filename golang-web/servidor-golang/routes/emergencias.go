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

func GetEntradaEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
	var entrada util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entrada)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(entrada.UserToken.UserId, entrada.UserToken.Token, models.Rol_emergencias.Id)
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

func GetAnaliticaEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
	var analitica util.AnaliticaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&analitica)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(analitica.UserToken.UserId, analitica.UserToken.Token, models.Rol_emergencias.Id)
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

func AddEntradaEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
	var entradaHistorial util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entradaHistorial)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(entradaHistorial.UserToken.UserId, entradaHistorial.UserToken.Token, models.Rol_emergencias.Id)
	if authorized == true {
		var returnJSON util.JSON_Return
		//Insertamos la entrada
		result, err := models.InsertEntradaHistorial(entradaHistorial)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if result != -1 {
			returnJSON.Result = strconv.Itoa(result)
		} else {
			returnJSON.Error = "Error insertando la entrada"
		}

		js, err := json.Marshal(returnJSON)
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

func AddAnaliticaEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
	var analiticaHistorial util.AnaliticaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&analiticaHistorial)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(analiticaHistorial.UserToken.UserId, analiticaHistorial.UserToken.Token, models.Rol_emergencias.Id)
	if authorized == true {
		var returnJSON util.JSON_Return
		//Insertamos la analítica
		result, err := models.InsertAnaliticaHistorial(analiticaHistorial)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if result != -1 {
			returnJSON.Result = strconv.Itoa(result)
		} else {
			returnJSON.Error = "Error insertando la analítica"
		}

		js, err := json.Marshal(returnJSON)
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

func AddEstadisticaAnaliticaEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
	var analiticaHistorial util.AnaliticaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&analiticaHistorial)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	//MEDICO,ENFERMERO Y EMERGENCIAS
	authorized, _ := models.GetAuthorizationbyUserId(analiticaHistorial.UserToken.UserId, analiticaHistorial.UserToken.Token, models.Rol_emergencias.Id)
	if authorized == true {
		var returnJSON util.JSON_Return
		//Insertamos la analítica
		result, err := models.InsertEstadisticaAnaliticaHistorial(analiticaHistorial)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if result == true {

		} else {
			returnJSON.Error = "Error insertando la analítica"
		}

		js, err := json.Marshal(returnJSON)
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
