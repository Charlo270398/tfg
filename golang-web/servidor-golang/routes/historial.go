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

func solicitarPermisoTotal(w http.ResponseWriter, req *http.Request) {
	var historial util.Historial_JSON
	json.NewDecoder(req.Body).Decode(&historial)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_paciente.Id)
	authorizedEnfermero, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_enfermero.Id)
	authorizedMedico, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_medico.Id)
	authorizedEmergencias, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_emergencias.Id)
	authorized := authorizedPaciente || authorizedEnfermero || authorizedMedico || authorizedEmergencias
	if authorized == true {
		//Agregamos petición
		result, _ := models.SolicitarPermisoTotalHistorial(historial)
		if result == true {
			js, err := json.Marshal(util.JSON_Return{Result: "OK"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else {
			js, err := json.Marshal(util.JSON_Return{Error: "Error insertando solicitud"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}
	http.Error(w, "No estas autorizado", http.StatusInternalServerError)
	return
}

func solicitarPermisoBasico(w http.ResponseWriter, req *http.Request) {
	var historial util.Historial_JSON
	json.NewDecoder(req.Body).Decode(&historial)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_paciente.Id)
	authorizedEnfermero, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_enfermero.Id)
	authorizedMedico, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_medico.Id)
	authorizedEmergencias, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_emergencias.Id)
	authorized := authorizedPaciente || authorizedEnfermero || authorizedMedico || authorizedEmergencias
	if authorized == true {
		//Agregamos petición
		result, _ := models.SolicitarPermisoBasicoHistorial(historial)
		if result == true {
			js, err := json.Marshal(util.JSON_Return{Result: "OK"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else {
			js, err := json.Marshal(util.JSON_Return{Error: "Error insertando solicitud"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}
	http.Error(w, "No estas autorizado", http.StatusInternalServerError)
	return
}

func solicitarPermisoEntrada(w http.ResponseWriter, req *http.Request) {
	var historial util.Historial_JSON
	json.NewDecoder(req.Body).Decode(&historial)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_paciente.Id)
	authorizedEnfermero, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_enfermero.Id)
	authorizedMedico, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_medico.Id)
	authorizedEmergencias, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_emergencias.Id)
	authorized := authorizedPaciente || authorizedEnfermero || authorizedMedico || authorizedEmergencias
	if authorized == true {
		//Agregamos petición
		result, _ := models.SolicitarPermisoTotalHistorial(historial)
		if result == true {
			js, err := json.Marshal(util.JSON_Return{Result: "OK"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else {
			js, err := json.Marshal(util.JSON_Return{Error: "Error insertando solicitud"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}
	http.Error(w, "No estas autorizado", http.StatusInternalServerError)
	return
}

func solicitarPermisoAnalitica(w http.ResponseWriter, req *http.Request) {
	var historial util.Historial_JSON
	json.NewDecoder(req.Body).Decode(&historial)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_paciente.Id)
	authorizedEnfermero, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_enfermero.Id)
	authorizedMedico, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_medico.Id)
	authorizedEmergencias, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_emergencias.Id)
	authorized := authorizedPaciente || authorizedEnfermero || authorizedMedico || authorizedEmergencias
	if authorized == true {
		//Agregamos petición
		result, _ := models.SolicitarPermisoTotalHistorial(historial)
		if result == true {
			js, err := json.Marshal(util.JSON_Return{Result: "OK"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else {
			js, err := json.Marshal(util.JSON_Return{Error: "Error insertando solicitud"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}
	http.Error(w, "No estas autorizado", http.StatusInternalServerError)
	return
}

func listarSolicitudesPermiso(w http.ResponseWriter, req *http.Request) {
	var userToken util.UserToken_JSON
	json.NewDecoder(req.Body).Decode(&userToken)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_paciente.Id)
	if authorizedPaciente == true {
		result, _ := models.ListarSolicitudesPermiso(userToken.UserId)
		js, err := json.Marshal(result)
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

func comprobarSolicitudesPermiso(w http.ResponseWriter, req *http.Request) {
	var userToken util.UserToken_JSON
	json.NewDecoder(req.Body).Decode(&userToken)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(userToken.UserId, userToken.Token, models.Rol_paciente.Id)
	if authorizedPaciente == true {
		//Agregamos petición
		result, _ := models.ComprobarSolicitudesPermiso(userToken.UserId)
		if result == true {
			js, err := json.Marshal(util.JSON_Return{Result: "OK"})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		} else {
			js, err := json.Marshal(util.JSON_Return{Result: ""})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}
	http.Error(w, "No estas autorizado", http.StatusInternalServerError)
	return
}
