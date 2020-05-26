package routes

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	util "../utils"
)

//GET
func menuMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	locJson, err := json.Marshal(prepareUserToken(req))

	//Certificado
	client := GetTLSClient()
	var cita util.CitaJSON

	//Request al servidor para obtener citas futuras
	response, err := client.Post(SERVER_URL+"/user/doctor/citas/actual", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&cita)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/index.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.PageMenuMedico{Title: "Menú médico", Body: "body", CitaActual: cita.Id}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func solicitarHistorialMedicoFormHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/user/medico/historial/solicitar.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Solicitar historial", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func solicitarHistorialMedicoHandler(w http.ResponseWriter, req *http.Request) {
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
	var user util.User
	json.NewDecoder(req.Body).Decode(&user)
	var solicitarHistorial util.SolicitarHistorial_JSON
	solicitarHistorial.UserDNI = user.Identificacion
	solicitarHistorial.UserToken = prepareUserToken(req)
	locJson, err := json.Marshal(solicitarHistorial)

	//Certificado
	client := GetTLSClient()

	//Request al servidor para recibir lista de roles
	response, err := client.Post(SERVER_URL+"/user/doctor/historial/solicitar", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.SolicitarHistorial_JSON
		err := json.NewDecoder(response.Body).Decode(&result)

		//Recuperamos nuestra clave privada cifrada
		userId, _ := session.Values["userId"].(string)
		userPairkeys := getUserPairKeys(userId)
		userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)

		//Desciframos nuestra clave privada cifrada con AES
		userPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(userPairkeys.PrivateKey))
		userPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(userPrivateKeyString)))

		//DESCIFRADO DE DATOS
		claveAEShistorial := util.RSADecryptOAEP(result.HistorialPermitido.Clave, *userPrivateKey)
		claveAEShistorialByte := util.Base64Decode([]byte(claveAEShistorial))
		//Desciframos los datos del historial con AES
		result.HistorialPermitido.NombrePaciente, _ = util.AESdecrypt(claveAEShistorialByte, result.HistorialPermitido.NombrePaciente)
		result.HistorialPermitido.ApellidosPaciente, _ = util.AESdecrypt(claveAEShistorialByte, result.HistorialPermitido.ApellidosPaciente)
		result.HistorialPermitido.Sexo, _ = util.AESdecrypt(claveAEShistorialByte, result.HistorialPermitido.Sexo)
		result.HistorialPermitido.Alergias, _ = util.AESdecrypt(claveAEShistorialByte, result.HistorialPermitido.Alergias)
		for index, entrada := range result.HistorialPermitido.Entradas {
			if entrada.Clave != "" {
				//Desciframos la clave AES de los datos cifrados
				claveAESentrada := util.RSADecryptOAEP(entrada.Clave, *userPrivateKey)
				claveAESentradaByte := util.Base64Decode([]byte(claveAESentrada))
				//Desciframos los datos de la entrada con AES
				result.HistorialPermitido.Entradas[index].Tipo, _ = util.AESdecrypt(claveAESentradaByte, entrada.Tipo)
			} else {
				result.HistorialPermitido.Entradas[index].Tipo = ""
			}
		}

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
}

func solicitarPermisoTotalHistorialMedicoHandler(w http.ResponseWriter, req *http.Request) {
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
	var user util.User
	json.NewDecoder(req.Body).Decode(&user)
	var solicitarHistorial util.SolicitarHistorial_JSON
	solicitarHistorial.UserDNI = user.Identificacion
	solicitarHistorial.UserToken = prepareUserToken(req)
	locJson, err := json.Marshal(solicitarHistorial)

	//Certificado
	client := GetTLSClient()

	//Request al servidor para recibir lista de roles
	response, err := client.Post(SERVER_URL+"/user/doctor/historial/solicitar", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.SolicitarHistorial_JSON
		err := json.NewDecoder(response.Body).Decode(&result)

		//Recuperamos nuestra clave privada cifrada
		userId, _ := session.Values["userId"].(string)
		userPairkeys := getUserPairKeys(userId)
		userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)

		//Desciframos nuestra clave privada cifrada con AES
		userPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(userPairkeys.PrivateKey))
		userPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(userPrivateKeyString)))

		//DESCIFRADO DE DATOS
		claveAEShistorial := util.RSADecryptOAEP(result.HistorialPermitido.Clave, *userPrivateKey)
		claveAEShistorialByte := util.Base64Decode([]byte(claveAEShistorial))
		//Desciframos los datos del historial con AES
		result.HistorialPermitido.NombrePaciente, _ = util.AESdecrypt(claveAEShistorialByte, result.HistorialPermitido.NombrePaciente)
		result.HistorialPermitido.ApellidosPaciente, _ = util.AESdecrypt(claveAEShistorialByte, result.HistorialPermitido.ApellidosPaciente)
		result.HistorialPermitido.Sexo, _ = util.AESdecrypt(claveAEShistorialByte, result.HistorialPermitido.Sexo)
		result.HistorialPermitido.Alergias, _ = util.AESdecrypt(claveAEShistorialByte, result.HistorialPermitido.Alergias)
		for index, entrada := range result.HistorialPermitido.Entradas {
			if entrada.Clave != "" {
				//Desciframos la clave AES de los datos cifrados
				claveAESentrada := util.RSADecryptOAEP(entrada.Clave, *userPrivateKey)
				claveAESentradaByte := util.Base64Decode([]byte(claveAESentrada))
				//Desciframos los datos de la entrada con AES
				result.HistorialPermitido.Entradas[index].Tipo, _ = util.AESdecrypt(claveAESentradaByte, entrada.Tipo)
			} else {
				result.HistorialPermitido.Entradas[index].Tipo = ""
			}
		}

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
}

//Listar medicos dada una especialidad de la clinica
func getMedicoDiasDisponiblesHandler(w http.ResponseWriter, req *http.Request) {
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
	medicoId, _ := req.URL.Query()["doctorId"]
	//Certificado
	client := GetTLSClient()

	response, err := client.Get(SERVER_URL + "/user/doctor/disponible/dia?doctorId=" + medicoId[0])
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		var citasListJSON []util.Cita
		err := json.NewDecoder(response.Body).Decode(&citasListJSON)
		js, err := json.Marshal(citasListJSON)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

//Listar medicos dada una especialidad de la clinica
func getMedicoHorasDiaDisponiblesHandler(w http.ResponseWriter, req *http.Request) {
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
	medicoId, _ := req.URL.Query()["doctorId"]
	dia, _ := req.URL.Query()["dia"]
	//Certificado
	client := GetTLSClient()
	diaTratado := strings.Replace(dia[0], " ", "", -1)
	response, err := client.Get(SERVER_URL + "/user/doctor/disponible/hora?doctorId=" + medicoId[0] + "&dia=" + diaTratado)
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		var citasListJSON []util.Cita
		err := json.NewDecoder(response.Body).Decode(&citasListJSON)
		js, err := json.Marshal(citasListJSON)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func medicoCitaListHandler(w http.ResponseWriter, req *http.Request) {

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

	locJson, err := json.Marshal(prepareUserToken(req))

	//Certificado
	client := GetTLSClient()
	var citasList []util.CitaJSON

	//Request al servidor para obtener citas futuras
	response, err := client.Post(SERVER_URL+"/user/doctor/citas/list", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&citasList)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Recuperamos nuestra clave privada cifrada
	userId, _ := session.Values["userId"].(string)
	userPairkeys := getUserPairKeys(userId)
	userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)

	//Desciframos nuestra clave privada cifrada con AES
	userPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(userPairkeys.PrivateKey))
	userPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(userPrivateKeyString)))

	//DESCIFRADO DE DATOS
	for index, cita := range citasList {
		//Desciframos la clave AES de los datos cifrados
		claveAEShistorial := util.RSADecryptOAEP(cita.Historial.Clave, *userPrivateKey)
		claveAEShistorialByte := util.Base64Decode([]byte(claveAEShistorial))
		//Desciframos los datos del historial con AES
		citasList[index].Historial.NombrePaciente, _ = util.AESdecrypt(claveAEShistorialByte, cita.Historial.NombrePaciente)
		citasList[index].Historial.ApellidosPaciente, _ = util.AESdecrypt(claveAEShistorialByte, cita.Historial.ApellidosPaciente)
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/citas/list.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.CitaListPage{Title: "Citas pendientes", Body: "body", Citas: citasList}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//GET
func getCitaFormMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	//Preparamos Id cita
	citaId, _ := req.URL.Query()["citaId"]
	userId_int, _ := strconv.Atoi(citaId[0])

	//Certificado
	client := GetTLSClient()
	var cita util.CitaJSON
	cita.UserToken = prepareUserToken(req)
	cita.Id = userId_int
	locJson, err := json.Marshal(cita)

	//Request al servidor para obtener la cita
	response, err := client.Post(SERVER_URL+"/user/doctor/citas", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&cita)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//DESCIFRAMOS DATOS HISTORIAL
	//Recuperamos nuestra clave privada cifrada
	userId, _ := session.Values["userId"].(string)
	userPairkeys := getUserPairKeys(userId)
	userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)
	//Desciframos nuestra clave privada cifrada con AES
	userPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(userPairkeys.PrivateKey))
	userPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(userPrivateKeyString)))
	//Desciframos la clave AES de los datos cifrados
	claveAEShistorial := util.RSADecryptOAEP(cita.Historial.Clave, *userPrivateKey)
	claveAEShistorialByte := util.Base64Decode([]byte(claveAEShistorial))
	//Desciframos los datos del historial con AES
	cita.Historial.NombrePaciente, _ = util.AESdecrypt(claveAEShistorialByte, cita.Historial.NombrePaciente)
	cita.Historial.ApellidosPaciente, _ = util.AESdecrypt(claveAEShistorialByte, cita.Historial.ApellidosPaciente)
	cita.Historial.Sexo, _ = util.AESdecrypt(claveAEShistorialByte, cita.Historial.Sexo)
	cita.Historial.Alergias, _ = util.AESdecrypt(claveAEShistorialByte, cita.Historial.Alergias)

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/citas/index.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", util.ConsultaPage{Title: "Pasar consulta", Body: "body", Cita: cita}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getListHistorialMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	locJson, err := json.Marshal(prepareUserToken(req))

	//Certificado
	client := GetTLSClient()
	var historialList []util.Historial_JSON

	//Request al servidor para obtener historiales compartidos
	response, err := client.Post(SERVER_URL+"/user/doctor/historial/list", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&historialList)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Recuperamos nuestra clave privada cifrada
	userId, _ := session.Values["userId"].(string)
	userPairkeys := getUserPairKeys(userId)
	userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)

	//Desciframos nuestra clave privada cifrada con AES
	userPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(userPairkeys.PrivateKey))
	userPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(userPrivateKeyString)))

	//DESCIFRADO DE DATOS
	for index, historial := range historialList {
		//Desciframos la clave AES de los datos cifrados
		claveAEShistorial := util.RSADecryptOAEP(historial.Clave, *userPrivateKey)
		claveAEShistorialByte := util.Base64Decode([]byte(claveAEShistorial))
		//Desciframos los datos del historial con AES
		historialList[index].NombrePaciente, _ = util.AESdecrypt(claveAEShistorialByte, historial.NombrePaciente)
		historialList[index].ApellidosPaciente, _ = util.AESdecrypt(claveAEShistorialByte, historial.ApellidosPaciente)
		historialList[index].Sexo, _ = util.AESdecrypt(claveAEShistorialByte, historial.Sexo)
	}
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/historial/list.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.HistorialListPage{Title: "Historiales compartidos", Body: "body", Historiales: historialList}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getEntradaHistorialMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	//Certificado
	client := GetTLSClient()

	//Preparamos datos de la request
	entradaId, _ := req.URL.Query()["entradaId"]
	entradaIdInt, _ := strconv.Atoi(entradaId[0])
	entradaJSON := util.EntradaHistorial_JSON{Id: entradaIdInt, UserToken: prepareUserToken(req)}
	locJson, err := json.Marshal(entradaJSON)

	//Request para obtener entrada si existe
	response, err := client.Post(SERVER_URL+"/user/doctor/historial/entrada", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&entradaJSON)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//DESCIFRAMOS DATOS CON CLAVE
		//Recuperamos la CLAVE
		userId, _ := session.Values["userId"].(string)
		pairKeys := getUserPairKeys(userId)

		//Desciframos la clave privada CLAVE cifrada con AES
		userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)
		privateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(pairKeys.PrivateKey))
		privateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(privateKeyString)))

		//Desciframos la clave AES
		claveAESentrada := util.RSADecryptOAEP(entradaJSON.Clave, *privateKey)
		claveAESentradaByte := util.Base64Decode([]byte(claveAESentrada))

		//Desciframos los datos del historial con AES
		entradaJSON.JuicioDiagnostico, _ = util.AESdecrypt(claveAESentradaByte, entradaJSON.JuicioDiagnostico)
		entradaJSON.MotivoConsulta, _ = util.AESdecrypt(claveAESentradaByte, entradaJSON.MotivoConsulta)
		entradaJSON.Tipo, _ = util.AESdecrypt(claveAESentradaByte, entradaJSON.Tipo)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/historial/entrada.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.EntradaPage{Title: "Consultar entrada", Body: "body", Entrada: entradaJSON}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getAnaliticaHistorialMedicoHandler(w http.ResponseWriter, req *http.Request) {
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

	//Certificado
	client := GetTLSClient()

	//Preparamos datos de la request
	analiticaId, _ := req.URL.Query()["analiticaId"]
	analiticaIdInt, _ := strconv.Atoi(analiticaId[0])
	analiticaJSON := util.AnaliticaHistorial_JSON{Id: analiticaIdInt, UserToken: prepareUserToken(req)}
	locJson, err := json.Marshal(analiticaJSON)

	//Request para obtener analitica si existe
	response, err := client.Post(SERVER_URL+"/user/doctor/historial/analitica", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&analiticaJSON)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//DESCIFRAMOS DATOS CON CLAVE
		//Recuperamos la CLAVE
		userId, _ := session.Values["userId"].(string)
		pairKeys := getUserPairKeys(userId)

		//Desciframos la clave privada CLAVE cifrada con AES
		userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)
		privateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(pairKeys.PrivateKey))
		privateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(privateKeyString)))

		//Desciframos la clave AES
		claveAESentrada := util.RSADecryptOAEP(analiticaJSON.Clave, *privateKey)
		claveAESentradaByte := util.Base64Decode([]byte(claveAESentrada))

		//Desciframos los datos de la analítica con AES
		analiticaJSON.Leucocitos, _ = util.AESdecrypt(claveAESentradaByte, analiticaJSON.Leucocitos)
		analiticaJSON.Hematies, _ = util.AESdecrypt(claveAESentradaByte, analiticaJSON.Hematies)
		analiticaJSON.Hierro, _ = util.AESdecrypt(claveAESentradaByte, analiticaJSON.Hierro)
		analiticaJSON.Plaquetas, _ = util.AESdecrypt(claveAESentradaByte, analiticaJSON.Plaquetas)
		analiticaJSON.Glucosa, _ = util.AESdecrypt(claveAESentradaByte, analiticaJSON.Glucosa)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/historial/analitica.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.AnaliticaPage{Title: "Consultar analitica", Body: "body", Analitica: analiticaJSON}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getInvestigacionAnaliticasMedicoFormHandler(w http.ResponseWriter, req *http.Request) {
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
	locJson, err := json.Marshal(prepareUserToken(req))

	//Certificado
	client := GetTLSClient()
	var analiticas []util.AnaliticaHistorial_JSON

	//Request al servidor para obtener citas futuras
	response, err := client.Post(SERVER_URL+"/user/doctor/research/analiticas", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&analiticas)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/medico/investigacion/index.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", util.EstadisticasAnaliticaPage{Title: "Estadísticas analíticas", Body: "body", Analiticas: analiticas}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
