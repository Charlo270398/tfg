package models

import (
	"fmt"
	"strconv"
	"time"

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
		if time.Now().Local().Hour() > hora-1 && time.Now().Local().Day() == dia {
			return false
		}
	} else {
		fmt.Println(err)
	}
	return true
}

func GetCitasFuturasPaciente(pacienteId string) (citasList []util.CitaJSON, err error) {
	rows, err := db.Query("SELECT id, medico_id, tipo, anyo, mes, dia, hora FROM citas WHERE paciente_id = " + pacienteId + " ORDER BY anyo, mes, dia, hora")
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var c util.CitaJSON
			var nombreDoctor string
			rows.Scan(&c.Id, &c.MedicoId, &c.Tipo, &c.Anyo, &c.Mes, &c.Dia, &c.Hora)
			rowNombreMedico, _ := db.Query("SELECT nombre FROM empleados_nombres WHERE usuario_id = " + c.MedicoId)
			rowNombreMedico.Next()
			rowNombreMedico.Scan(&nombreDoctor)
			c.MedicoNombre = nombreDoctor
			mesFormato := fmt.Sprintf("%02d", c.Mes)
			diaFormato := fmt.Sprintf("%02d", c.Dia)
			diaCompleto := strconv.Itoa(c.Anyo) + "-" + mesFormato + "-" + diaFormato + "T" + strconv.Itoa(c.Hora-2) //TODO hacer esto mejor
			layout := "2006-01-02T15:04:05.000"
			fechaCita, _ := time.Parse(layout, diaCompleto+":00:00.000")
			if time.Now().Local().Before(fechaCita) {
				citasList = append(citasList, c)
			}
		}
		return citasList, err
	} else {
		fmt.Println(err)
		return nil, err
	}
}

func GetCitasFuturasMedico(medicoId string) (citasList []util.CitaJSON, err error) {
	rows, err := db.Query("SELECT id, paciente_id, tipo, anyo, mes, dia, hora FROM citas WHERE medico_id = " + medicoId)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var c util.CitaJSON
			rows.Scan(&c.Id, &c.PacienteId, &c.Tipo, &c.Anyo, &c.Mes, &c.Dia, &c.Hora)
			//Deberiamos obtener nombre del paciente y asignarlo
			mesFormato := fmt.Sprintf("%02d", c.Mes)
			diaFormato := fmt.Sprintf("%02d", c.Dia)
			diaCompleto := strconv.Itoa(c.Anyo) + "-" + mesFormato + "-" + diaFormato + "T" + strconv.Itoa(c.Hora-2) //TODO hacer esto mejor
			layout := "2006-01-02T15:04:05.000"
			fechaCita := time.Now().Local()
			fechaCita, _ = time.Parse(layout, diaCompleto+":00:00.000")
			if time.Now().Local().Before(fechaCita) {
				citasList = append(citasList, c)
			}
		}
		return citasList, err
	} else {
		fmt.Println(err)
		return nil, err
	}
}

func GetCitaActualMedico(medicoId string) (citaId int, err error) {
	citaId = -1
	day := time.Now().Local()
	anyoString := strconv.Itoa(day.Year())
	mesString := strconv.Itoa(int(day.Month()))
	diaString := strconv.Itoa(day.Day())
	rows, err := db.Query("SELECT id, anyo, mes, dia, hora FROM citas WHERE medico_id = " + medicoId + " AND " +
		" anyo = " + anyoString + " AND mes = " + mesString + " AND dia = " + diaString)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var c util.CitaJSON
			rows.Scan(&c.Id, &c.Anyo, &c.Mes, &c.Dia, &c.Hora)
			//Deberiamos obtener nombre del paciente y asignarlo
			mesFormato := fmt.Sprintf("%02d", c.Mes)
			diaFormato := fmt.Sprintf("%02d", c.Dia)
			diaCompleto := strconv.Itoa(c.Anyo) + "-" + mesFormato + "-" + diaFormato + "T" + strconv.Itoa(c.Hora-2) //TODO hacer esto mejor
			layout := "2006-01-02T15:04:05.000"
			fechaCita := time.Now().Local()
			fechaCita, _ = time.Parse(layout, diaCompleto+":00:00.000")
			if time.Now().Local().YearDay() == fechaCita.YearDay() && time.Now().Local().Year() == fechaCita.Year() {
				//Es la misma hora
				if day.Hour() == c.Hora {
					return c.Id, nil
				}
				//Quedan 10 min para la cita
				if day.Hour() == c.Hora-1 && day.Minute() >= 50 {
					citaId = c.Id
				}
			}
		}
	} else {
		fmt.Println(err)
	}
	return citaId, err
}

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

func GetCitaUserIdByCitaId(citaId int) (pacienteId int, err error) {
	row, err := db.Query(`SELECT paciente_id FROM citas WHERE id = ` + strconv.Itoa(citaId)) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&pacienteId)
		return pacienteId, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return pacienteId, err
	}
}

func GetCitaById(citaId int) (cita util.CitaJSON, err error) {
	row, err := db.Query(`SELECT * FROM citas WHERE id = ` + strconv.Itoa(citaId)) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&cita.Id, &cita.MedicoId, &cita.PacienteId, &cita.Anyo, &cita.Mes, &cita.Dia, &cita.Hora, &cita.Tipo)
		return cita, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return cita, err
	}
}
