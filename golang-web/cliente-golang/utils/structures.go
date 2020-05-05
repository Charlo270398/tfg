package utils

import (
	"time"
)

//User

type User struct {
	Id             int
	Identificacion string
	Nombre         string
	Apellidos      string
	Email          string
	Password       string
	CreatedAt      time.Time
}

type User_JSON_AddUsers struct {
	Id                 int       `json:"id"`
	Identificacion     string    `json:"identificacion"`
	IdentificacionHash string    `json:"identificacionHash"`
	Nombre             string    `json:"nombre"`
	Apellidos          string    `json:"apellidos"`
	Email              string    `json:"email"`
	CreatedAt          time.Time `json:"createdAt"`
	Password           string    `json:"password"`
	Roles              []int     `json:"roles"`
	EnfermeroClinica   string    `json:"enfermeroClinica"`
	MedicoClinica      string    `json:"medicoClinica"`
	AdminClinica       string    `json:"adminClinica"`
	MedicoEspecialidad string    `json:"medicoEspecialidad"`
	UserToken          UserToken `json:"userToken"`
	PairKeys           PairKeys  `json:"pairKeys"`
	Sexo               string    `json:"sexo"`
	Alergias           string    `json:"alergias"`
	Clave              string    `json:"clave"`
	NombreDoctor       string    `json:nombreDoctor`
}

type User_JSON struct {
	Id                 int       `json:"id"`
	Identificacion     string    `json:"identificacion"`
	IdentificacionHash string    `json:"identificacionHash"`
	Nombre             string    `json:"nombre"`
	Apellidos          string    `json:"apellidos"`
	Email              string    `json:"email"`
	CreatedAt          time.Time `json:"createdAt"`
	Password           []byte    `json:"password"`
	Roles              []int     `json:"roles"`
	EnfermeroClinica   string    `json:"enfermeroClinica"`
	MedicoClinica      string    `json:"medicoClinica"`
	AdminClinica       string    `json:"adminClinica"`
	MedicoEspecialidad string    `json:"medicoEspecialidad"`
	UserToken          UserToken `json:"userToken"`
	PairKeys           PairKeys  `json:"pairKeys"`
	Sexo               string    `json:"sexo"`
	Alergias           string    `json:"alergias"`
	Clave              string    `json:"clave"`
	NombreDoctor       string    `json:nombreDoctor`
}

type JSON_Credentials_CLIENTE struct {
	Password       string `json:"password"`
	Identificacion string `json:"identificacion"`
}

type JSON_Credentials_SERVIDOR struct {
	Password       []byte `json:"password"`
	Identificacion string `json:"identificacion"`
}

type JSON_user_SERVIDOR struct {
	Identificacion string `json:"identificacion"`
	Nombre         string `json:"nombre"`
	Apellidos      string `json:"apellidos"`
	Email          string `json:"email"`
	Password       []byte `json:"password"`
	Roles          []int  `json:"roles"`
}

type User_id_JSON struct {
	Id        int       `json:"user_id"`
	UserToken UserToken `json:"userToken"`
}

type Params_argon2 struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

type PaginationPage_JSON struct {
	Page int `json:"Page"`
}

type UserList_JSON_Pagination struct {
	Page       int    `json:"Page"`
	NextPage   int    `json:"NextPage"`
	BeforePage int    `json:"BeforePage"`
	UserList   []User `json:"UserList"`
}

type UserList_Page struct {
	Title      string
	Body       string
	Page       int
	NextPage   int
	BeforePage int
	UserList   []User
}

//Rol

type Rol struct {
	Id          int
	Nombre      string
	Descripcion string
}

type Rol_json struct {
	Id          int    `json:"id"`
	Nombre      string `json:"nombre"`
	Descripcion string `json:"descripcion"`
}

type Roles_List_json struct {
	Roles []Rol `json:"roles"`
}

//Clinica

type Clinica struct {
	Id                    int
	Nombre                string
	Direccion             string
	Telefono              string
	NumeroEnfermeros      int
	NumeroMedicos         int
	NumeroAdministradores int
}

type Clinica_JSON struct {
	Id        int       `json:"id"`
	Nombre    string    `json:"nombre"`
	Direccion string    `json:"direccion"`
	Telefono  string    `json:"telefono"`
	UserToken UserToken `json:"userToken"`
}

type Clinica_JSON_Pagination struct {
	Page        int       `json:"Page"`
	NextPage    int       `json:"NextPage"`
	BeforePage  int       `json:"BeforePage"`
	ClinicaList []Clinica `json:"ClinicaList"`
}

type ClinicaList_Page struct {
	Title       string
	Body        string
	Page        int
	NextPage    int
	BeforePage  int
	ClinicaList []Clinica
}

//Especialidad

type Especialidad struct {
	Id     int
	Nombre string
}

type Especialidad_JSON struct {
	Id        int       `json:"id"`
	Nombre    string    `json:"nombre"`
	UserToken UserToken `json:"userToken"`
}

type Especialidad_JSON_Pagination struct {
	Page             int            `json:"Page"`
	NextPage         int            `json:"NextPage"`
	BeforePage       int            `json:"BeforePage"`
	EspecialidadList []Especialidad `json:"EspecialidadList"`
}

type EspecialidadList_Page struct {
	Title            string
	Body             string
	Page             int
	NextPage         int
	BeforePage       int
	EspecialidadList []Especialidad
}

//Historial

type Historial_JSON struct {
	Id             int       `json:"id"`
	PacienteId     int       `json:"pacienteId"`
	MedicoId       int       `json:"medicoId"`
	Sexo           string    `json:"sexo"`
	Alergias       string    `json:"alergias"`
	NombrePaciente string    `json:"nombrePaciente"`
	Clave          string    `json:"clave"`
	UserToken      UserToken `json:"userToken"`
}

type SolicitarHistorial_JSON struct {
	UserDNI   string    `json:"userDNI"`
	UserToken UserToken `json:"userToken"`
}

type EntradaHistorial_JSON struct {
	Id                int       `json:"id"`
	PacienteId        int       `json:"pacienteId"`
	CitaId            int       `json:"citaId"`
	MotivoConsulta    string    `json:"motivoConsulta"`
	JuicioDiagnostico string    `json:"juicioDiagnostico"`
	Tipo              string    `json:"tipo"`
	Clave             string    `json:"clave"`
	UserToken         UserToken `json:"userToken"`
}

//Citas

type Cita struct {
	Id             int
	PacienteId     int
	PacienteNombre string
	MedicoId       int
	MedicoNombre   string
	Hora           int
	Dia            int
	Mes            int
	Anyo           int
	Tipo           string
	Fecha          time.Time
	UserToken      UserToken
}

type CitaJSON struct {
	Id             int       `json:"id"`
	PacienteId     string    `json:"pacienteId"`
	PacienteNombre string    `json:"pacienteNombre"`
	MedicoId       string    `json:"medicoId"`
	MedicoNombre   string    `json:"medicoNombre"`
	Hora           int       `json:"hora"`
	Dia            int       `json:"dia"`
	Mes            int       `json:"mes"`
	Anyo           int       `json:"anyo"`
	Tipo           string    `json:"tipo"`
	Fecha          time.Time `json:"fecha"`
	FechaString    string    `json:"fechaString"`
	UserToken      UserToken `json:"userToken"`
}

//SEGURIDAD
//Token

type UserToken struct {
	UserId string
	Token  string
}

type UserToken_JSON struct {
	UserId string `json:"UserId"`
	Token  string `json:"Token"`
}

//PairKeys
type PairKeys struct {
	PublicKey  []byte
	PrivateKey []byte
}

//Response

type JSON_Login_Return struct {
	UserId    string
	Nombre    string
	Apellidos string
	Email     string
	Error     string
	Token     string
	PairKeys  PairKeys
	Clave     string `json:"clave"`
}

type JSON_Return struct {
	Result string
	Error  string
}

//NECESARIO EL TOKEN

//Admin
type JSON_Admin_Menu struct {
	Clinica   Clinica
	UserToken UserToken
	Error     string
}

//PAGES

type PageMenuAdmin struct {
	Title   string
	Body    string
	Clinica Clinica
}

type PageMenuMedico struct {
	Title      string
	Body       string
	CitaActual int
}

type PageAdminAddMedico struct {
	Title          string
	Body           string
	Clinica        Clinica
	Especialidades []Especialidad_JSON
}

type CitaPage struct {
	Title          string
	Body           string
	Clinicas       []Clinica_JSON
	Especialidades []Especialidad_JSON
}

type ConsultaPage struct {
	Title string
	Body  string
	Cita  CitaJSON
}

type CitaListPage struct {
	Title string
	Body  string
	Citas []CitaJSON
}

type HistorialListPage struct {
	Title       string
	Body        string
	Historiales []Historial_JSON
}

type CambiarDatosPage struct {
	Title     string
	Body      string
	Nombre    string
	Apellido1 string
	Apellido2 string
	Email     string
}
