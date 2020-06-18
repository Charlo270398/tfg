package routes

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	util "../utils"
	"github.com/gorilla/mux"
)

type Page struct {
	Title     string
	Body      string
	UserRoles []int
}

func prepareUserToken(req *http.Request) util.UserToken {
	var userToken util.UserToken
	session, _ := store.Get(req, "userSession")
	userId, _ := session.Values["userId"].(string)
	token, _ := session.Values["userToken"].(string)
	userToken.Token = token
	userToken.UserId = userId
	return userToken
}

func proveToken(req *http.Request) bool {
	session, _ := store.Get(req, "userSession")
	userId, _ := session.Values["userId"].(string)
	token, _ := session.Values["userToken"].(string)
	var responseJSON util.JSON_Return
	locJson, err := json.Marshal(util.UserToken_JSON{Token: token, UserId: userId})
	if err != nil {
		return false
	}
	//Certificado
	client := GetTLSClient()

	//Request al servidor para registrar usuario
	response, err := client.Post(SERVER_URL+"/provetoken", "application/json", bytes.NewBuffer(locJson))
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&responseJSON)
		if err != nil {
			return false
		}
		if responseJSON.Result == "OK" {
			return true
		}
	}
	return false
}

func GetTLSClient() *http.Client {
	//Certificado
	// Read the key pair to create certificate
	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		log.Fatal(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create a HTTPS client and supply the created CA pool and certificate
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{cert},
			},
		},
	}
	return client
}

func getUserPairKeys(userId string) util.PairKeys {
	//Certificado
	client := GetTLSClient()
	var user util.User_JSON
	//Recuperamos la clave publica del medico
	response, _ := client.Get(SERVER_URL + "/user/pairkeys?userId=" + userId)
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&user)
		if err != nil {
			return user.PairKeys
		}
	} else {
		return user.PairKeys
	}
	return user.PairKeys
}

func getUserPublicKeyByHistorialId(historialId string) util.PairKeys {
	//Certificado
	client := GetTLSClient()
	var user util.User_JSON
	//Recuperamos la clave publica del usuario al que pertenece ese historial
	response, _ := client.Get(SERVER_URL + "/user/pairkeysByHistorialId?historialId=" + historialId)
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&user)
		if err != nil {
			return user.PairKeys
		}
	} else {
		return user.PairKeys
	}
	return user.PairKeys
}

func getUserMasterPairKeys(userId string) util.PairKeys {
	//Certificado
	client := GetTLSClient()
	var user util.User_JSON
	//Recuperamos la clave publica del medico
	response, _ := client.Get(SERVER_URL + "/user/masterPairkeys?userId=" + userId)
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&user)
		if err != nil {
			return user.MasterPairKeys
		}
	} else {
		return user.MasterPairKeys
	}
	return user.MasterPairKeys
}

func getPublicMasterKey() util.PairKeys {
	//Certificado
	client := GetTLSClient()
	var user util.User_JSON
	//Recuperamos la clave publica del medico
	response, _ := client.Get(SERVER_URL + "/user/publicMasterKey")
	if response != nil {
		err := json.NewDecoder(response.Body).Decode(&user)
		if err != nil {
			return user.MasterPairKeys
		}
	} else {
		return user.MasterPairKeys
	}
	return user.MasterPairKeys
}

//URL DEL SERVIDOR AL QUE NOS CONECTAMOS
const SERVER_URL = "https://localhost:5001"

func LoadRouter() {
	router := mux.NewRouter()

	//STATIC RESOURCES
	http.Handle("/", router)
	router.
		PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("."+"/public/"))))
	router.HandleFunc("/forbidden", forbiddenHandler).Methods("GET")

	//ROLES
	router.HandleFunc("/rol/list", rolesListHandler).Methods("GET")
	router.HandleFunc("/rol/list/user", rolesListByUserHandler).Methods("GET")

	//TAGS
	router.HandleFunc("/tag/list", tagsListHandler).Methods("GET")

	//CLINICA
	router.HandleFunc("/clinica/add", addClinicaFormGadminHandler).Methods("GET")
	router.HandleFunc("/clinica/add", addClinicaGadminHandler).Methods("POST")
	router.HandleFunc("/clinica/especialidad/add", addClinicaEspecialidadFormGadminHandler).Methods("GET")
	router.HandleFunc("/clinica/especialidad/add", addClinicaEspecialidadFormGadminHandler).Methods("POST")
	router.HandleFunc("/clinica/list", getClinicaListGadminHandler).Methods("GET")
	router.HandleFunc("/clinica/especialidad/list", getClinicaEspecialidadListHandler).Methods("GET")
	router.HandleFunc("/clinica/especialidad/doctor/list", getMedicosClinicaByEspecialidadListHandler).Methods("GET")

	//ESPECIALIDAD
	router.HandleFunc("/especialidad/add", addEspecialidadFormGadminHandler).Methods("GET")
	router.HandleFunc("/especialidad/add", addEspecialidadGadminHandler).Methods("POST")
	router.HandleFunc("/especialidad/list", getEspecialidadListGadminHandler).Methods("GET")

	//LOGIN
	router.HandleFunc("/login", loginIndexHandler).Methods("GET")
	router.HandleFunc("/login", loginUserHandler).Methods("POST")
	router.HandleFunc("/register", registerIndexHandler).Methods("GET")
	router.HandleFunc("/register", registerUserHandler).Methods("POST")
	router.HandleFunc("/logout", logoutUserHandler).Methods("GET")
	router.HandleFunc("/logout", logoutUserHandler).Methods("POST")

	//HOME
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/home", homeHandler).Methods("GET")

	//USER(GLOBAL)
	router.HandleFunc("/user/menu", menuUserHandler).Methods("GET")
	router.HandleFunc("/user/{userId}/delete", deleteUserHandler).Methods("DELETE")
	router.HandleFunc("/user/menu/edit", menuEditUserFormHandler).Methods("GET")
	router.HandleFunc("/user/menu/edit", menuEditUserHandler).Methods("POST")

	//SOLICITAR PERMISOS
	router.HandleFunc("/permisos/historial/total/solicitar", solicitarPermisoTotal).Methods("POST")
	router.HandleFunc("/permisos/historial/basico/solicitar", solicitarPermisoBasico).Methods("POST")
	router.HandleFunc("/permisos/entrada/solicitar", solicitarPermisoEntrada).Methods("POST")
	router.HandleFunc("/permisos/analitica/solicitar", solicitarPermisoAnalitica).Methods("POST")
	router.HandleFunc("/permisos/solicitudes/denegar", denegarPermiso).Methods("POST")
	router.HandleFunc("/permisos/solicitudes/permitir", permitirPermiso).Methods("POST")

	//USER(PACIENTE)
	router.HandleFunc("/user/patient", menuPatientHandler).Methods("GET")
	router.HandleFunc("/user/patient/edit", editUserPatientHandler).Methods("GET")
	router.HandleFunc("/user/patient/historial", historialPatientHandler).Methods("GET")
	router.HandleFunc("/user/patient/historial/entrada", historialEntradaPatientHandler).Methods("GET")
	router.HandleFunc("/user/patient/historial/analitica", historialAnaliticaPatientHandler).Methods("GET")
	router.HandleFunc("/user/patient/citas", patientCitaListHandler).Methods("GET")
	router.HandleFunc("/user/patient/citas/add", addPatientCitaFormHandler).Methods("GET")
	router.HandleFunc("/user/patient/citas/add", addCitaPacienteHandler).Methods("POST")
	router.HandleFunc("/user/patient/historial/share", patientAutorizationsHandler).Methods("GET")

	//USER(ENFERMERO)
	router.HandleFunc("/user/nurse", menuEnfermeroHandler).Methods("GET")
	router.HandleFunc("/user/nurse/historial/solicitar", solicitarHistorialEnfermeroHandler).Methods("GET")

	//USER(MEDICO)
	router.HandleFunc("/user/doctor", menuMedicoHandler).Methods("GET")
	router.HandleFunc("/user/doctor/historial/solicitar", solicitarHistorialMedicoFormHandler).Methods("GET")
	router.HandleFunc("/user/doctor/historial/solicitar", solicitarHistorialMedicoHandler).Methods("POST")
	router.HandleFunc("/user/doctor/historial/addEntrada", addEntradaHistorialFormMedicoHandler).Methods("GET")
	router.HandleFunc("/user/doctor/historial/addEntrada", addEntradaHistorialMedicoHandler).Methods("POST")
	router.HandleFunc("/user/doctor/historial/addAnalitica", addAnaliticaHistorialFormMedicoHandler).Methods("GET")
	router.HandleFunc("/user/doctor/historial/addAnalitica", addAnaliticaHistorialMedicoHandler).Methods("POST")
	router.HandleFunc("/user/doctor/historial/list", getListHistorialMedicoHandler).Methods("GET")
	router.HandleFunc("/user/medico/historial/entrada", getEntradaHistorialMedicoHandler).Methods("GET")
	router.HandleFunc("/user/medico/historial/analitica", getAnaliticaHistorialMedicoHandler).Methods("GET")
	router.HandleFunc("/user/doctor/disponible/dia", getMedicoDiasDisponiblesHandler).Methods("GET")
	router.HandleFunc("/user/doctor/disponible/hora", getMedicoHorasDiaDisponiblesHandler).Methods("GET")
	router.HandleFunc("/user/doctor/citas", getCitaFormMedicoHandler).Methods("GET")
	router.HandleFunc("/user/doctor/citas/list", medicoCitaListHandler).Methods("GET")
	router.HandleFunc("/user/doctor/citas/addEntrada", addEntradaHistorialConsultaMedicoHandler).Methods("POST")
	router.HandleFunc("/user/doctor/research/analiticas", getInvestigacionAnaliticasMedicoFormHandler).Methods("GET")

	//USER(ADMIN-CLINICA)
	router.HandleFunc("/user/admin", menuAdminHandler).Methods("GET")
	router.HandleFunc("/user/admin/nurse/add", adminAddEnfermeroFormHandler).Methods("GET")
	router.HandleFunc("/user/admin/nurse/add", addEnfermeroAdminHandler).Methods("POST")
	router.HandleFunc("/user/admin/doctor/add", adminAddMedicoFormHandler).Methods("GET")
	router.HandleFunc("/user/admin/doctor/add", addMedicoAdminHandler).Methods("POST")

	//USER(ADMIN-GLOBAL)
	router.HandleFunc("/user/adminG", menuAdminGHandler).Methods("GET")
	router.HandleFunc("/user/adminG/userList", getUserListAdminGHandler).Methods("GET")
	router.HandleFunc("/user/adminG/userList/add", addUserFormGadminHandler).Methods("GET")
	router.HandleFunc("/user/adminG/userList/add", addUserGadminHandler).Methods("POST")

	//USER(EMERGENCIAS)
	router.HandleFunc("/user/emergency", menuEmergenciasHandler).Methods("GET")
	router.HandleFunc("/user/emergency/historial", GetHistorialEmergenciasHandler).Methods("POST")
	router.HandleFunc("/user/emergency/historial/entrada", GetEntradaEmergenciasHandler).Methods("GET")
	router.HandleFunc("/user/emergency/historial/analitica", GetAnaliticaEmergenciasHandler).Methods("GET")
	router.HandleFunc("/user/emergency/historial/addEntrada", AddEntradaEmergenciasFormHandler).Methods("GET")
	router.HandleFunc("/user/emergency/historial/addEntrada", AddEntradaEmergenciasHandler).Methods("POST")
	router.HandleFunc("/user/emergency/historial/addAnalitica", AddAnaliticaEmergenciasFormHandler).Methods("GET")
	router.HandleFunc("/user/emergency/historial/addAnalitica", AddAnaliticaEmergenciasHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	fmt.Println("Servidor cliente escuchando en el puerto ", port)
	err := http.ListenAndServeTLS(":5000", "cert.pem", "key.pem", nil)
	if err != nil {
		util.PrintErrorLog(err)
	}
}
