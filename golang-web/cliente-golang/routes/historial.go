package routes

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"

	util "../utils"
)

//GET
func addEntradaHistorialFormMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/historial/addEntrada.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", Page{Title: "Añadir entrada historial", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func addAnaliticaHistorialFormMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/historial/addAnalitica.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", Page{Title: "Añadir analitica historial", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//POST
func addEntradaHistorialConsultaMedicoHandler(w http.ResponseWriter, req *http.Request) {
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
	var entradaHistorial util.EntradaHistorial_JSON
	var entradaHistorialCompartida util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entradaHistorial)

	//Certificado
	client := GetTLSClient()

	//Recuperamos clave publica del paciente
	var user util.User_JSON
	response, err := client.Get(SERVER_URL + "/user/pairkeys?userId=" + strconv.Itoa(entradaHistorial.PacienteId))
	if response != nil {
		json.NewDecoder(response.Body).Decode(&user)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Recuperamos CLAVE PUBLICA MAESTRA
	masterPairKeys := getPublicMasterKey()

	//Generamos una clave AES aleatoria de 256 bits para cifrar los datos sensibles
	AESkeyDatos := util.AEScreateKey()

	//Ciframos los datos sensibles con la clave
	entradaHistorial.JuicioDiagnostico, _ = util.AESencrypt(AESkeyDatos, entradaHistorial.JuicioDiagnostico)
	entradaHistorial.MotivoConsulta, _ = util.AESencrypt(AESkeyDatos, entradaHistorial.MotivoConsulta)
	entradaHistorial.Tipo, _ = util.AESencrypt(AESkeyDatos, "Consulta")

	//Pasamos la clave a base 64
	AESkeyBase64String := string(util.Base64Encode(AESkeyDatos))

	//CIFRADO PARA PACIENTE
	//Recuperamos nuestra clave publica del paciente
	pacienteIdString := strconv.Itoa(entradaHistorial.PacienteId)
	pacientePairkeys := getUserPairKeys(pacienteIdString)

	//Ciframos la clave AES usada con nuestra clave pública
	claveAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(pacientePairkeys.PublicKey))
	claveMaestraAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(masterPairKeys.PublicKey))

	//Preparamos los datos para enviar
	entradaHistorial.UserToken = prepareUserToken(req)
	entradaHistorial.Clave = claveAEScifrada
	entradaHistorial.ClaveMaestra = claveMaestraAEScifrada
	locJson, err := json.Marshal(entradaHistorial)

	//Request al servidor para añadir entrada paciente
	response, err = client.Post(SERVER_URL+"/user/doctor/citas/addEntrada", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			if result.Error == "" {
				//CIFRADO PARA MEDICO
				//Recuperamos la clave publica del medico
				userId, _ := session.Values["userId"].(string)
				medicoPairkeys := getUserPairKeys(userId)
				//Ciframos la clave AES usada con nuestra clave pública
				claveAEScifrada = util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(medicoPairkeys.PublicKey))
				//Preparamos los datos para enviar
				entradaHistorialCompartida = entradaHistorial
				entradaHistorialCompartida.Id, _ = strconv.Atoi(result.Result)
				entradaHistorialCompartida.Clave = claveAEScifrada
				locJson, err = json.Marshal(entradaHistorialCompartida)
				//Request al servidor para añadir entrada compartida
				response, err = client.Post(SERVER_URL+"/user/doctor/citas/addEntradaCompartida", "application/json", bytes.NewBuffer(locJson))
				if response != nil {
					err := json.NewDecoder(response.Body).Decode(&result)
					if err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				} else {
					util.PrintErrorLog(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		js, err := json.Marshal(result)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addEntradaHistorialMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	//Preparamos datos request
	var entrada util.EntradaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&entrada)
	//Generamos una clave AES aleatoria de 256 bits para cifrar los datos sensibles
	AESkeyDatos := util.AEScreateKey()

	//Ciframos los datos sensibles con la clave
	entrada.JuicioDiagnostico, _ = util.AESencrypt(AESkeyDatos, entrada.JuicioDiagnostico)
	entrada.MotivoConsulta, _ = util.AESencrypt(AESkeyDatos, entrada.MotivoConsulta)
	entrada.Tipo, _ = util.AESencrypt(AESkeyDatos, entrada.Tipo)

	//Pasamos la clave a base 64
	AESkeyBase64String := string(util.Base64Encode(AESkeyDatos))

	//CIFRAMOS LA CLAVE AES CON LA CLAVE PUBLICA DEL PACIENTE
	pacienteIdString := strconv.Itoa(entrada.HistorialId)
	pacientePublicKey := getUserPublicKeyByHistorialId(pacienteIdString)

	//Ciframos la clave AES usada con nuestra clave pública
	claveAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(pacientePublicKey.PublicKey))

	//CIFRAMOS LOS DATOS CON LA CLAVE MAESTRA
	//Recuperamos CLAVE PUBLICA MAESTRA
	masterPairKeys := getPublicMasterKey()
	claveMaestraAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(masterPairKeys.PublicKey))

	//Preparamos los datos para enviar
	entrada.UserToken = prepareUserToken(req)
	entrada.Clave = claveAEScifrada
	entrada.ClaveMaestra = claveMaestraAEScifrada
	locJson, err := json.Marshal(entrada)

	//Certificado
	client := GetTLSClient()

	//Request al servidor para añadir entrada paciente
	response, err := client.Post(SERVER_URL+"/user/doctor/historial/addEntrada", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			//CIFRADO PARA MEDICO
			//Recuperamos la clave publica del medico
			userId, _ := session.Values["userId"].(string)
			medicoPairkeys := getUserPairKeys(userId)
			//Ciframos la clave AES usada con nuestra clave pública
			claveAEScifrada = util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(medicoPairkeys.PublicKey))
			//Preparamos los datos para enviar
			entradaHistorialCompartida := entrada
			entradaHistorialCompartida.Id, _ = strconv.Atoi(result.Result)
			entradaHistorialCompartida.Clave = claveAEScifrada
			locJson, err = json.Marshal(entradaHistorialCompartida)
			//Request al servidor para añadir entrada compartida
			response, err = client.Post(SERVER_URL+"/user/doctor/citas/addEntradaCompartida", "application/json", bytes.NewBuffer(locJson))
			if response != nil {
				err := json.NewDecoder(response.Body).Decode(&result)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				util.PrintErrorLog(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			result = util.JSON_Return{Result: "OK"}
		}
		js, err := json.Marshal(result)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addAnaliticaHistorialMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	//Preparamos datos request
	var analitica util.AnaliticaHistorial_JSON
	var analiticaEstadistica util.AnaliticaHistorial_JSON
	json.NewDecoder(req.Body).Decode(&analitica)
	analiticaEstadistica = analitica
	analiticaEstadistica.HistorialId = -1

	//Generamos una clave AES aleatoria de 256 bits para cifrar los datos sensibles
	AESkeyDatos := util.AEScreateKey()

	//Ciframos los datos sensibles con la clave
	analitica.Hematies, _ = util.AESencrypt(AESkeyDatos, analitica.Hematies)
	analitica.Glucosa, _ = util.AESencrypt(AESkeyDatos, analitica.Glucosa)
	analitica.Hierro, _ = util.AESencrypt(AESkeyDatos, analitica.Hierro)
	analitica.Leucocitos, _ = util.AESencrypt(AESkeyDatos, analitica.Leucocitos)
	analitica.Plaquetas, _ = util.AESencrypt(AESkeyDatos, analitica.Plaquetas)

	//Pasamos la clave a base 64
	AESkeyBase64String := string(util.Base64Encode(AESkeyDatos))

	//CIFRAMOS LA CLAVE AES CON LA CLAVE PUBLICA DEL PACIENTE
	pacienteIdString := strconv.Itoa(analitica.HistorialId)
	pacientePublicKey := getUserPublicKeyByHistorialId(pacienteIdString)

	//Ciframos la clave AES usada con nuestra clave pública
	claveAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(pacientePublicKey.PublicKey))

	//CIFRAMOS LOS DATOS CON LA CLAVE MAESTRA
	//Recuperamos CLAVE PUBLICA MAESTRA
	masterPairKeys := getPublicMasterKey()
	claveMaestraAEScifrada := util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(masterPairKeys.PublicKey))

	//Preparamos los datos para enviar
	analiticaEstadistica.UserToken = prepareUserToken(req)
	analitica.UserToken = prepareUserToken(req)
	analitica.Clave = claveAEScifrada
	analitica.ClaveMaestra = claveMaestraAEScifrada
	locJson, err := json.Marshal(analitica)

	//Certificado
	client := GetTLSClient()

	//Request al servidor para añadir analitica paciente
	response, err := client.Post(SERVER_URL+"/user/doctor/historial/addAnalitica", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			//CIFRADO PARA MEDICO
			//Recuperamos la clave publica del medico
			userId, _ := session.Values["userId"].(string)
			medicoPairkeys := getUserPairKeys(userId)
			//Ciframos la clave AES usada con nuestra clave pública
			claveAEScifrada = util.RSAEncryptOAEP(AESkeyBase64String, *util.RSABytesToPublicKey(medicoPairkeys.PublicKey))
			//Preparamos los datos para enviar
			analiticaCompartida := analitica
			analiticaCompartida.Id, _ = strconv.Atoi(result.Result)
			analiticaCompartida.Clave = claveAEScifrada
			analiticaCompartida.EmpleadoId, _ = strconv.Atoi(analitica.UserToken.UserId)
			locJson, err := json.Marshal(analiticaCompartida)

			//Request al servidor para compartir analitica
			response, err = client.Post(SERVER_URL+"/user/doctor/historial/addAnaliticaCompartida", "application/json", bytes.NewBuffer(locJson))
			if response != nil {
				//Request al servidor para añadir estadistica analitica
				response, err = client.Post(SERVER_URL+"/user/doctor/historial/addEstadisticaAnalitica", "application/json", bytes.NewBuffer(locJson))
				if response != nil {
					result = util.JSON_Return{Result: "OK"}
					js, err := json.Marshal(result)
					if err != nil {
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
			} else {
				util.PrintErrorLog(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
