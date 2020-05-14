package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "../models"
	util "../utils"
)

//POST
func addUserHandler(w http.ResponseWriter, req *http.Request) {
	var user util.User_JSON
	json.NewDecoder(req.Body).Decode(&user)
	jsonReturn := util.JSON_Return{"", ""}
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(user.UserToken.UserId, user.UserToken.Token, models.Rol_administradorG.Id)
	if authorized == true {
		util.PrintLog("Insertando usuario " + user.Email)
		//INSERTAMOS EL USUARIO
		userId, err := models.InsertUser(user)
		if err == nil {
			user.Id = userId
			//INSERTAMOS CLAVES RSA
			_, err := models.InsertUserPairKeys(userId, user.PairKeys)
			if err != nil {
				jsonReturn = util.JSON_Return{"", "Error insertando las claves del usuario en la base de datos"}
			} else {
				//INSERTAMOS LA CLAVE MAESTRA SI TIENE
				if user.MasterPairKeys.PrivateKey != nil {
					_, err = models.InsertUserMasterPairKeys(userId, user.MasterPairKeys)
					if err != nil {
						util.PrintErrorLog(err)
						jsonReturn = util.JSON_Return{Error: "Las claves maestras no se han podido insertar"}
					}
					//Insertamos nombre medico
					models.InsertNombresEmpleado(user)
				}
				//Insertamos DNI hasheado
				_, err := models.InsertUserDniHash(userId, user.IdentificacionHash)
				if err != nil {
					jsonReturn = util.JSON_Return{"", "El documento de identificación ya existe en la base de datos"}
					models.DeleteUser(user.Id)
				} else {
					//INSERTAMOS LOS ROLES DEL USUARIO
					inserted, err := models.InsertUserAndRole(userId, user.Roles)
					if err == nil && inserted == true {
						jsonReturn = util.JSON_Return{"OK", ""}
						//INSERTAR EL USUARIO EN LAS CLINICAS
						clinicaId, _ := strconv.Atoi(user.EnfermeroClinica)
						if clinicaId != -1 {
							result, err := models.InsertarUserClinica(clinicaId, userId, models.Rol_enfermero.Id)
							if err != nil || result == false {
								jsonReturn = util.JSON_Return{"", "Error insertando el usuario en la clínica"}
							}
							//Insertamos nombre medico
							models.InsertNombresEmpleado(user)
						}
						clinicaId, _ = strconv.Atoi(user.MedicoClinica)
						if clinicaId != -1 {
							result, err := models.InsertarUserClinica(clinicaId, userId, models.Rol_medico.Id)
							if err != nil || result == false {
								jsonReturn = util.JSON_Return{"", "Error insertando el usuario en la clínica"}
							}
							//Insertamos nombre medico
							models.InsertNombresEmpleado(user)
						}
						clinicaId, _ = strconv.Atoi(user.AdminClinica)
						if clinicaId != -1 {
							result, err := models.InsertarUserClinica(clinicaId, userId, models.Rol_administradorC.Id)
							if err != nil || result == false {
								jsonReturn = util.JSON_Return{"", "Error insertando el usuario en la clínica"}
							}
							//Insertamos nombre medico
							models.InsertNombresEmpleado(user)
						}
						especialidadId, _ := strconv.Atoi(user.MedicoEspecialidad)
						if especialidadId != -1 {
							result, err := models.InsertEspecialidadMedico(userId, especialidadId)
							if err != nil || result == false {
								jsonReturn = util.JSON_Return{"", "Error insertando el usuario en la clínica"}
							}
						}
					} else {
						jsonReturn = util.JSON_Return{"", "Error insertando los roles del usuario en la base de datos"}
					}
				}
			}
		} else {
			jsonReturn = util.JSON_Return{"", err.Error()}
		}
	} else {
		jsonReturn = util.JSON_Return{"", "No dispones de permisos para realizar esa acción"}
	}
	js, err := json.Marshal(jsonReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func menuUserEditHandler(w http.ResponseWriter, req *http.Request) {
	var user util.User_JSON
	json.NewDecoder(req.Body).Decode(&user)
	util.PrintLog("Editando datos del usuario " + user.Email)
	inserted, err := models.EditUserData(user)
	jsonReturn := util.JSON_Return{"", ""}
	if inserted == true {
		jsonReturn = util.JSON_Return{"OK", ""}
	} else {
		jsonReturn = util.JSON_Return{"", "El usuario no se ha podido editar"}
	}

	js, err := json.Marshal(jsonReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

//DELETE
func deleteUserHandler(w http.ResponseWriter, req *http.Request) {
	var user util.User_id_JSON
	json.NewDecoder(req.Body).Decode(&user)
	jsonReturn := util.JSON_Return{"", ""}
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(user.UserToken.UserId, user.UserToken.Token, models.Rol_administradorG.Id)
	if authorized == true {
		util.PrintLog("Borrando Usuario con ID " + strconv.Itoa(user.Id))
		result, err := models.DeleteUser(user.Id)
		if err == nil {
			if result == true {
				util.PrintLog("El usuario con ID " + strconv.Itoa(user.Id) + " ha sido borrado correctamente")
				jsonReturn = util.JSON_Return{"OK", ""}
			} else {
				util.PrintLog("El usuario con ID " + strconv.Itoa(user.Id) + "no se ha podido borrar")
				jsonReturn = util.JSON_Return{"", "El usuario con ID " + strconv.Itoa(user.Id) + " no se ha podido borrar"}
			}
		} else {
			jsonReturn = util.JSON_Return{"", err.Error()}
		}
	} else {
		jsonReturn = util.JSON_Return{"", "No dispones de permisos para realizar esa acción"}
	}
	js, err := json.Marshal(jsonReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
