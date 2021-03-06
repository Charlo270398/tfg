package routes

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	util "../utils"
)

type Page struct {
	Title string
	Body  string
}

func LoadRouter(port string) {

	//STATIC RESOURCES
	http.HandleFunc("/inicio", getInicioHandler)

	//LOGIN
	http.HandleFunc("/login", loginUserHandler)
	http.HandleFunc("/register", registerUserHandler)
	http.HandleFunc("/provetoken", proveUserTokenHandler)

	//CLINICA
	http.HandleFunc("/clinica/add", addClinicaHandler)
	http.HandleFunc("/clinica/list", getClinicaPaginationHandler)
	http.HandleFunc("/clinica/list/query", getClinicaListHandler)

	//ESPECIALIDAD
	http.HandleFunc("/especialidad/add", addEspecialidadHandler)
	http.HandleFunc("/especialidad/list", getEspecialidadPaginationHandler)
	http.HandleFunc("/especialidad/list/query", getEspecialidadListHandler)

	//ROLES
	http.HandleFunc("/rol/list", getRolesListHandler)
	http.HandleFunc("/rol/list/user", getRolesByUserHandler)

	//USER(GLOBAL)
	http.HandleFunc("/user/menu/edit", menuUserEditHandler)
	http.HandleFunc("/user/delete", deleteUserHandler)
	http.HandleFunc("/user", getUserHandler)

	//USER(ADMIN)
	http.HandleFunc("/user/admin", getAdminMenuDataHandler)

	//USER(ADMING)
	http.HandleFunc("/user/adminG/userList/add", addUserHandler)
	http.HandleFunc("/user/adminG/userList", getUsersPaginationHandler)
	if port == "" {
		port = "5001"
	}
	// Create a CA certificate pool and add cert.pem to it
	caCert, err := ioutil.ReadFile("cert.pem")
	if err != nil {
		util.PrintErrorLog(err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Create the TLS Config with the CA pool and enable Client certificate validation
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()

	// Create a Server instance to listen on port 8443 with the TLS config
	server := &http.Server{
		Addr:      ":" + port,
		TLSConfig: tlsConfig,
	}
	fmt.Println("Servidor escuchando en el puerto ", port)

	//log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
	err = server.ListenAndServeTLS("cert.pem", "key.pem")
	if err != nil {
		util.PrintErrorLog(err)
	}
}
