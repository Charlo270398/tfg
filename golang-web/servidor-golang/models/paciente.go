package models

import (
	"fmt"
	"time"

	util "../utils"
)

func InsertCita(cita util.CitaJSON) (result bool, err error) {
	layout := "2006-01-02T15:04:05.000Z"
	fechaCita, err := time.Parse(layout, cita.FechaString+".000Z")
	if fechaCita.Hour() != 9 && fechaCita.Hour() != 10 && fechaCita.Hour() != 11 && fechaCita.Hour() != 12 && fechaCita.Hour() != 13 {
		return false, nil
	} else {
		//INSERT
		_, err = db.Exec(`INSERT INTO citas (medico_id, paciente_id, anyo, mes, dia, hora, tipo) VALUES (?, ?, ?, ?, ?, ?, ?)`, cita.MedicoId,
			cita.UserToken.UserId, fechaCita.Year(), int(fechaCita.Month()), fechaCita.Day(), fechaCita.Hour(), "Consulta")
		if err == nil {
			return true, nil
		} else {
			fmt.Println(err)
			util.PrintErrorLog(err)
			return false, err
		}
	}
}
