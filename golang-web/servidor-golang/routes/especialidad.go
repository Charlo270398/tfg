package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "../models"
	util "../utils"
)

//POST
func addEspecialidadHandler(w http.ResponseWriter, req *http.Request) {
	var especialidad util.Especialidad_JSON
	json.NewDecoder(req.Body).Decode(&especialidad)
	if especialidad.Nombre == "" {
		http.Error(w, "Nombre incompatible", http.StatusInternalServerError)
		return
	}

	jsonReturn := util.JSON_Return{"", ""}
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(especialidad.UserToken.UserId, especialidad.UserToken.Token, models.Rol_administradorG.Id)
	if authorized == true {
		util.PrintLog("Insertando especialidad " + especialidad.Nombre)
		_, err := models.InsertEspecialidad(especialidad)
		if err == nil {
			util.PrintLog("Especialidad " + especialidad.Nombre + " INSERTADA")
			jsonReturn = util.JSON_Return{"OK", ""}
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

//GET

func getEspecialidadListHandler(w http.ResponseWriter, req *http.Request) {
	var especialidadList []util.Especialidad
	especialidadList, err := models.GetEspecialidadList()
	js, err := json.Marshal(especialidadList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getEspecialidadPaginationHandler(w http.ResponseWriter, req *http.Request) {
	page, ok := req.URL.Query()["page"]
	var especialidadListReturn util.Especialidad_JSON_Pagination
	var especialidadList []util.Especialidad
	if !ok || len(page[0]) < 1 {
		especialidadList = models.GetEspecialidadPagination(0) //Devolvemos primera pagina
	} else {
		pageNumber, err := strconv.Atoi(page[0])
		especialidadListReturn.Page = pageNumber
		especialidadListReturn.BeforePage = pageNumber - 1
		especialidadListReturn.NextPage = pageNumber + 1
		if err != nil {
			especialidadList = models.GetEspecialidadPagination(0) //Devolvemos primera pagina
		} else {
			especialidadList = models.GetEspecialidadPagination(pageNumber)
		}
	}
	especialidadListReturn.EspecialidadList = especialidadList

	js, err := json.Marshal(especialidadListReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}
