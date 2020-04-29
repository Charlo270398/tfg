package models

import (
	"database/sql"
	"fmt"
	"log"

	"strconv"

	util "../utils"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB //variable db común a todos

func ConnectDB() {
	var err error
	db, err = sql.Open("mysql", "golang:golang@(127.0.0.1:3306)/golang?parseTime=true")
	if err != nil {
		util.PrintErrorLog(err)
		log.Panic(err)
	}

	if err = db.Ping(); err != nil {
		util.PrintErrorLog(err)
		log.Panic(err)
	}
}
func query(query string) bool {

	// Executes the SQL query in our database. Check err to ensure there was no error.
	if _, err := db.Exec(query); err != nil {
		util.PrintErrorLog(err)
		return false
	}
	return true
}

func CreateDB() {
	ConnectDB()
	//CREATE TABLES
	query(CLINICAS_TABLE)
	query(ESPECIALIDADES_TABLE)
	query(USUARIOS_TABLE)
	query(ROLES_TABLE)
	query(USERS_CLINICAS_TABLE)
	query(USERS_ESPECIALIDADES_TABLE)
	query(USERS_ROLES_TABLE)
	query(USERS_TOKENS_TABLE)
	query(USERS_PAIRKEYS_TABLE)
	query(USERS_DNIHASHES_TABLE)
	query(MEDICOS_NOMBRES_TABLE)
	query(CITAS_TABLE)
	query(USERS_HISTORIAL_TABLE)
	query(USERS_PERMISOS_HISTORIAL_TABLE)
	query(USERS_ENTRADAS_HISTORIAL_TABLE)
	query(USERS_PERMISOS_ENTRADAS_HISTORIAL_TABLE)
	query(USERS_ANALITICAS_TABLE)
	query(USERS_PERMISOS_ANALITICAS_TABLE)
	query(TAGS_TABLE)
	query(ANALITICAS_TAGS_TABLE)
	query(ESTADISTICAS_ANALITICAS_TABLE)

	//SEEDERS
	//Roles
	query("INSERT IGNORE INTO roles (id,nombre,descripcion) VALUES (" + strconv.Itoa(Rol_paciente.Id) + ",'" + Rol_paciente.Nombre + "', '" + Rol_paciente.Descripcion + "');")
	query("INSERT IGNORE INTO roles (id,nombre,descripcion) VALUES (" + strconv.Itoa(Rol_enfermero.Id) + ",'" + Rol_enfermero.Nombre + "', '" + Rol_enfermero.Descripcion + "');")
	query("INSERT IGNORE INTO roles (id,nombre,descripcion) VALUES (" + strconv.Itoa(Rol_medico.Id) + ",'" + Rol_medico.Nombre + "', '" + Rol_medico.Descripcion + "');")
	query("INSERT IGNORE INTO roles (id,nombre,descripcion) VALUES (" + strconv.Itoa(Rol_administradorC.Id) + ",'" + Rol_administradorC.Nombre + "', '" + Rol_administradorC.Descripcion + "');")
	query("INSERT IGNORE INTO roles (id,nombre,descripcion) VALUES (" + strconv.Itoa(Rol_administradorG.Id) + ",'" + Rol_administradorG.Nombre + "', '" + Rol_administradorG.Descripcion + "');")

	//Tags

	fmt.Println("Database OK")
}

var CLINICAS_TABLE string = `
CREATE TABLE IF NOT EXISTS clinicas (
	id INT AUTO_INCREMENT,
	nombre VARCHAR(20) UNIQUE,
	direccion VARCHAR(50),
	telefono VARCHAR(16),
	PRIMARY KEY (id)
);`

var ESPECIALIDADES_TABLE string = `
CREATE TABLE IF NOT EXISTS especialidades (
	id INT AUTO_INCREMENT,
	nombre VARCHAR(30) UNIQUE,
	PRIMARY KEY (id)
);`

var USUARIOS_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios (
	id INT AUTO_INCREMENT,
	dni VARCHAR(36) UNIQUE,
	nombre VARCHAR(100) NOT NULL,
	apellidos VARCHAR(150) NOT NULL,
	email VARCHAR(100) UNIQUE,
	password VARCHAR(100) NOT NULL,
	created_at DATETIME,
	clave VARCHAR(344) NOT NULL,
	PRIMARY KEY (id)
);`

var ROLES_TABLE string = `
CREATE TABLE IF NOT EXISTS roles (
	id INT AUTO_INCREMENT,
	nombre VARCHAR(20) UNIQUE,
	descripcion VARCHAR(50),
	PRIMARY KEY (id)
);`

var TAGS_TABLE string = `
CREATE TABLE IF NOT EXISTS tags (
	id INT AUTO_INCREMENT,
	nombre VARCHAR(30) UNIQUE,
	PRIMARY KEY (id)
);`

//Relaciones

var USERS_CLINICAS_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_clinicas (
	usuario_id INT,
	clinica_id INT,
	rol_id INT,
	PRIMARY KEY (usuario_id, clinica_id, rol_id),
	FOREIGN KEY(usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
	FOREIGN KEY(clinica_id) REFERENCES clinicas(id) ON DELETE CASCADE,
	FOREIGN KEY(rol_id) REFERENCES roles(id) ON DELETE CASCADE
);`

var USERS_ESPECIALIDADES_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_especialidades (
	usuario_id INT,
	especialidad_id INT,
	PRIMARY KEY(usuario_id, especialidad_id),
	FOREIGN KEY(usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
	FOREIGN KEY(especialidad_id) REFERENCES especialidades(id) ON DELETE CASCADE
);`

var USERS_ROLES_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_roles (
	usuario_id INT,
	rol_id INT,
	PRIMARY KEY (usuario_id, rol_id),
	FOREIGN KEY(usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
	FOREIGN KEY(rol_id) REFERENCES roles(id) ON DELETE CASCADE
);`

var USERS_TOKENS_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_tokens (
	id INT AUTO_INCREMENT,
	usuario_id INT UNIQUE,
	token VARCHAR(156),
	fecha_expiracion DATETIME,
	PRIMARY KEY (id),
	FOREIGN KEY(usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE
);`

var USERS_PAIRKEYS_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_pairkeys (
	id INT AUTO_INCREMENT,
	usuario_id INT UNIQUE,
	public_key BLOB,
	private_key BLOB,
	PRIMARY KEY (id),
	FOREIGN KEY(usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE
);`

var USERS_DNIHASHES_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_dnihashes(
	usuario_id INT,
	dni_hash VARCHAR(64),
	FOREIGN KEY(usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
	PRIMARY KEY (usuario_id,dni_hash)
);`

var MEDICOS_NOMBRES_TABLE string = `
CREATE TABLE IF NOT EXISTS medicos_nombres (
	usuario_id INT,
	nombreDoctor VARCHAR(150) NOT NULL,
	FOREIGN KEY(usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE,
	PRIMARY KEY (usuario_id)
);`

var CITAS_TABLE string = `
CREATE TABLE IF NOT EXISTS citas (
	id INT AUTO_INCREMENT,
	medico_id INT,
	paciente_id INT,
	anyo INT,
	mes INT,
	dia INT,
	hora INT,
	tipo VARCHAR(30) NOT NULL, 
	FOREIGN KEY(medico_id) REFERENCES usuarios(id) ON DELETE CASCADE,
	FOREIGN KEY(paciente_id) REFERENCES usuarios(id) ON DELETE CASCADE,
	UNIQUE (medico_id, anyo, mes, dia, hora),
	PRIMARY KEY (id)
);`

var USERS_HISTORIAL_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_historial (
	id INT AUTO_INCREMENT,
	usuarios_id INT,
	PRIMARY KEY (id),
	FOREIGN KEY(usuarios_id) REFERENCES usuarios(id) ON DELETE CASCADE
);`

var USERS_ENTRADAS_HISTORIAL_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_entradas_historial (
	id INT AUTO_INCREMENT,
	historial_id INT,
	motivo_consulta varchar(500), 
	juicio_diagnostico varchar(500),
	clave VARCHAR(344) NOT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY(historial_id) REFERENCES usuarios_historial(id) ON DELETE CASCADE
);`

var USERS_PERMISOS_HISTORIAL_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_permisos_historial (
	historial_id INT,
	medico_id INT,
	fecha_expiracion DATETIME,
	PRIMARY KEY (historial_id, medico_id),
	FOREIGN KEY(historial_id) REFERENCES usuarios_historial(id) ON DELETE CASCADE,
	FOREIGN KEY(medico_id) REFERENCES usuarios(id) ON DELETE CASCADE
);`

var USERS_PERMISOS_ENTRADAS_HISTORIAL_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_permisos_entradas_historial (
	entrada_id INT,
	medico_id INT,
	clave VARCHAR(344) NOT NULL,
	fecha_expiracion DATETIME,
	PRIMARY KEY (entrada_id, medico_id),
	FOREIGN KEY(entrada_id) REFERENCES usuarios_entradas_historial(id) ON DELETE CASCADE,
	FOREIGN KEY(medico_id) REFERENCES usuarios(id) ON DELETE CASCADE
);`

var USERS_ANALITICAS_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_analiticas (
	id INT AUTO_INCREMENT,
	usuario_id INT,
	leucocitos VARCHAR(100),
	hematies VARCHAR(100),
	plaquetas VARCHAR(100),
	glucosa VARCHAR(100),
	hierro VARCHAR(100),
	clave VARCHAR(344) NOT NULL,
	PRIMARY KEY (id),
	FOREIGN KEY(usuario_id) REFERENCES usuarios(id) ON DELETE CASCADE
);`

var USERS_PERMISOS_ANALITICAS_TABLE string = `
CREATE TABLE IF NOT EXISTS usuarios_permisos_analiticas (
	analitica_id INT,
	medico_id INT,
	clave VARCHAR(344) NOT NULL,
	fecha_expiracion DATETIME,
	PRIMARY KEY (analitica_id, medico_id),
	FOREIGN KEY(analitica_id) REFERENCES usuarios_analiticas(id) ON DELETE CASCADE,
	FOREIGN KEY(medico_id) REFERENCES usuarios(id) ON DELETE CASCADE
);`

var ANALITICAS_TAGS_TABLE string = `
CREATE TABLE IF NOT EXISTS analiticas_tags (
	analitica_id INT,
	tag_id INT,
	PRIMARY KEY (analitica_id, tag_id),
	FOREIGN KEY(tag_id) REFERENCES tags(id) ON DELETE CASCADE,
	FOREIGN KEY(analitica_id) REFERENCES usuarios_analiticas(id) ON DELETE CASCADE
);`

var ESTADISTICAS_ANALITICAS_TABLE string = `
CREATE TABLE IF NOT EXISTS estadisticas_analiticas (
	leucocitos FLOAT,
	hematies FLOAT,
	plaquetas FLOAT,
	glucosa FLOAT,
	hierro FLOAT
);`
