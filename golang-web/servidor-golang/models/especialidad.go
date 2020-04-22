package models

import (
	"fmt"
	"strconv"

	util "../utils"
)

func InsertEspecialidad(especialidad util.Especialidad_JSON) (especialidadId int, err error) {
	//INSERT
	res, err := db.Exec(`INSERT INTO especialidades (nombre) VALUES (?)`, especialidad.Nombre)
	if err == nil {
		especialidadId, _ := res.LastInsertId()
		return int(especialidadId), nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return -1, err
	}
}

func EditEspecialidadData(especialidad util.Especialidad_JSON) (edited bool, err error) {
	//UPDATE
	_, err = db.Exec(`UPDATE especialidades set nombre = ? where id = ?`, especialidad.Nombre, especialidad.Id)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return false, nil
}

func DeleteEspecialidad(especialidad_id int) bool {
	_, err := db.Exec(`DELETE FROM especialidades where id = ` + strconv.Itoa(especialidad_id))
	if err == nil {
		return true
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return false
}

func GetEspecialidadList() (especialidadList []util.Especialidad, err error) {
	rows, err := db.Query("SELECT id, nombre FROM especialidades")
	if err == nil {
		defer rows.Close()
		var especialidades []util.Especialidad
		for rows.Next() {
			var e util.Especialidad
			rows.Scan(&e.Id, &e.Nombre)
			especialidades = append(especialidades, e)
		}
		return especialidades, err
	} else {
		fmt.Println(err)
		return nil, err
	}
}

func GetEspecialidadPagination(page int) []util.Especialidad {
	firstRow := strconv.Itoa(page * 10)
	lastRow := strconv.Itoa((page * 10) + 10)
	rows, err := db.Query("SELECT id, nombre FROM especialidades LIMIT " + firstRow + "," + lastRow)
	if err == nil {
		defer rows.Close()
		var especialidades []util.Especialidad
		for rows.Next() {
			var e util.Especialidad
			rows.Scan(&e.Id, &e.Nombre)
			especialidades = append(especialidades, e)
		}
		return especialidades
	} else {
		fmt.Println(err)
		return nil
	}
}
