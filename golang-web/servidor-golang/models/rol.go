package models

import (
	"fmt"
	"strconv"

	util "../utils"
)

var Rol_paciente util.Rol
var Rol_enfermero util.Rol
var Rol_medico util.Rol
var Rol_administradorC util.Rol
var Rol_administradorG util.Rol
var Rol_emergencias util.Rol

func LoadRoles() {
	//Definimos los roles basicos
	Rol_paciente = util.Rol{Id: 1, Nombre: "paciente", Descripcion: "Paciente"}
	Rol_enfermero = util.Rol{Id: 2, Nombre: "enfermero", Descripcion: "Enfermero"}
	Rol_medico = util.Rol{Id: 3, Nombre: "medico", Descripcion: "Medico"}
	Rol_administradorC = util.Rol{Id: 4, Nombre: "administradorC", Descripcion: "Administrador clinica"}
	Rol_administradorG = util.Rol{Id: 5, Nombre: "administradorG", Descripcion: "Administrador global"}
	Rol_emergencias = util.Rol{Id: 6, Nombre: "emergencias", Descripcion: "Emergencias"}
}

func InsertUserAndRole(userid int, roles []int) (inserted bool, err error) {
	//INSERT
	for _, rolId := range roles {
		_, err = db.Exec(`INSERT INTO usuarios_roles (usuario_id, rol_id) VALUES (?, ?)`, userid,
			rolId)
		if err != nil {
			fmt.Println(err)
			util.PrintErrorLog(err)
			return false, nil
		}
	}
	return true, nil
}

func GetRolesList() (roles []util.Rol, err error) {
	row, err := db.Query(`SELECT id, nombre, descripcion FROM roles`) // check err
	if err == nil {
		defer row.Close()
		var roles []util.Rol
		for row.Next() {
			var r util.Rol
			row.Scan(&r.Id, &r.Nombre, &r.Descripcion)
			roles = append(roles, r)
		}
		return roles, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return nil, err
	}
}

func GetRoleById(rol_id string) (rol util.Rol, err error) {
	row, err := db.Query(`SELECT id, nombre, descripcion FROM usuarios where id = '` + rol_id + `'`) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&rol.Id, &rol.Nombre, &rol.Descripcion)
		return rol, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return rol, err
	}
}

func GetRolesbyUserId(user_id string) (roles []int, err error) {
	row, err := db.Query(`SELECT rol_id FROM usuarios_roles where usuario_id = '` + user_id + `'`) // check err
	if err == nil {
		defer row.Close()
		var roles []int
		for row.Next() {
			var r int
			row.Scan(&r)
			roles = append(roles, r)
		}
		return roles, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return nil, err
	}
}

//GESTION AUTORIZACION DEL USUARIO
//Indica si ese usuario tiene permisos para ejercer el  rol indicado
func GetAuthorizationbyUserId(user_id string, token string, rol_id int) (result bool, err error) {
	user_id_int, _ := strconv.Atoi(user_id)
	//Comprobamos token del usuario
	proved, err := ProveUserToken(user_id_int, token)
	if err != nil {
		return false, err
	}
	if proved == false {
		return false, nil
	}
	//Comprobamos que el rol es correcto
	row, err := db.Query(`SELECT count(*) FROM usuarios_roles where usuario_id = '` + user_id + `' and rol_id = '` + strconv.Itoa(rol_id) + `'`) // check err
	if err == nil {
		defer row.Close()
		var count int
		row.Next()
		row.Scan(&count)
		if count == 1 {
			return true, nil
		}
		return false, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}
