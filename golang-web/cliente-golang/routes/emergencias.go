package routes

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	util "../utils"
)

//GET
func menuEmergenciasHandler(w http.ResponseWriter, req *http.Request) {
	session, _ := store.Get(req, "userSession")
	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Redirect(w, req, "/login", http.StatusSeeOther)
		return
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
