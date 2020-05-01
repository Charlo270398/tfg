package models

import (
	"fmt"
	"time"

	util "../utils"
)

func InsertHistorial(user util.User_JSON) (result bool, err error) {
	//INSERT
	createdAt := time.Now()
	_, err = db.Exec(`INSERT INTO usuarios_historial (sexo,alergias,usuario_id,ultima_actualizacion, clave) VALUES (?, ?, ?, ?, ?)`, user.Sexo,
		user.Alergias, user.Id, createdAt, user.Clave)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func InsertEntradaHistorial(entrada util.EntradaHistorial_JSON) (result bool, err error) {
	//createdAt := time.Now()
	fmt.Println(entrada.JuicioDiagnostico)
	fmt.Println(entrada.MotivoConsulta)
	return true, nil
	//INSERT
	/*
		_, err = db.Exec(`INSERT INTO citas (medico_id, paciente_id, anyo, mes, dia, hora, tipo) VALUES (?, ?, ?, ?, ?, ?, ?)`, cita.MedicoId,
			cita.UserToken.UserId, fechaCita.Year(), int(fechaCita.Month()), fechaCita.Day(), fechaCita.Hour(), "Consulta")
		if err == nil {
			return true, nil
		} else {
			fmt.Println(err)
			util.PrintErrorLog(err)
			return false, err
		}*/
}
