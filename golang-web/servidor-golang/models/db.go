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

	//SEEDERS
	//Roles
	query("INSERT IGNORE INTO roles (id,nombre,descripcion) VALUES (" + strconv.Itoa(Rol_paciente.Id) + ",'" + Rol_paciente.Nombre + "', '" + Rol_paciente.Descripcion + "');")
	query("INSERT IGNORE INTO roles (id,nombre,descripcion) VALUES (" + strconv.Itoa(Rol_enfermero.Id) + ",'" + Rol_enfermero.Nombre + "', '" + Rol_enfermero.Descripcion + "');")
	query("INSERT IGNORE INTO roles (id,nombre,descripcion) VALUES (" + strconv.Itoa(Rol_medico.Id) + ",'" + Rol_medico.Nombre + "', '" + Rol_medico.Descripcion + "');")
	query("INSERT IGNORE INTO roles (id,nombre,descripcion) VALUES (" + strconv.Itoa(Rol_administradorC.Id) + ",'" + Rol_administradorC.Nombre + "', '" + Rol_administradorC.Descripcion + "');")
	query("INSERT IGNORE INTO roles (id,nombre,descripcion) VALUES (" + strconv.Itoa(Rol_administradorG.Id) + ",'" + Rol_administradorG.Nombre + "', '" + Rol_administradorG.Descripcion + "');")

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
	dni VARCHAR(9) UNIQUE,
	nombre VARCHAR(25) NOT NULL,
	apellidos VARCHAR(50) NOT NULL,
	email VARCHAR(59) UNIQUE,
	password VARCHAR(100) NOT NULL,
	created_at DATETIME,
	PRIMARY KEY (id)
);`

var ROLES_TABLE string = `
CREATE TABLE IF NOT EXISTS roles (
	id INT AUTO_INCREMENT,
	nombre VARCHAR(20) UNIQUE,
	descripcion VARCHAR(50),
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
	created_at DATETIME,
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
