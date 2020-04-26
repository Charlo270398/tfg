package models

import (
	"fmt"

	util "../utils"
)

func InsertEspecialidadMedico(user_id int, especialidad_id int) (result bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT INTO usuarios_especialidades (usuario_id, especialidad_id) VALUES (?, ?)`, user_id,
		especialidad_id)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func InsertNombresMedico(user util.User_JSON) (result bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT INTO medicos_nombres (usuario_id, nombreDoctor) VALUES (?, ?)`, user.Id, user.NombreDoctor)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}
