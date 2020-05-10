package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "../models"
	util "../utils"
)

//POST
func getAdminMenuDataHandler(w http.ResponseWriter, req *http.Request) {
	var userToken util.JSON_Admin_Menu
	json.NewDecoder(req.Body).Decode(&userToken)
	jsonReturn := util.JSON_Admin_Menu{}
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(userToken.UserToken.UserId, userToken.UserToken.Token, models.Rol_administradorC.Id)
	if authorized == true {
		clinica, err := models.GetClinicaByAdmin(userToken.UserToken.UserId)
		if err != nil {
			jsonReturn = util.JSON_Admin_Menu{Error: "Error cargando los datos de la clínica"}
		} else {
			jsonReturn = util.JSON_Admin_Menu{Clinica: clinica, Error: ""}
		}

	} else {
		jsonReturn = util.JSON_Admin_Menu{Error: "No dispones de permisos para realizar esa acción"}
	}
	js, err := json.Marshal(jsonReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func addEnfermeroAdminHandler(w http.ResponseWriter, req *http.Request) {
	var user util.User_JSON
	json.NewDecoder(req.Body).Decode(&user)
	jsonReturn := util.JSON_Return{"", ""}
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.VerifyAdmin(user.UserToken.UserId, user.EnfermeroClinica, user.UserToken.Token)
	if authorized == true {
		util.PrintLog("Insertando usuario " + user.Email)
		//INSERTAMOS EL USUARIO
		userId, err := models.InsertUser(user)
		if err == nil {
			//Insertamos DNI hasheado
			_, err = models.InsertUserDniHash(userId, user.IdentificacionHash)
			if err != nil {
				util.PrintErrorLog(err)
				jsonReturn = util.JSON_Return{"", "El documento de identificación ya existe en la base de datos"}
				models.DeleteUser(user.Id)
			} else {
				//INSERTAMOS CLAVES RSA
				_, err := models.InsertUserPairKeys(userId, user.PairKeys)
				if err != nil {
					jsonReturn = util.JSON_Return{"", "Error insertando las claves del usuario en la base de datos"}
				} else {
					//INSERTAMOS LOS ROLES DEL USUARIO
					inserted, err := models.InsertUserAndRole(userId, user.Roles)
					if err == nil && inserted == true {
						jsonReturn = util.JSON_Return{"OK", ""}
						//INSERTAR EL USUARIO EN LAS CLINICA
						clinicaId, _ := strconv.Atoi(user.EnfermeroClinica)
						if clinicaId != -1 {
							result, err := models.InsertarUserClinica(clinicaId, userId, models.Rol_enfermero.Id)
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

func addMedicoAdminHandler(w http.ResponseWriter, req *http.Request) {
	var user util.User_JSON
	json.NewDecoder(req.Body).Decode(&user)
	jsonReturn := util.JSON_Return{"", ""}
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.VerifyAdmin(user.UserToken.UserId, user.MedicoClinica, user.UserToken.Token)
	if authorized == true {
		util.PrintLog("Insertando usuario " + user.Email)
		//INSERTAMOS EL USUARIO
		userId, err := models.InsertUser(user)
		user.Id = userId
		if err == nil {
			//Insertamos DNI hasheado
			_, err = models.InsertUserDniHash(userId, user.IdentificacionHash)
			if err != nil {
				util.PrintErrorLog(err)
				jsonReturn = util.JSON_Return{"", "El documento de identificación ya existe en la base de datos"}
				models.DeleteUser(user.Id)
			} else {
				//INSERTAMOS CLAVES RSA
				_, err := models.InsertUserPairKeys(userId, user.PairKeys)
				if err != nil {
					jsonReturn = util.JSON_Return{"", "Error insertando las claves del usuario en la base de datos"}
				} else {
					//INSERTAMOS LOS ROLES DEL USUARIO
					inserted, err := models.InsertUserAndRole(userId, user.Roles)
					if err == nil && inserted == true {
						jsonReturn = util.JSON_Return{"OK", ""}
						//INSERTAR EL USUARIO EN LAS CLINICA
						clinicaId, _ := strconv.Atoi(user.MedicoClinica)
						if clinicaId != -1 {
							result, err := models.InsertarUserClinica(clinicaId, userId, models.Rol_medico.Id)
							//Insertamos nombre medico
							models.InsertNombresEmpleado(user)
							if err != nil || result == false {
								jsonReturn = util.JSON_Return{"", "Error insertando el usuario en la clínica"}
							}
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
