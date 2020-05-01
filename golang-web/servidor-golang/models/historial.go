package models

import (
	"fmt"

	util "../utils"
)

func InsertHistorial(entrada util.User_JSON) (result bool, err error) {
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
