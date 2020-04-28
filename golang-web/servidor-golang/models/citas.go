package models

import (
	"fmt"
	"strconv"

	util "../utils"
)

func ComprobarDiaDisponible(doctor_id string, anyo int, mes int, dia int) bool {
	anyoString := strconv.Itoa(anyo)
	mesString := strconv.Itoa(mes)
	diaString := strconv.Itoa(dia)
	rows, err := db.Query("SELECT count(*) FROM citas WHERE medico_id = " + doctor_id + " AND " +
		" anyo = " + anyoString + " AND mes = " + mesString + " AND dia = " + diaString)
	if err == nil {
		var horasNumber int
		defer rows.Close()
		rows.Next()
		rows.Scan(&horasNumber)
		if horasNumber >= 5 {
			return false
		}
	} else {
		fmt.Println(err)
	}
	return true
}

func ComprobarHoraDisponible(doctor_id string, anyo int, mes int, dia int, hora int) bool {
	anyoString := strconv.Itoa(anyo)
	mesString := strconv.Itoa(mes)
	diaString := strconv.Itoa(dia)
	horaString := strconv.Itoa(hora)
	rows, err := db.Query("SELECT count(*) FROM citas WHERE medico_id = " + doctor_id + " AND " +
		" anyo = " + anyoString + " AND mes = " + mesString + " AND dia = " + diaString + " AND hora = " + horaString)
	if err == nil {
		var horasNumber int
		defer rows.Close()
		rows.Next()
		rows.Scan(&horasNumber)
		if horasNumber >= 1 {
			return false
		}
	} else {
		fmt.Println(err)
	}
	return true
}

func GetCitasFuturas(userId string) (citasList []util.CitaJSON, err error) {
	rows, err := db.Query("SELECT id, medico_id, tipo, anyo, mes, dia, hora FROM citas WHERE paciente_id = " + userId)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var c util.CitaJSON
			var nombreDoctor string
			rows.Scan(&c.Id, &c.MedicoId, &c.Tipo, &c.Anyo, &c.Mes, &c.Dia, &c.Hora)
			rowNombreMedico, _ := db.Query("SELECT nombreDoctor FROM medicos_nombres WHERE usuario_id = " + c.MedicoId)
			rowNombreMedico.Next()
			rowNombreMedico.Scan(&nombreDoctor)
			c.MedicoNombre = nombreDoctor
			citasList = append(citasList, c)
		}
		return citasList, err
	} else {
		fmt.Println(err)
		return nil, err
	}
}
