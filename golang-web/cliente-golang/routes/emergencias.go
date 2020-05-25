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
func menuEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/user/emergencias/index.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Menú emergencias", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//GET
func GetHistorialEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
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
	var user util.User_JSON
	json.NewDecoder(req.Body).Decode(&user)
	user.UserToken = prepareUserToken(req)
	locJson, err := json.Marshal(user)

	//Certificado
	client := GetTLSClient()

	//Request para obtener historial si existe
	response, err := client.Post(SERVER_URL+"/user/emergency/historial", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.Historial_JSON
		err := json.NewDecoder(response.Body).Decode(&result)

		//DESCIFRAMOS DATOS CON CLAVE MAESTRA
		//Recuperamos la CLAVE MAESTRA
		userId, _ := session.Values["userId"].(string)
		masterPairKeys := getUserMasterPairKeys(userId)

		//Desciframos la clave privada CLAVE MAESTRA cifrada con AES
		userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)
		masterPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(util.Base64Decode(masterPairKeys.PrivateKey)))
		masterPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(masterPrivateKeyString)))

		//Desciframos la clave AES maestra
		claveAEShistorial := util.RSADecryptOAEP(result.ClaveMaestra, *masterPrivateKey)
		claveAEShistorialByte := util.Base64Decode([]byte(claveAEShistorial))

		//Desciframos los datos del historial con AES
		result.NombrePaciente, _ = util.AESdecrypt(claveAEShistorialByte, result.NombrePaciente)
		result.ApellidosPaciente, _ = util.AESdecrypt(claveAEShistorialByte, result.ApellidosPaciente)
		result.Sexo, _ = util.AESdecrypt(claveAEShistorialByte, result.Sexo)
		result.Alergias, _ = util.AESdecrypt(claveAEShistorialByte, result.Alergias)

		//Desciframos el Tipo de las entradas del Historial
		for index, _ := range result.Entradas {
			//Desciframos la clave AES maestra
			claveAESentrada := util.RSADecryptOAEP(result.Entradas[index].ClaveMaestra, *masterPrivateKey)
			claveAESentradaByte := util.Base64Decode([]byte(claveAESentrada))
			//Desciframos los datos
			result.Entradas[index].Tipo, _ = util.AESdecrypt(claveAESentradaByte, result.Entradas[index].Tipo)
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

//GET
func GetEntradaEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
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
	response, err := client.Post(SERVER_URL+"/user/emergency/historial/entrada", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&entradaJSON)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//DESCIFRAMOS DATOS CON CLAVE MAESTRA
		//Recuperamos la CLAVE MAESTRA
		userId, _ := session.Values["userId"].(string)
		masterPairKeys := getUserMasterPairKeys(userId)

		//Desciframos la clave privada CLAVE MAESTRA cifrada con AES
		userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)
		masterPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(util.Base64Decode(masterPairKeys.PrivateKey)))
		masterPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(masterPrivateKeyString)))

		//Desciframos la clave AES maestra
		claveAESentrada := util.RSADecryptOAEP(entradaJSON.ClaveMaestra, *masterPrivateKey)
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
		template.New("").ParseFiles("public/templates/user/emergencias/entrada.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.EntradaPage{Title: "Consultar entrada", Body: "body", Entrada: entradaJSON}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func GetAnaliticaEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
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

	//Request para obtener historial si existe
	response, err := client.Post(SERVER_URL+"/user/emergency/historial/analitica", "application/json", bytes.NewBuffer(locJson))
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
		masterPairKeys := getUserMasterPairKeys(userId)

		//Desciframos la clave privada CLAVE MAESTRA cifrada con AES
		userPrivateKeyHash, _ := session.Values["userPrivateKeyHash"].([]byte)
		masterPrivateKeyString, _ := util.AESdecrypt(userPrivateKeyHash, string(util.Base64Decode(masterPairKeys.PrivateKey)))
		masterPrivateKey := util.RSABytesToPrivateKey(util.Base64Decode([]byte(masterPrivateKeyString)))

		//Desciframos la clave AES maestra
		claveAESentrada := util.RSADecryptOAEP(analiticaJSON.ClaveMaestra, *masterPrivateKey)
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
		template.New("").ParseFiles("public/templates/user/emergencias/analitica.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &util.AnaliticaPage{Title: "Consultar analítica", Body: "body", Analitica: analiticaJSON}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func AddEntradaEmergenciasFormHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/user/emergencias/addEntrada.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Añadir entrada", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func AddAnaliticaEmergenciasFormHandler(w http.ResponseWriter, req *http.Request) {
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
		template.New("").ParseFiles("public/templates/user/emergencias/addAnalitica.html", "public/templates/layouts/base.html"),
	)
	if err := tmp.ExecuteTemplate(w, "base", &Page{Title: "Añadir analítica", Body: "body"}); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

//POST

func AddEntradaEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
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
	response, err := client.Post(SERVER_URL+"/user/emergency/addEntrada", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
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

func AddAnaliticaEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
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
	response, err := client.Post(SERVER_URL+"/user/emergency/addAnalitica", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		var result util.JSON_Return
		err := json.NewDecoder(response.Body).Decode(&result)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			locJson, err := json.Marshal(analiticaEstadistica)
			//Request al servidor para añadir estadistica analitica
			response, err = client.Post(SERVER_URL+"/user/emergency/addEstadisticaAnalitica", "application/json", bytes.NewBuffer(locJson))
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
		}
	} else {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
