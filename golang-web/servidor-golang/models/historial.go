package models

import (
	"fmt"
	"strconv"
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
	_, err = db.Exec(`INSERT IGNORE INTO usuarios_permisos_historial (historial_id, empleado_id,sexo,alergias, nombrePaciente, clave) VALUES (?, ?, ?, ?, ?, ?)`, historial.Id,
		historial.MedicoId, historial.Sexo, historial.Alergias, historial.NombrePaciente, historial.Clave)
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
	entradaId, err := db.Exec(`INSERT INTO usuarios_entradas_historial (empleado_id, historial_id, motivo_consulta, juicio_diagnostico, clave, created_at) VALUES (?, ?, ?, ?, ?, ?)`, entrada.UserToken.UserId,
		historialPacienteIdString, entrada.MotivoConsulta, entrada.JuicioDiagnostico, entrada.Clave, createdAt.Local())
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

func GetHistorialCompartidosByCita(medicoId string) (historiales []util.Historial_JSON, err error) {
	rows, err := db.Query(`SELECT historial_id, sexo, alergias, nombrePaciente, clave FROM usuarios_permisos_historial where empleado_id = ` + medicoId) // check err
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

func GetHistorialesCompartidosByMedicoId(medicoId string) (historiales []util.Historial_JSON, err error) {
	rows, err := db.Query(`SELECT historial_id, sexo, alergias, nombrePaciente, clave FROM usuarios_permisos_historial where empleado_id = ` + medicoId) // check err
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

func GetHistorialCompartidoByMedicoIdPacienteId(medicoId string, pacienteId string) (historial util.Historial_JSON, err error) {
	historialPaciente, _ := GetHistorialByUserId(pacienteId)
	historialPacienteIdString := strconv.Itoa(historialPaciente.Id)
	rows, err := db.Query(`SELECT historial_id, sexo, alergias, nombrePaciente, clave FROM usuarios_permisos_historial where empleado_id = ` + medicoId + ` and historial_id = ` + historialPacienteIdString) // check err
	if err == nil {
		defer rows.Close()
		rows.Next()
		rows.Scan(&historial.Id, &historial.Sexo, &historial.Alergias, &historial.NombrePaciente, &historial.Clave)
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return historial, err
	}
	return historial, nil
}
