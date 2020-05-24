package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "../models"
	util "../utils"
)

//POST
func MedicoSolicitarHistorialHandler(w http.ResponseWriter, req *http.Request) {
	var solicitarHistorial util.SolicitarHistorial_JSON
	json.NewDecoder(req.Body).Decode(&solicitarHistorial)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(solicitarHistorial.UserToken.UserId, solicitarHistorial.UserToken.Token, models.Rol_medico.Id)
	if authorized == true {
		userId, err := models.CheckUserDniHash(solicitarHistorial.UserDNI)
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
		solicitarHistorial.Historial = historialJSON
		//Historial permitido
		historialPermitidoJSON, _ := models.GetHistorialCompartidoByMedicoIdPacienteId(solicitarHistorial.UserToken.UserId, userIdString)
		solicitarHistorial.HistorialPermitido = historialPermitidoJSON
		js, err := json.Marshal(solicitarHistorial)
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

func MedicoDiasDisponiblesHandler(w http.ResponseWriter, req *http.Request) {
	urlParam, ok := req.URL.Query()["doctorId"]
	if !ok || len(urlParam[0]) < 1 {
		http.Error(w, "¡No hay parametro doctorId!", http.StatusInternalServerError)
		return
	}
	doctorId := urlParam[0]

	diasDisponibles, err := models.GetDiasDisponiblesMedico(doctorId)
	js, err := json.Marshal(diasDisponibles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func MedicoHorasDiaDisponiblesHandler(w http.ResponseWriter, req *http.Request) {
	urlParam, ok := req.URL.Query()["doctorId"]
	if !ok || len(urlParam[0]) < 1 {
		http.Error(w, "¡No hay parametro doctorId!", http.StatusInternalServerError)
		return
	}
	doctorId := urlParam[0]

	urlParam, ok = req.URL.Query()["dia"]
	if !ok || len(urlParam[0]) < 1 {
		http.Error(w, "¡No hay parametro doctorId!", http.StatusInternalServerError)
		return
	}
	dia := urlParam[0]
	horasDisponibles, err := models.GetHorasDiaDisponiblesMedico(doctorId, dia)
	js, err := json.Marshal(horasDisponibles)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func MedicoGetCitasFuturasList(w http.ResponseWriter, req *http.Request) {
	var userToken util.UserToken_JSON
	json.NewDecoder(req.Body).Decode(&userToken)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_medico.Id)
	if authorized == true {
		jsonReturn, _ := models.GetCitasFuturasMedico(userToken.UserId)
		for index, cita := range jsonReturn {
			jsonReturn[index].Historial, _ = models.GetHistorialCompartidoByMedicoIdPacienteId(userToken.UserId, cita.PacienteId)
			jsonReturn[index].Historial.Sexo = ""
			jsonReturn[index].Historial.Alergias = ""
		}
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

func MedicoGetCitaActual(w http.ResponseWriter, req *http.Request) {
	var userToken util.UserToken_JSON
	json.NewDecoder(req.Body).Decode(&userToken)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_medico.Id)
	if authorized == true {
		citaId, _ := models.GetCitaActualMedico(userToken.UserId)
		var citaJson util.CitaJSON
		citaJson.Id = citaId
		js, err := json.Marshal(citaJson)
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

func MedicoGetCita(w http.ResponseWriter, req *http.Request) {
	var cita util.CitaJSON
	json.NewDecoder(req.Body).Decode(&cita)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(cita.UserToken.UserId, cita.UserToken.Token, models.Rol_medico.Id)
	if authorized == true {
		cita, _ := models.GetCitaById(cita.Id)
		historialPaciente, _ := models.GetHistorialCompartidoByMedicoIdPacienteId(cita.MedicoId, cita.PacienteId)
		cita.Historial = historialPaciente
		js, err := json.Marshal(cita)
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

func MedicoAddEntradaHistorialConsulta(w http.ResponseWriter, req *http.Request) {
	var entradaHistorial util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entradaHistorial)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(entradaHistorial.UserToken.UserId, entradaHistorial.UserToken.Token, models.Rol_medico.Id)
	if authorized == true {
		var returnJSON util.JSON_Return
		//Insertamos la entrada
		result, err := models.InsertEntradaHistorialPacienteId(entradaHistorial)
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

func MedicoAddEntradaHistorialCompartidaConsulta(w http.ResponseWriter, req *http.Request) {
	var entradaHistorial util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entradaHistorial)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(entradaHistorial.UserToken.UserId, entradaHistorial.UserToken.Token, models.Rol_medico.Id)
	if authorized == true {
		var returnJSON util.JSON_Return
		//Insertamos la entrada
		result, err := models.InsertEntradaCompartidaHistorial(entradaHistorial)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if result == true {
			returnJSON.Result = "OK"
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

func MedicoGetHistorialesCompartidos(w http.ResponseWriter, req *http.Request) {
	var userToken util.UserToken_JSON
	json.NewDecoder(req.Body).Decode(&userToken)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_medico.Id)
	if authorized == true {
		historiales, _ := models.GetHistorialesCompartidosByMedicoId(userToken.UserId)
		js, err := json.Marshal(historiales)
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

func GetHistorialCompartido(w http.ResponseWriter, req *http.Request) {
	var historial util.Historial_JSON
	json.NewDecoder(req.Body).Decode(&historial)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_medico.Id)
	if authorized == true {
		medicoIdString := strconv.Itoa(historial.MedicoId)
		pacienteIdString := strconv.Itoa(historial.PacienteId)
		historialPaciente, _ := models.GetHistorialCompartidoByMedicoIdPacienteId(medicoIdString, pacienteIdString)
		js, err := json.Marshal(historialPaciente)
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
