package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	util "../utils"
)

func solicitarPermisoTotal(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}

	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}

	//Recuperamos datos del form
	var historial util.Historial_JSON
	json.NewDecoder(req.Body).Decode(&historial)
	historial.UserToken = prepareUserToken(req)

	//Certificado
	client := GetTLSClient()

	locJson, err := json.Marshal(historial)
	response, err := client.Post(SERVER_URL+"/permisos/historial/total/solicitar", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		json.NewDecoder(response.Body).Decode(&result)
		js, err := json.Marshal(result)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func solicitarPermisoBasico(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}

	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}

	//Recuperamos datos del form
	var historial util.Historial_JSON
	json.NewDecoder(req.Body).Decode(&historial)
	historial.UserToken = prepareUserToken(req)

	//Certificado
	client := GetTLSClient()

	locJson, err := json.Marshal(historial)
	response, err := client.Post(SERVER_URL+"/permisos/historial/basico/solicitar", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		json.NewDecoder(response.Body).Decode(&result)
		js, err := json.Marshal(result)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func solicitarPermisoEntrada(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}

	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}

	//Recuperamos datos del form
	var entrada util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entrada)
	entrada.UserToken = prepareUserToken(req)

	//Certificado
	client := GetTLSClient()
	locJson, err := json.Marshal(entrada)
	response, err := client.Post(SERVER_URL+"/permisos/entrada/solicitar", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		json.NewDecoder(response.Body).Decode(&result)
		js, err := json.Marshal(result)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func solicitarPermisoAnalitica(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}

	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}

	//Recuperamos datos del form
	var analitica util.AnaliticaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&analitica)
	analitica.UserToken = prepareUserToken(req)

	//Certificado
	client := GetTLSClient()

	//Recuperamos clave publica del paciente
	locJson, err := json.Marshal(analitica)
	response, err := client.Post(SERVER_URL+"/permisos/analitica/solicitar", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		json.NewDecoder(response.Body).Decode(&result)
		js, err := json.Marshal(result)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func denegarPermiso(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}

	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}

	//Recuperamos datos del form
	var solicitud util.Solicitud_JSON
	json.NewDecoder(req.Body).Decode(&solicitud)
	solicitud.UserToken = prepareUserToken(req)
	//Certificado
	client := GetTLSClient()
	locJson, err := json.Marshal(solicitud)
	response, err := client.Post(SERVER_URL+"/permisos/solicitudes/eliminar", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		json.NewDecoder(response.Body).Decode(&result)
		js, err := json.Marshal(result)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func permitirPermiso(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
	} else {
		//Refrescar sesión
		session.Options.MaxAge = 60 * 30
		session.Save(req, w)
	}

	// Check user Token
	if !proveToken(req) {
		http.Redirect(w, req, "/forbidden", http.StatusSeeOther)
		return
	}

	//Recuperamos datos de la solicitud
	var solicitud util.Solicitud_JSON
	json.NewDecoder(req.Body).Decode(&solicitud)
	solicitud.UserToken = prepareUserToken(req)

	//Certificado
	client := GetTLSClient()

	//Recuperamos nuestra clave privada cifrada
	userId, _ := session.Values["userId"].(string)
	userPairkeys := getUserPairKeys(userId)
	userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)
	//Desciframos nuestra clave privada cifrada con AES
	userPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(userPairkeys.PrivateKey))
	userPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(userPrivateKeyString)))
	//Recuperamos la clave publica del empleado
	empleadoPairkeys := getUserPairKeys(strconv.Itoa(solicitud.EmpleadoId))
	empleadoPublicKey := *util.RSABytesToPublicKey(empleadoPairkeys.PublicKey)

	if solicitud.TipoHistorial != "" {
		//Dar permisos historial
		var historial util.Historial_JSON
		//Recuperamos historial
		locJson, _ := json.Marshal(prepareUserToken(req))
		response, _ := client.Post(SERVER_URL+"/user/patient/historial", "application/json", bytes.NewBuffer(locJson))
		if response != nil {
			json.NewDecoder(response.Body).Decode(&historial)
			//Desciframos la clave AES de los datos del usuario
			userDataKey, _ := session.Values["userDataKey"].(string)
			claveAESuserData := util.RSADecryptOAEP(userDataKey, *userPrivateKey)
			claveAESuserDataByte := util.Base64Decode([]byte(claveAESuserData))
			//Ciframos los datos del historial con AES y terminamos de rellenar el historial
			var historialCompartido util.Historial_JSON
			historialCompartido.Id = historial.Id
			historialCompartido.MedicoId = solicitud.EmpleadoId
			historialCompartido.PacienteId, _ = strconv.Atoi(prepareUserToken(req).UserId)
			historialCompartido.UserToken = prepareUserToken(req)
			//Ciframos la clave AES con la clave publica del Medico
			historialCompartido.Clave = util.RSAEncryptOAEP(string(util.Base64Encode(claveAESuserDataByte)), empleadoPublicKey)
			if solicitud.TipoHistorial == "TOTAL" {
				//PERMISO TOTAL
				//Permiso historial basico
				locJson, _ = json.Marshal(historialCompartido)
				response, _ = client.Post(SERVER_URL+"/permisos/historial/permitir", "application/json", bytes.NewBuffer(locJson))
				//Permiso entradas
				for _, entrada := range historial.Entradas {
					var entradaCompartir util.EntradaHistorial_JSON
					//Desciframos la clave AES de los datos cifrados
					claveAESentrada := util.RSADecryptOAEP(entrada.Clave, *userPrivateKey)
					claveAESentradaByte := util.Base64Decode([]byte(claveAESentrada))
					//Asignamos los datos para la request
					entradaCompartir.UserToken = prepareUserToken(req)
					entradaCompartir.EmpleadoId = solicitud.EmpleadoId
					entradaCompartir.Id = entrada.Id
					entradaCompartir.Clave = util.RSAEncryptOAEP(string(util.Base64Encode(claveAESentradaByte)), empleadoPublicKey)
					//Mandamos la petición
					locJson, _ = json.Marshal(entradaCompartir)
					response, _ = client.Post(SERVER_URL+"/permisos/entrada/permitir", "application/json", bytes.NewBuffer(locJson))
				}
				//Permiso analiticas
				for _, analitica := range historial.Analiticas {
					var analiticaCompartir util.AnaliticaHistorial_JSON
					//Desciframos la clave AES de los datos cifrados
					claveAESanalitica := util.RSADecryptOAEP(analitica.Clave, *userPrivateKey)
					claveAESanaliticaByte := util.Base64Decode([]byte(claveAESanalitica))
					//Asignamos los datos para la request
					analiticaCompartir.UserToken = prepareUserToken(req)
					analiticaCompartir.EmpleadoId = solicitud.EmpleadoId
					analiticaCompartir.Id = analitica.Id
					analiticaCompartir.Clave = util.RSAEncryptOAEP(string(util.Base64Encode(claveAESanaliticaByte)), empleadoPublicKey)
					//Mandamos la petición
					locJson, _ = json.Marshal(analiticaCompartir)
					response, _ = client.Post(SERVER_URL+"/permisos/analitica/permitir", "application/json", bytes.NewBuffer(locJson))
				}
			} else {
				//PERMISO BASICO
				locJson, _ = json.Marshal(historialCompartido)
				response, _ = client.Post(SERVER_URL+"/permisos/historial/permitir", "application/json", bytes.NewBuffer(locJson))
			}
		}
	} else {
		//Dar permisos entradas/analiticas
		if solicitud.EntradaId != 0 {
			//Obtenemos la entrada
			entrada := util.EntradaHistorial_JSON{Id: solicitud.EntradaId, UserToken: prepareUserToken(req)}
			locJson, _ := json.Marshal(entrada)
			//Request para obtener historial si existe
			response, _ := client.Post(SERVER_URL+"/user/patient/historial/entrada", "application/json", bytes.NewBuffer(locJson))
			if response != nil {
				err := json.NewDecoder(response.Body).Decode(&entrada)
				if err != nil {
					util.PrintErrorLog(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				var entradaCompartir util.EntradaHistorial_JSON
				//Desciframos la clave AES de los datos cifrados
				claveAESentrada := util.RSADecryptOAEP(entrada.Clave, *userPrivateKey)
				claveAESentradaByte := util.Base64Decode([]byte(claveAESentrada))
				//Asignamos los datos para la request
				entradaCompartir.UserToken = prepareUserToken(req)
				entradaCompartir.EmpleadoId = solicitud.EmpleadoId
				entradaCompartir.Id = entrada.Id
				entradaCompartir.Clave = util.RSAEncryptOAEP(string(util.Base64Encode(claveAESentradaByte)), empleadoPublicKey)
				//Mandamos la petición
				locJson, _ = json.Marshal(entradaCompartir)
				response, _ = client.Post(SERVER_URL+"/permisos/entrada/permitir", "application/json", bytes.NewBuffer(locJson))
			}
		} else {
			//Obtenemos la analítica
			analitica := util.AnaliticaHistorial_JSON{Id: solicitud.AnaliticaId, UserToken: prepareUserToken(req)}
			locJson, _ := json.Marshal(analitica)
			//Request para obtener analitica si existe
			response, _ := client.Post(SERVER_URL+"/user/patient/historial/analitica", "application/json", bytes.NewBuffer(locJson))
			if response != nil {
				err := json.NewDecoder(response.Body).Decode(&analitica)
				if err != nil {
					util.PrintErrorLog(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				var analiticaCompartir util.AnaliticaHistorial_JSON
				//Desciframos la clave AES de los datos cifrados
				claveAESanalitica := util.RSADecryptOAEP(analitica.Clave, *userPrivateKey)
				claveAESanaliticaByte := util.Base64Decode([]byte(claveAESanalitica))
				//Asignamos los datos para la request
				analiticaCompartir.UserToken = prepareUserToken(req)
				analiticaCompartir.EmpleadoId = solicitud.EmpleadoId
				analiticaCompartir.Id = analitica.Id
				analiticaCompartir.Clave = util.RSAEncryptOAEP(string(util.Base64Encode(claveAESanaliticaByte)), empleadoPublicKey)
				//Mandamos la petición
				locJson, _ = json.Marshal(analiticaCompartir)
				response, _ = client.Post(SERVER_URL+"/permisos/analitica/permitir", "application/json", bytes.NewBuffer(locJson))
			}
		}
	}

	//Actualizamos el estado de la solicitud
	locJson, err := json.Marshal(solicitud)
	response, err := client.Post(SERVER_URL+"/permisos/solicitudes/eliminar", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		json.NewDecoder(response.Body).Decode(&result)
		js, err := json.Marshal(result)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
