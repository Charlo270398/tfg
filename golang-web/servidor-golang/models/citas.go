package models

import (
	"fmt"
	"strconv"
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
