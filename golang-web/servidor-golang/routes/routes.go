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
	http.HandleFunc("/clinica/especialidad/list", getEspecialidadesListClinicaHandler)
	http.HandleFunc("/clinica/especialidad/doctor/list", getMedicosByEspecialidadListClinicaHandler)

	//ESPECIALIDAD
	http.HandleFunc("/especialidad/add", addEspecialidadHandler)
	http.HandleFunc("/especialidad/list", getEspecialidadPaginationHandler)
	http.HandleFunc("/especialidad/list/query", getEspecialidadListHandler)

	//ROLES
	http.HandleFunc("/rol/list", getRolesListHandler)
	http.HandleFunc("/rol/list/user", getRolesByUserHandler)

	//TAGS
	http.HandleFunc("/tag/list", getTagsListHandler)

	//USER(GLOBAL)
	http.HandleFunc("/user", getUserHandler)
	http.HandleFunc("/user/menu/edit", menuUserEditHandler)
	http.HandleFunc("/user/delete", deleteUserHandler)
	http.HandleFunc("/user/pairkeys", getUserPairKeysHandler)
	http.HandleFunc("/user/masterPairkeys", getUserMasterPairKeysHandler)
	http.HandleFunc("/user/publicMasterKey", getPublicMasterKeyHandler)
	http.HandleFunc("/user/pairkeysByHistorialId", getUserPairKeysByHistorialIdHandler)

	//SOLICITAR PERMISOS
	http.HandleFunc("/permisos/historial/total/solicitar", solicitarPermisoTotal)
	http.HandleFunc("/permisos/historial/basico/solicitar", solicitarPermisoBasico)
	http.HandleFunc("/permisos/entrada/solicitar", solicitarPermisoEntrada)
	http.HandleFunc("/permisos/analitica/solicitar", solicitarPermisoAnalitica)
	http.HandleFunc("/permisos/solicitudes/listar", listarSolicitudesPermiso)
	http.HandleFunc("/permisos/solicitudes/comprobar", comprobarSolicitudesPermiso)
	http.HandleFunc("/permisos/solicitudes/eliminar", eliminarSolicitudPermiso)
	http.HandleFunc("/permisos/historial/permitir", permitirHistorial)
	http.HandleFunc("/permisos/entrada/permitir", permitirEntrada)
	http.HandleFunc("/permisos/analitica/permitir", permitirAnalitica)

	//USER(PACIENTE)
	http.HandleFunc("/user/patient/citas/add", PacienteInsertCita)
	http.HandleFunc("/user/patient/citas/list", PacienteGetCitasFuturasList)
	http.HandleFunc("/user/patient/historial", GetHistorialPaciente)
	http.HandleFunc("/user/patient/historial/share", ShareHistorialPaciente)
	http.HandleFunc("/user/patient/historial/entrada", PacienteGetEntradaHandler)
	http.HandleFunc("/user/patient/historial/analitica", PacienteGetAnaliticaHandler)

	//USER(MEDICO)
	http.HandleFunc("/user/doctor/historial", GetHistorialCompartido)
	http.HandleFunc("/user/doctor/historial/entrada", GetEntradaHistorialCompartido)
	http.HandleFunc("/user/doctor/historial/analitica", GetAnaliticaHistorialCompartido)
	http.HandleFunc("/user/doctor/historial/solicitar", MedicoSolicitarHistorialHandler)
	http.HandleFunc("/user/doctor/historial/list", MedicoGetHistorialesCompartidos)
	http.HandleFunc("/user/doctor/disponible/dia", MedicoDiasDisponiblesHandler)
	http.HandleFunc("/user/doctor/disponible/hora", MedicoHorasDiaDisponiblesHandler)
	http.HandleFunc("/user/doctor/citas/list", MedicoGetCitasFuturasList)
	http.HandleFunc("/user/doctor/citas/actual", MedicoGetCitaActual)
	http.HandleFunc("/user/doctor/citas", MedicoGetCita)
	http.HandleFunc("/user/doctor/citas/addEntrada", MedicoAddEntradaHistorialConsulta)
	http.HandleFunc("/user/doctor/citas/addEntradaCompartida", MedicoAddEntradaHistorialCompartidaConsulta)

	//USER(ADMIN)
	http.HandleFunc("/user/admin", getAdminMenuDataHandler)
	http.HandleFunc("/user/admin/nurse/add", addEnfermeroAdminHandler)
	http.HandleFunc("/user/admin/doctor/add", addMedicoAdminHandler)

	//USER(ADMING)
	http.HandleFunc("/user/adminG/userList/add", addUserHandler)
	http.HandleFunc("/user/adminG/userList", getUsersPaginationHandler)

	//USER(EMERGENCIAS)
	http.HandleFunc("/user/emergency/historial", GetHistorialEmergencias)
	http.HandleFunc("/user/emergency/historial/entrada", GetEntradaEmergenciasHandler)
	http.HandleFunc("/user/emergency/historial/analitica", GetAnaliticaEmergenciasHandler)
	http.HandleFunc("/user/emergency/addEntrada", AddEntradaEmergenciasHandler)
	http.HandleFunc("/user/emergency/addAnalitica", AddAnaliticaEmergenciasHandler)
	http.HandleFunc("/user/emergency/addEstadisticaAnalitica", AddEstadisticaAnaliticaEmergenciasHandler)

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
