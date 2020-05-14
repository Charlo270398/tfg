package models

import (
	"fmt"
	"time"

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

func InsertNombresEmpleado(user util.User_JSON) (result bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT IGNORE INTO empleados_nombres (usuario_id, nombre) VALUES (?, ?)`, user.Id, user.NombreDoctor)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func GetDiasDisponiblesMedico(doctor_id string) (diasDisponibles []util.Cita, err error) {
	day := time.Now()
	proximosDias := 30
	//Comprobamos que la hora no esta por detras de la ultima hora laboral (13)
	if day.Hour() >= 13 {
		day = day.Add(time.Hour * 24)
		proximosDias = proximosDias - 1
	}
	//Proximos 30 dias
	for i := 0; i < proximosDias; i++ {
		//Sabado y domingo no se abre, son el dia 6 y 7 de la semana
		if int(day.Weekday()) != 6 && int(day.Weekday()) != 7 {
			var cita util.Cita
			cita.Anyo = day.Year()
			cita.Mes = int(day.Month())
			cita.Dia = day.Day()
			cita.Fecha = day
			//Comprobar en BD que el dia está disponible (quedan horas)
			if ComprobarDiaDisponible(doctor_id, cita.Anyo, cita.Mes, cita.Dia) == true {
				diasDisponibles = append(diasDisponibles, cita)
			}
		}
		day = day.Add(time.Hour * 24)
	}
	return diasDisponibles, nil
}

func GetHorasDiaDisponiblesMedico(doctor_id string, dia string) (horasDisponibles []util.Cita, err error) {
	layout := "2006-01-02T15:04:05.000Z"
	day, err := time.Parse(layout, dia+"T09:00:00.000Z")
	//Se trabaja de 9 a 14
	for i := 0; i < 5; i++ {
		var cita util.Cita
		cita.Anyo = day.Year()
		cita.Mes = int(day.Month())
		cita.Dia = day.Day()
		cita.Hora = day.Hour()
		cita.Fecha = day
		if ComprobarHoraDisponible(doctor_id, cita.Anyo, cita.Mes, cita.Dia, cita.Hora) == true {
			horasDisponibles = append(horasDisponibles, cita)
		}
		day = day.Add(time.Hour)
	}
	return horasDisponibles, nil
}
