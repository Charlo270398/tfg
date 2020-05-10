package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	util "../utils"
)

func GetHistorialById(historialId int) (historial util.Historial_JSON, err error) {
	historialIdString := strconv.Itoa(historialId)
	row, err := db.Query(`SELECT id, usuario_id, sexo, alergias, clave_maestra FROM usuarios_historial where id = ` + historialIdString) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&historial.Id, &historial.PacienteId, &historial.Sexo, &historial.Alergias, &historial.Clave, &historial.ClaveMaestra)
		return historial, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return historial, err
	}
}

func GetHistorialByUserId(userId string) (historial util.Historial_JSON, err error) {
	row, err := db.Query(`SELECT id, sexo, alergias, clave, clave_maestra FROM usuarios_historial where usuario_id = ` + userId) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&historial.Id, &historial.Sexo, &historial.Alergias, &historial.Clave, &historial.ClaveMaestra)
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
	_, err = db.Exec(`INSERT INTO usuarios_historial (sexo,alergias,usuario_id,ultima_actualizacion, clave, clave_maestra) VALUES (?, ?, ?, ?, ?, ?)`, user.Sexo,
		user.Alergias, user.Id, createdAt, user.Clave, user.ClaveMaestra)
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
	_, err = db.Exec(`INSERT IGNORE INTO usuarios_permisos_historial (historial_id, empleado_id, clave) VALUES (?, ?, ?)`, historial.Id,
		historial.MedicoId, historial.Clave)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func InsertEntradaHistorial(entrada util.EntradaHistorial_JSON) (inserted_id int, err error) {
	createdAt := time.Now()
	pacienteIdString := strconv.Itoa(entrada.PacienteId)
	historialPaciente, _ := GetHistorialByUserId(pacienteIdString)
	historialPacienteIdString := strconv.Itoa(historialPaciente.Id)
	//INSERT
	entradaId, err := db.Exec(`INSERT INTO usuarios_entradas_historial (empleado_id, historial_id, motivo_consulta, juicio_diagnostico, clave, created_at, tipo) VALUES (?, ?, ?, ?, ?, ?, ?)`, entrada.UserToken.UserId,
		historialPacienteIdString, entrada.MotivoConsulta, entrada.JuicioDiagnostico, entrada.Clave, createdAt.Local(), entrada.Tipo)
	if err == nil {
		id, _ := entradaId.LastInsertId()
		inserted_id = int(id)
		return inserted_id, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return -1, err
	}
}

func InsertEntradaCompartidaHistorial(entrada util.EntradaHistorial_JSON) (result bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT INTO usuarios_permisos_entradas_historial (entrada_id, empleado_id, clave) VALUES (?, ?, ?)`,
		entrada.Id, entrada.UserToken.UserId, entrada.Clave)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func GetHistorialesCompartidosByMedicoId(medicoId string) (historiales []util.Historial_JSON, err error) {
	rows, err := db.Query(`SELECT historial_id, clave FROM usuarios_permisos_historial where empleado_id = ` + medicoId) // check err
	if err == nil {
		var h util.Historial_JSON
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&h.Id, &h.Clave)
			historial, _ := GetHistorialById(h.Id)
			h.Sexo = historial.Sexo
			userData, _ := GetUserById(historial.PacienteId)
			h.NombrePaciente = userData.Nombre
			h.ApellidosPaciente = userData.Apellidos
			historiales = append(historiales, h)
		}
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return historiales, err
	}
	return historiales, nil
}

func GetHistorialCompartidoByMedicoIdPacienteId(medicoId string, pacienteId string) (historial util.Historial_JSON, err error) {
	historialPaciente, _ := GetHistorialByUserId(pacienteId)
	pacienteIdInt, _ := strconv.Atoi(pacienteId)
	userData, _ := GetUserById(pacienteIdInt)
	historialPacienteIdString := strconv.Itoa(historialPaciente.Id)
	rows, err := db.Query(`SELECT historial_id, clave FROM usuarios_permisos_historial where empleado_id = ` + medicoId + ` and historial_id = ` + historialPacienteIdString) // check err
	if err == nil {
		defer rows.Close()
		rows.Next()
		rows.Scan(&historial.Id, &historial.Clave)
		historial.Alergias = historialPaciente.Alergias
		historial.Sexo = historialPaciente.Sexo
		historial.NombrePaciente = userData.Nombre
		historial.ApellidosPaciente = userData.Apellidos
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return historial, err
	}
	return historial, nil
}

func GetEntradasHistorialByHistorialId(historialId int) (entradas []util.EntradaHistorial_JSON, err error) {
	historialPacienteIdString := strconv.Itoa(historialId)
	rows, err := db.Query(`SELECT id, empleado_id, historial_id, motivo_consulta, juicio_diagnostico, clave, created_at, tipo FROM usuarios_entradas_historial where historial_id = ` + historialPacienteIdString + " order by created_at desc") // check err
	if err == nil {
		var e util.EntradaHistorial_JSON
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&e.Id, &e.EmpleadoId, &e.HistorialId, &e.MotivoConsulta, &e.JuicioDiagnostico, &e.Clave, &e.CreatedAt, &e.Tipo)
			//Cambio horario y formato
			words := strings.Fields(e.CreatedAt)
			day := words[0] + "T" + words[1] + "Z"
			layout := "2006-01-02T15:04:05.000000Z"
			t, err := time.Parse(layout, day)
			if err != nil {
				fmt.Println(err)
			}
			t = t.Local()
			e.CreatedAt = fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d",
				t.Day(), t.Month(), t.Year(),
				t.Hour(), t.Minute(), t.Second())
			e.EmpleadoNombre, _ = GetNombreEmpleado(e.EmpleadoId)
			entradas = append(entradas, e)
		}
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return entradas, err
	}
	return entradas, nil
}
