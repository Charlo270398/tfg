package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	models "../models"
	util "../utils"
)

//GET

func getClinicaListHandler(w http.ResponseWriter, req *http.Request) {
	var clinicaList []util.Clinica
	clinicaList, err := models.GetClinicaList()
	js, err := json.Marshal(clinicaList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getEspecialidadesListClinicaHandler(w http.ResponseWriter, req *http.Request) {
	clinicaId, ok := req.URL.Query()["clinicaId"]
	if !ok || len(clinicaId[0]) < 1 {
		http.Error(w, "No existe el parámetro clinicaId", http.StatusInternalServerError)
		return
	} else {
		var especialidadList []util.Especialidad
		especialidadList, err := models.GetEspecialidadesClinica(clinicaId[0])
		js, err := json.Marshal(especialidadList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func getMedicosByEspecialidadListClinicaHandler(w http.ResponseWriter, req *http.Request) {
	clinicaId, ok := req.URL.Query()["clinicaId"]
	if !ok || len(clinicaId[0]) < 1 {
		http.Error(w, "No existe el parámetro clinicaId", http.StatusInternalServerError)
		return
	}
	especialidadId, ok := req.URL.Query()["especialidadId"]
	if !ok || len(especialidadId[0]) < 1 {
		http.Error(w, "No existe el parámetro especialidadId", http.StatusInternalServerError)
		return
	}

	var medicosList []util.User_JSON
	medicosList, err := models.GetMedicosClinicaByEspecialidad(clinicaId[0], especialidadId[0])
	js, err := json.Marshal(medicosList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func getClinicaPaginationHandler(w http.ResponseWriter, req *http.Request) {
	page, ok := req.URL.Query()["page"]
	var clinicaListReturn util.Clinica_JSON_Pagination
	var clinicaList []util.Clinica
	if !ok || len(page[0]) < 1 {
		clinicaList = models.GetClinicaPagination(0) //Devolvemos primera pagina
	} else {
		pageNumber, err := strconv.Atoi(page[0])
		clinicaListReturn.Page = pageNumber
		clinicaListReturn.BeforePage = pageNumber - 1
		clinicaListReturn.NextPage = pageNumber + 1
		if err != nil {
			clinicaList = models.GetClinicaPagination(0) //Devolvemos primera pagina
		} else {
			clinicaList = models.GetClinicaPagination(pageNumber)
		}
	}
	clinicaListReturn.ClinicaList = clinicaList

	js, err := json.Marshal(clinicaListReturn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

//POST
func addClinicaHandler(w http.ResponseWriter, req *http.Request) {
	var clinica util.Clinica_JSON
	json.NewDecoder(req.Body).Decode(&clinica)
	if clinica.Nombre == "" {
		http.Error(w, "Nombre incompatible", http.StatusInternalServerError)
		return
	}

	jsonReturn := util.JSON_Return{"", ""}
	//Comprobamos que el usuario esta autorizado y el token es correcto
	authorized, _ := models.GetAuthorizationbyUserId(clinica.UserToken.UserId, clinica.UserToken.Token, models.Rol_administradorG.Id)
	if authorized == true {
		util.PrintLog("Insertando especialidad " + clinica.Nombre)
		_, err := models.InsertClinica(clinica)
		if err == nil {
			util.PrintLog("Clinica " + clinica.Nombre + " INSERTADA")
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
