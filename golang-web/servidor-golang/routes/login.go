package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	models "../models"
	util "../utils"
)

type JSON_Credentials struct {
	Password []byte `json:"password"`
	Email    string `json:"email"`
}

func getInicioHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", req.URL.Path)
}

//POST
func loginUserHandler(w http.ResponseWriter, req *http.Request) {
	util.PrintLog("Intentando iniciar sesión...")
	var creds JSON_Credentials
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
	correctLogin := models.LoginUser(creds.Email, creds.Password)
	user, err := models.GetUserByEmail(creds.Email)
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
		jsonReturn = util.JSON_Login_Return{UserId: strconv.Itoa(user.Id), Nombre: user.Nombre, Apellidos: user.Apellidos, Token: token}
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
			//SI ES EL PRIMER USUARIO DE LA BD LE DAMOS TODOS LOS ROLES SALVO ADMIN DE UNA CLINICA CONCRETA
			rolesList = []int{models.Rol_paciente.Id, models.Rol_enfermero.Id, models.Rol_medico.Id, models.Rol_administradorG.Id}
		} else {
			rolesList = []int{models.Rol_paciente.Id}
		}
		//INSERTAMOS CLAVES RSA
		_, err = models.InsertUserPairKeys(userId, user.PairKeys)
		if err != nil {
			util.PrintErrorLog(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
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
			jsonReturn = util.JSON_Login_Return{UserId: strconv.Itoa(userId), Nombre: user.Nombre, Apellidos: user.Apellidos, Token: token}
		} else {
			jsonReturn = util.JSON_Login_Return{Error: "Los roles no se han podido registrar"}
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
