package models

import (
	"fmt"
	"time"

	util "../utils"
)

func GetHistorialByUserId(userId string) (historial util.Historial_JSON, err error) {
	row, err := db.Query(`SELECT id, sexo, alergias, clave FROM usuarios_historial where usuario_id = ` + userId) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&historial.Id, &historial.Sexo, &historial.Alergias, &historial.Clave)
		return historial, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return historial, err
	}
}

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

func InsertShareHistorial(historial util.Historial_JSON) (result bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT IGNORE INTO usuarios_permisos_historial (historial_id, medico_id,sexo,alergias, nombrePaciente, clave) VALUES (?, ?, ?, ?, ?, ?)`, historial.Id,
		historial.MedicoId, historial.Sexo, historial.Alergias, historial.NombrePaciente, historial.Clave)
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

func GetHistorialesCompartidosByMedicoId(medicoId string) (historiales []util.Historial_JSON, err error) {
	rows, err := db.Query(`SELECT historial_id, sexo, alergias, nombrePaciente, clave FROM usuarios_permisos_historial where medico_id = ` + medicoId) // check err
	if err == nil {
		var h util.Historial_JSON
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&h.Id, &h.Sexo, &h.Alergias, &h.NombrePaciente, &h.Clave)
			historiales = append(historiales, h)
		}
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return historiales, err
	}
	return historiales, nil
}
