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
func menuPatientHandler(w http.ResponseWriter, req *http.Request) {
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

	//VER SOLICITUDES
	//Certificado
	client := GetTLSClient()
	permisos := false
	//Recuperamos datos de nuestro historial
	locJson, err := json.Marshal(util.UserToken_JSON{UserId: prepareUserToken(req).UserId, Token: prepareUserToken(req).Token})
	response, err := client.Post(SERVER_URL+"/permisos/solicitudes/comprobar", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if result.Result == "OK" {
			permisos = true
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/paciente/index.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", util.PacientePage{Title: "Menú paciente", Body: "body", Permisos: permisos}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func patientAutorizationsHandler(w http.ResponseWriter, req *http.Request) {
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

	//LISTAR SOLICITUDES
	//Certificado
	client := GetTLSClient()
	//Recuperamos datos de las solicitudes
	var solicitudes []util.Solicitud_JSON
	locJson, err := json.Marshal(util.UserToken_JSON{UserId: prepareUserToken(req).UserId, Token: prepareUserToken(req).Token})
	response, err := client.Post(SERVER_URL+"/permisos/solicitudes/listar", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&solicitudes)
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
		template.New("").ParseFiles("public/templates/user/paciente/historial/autorizar.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", util.PermisosPage{Title: "Gestión de permisos", Body: "body", Solicitudes: solicitudes}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func editUserPatientHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/user/paciente/edit.html", "public/templates/layouts/menuPaciente.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Mis datos", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func historialPatientHandler(w http.ResponseWriter, req *http.Request) {

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
	var historial util.Historial_JSON

	//Recuperamos datos de nuestro historial
	locJson, err := json.Marshal(util.UserToken_JSON{UserId: prepareUserToken(req).UserId, Token: prepareUserToken(req).Token})
	response, err := client.Post(SERVER_URL+"/user/patient/historial", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&historial)
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

	//Desciframos la clave AES de los datos cifrados
	claveAEShistorial := util.RSADecryptOAEP(historial.Clave, *userPrivateKey)
	claveAEShistorialByte := util.Base64Decode([]byte(claveAEShistorial))

	//Desciframos los datos del historial con AES
	historial.Alergias, _ = util.AESdecrypt(claveAEShistorialByte, historial.Alergias)
	historial.Sexo, _ = util.AESdecrypt(claveAEShistorialByte, historial.Sexo)

	for index, entrada := range historial.Entradas {
		//Desciframos la clave AES de los datos cifrados
		claveAESentrada := util.RSADecryptOAEP(entrada.Clave, *userPrivateKey)
		claveAESentradaByte := util.Base64Decode([]byte(claveAESentrada))

		//Desciframos los datos de la entrada con AES
		historial.Entradas[index].MotivoConsulta, _ = util.AESdecrypt(claveAESentradaByte, entrada.MotivoConsulta)
		historial.Entradas[index].JuicioDiagnostico, _ = util.AESdecrypt(claveAESentradaByte, entrada.JuicioDiagnostico)
		historial.Entradas[index].Tipo, _ = util.AESdecrypt(claveAESentradaByte, entrada.Tipo)
	}

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/paciente/historial/index.html", "public/templates/layouts/menuPaciente.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", util.HistorialPage{Title: "Mi historia clínica", Body: "body", Historial: historial}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func patientCitaListHandler(w http.ResponseWriter, req *http.Request) {

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
	response, err := client.Post(SERVER_URL+"/user/patient/citas/list", "application/json", bytes.NewBuffer(locJson))
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
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/paciente/citas/list.html", "public/templates/layouts/menuPaciente.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.CitaListPage{Title: "Citas pendientes", Body: "body", Citas: citasList}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func addPatientCitaFormHandler(w http.ResponseWriter, req *http.Request) {

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

	// Request /hello via the created HTTPS client over port 5001 via GET
	response, err := client.Get(SERVER_URL + "/clinica/list/query")
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//Request al servidor para comprobar usuario/pass
	var serverReq []util.Clinica_JSON
	json.NewDecoder(response.Body).Decode(&serverReq)

	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/paciente/citas/add.html", "public/templates/layouts/menuPaciente.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.CitaPage{Title: "Solicitar cita", Body: "body", Clinicas: serverReq}); err != nil {
		log.Printf("Error executing template: %v", err)
		util.PrintErrorLog(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func historialEntradaPatientHandler(w http.ResponseWriter, req *http.Request) {
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

	//Request para obtener historial si existe
	response, err := client.Post(SERVER_URL+"/user/patient/historial/entrada", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&entradaJSON)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//DESCIFRAMOS DATOS CON NUESTRA CLAVE
		//Recuperamos NUESTRA CLAVE
		userId, _ := session.Values["userId"].(string)
		pairKeys := getUserPairKeys(userId)

		//Desciframos la clave privada CLAVE MAESTRA cifrada con AES
		userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)
		PrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(pairKeys.PrivateKey))
		PrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(PrivateKeyString)))

		//Desciframos la clave AES maestra
		claveAESentrada := util.RSADecryptOAEP(entradaJSON.Clave, *PrivateKey)
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
		template.New("").ParseFiles("public/templates/user/paciente/historial/entrada.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.EntradaPage{Title: "Consultar entrada", Body: "body", Entrada: entradaJSON}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func historialAnaliticaPatientHandler(w http.ResponseWriter, req *http.Request) {
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
	response, err := client.Post(SERVER_URL+"/user/patient/historial/analitica", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&analiticaJSON)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//DESCIFRAMOS DATOS CON CLAVE MAESTRA
		//Recuperamos la CLAVE MAESTRA
		userId, _ := session.Values["userId"].(string)
		PairKeys := getUserPairKeys(userId)

		//Desciframos la clave privada CLAVE MAESTRA cifrada con AES
		userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)
		PrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(PairKeys.PrivateKey))
		PrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(PrivateKeyString)))

		//Desciframos la clave AES maestra
		claveAESentrada := util.RSADecryptOAEP(analiticaJSON.Clave, *PrivateKey)
		claveAESentradaByte := util.Base64Decode([]byte(claveAESentrada))

		//Desciframos los datos del historial con AES
		analiticaJSON.Leucocitos, _ = util.AESdecrypt(claveAESentradaByte, analiticaJSON.Leucocitos)
		analiticaJSON.Hematies, _ = util.AESdecrypt(claveAESentradaByte, analiticaJSON.Hematies)
		analiticaJSON.Hierro, _ = util.AESdecrypt(claveAESentradaByte, analiticaJSON.Hierro)
		analiticaJSON.Glucosa, _ = util.AESdecrypt(claveAESentradaByte, analiticaJSON.Glucosa)
		analiticaJSON.Plaquetas, _ = util.AESdecrypt(claveAESentradaByte, analiticaJSON.Plaquetas)
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var tmp = template.Must(
		template.New("").ParseFiles("public/templates/user/paciente/historial/analitica.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.AnaliticaPage{Title: "Consultar analítica", Body: "body", Analitica: analiticaJSON}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//POST

//reservar cita
func addCitaPacienteHandler(w http.ResponseWriter, req *http.Request) {
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
	var cita util.CitaJSON
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&cita)
	if err != nil {
		util.PrintErrorLog(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//Certificado
	client := GetTLSClient()
	var historial util.Historial_JSON

	//Recuperamos datos de nuestro historial
	locJson, err := json.Marshal(util.UserToken_JSON{UserId: prepareUserToken(req).UserId, Token: prepareUserToken(req).Token})
	response, err := client.Post(SERVER_URL+"/user/patient/historial", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&historial)
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

	//Desciframos la clave AES de los datos del usuario
	userDataKey, _ := session.Values["userDataKey"].(string)
	claveAESuserData := util.RSADecryptOAEP(userDataKey, *userPrivateKey)
	claveAESuserDataByte := util.Base64Decode([]byte(claveAESuserData))

	//Recuperamos la clave publica del medico
	medicoPairkeys := getUserPairKeys(cita.MedicoId)
	medicoPublicKey := *util.RSABytesToPublicKey(medicoPairkeys.PublicKey)

	//Ciframos los datos del historial con AES y terminamos de rellenar el historial
	var historialCompartido util.Historial_JSON
	historialCompartido.Id = historial.Id
	historialCompartido.MedicoId, _ = strconv.Atoi(cita.MedicoId)
	historialCompartido.PacienteId, _ = strconv.Atoi(prepareUserToken(req).UserId)
	historialCompartido.UserToken = prepareUserToken(req)

	//Ciframos la clave AES con la clave publica del Medico
	historialCompartido.Clave = util.RSAEncryptOAEP(string(util.Base64Encode(claveAESuserDataByte)), medicoPublicKey)

	//Enviamos historial compartido al servidor
	locJson, err = json.Marshal(historialCompartido)
	response, err = client.Post(SERVER_URL+"/user/patient/historial/share", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&historial)
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

	//Request al servidor para registrar cita
	locJson, err = json.Marshal(util.CitaJSON{FechaString: cita.FechaString, MedicoId: cita.MedicoId, UserToken: prepareUserToken(req)})
	response, err = client.Post(SERVER_URL+"/user/patient/citas/add", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var responseJSON JSON_Return
		err := json.NewDecoder(response.Body).Decode(&responseJSON)
		js, err := json.Marshal(responseJSON)
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
