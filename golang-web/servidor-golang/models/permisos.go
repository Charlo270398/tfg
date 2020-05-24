package models

import (
	"fmt"

	util "../utils"
)

func SolicitarPermisoTotalHistorial(historial util.Historial_JSON) (result bool, err error) {
	loadHistorial, _ := GetHistorialById(historial.Id)
	//INSERT
	_, err = db.Exec(`INSERT IGNORE INTO solicitar_historial_total (paciente_id, empleado_id) VALUES (?, ?)`, loadHistorial.PacienteId, historial.UserToken.UserId)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func SolicitarPermisoBasicoHistorial(historial util.Historial_JSON) (result bool, err error) {
	loadHistorial, _ := GetHistorialById(historial.Id)
	//INSERT
	_, err = db.Exec(`INSERT IGNORE INTO solicitar_historial (paciente_id, empleado_id) VALUES (?, ?)`, loadHistorial.PacienteId, historial.UserToken.UserId)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func SolicitarPermisoEntrada(historial util.Historial_JSON) (result bool, err error) {
	loadHistorial, _ := GetHistorialById(historial.Id)
	//INSERT
	_, err = db.Exec(`INSERT IGNORE INTO solicitar_historial (paciente_id, empleado_id) VALUES (?, ?)`, loadHistorial.PacienteId, historial.UserToken.UserId)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func ComprobarSolicitudesPermiso(userId string) (result bool, err error) {
	var number int
	row, _ := db.Query(`SELECT count(*) FROM solicitar_historial where paciente_id = ` + userId)
	defer row.Close()
	row.Next()
	row.Scan(&number)
	if number == 0 {
		row, _ = db.Query(`SELECT count(*) FROM solicitar_historial_total where paciente_id = ` + userId)
		defer row.Close()
		row.Next()
		row.Scan(&number)
		if number == 0 {
			row, _ = db.Query(`SELECT count(*) FROM solicitar_entradas_historial where paciente_id = ` + userId)
			defer row.Close()
			row.Next()
			row.Scan(&number)
			if number == 0 {
				row, _ = db.Query(`SELECT count(*) FROM solicitar_analiticas where paciente_id = ` + userId)
				defer row.Close()
				row.Next()
				row.Scan(&number)
				if number == 0 {
					return false, nil
				} else {
					return true, nil
				}
			} else {
				return true, nil
			}
		} else {
			return true, nil
		}
	} else {
		return true, nil
	}
}

func ListarSolicitudesPermiso(userId string) (solicitudes []util.Solicitud_JSON, err error) {
	historial, _ := GetHistorialByUserId(userId)
	//HISTORIAL BASICO
	rows, err := db.Query(`SELECT empleado_id FROM solicitar_historial where paciente_id = ` + userId) // check err
	if err == nil {
		var s util.Solicitud_JSON
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&s.EmpleadoId)
			s.HistorialId = historial.Id
			empleadoNombre, _ := GetNombreEmpleado(s.EmpleadoId)
			s.NombreEmpleado = empleadoNombre
			s.TipoHistorial = "BASICO"
			solicitudes = append(solicitudes, s)
		}
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return solicitudes, err
	}
	//HISTORIAL TOTAL
	rows, err = db.Query(`SELECT empleado_id FROM solicitar_historial_total where paciente_id = ` + userId) // check err
	if err == nil {
		var s util.Solicitud_JSON
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&s.EmpleadoId)
			s.HistorialId = historial.Id
			empleadoNombre, _ := GetNombreEmpleado(s.EmpleadoId)
			s.NombreEmpleado = empleadoNombre
			s.TipoHistorial = "TOTAL"
			solicitudes = append(solicitudes, s)
		}
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return solicitudes, err
	}
	//ENTRADAS
	rows, err = db.Query(`SELECT empleado_id, entrada_id FROM solicitar_entradas_historial where paciente_id = ` + userId) // check err
	if err == nil {
		var s util.Solicitud_JSON
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&s.EmpleadoId, &s.EntradaId)
			empleadoNombre, _ := GetNombreEmpleado(s.EmpleadoId)
			s.NombreEmpleado = empleadoNombre
			solicitudes = append(solicitudes, s)
		}
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return solicitudes, err
	}
	//ANALITICAS
	rows, err = db.Query(`SELECT empleado_id, analitica_id  FROM solicitar_analiticas where paciente_id = ` + userId) // check err
	if err == nil {
		var s util.Solicitud_JSON
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&s.EmpleadoId, &s.AnaliticaId)
			empleadoNombre, _ := GetNombreEmpleado(s.EmpleadoId)
			s.NombreEmpleado = empleadoNombre
			solicitudes = append(solicitudes, s)
		}
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return solicitudes, err
	}
	return solicitudes, nil
}
