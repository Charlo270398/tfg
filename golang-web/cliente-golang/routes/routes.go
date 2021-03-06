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

//URL DEL SERVIDOR AL QUE NOS CONECTAMOS
const SERVER_URL = "https://localhost:5001"

func LoadRouter() {
	router := mux.NewRouter()

	//STATIC RESOURCES
	http.Handle("/", router)
	router.
		PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("."+"/public/"))))

	//ROLES
	router.HandleFunc("/rol/list", rolesListHandler).Methods("GET")
	router.HandleFunc("/rol/list/user", rolesListByUserHandler).Methods("GET")

	//CLINICA
	router.HandleFunc("/clinica/add", addClinicaFormGadminHandler).Methods("GET")
	router.HandleFunc("/clinica/add", addClinicaGadminHandler).Methods("POST")
	router.HandleFunc("/clinica/especialidad/add", addClinicaEspecialidadFormGadminHandler).Methods("GET")
	router.HandleFunc("/clinica/especialidad/add", addClinicaEspecialidadFormGadminHandler).Methods("POST")
	router.HandleFunc("/clinica/list", getClinicaListGadminHandler).Methods("GET")

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

	//USER(PACIENTE)
	router.HandleFunc("/user/patient", menuPatientHandler).Methods("GET")
	router.HandleFunc("/user/patient/edit", editUserPatientHandler).Methods("GET")
	router.HandleFunc("/user/patient/historial", historialPatientHandler).Methods("GET")
	router.HandleFunc("/user/patient/citas", patientCitaListHandler).Methods("GET")
	router.HandleFunc("/user/patient/citas/add", addPatientCitaFormHandler).Methods("GET")

	//USER(ENFERMERO)
	router.HandleFunc("/user/nurse", menuEnfermeroHandler).Methods("GET")

	//USER(MEDICO)
	router.HandleFunc("/user/doctor", menuMedicoHandler).Methods("GET")

	//USER(ADMIN-CLINICA)
	router.HandleFunc("/user/admin", menuAdminHandler).Methods("GET")

	//USER(ADMIN-GLOBAL)
	router.HandleFunc("/user/adminG", menuAdminGHandler).Methods("GET")
	router.HandleFunc("/user/adminG/userList", getUserListAdminGHandler).Methods("GET")
	router.HandleFunc("/user/adminG/userList/add", addUserFormGadminHandler).Methods("GET")
	router.HandleFunc("/user/adminG/userList/add", addUserGadminHandler).Methods("POST")

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
