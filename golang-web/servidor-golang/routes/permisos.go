package routes

import (
	"encoding/json"
	"net/http"

	models "../models"
	util "../utils"
)

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
	var entrada util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entrada)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedEnfermero, _ := models.GetAuthorizationbyUserId(entrada.UserToken.UserId, entrada.UserToken.Token, models.Rol_enfermero.Id)
	authorizedMedico, _ := models.GetAuthorizationbyUserId(entrada.UserToken.UserId, entrada.UserToken.Token, models.Rol_medico.Id)
	authorized := authorizedEnfermero || authorizedMedico
	if authorized == true {
		//Agregamos petición
		result, _ := models.SolicitarPermisoEntrada(entrada)
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
	var analitica util.AnaliticaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&analitica)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedEnfermero, _ := models.GetAuthorizationbyUserId(analitica.UserToken.UserId, analitica.UserToken.Token, models.Rol_enfermero.Id)
	authorizedMedico, _ := models.GetAuthorizationbyUserId(analitica.UserToken.UserId, analitica.UserToken.Token, models.Rol_medico.Id)
	authorized := authorizedEnfermero || authorizedMedico
	if authorized == true {
		//Agregamos petición
		result, _ := models.SolicitarPermisoAnalitica(analitica)
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

func eliminarSolicitudPermiso(w http.ResponseWriter, req *http.Request) {
	var solicitud util.Solicitud_JSON
	json.NewDecoder(req.Body).Decode(&solicitud)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(solicitud.UserToken.UserId, solicitud.UserToken.Token, models.Rol_paciente.Id)
	if authorizedPaciente == true {
		//Procesamos petición
		result := false
		if solicitud.TipoHistorial != "" {
			//Denegar permisos historial
			if solicitud.TipoHistorial == "TOTAL" {
				result, _ = models.BorrarSolicitudHistorialTotal(solicitud.UserToken.UserId, solicitud.EmpleadoId)
			} else {
				result, _ = models.BorrarSolicitudHistorialBasico(solicitud.UserToken.UserId, solicitud.EmpleadoId)
			}
		} else {
			//Denegar permisos entradas/analiticas
			if solicitud.EntradaId != 0 {
				result, _ = models.BorrarSolicitudEntrada(solicitud.UserToken.UserId, solicitud.EmpleadoId, solicitud.EntradaId)
			} else {
				result, _ = models.BorrarSolicitudAnalitica(solicitud.UserToken.UserId, solicitud.EmpleadoId, solicitud.AnaliticaId)
			}
		}
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

func permitirHistorial(w http.ResponseWriter, req *http.Request) {
	var historial util.Historial_JSON
	json.NewDecoder(req.Body).Decode(&historial)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(historial.UserToken.UserId, historial.UserToken.Token, models.Rol_paciente.Id)
	if authorizedPaciente == true {
		//Agregamos petición
		result := false
		result, _ = models.InsertShareHistorial(historial)
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

func permitirEntrada(w http.ResponseWriter, req *http.Request) {
	var entrada util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entrada)
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(entrada.UserToken.UserId, entrada.UserToken.Token, models.Rol_paciente.Id)
	if authorizedPaciente == true {
		//Agregamos petición
		result := false
		result, _ = models.InsertEntradaCompartidaHistorialPermisos(entrada)
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

func permitirAnalitica(w http.ResponseWriter, req *http.Request) {
	var analitica util.AnaliticaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&analitica)

	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorizedPaciente, _ := models.GetAuthorizationbyUserId(analitica.UserToken.UserId, analitica.UserToken.Token, models.Rol_paciente.Id)
	if authorizedPaciente == true {
		//Agregamos petición
		result := false
		result, _ = models.InsertAnaliticaCompartidaHistorial(analitica)
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
