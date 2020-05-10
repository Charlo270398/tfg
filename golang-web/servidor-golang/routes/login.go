package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	models "../models"
	util "../utils"
)

func getInicioHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", req.URL.Path)
}

//POST
func loginUserHandler(w http.ResponseWriter, req *http.Request) {
	util.PrintLog("Intentando iniciar sesión...")
	var creds util.JSON_Credentials_SERVIDOR
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(req.Body).Decode(&creds)
	if err != nil {
		util.PrintErrorLog(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//COMPROBAMOS USER Y PASS
	jsonReturn := util.JSON_Login_Return{}
	user, _ := models.GetUserByIdentificacion(creds.Identificacion)
	correctLogin := models.LoginUser(user.Id, creds.Password)
	if err != nil {
		util.PrintErrorLog(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//RECUPERAMOS CLAVE PUBLICA Y PRIVADA DEL USUARIO
	pairKeys, err := models.GetUserPairKeys(strconv.Itoa(user.Id))
	if err != nil {
		util.PrintErrorLog(err)
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if correctLogin == true {
		token, err := models.InsertUserToken(user.Id)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Example: this will give us a 64 byte output
		jsonReturn = util.JSON_Login_Return{UserId: strconv.Itoa(user.Id), Nombre: user.Nombre, Apellidos: user.Apellidos, Email: user.Email, Token: token, PairKeys: pairKeys, Clave: user.Clave}
	} else {
		jsonReturn = util.JSON_Login_Return{Error: "Usuario y contraseña incorrectos"}
	}
	js, err := json.Marshal(jsonReturn)
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//POST
func registerUserHandler(w http.ResponseWriter, req *http.Request) {
	var user util.User_JSON
	json.NewDecoder(req.Body).Decode(&user)
	util.PrintLog("Insertando usuario " + user.Email)
	userId, err := models.InsertUser(user)
	jsonReturn := util.JSON_Login_Return{}
	if err == nil {
		userlist, err := models.GetUsersList()
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		var rolesList []int
		if len(userlist) == 1 {
			//SI ES EL PRIMER USUARIO DE LA BD LE DAMOS PERMISO DE ADMINISTRADOR GLOBAL
			rolesList = []int{models.Rol_administradorG.Id}
			//INSERTAMOS CLAVES RSA MAESTRAS
			_, err = models.InsertUserMasterPairKeys(userId, user.MasterPairKeys)
			if err != nil {
				util.PrintErrorLog(err)
				jsonReturn = util.JSON_Login_Return{Error: "Las claves no se han podido insertar"}
			}
		} else {
			rolesList = []int{models.Rol_paciente.Id}
		}
		user.Id = userId
		//Insertamos DNI hasheado
		_, err = models.InsertUserDniHash(userId, user.IdentificacionHash)
		if err != nil {
			util.PrintErrorLog(err)
			jsonReturn = util.JSON_Login_Return{Error: "El documento de identificación ya existe en la base de datos"}
			models.DeleteUser(user.Id)
		} else {
			//INSERTAMOS HISTORIAL
			if len(userlist) != 1 {
				//Insertamos Historial
				_, err = models.InsertHistorial(user)
				if err != nil {
					util.PrintErrorLog(err)
					jsonReturn = util.JSON_Login_Return{Error: "El historial no se ha podido insertar"}
				}
			}
			//INSERTAMOS CLAVES RSA
			_, err = models.InsertUserPairKeys(userId, user.PairKeys)
			if err != nil {
				util.PrintErrorLog(err)
				jsonReturn = util.JSON_Login_Return{Error: "Las claves no se han podido insertar"}
			}
			//INSERTAMOS ROLES DEL USUARIO
			inserted, err := models.InsertUserAndRole(userId, rolesList)
			if err == nil && inserted == true {
				//INSERTAMOS EL TOKEN DE LA SESION DEL USUARIO
				token, err := models.InsertUserToken(userId)
				if err != nil {
					util.PrintErrorLog(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				//RECUPERAMOS CLAVE PUBLICA Y PRIVADA DEL USUARIO
				jsonReturn = util.JSON_Login_Return{UserId: strconv.Itoa(user.Id), Nombre: user.Nombre, Apellidos: user.Apellidos, Email: user.Email, Token: token, PairKeys: user.PairKeys, Clave: user.Clave}
			} else {
				jsonReturn = util.JSON_Login_Return{Error: "Los roles no se han podido registrar"}
			}
		}
	} else {
		jsonReturn = util.JSON_Login_Return{Error: "El usuario no se ha podido registrar"}
	}

	js, err := json.Marshal(jsonReturn)
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func proveUserTokenHandler(w http.ResponseWriter, req *http.Request) {
	var userToken util.UserToken_JSON
	json.NewDecoder(req.Body).Decode(&userToken)
	id, err := strconv.Atoi(userToken.UserId)
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	result, err := models.ProveUserToken(id, userToken.Token)
	if err != nil {
		util.PrintErrorLog(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		if result != true {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			jsonReturn := util.JSON_Return{Result: "OK"}
			js, err := json.Marshal(jsonReturn)
			if err != nil {
				util.PrintErrorLog(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}
	}
}
