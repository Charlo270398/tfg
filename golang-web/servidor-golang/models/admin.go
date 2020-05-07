package models

import (
	"fmt"
	"strconv"

	util "../utils"
)

//GESTION AUTORIZACION DEL USUARIO
//Verifica si el usuario es admin de la clinica
func VerifyAdmin(user_id string, clinica_id string, token string) (result bool, err error) {
	//Comprobamos token del usuario
	user_idInt, _ := strconv.Atoi(user_id)
	proved, err := ProveUserToken(user_idInt, token)
	if err != nil {
		return false, err
	}
	if proved == false {
		return false, nil
	}
	//Comprobamos que el rol es correcto
	row, err := db.Query(`SELECT count(*) FROM usuarios_clinicas where usuario_id = '` + user_id + `' and clinica_id = '` + clinica_id + `' and ` +
		`rol_id = ` + strconv.Itoa(Rol_administradorC.Id)) // check err
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

func GetNombreEmpleado(empleadoId int) (nombre string, err error) {
	empleadoIdString := strconv.Itoa(empleadoId)
	rows, err := db.Query(`SELECT nombre FROM empleados_nombres where usuario_id = ` + empleadoIdString) // check err
	if err == nil {
		defer rows.Close()
		rows.Next()
		rows.Scan(&nombre)
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return nombre, err
	}
	return nombre, nil
}
