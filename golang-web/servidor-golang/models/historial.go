package models

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	util "../utils"
)

func GetHistorialById(historialId int) (historial util.Historial_JSON, err error) {
	historialIdString := strconv.Itoa(historialId)
	row, err := db.Query(`SELECT id, usuario_id, sexo, alergias, clave, clave_maestra FROM usuarios_historial where id = ` + historialIdString) // check err
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
	historialPacienteIdString := strconv.Itoa(entrada.HistorialId)
	//INSERT
	entradaId, err := db.Exec(`INSERT INTO usuarios_entradas_historial (empleado_id, historial_id, motivo_consulta, juicio_diagnostico, clave, clave_maestra, created_at, tipo) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, entrada.UserToken.UserId,
		historialPacienteIdString, entrada.MotivoConsulta, entrada.JuicioDiagnostico, entrada.Clave, entrada.ClaveMaestra, createdAt.Local(), entrada.Tipo)
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

func InsertEntradaHistorialPacienteId(entrada util.EntradaHistorial_JSON) (inserted_id int, err error) {
	createdAt := time.Now()
	pacienteIdString := strconv.Itoa(entrada.PacienteId)
	historialPaciente, _ := GetHistorialByUserId(pacienteIdString)
	historialPacienteIdString := strconv.Itoa(historialPaciente.Id)
	//INSERT
	entradaId, err := db.Exec(`INSERT INTO usuarios_entradas_historial (empleado_id, historial_id, motivo_consulta, juicio_diagnostico, clave, clave_maestra, created_at, tipo) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, entrada.UserToken.UserId,
		historialPacienteIdString, entrada.MotivoConsulta, entrada.JuicioDiagnostico, entrada.Clave, entrada.ClaveMaestra, createdAt.Local(), entrada.Tipo)
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
	_, err = db.Exec(`INSERT IGNORE INTO usuarios_permisos_entradas_historial (entrada_id, empleado_id, clave) VALUES (?, ?, ?)`,
		entrada.Id, entrada.UserToken.UserId, entrada.Clave)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func InsertEntradaCompartidaHistorialPermisos(entrada util.EntradaHistorial_JSON) (result bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT IGNORE INTO usuarios_permisos_entradas_historial (entrada_id, empleado_id, clave) VALUES (?, ?, ?)`,
		entrada.Id, entrada.EmpleadoId, entrada.Clave)
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
		if historial.Id != 0 {
			historial.Alergias = historialPaciente.Alergias
			historial.Sexo = historialPaciente.Sexo
			historial.NombrePaciente = userData.Nombre
			historial.ApellidosPaciente = userData.Apellidos
			historial.Entradas, _ = GetEntradasHistorialByHistorialId(historial.Id)
			for index, entrada := range historial.Entradas {
				empleadoIdInt, _ := strconv.Atoi(medicoId)
				historial.Entradas[index].Clave, _ = GetClaveCompartidaEntradaHistorialByEntradaIdEmpleadoId(entrada.Id, empleadoIdInt)
			}
			historial.Analiticas, _ = GetAnaliticasHistorialByHistorialId(historial.Id)
			for index, analitica := range historial.Analiticas {
				empleadoIdInt, _ := strconv.Atoi(medicoId)
				historial.Analiticas[index].Clave, _ = GetClaveCompartidaAnaliticaHistorialByEntradaIdEmpleadoId(analitica.Id, empleadoIdInt)
			}
		}
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return historial, err
	}
	return historial, nil
}

func GetEntradasHistorialByHistorialId(historialId int) (entradas []util.EntradaHistorial_JSON, err error) {
	historialPacienteIdString := strconv.Itoa(historialId)
	rows, err := db.Query(`SELECT id, empleado_id, historial_id, motivo_consulta, juicio_diagnostico, clave, clave_maestra, created_at, tipo FROM usuarios_entradas_historial where historial_id = ` + historialPacienteIdString + " order by created_at desc") // check err
	if err == nil {
		var e util.EntradaHistorial_JSON
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&e.Id, &e.EmpleadoId, &e.HistorialId, &e.MotivoConsulta, &e.JuicioDiagnostico, &e.Clave, &e.ClaveMaestra, &e.CreatedAt, &e.Tipo)
			//Cambio horario y formato
			words := strings.Fields(e.CreatedAt)
			day := words[0] + "T" + words[1] + "Z"
			layout := "2006-01-02T15:04:05.000000Z"
			t, err := time.Parse(layout, day)
			if err != nil {
				layout = "2006-01-02T15:04:05.00000Z"
				t, err = time.Parse(layout, day)
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

func GetClaveCompartidaEntradaHistorialByEntradaIdEmpleadoId(entradaId int, empleadoId int) (clave string, err error) {
	entradaIdString := strconv.Itoa(entradaId)
	empleadoIdString := strconv.Itoa(empleadoId)
	row, err := db.Query(`SELECT clave FROM usuarios_permisos_entradas_historial WHERE entrada_id = ` + entradaIdString + " AND empleado_id = " + empleadoIdString) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&clave)
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return clave, err
	}
	return clave, nil
}

func GetClaveCompartidaAnaliticaHistorialByEntradaIdEmpleadoId(analiticaId int, empleadoId int) (clave string, err error) {
	analiticaIdString := strconv.Itoa(analiticaId)
	empleadoIdString := strconv.Itoa(empleadoId)
	row, err := db.Query(`SELECT clave FROM usuarios_permisos_analiticas WHERE analitica_id = ` + analiticaIdString + " AND empleado_id = " + empleadoIdString) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&clave)
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return clave, err
	}
	return clave, nil
}

func GetAnaliticasHistorialByHistorialId(historialId int) (analiticas []util.AnaliticaHistorial_JSON, err error) {
	historialPacienteIdString := strconv.Itoa(historialId)
	rows, err := db.Query(`SELECT id, empleado_id, historial_id, leucocitos, hematies, plaquetas, glucosa, hierro, created_at, clave, clave_maestra FROM usuarios_analiticas where historial_id = ` + historialPacienteIdString + " order by created_at desc") // check err
	if err == nil {
		var a util.AnaliticaHistorial_JSON
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&a.Id, &a.EmpleadoId, &a.HistorialId, &a.Leucocitos, &a.Hematies, &a.Plaquetas, &a.Glucosa, &a.Hierro, &a.CreatedAt, &a.Clave, &a.ClaveMaestra)
			//Cambio horario y formato
			words := strings.Fields(a.CreatedAt)
			day := words[0] + "T" + words[1] + "Z"
			layout := "2006-01-02T15:04:05.000000Z"
			t, err := time.Parse(layout, day)
			if err != nil {
				layout = "2006-01-02T15:04:05.00000Z"
				t, err = time.Parse(layout, day)
			}
			t = t.Local()
			a.CreatedAt = fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d",
				t.Day(), t.Month(), t.Year(),
				t.Hour(), t.Minute(), t.Second())
			a.EmpleadoNombre, _ = GetNombreEmpleado(a.EmpleadoId)
			analiticas = append(analiticas, a)
		}
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return analiticas, err
	}
	return analiticas, nil
}

func GetEntradaById(entradaId int) (entrada util.EntradaHistorial_JSON, err error) {
	entradaIdString := strconv.Itoa(entradaId)
	row, err := db.Query(`SELECT id, empleado_id, historial_id, motivo_consulta, juicio_diagnostico, clave, clave_maestra, created_at, tipo FROM usuarios_entradas_historial where id = ` + entradaIdString) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&entrada.Id, &entrada.EmpleadoId, &entrada.HistorialId, &entrada.MotivoConsulta, &entrada.JuicioDiagnostico, &entrada.Clave, &entrada.ClaveMaestra, &entrada.CreatedAt, &entrada.Tipo)
		//Cambio horario y formato
		words := strings.Fields(entrada.CreatedAt)
		day := words[0] + "T" + words[1] + "Z"
		layout := "2006-01-02T15:04:05.000000Z"
		t, err := time.Parse(layout, day)
		if err != nil {
			layout = "2006-01-02T15:04:05.00000Z"
			t, err = time.Parse(layout, day)
		}
		t = t.Local()
		entrada.CreatedAt = fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d",
			t.Day(), t.Month(), t.Year(),
			t.Hour(), t.Minute(), t.Second())
		entrada.EmpleadoNombre, _ = GetNombreEmpleado(entrada.EmpleadoId)
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return entrada, err
	}
	return entrada, nil
}

//ANALITICAS

func InsertAnaliticaHistorial(analitica util.AnaliticaHistorial_JSON) (inserted_id int, err error) {
	createdAt := time.Now()
	historialPacienteIdString := strconv.Itoa(analitica.HistorialId)
	//INSERT
	analiticaId, err := db.Exec(`INSERT INTO usuarios_analiticas (empleado_id, historial_id, leucocitos, hematies, plaquetas, glucosa, hierro, clave, clave_maestra, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, analitica.UserToken.UserId,
		historialPacienteIdString, analitica.Leucocitos, analitica.Hematies, analitica.Plaquetas, analitica.Glucosa, analitica.Hierro, analitica.Clave, analitica.ClaveMaestra,
		createdAt.Local())
	if err == nil {
		id, _ := analiticaId.LastInsertId()
		inserted_id = int(id)
		return inserted_id, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return -1, err
	}
}

func InsertAnaliticaCompartidaHistorial(analitica util.AnaliticaHistorial_JSON) (inserted bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT IGNORE INTO usuarios_permisos_analiticas (analitica_id, empleado_id, clave) VALUES (?, ?, ?)`, analitica.Id, analitica.EmpleadoId, analitica.Clave)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func InsertEstadisticaAnaliticaHistorial(analitica util.AnaliticaHistorial_JSON) (inserted bool, err error) {
	hematies, _ := strconv.ParseFloat(analitica.Hematies, 32)
	hierro, _ := strconv.ParseFloat(analitica.Hierro, 32)
	leucocitos, _ := strconv.ParseFloat(analitica.Leucocitos, 32)
	plaquetas, _ := strconv.ParseFloat(analitica.Plaquetas, 32)
	glucosa, _ := strconv.ParseFloat(analitica.Glucosa, 32)

	//ID ALEATORIO DISTINTO AL ID DE LA ANALITICA
	id, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("uuid.NewV4() failed with %s\n", err)
	}
	uid := fmt.Sprint(id)

	//INSERT
	_, err = db.Exec(`INSERT INTO estadisticas_analiticas (id, leucocitos, hematies, plaquetas, glucosa, hierro) VALUES (?, ?, ?, ?, ?, ?)`, uid, leucocitos,
		hematies, plaquetas, glucosa, hierro)
	if err == nil {
		//AÑADIMOS TAGS DE LAS ANALITICAS
		for _, tagId := range analitica.Tags {
			_, err = db.Exec(`INSERT INTO estadisticas_analiticas_tags (analitica_id, tag_id) VALUES (?, ?)`, uid, tagId)
			if err != nil {

			}
		}
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func GetAnaliticaById(analiticaId int) (analitica util.AnaliticaHistorial_JSON, err error) {
	analiticaIdString := strconv.Itoa(analiticaId)
	row, err := db.Query(`SELECT id, empleado_id, historial_id, leucocitos, hematies, plaquetas, glucosa, hierro, clave, clave_maestra, created_at FROM usuarios_analiticas where id = ` + analiticaIdString) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&analitica.Id, &analitica.EmpleadoId, &analitica.HistorialId, &analitica.Leucocitos, &analitica.Hematies, &analitica.Plaquetas, &analitica.Glucosa, &analitica.Hierro, &analitica.Clave, &analitica.ClaveMaestra, &analitica.CreatedAt)
		//Cambio horario y formato
		words := strings.Fields(analitica.CreatedAt)
		day := words[0] + "T" + words[1] + "Z"
		layout := "2006-01-02T15:04:05.000000Z"
		t, err := time.Parse(layout, day)
		if err != nil {
			layout = "2006-01-02T15:04:05.00000Z"
			t, err = time.Parse(layout, day)
		}
		t = t.Local()
		analitica.CreatedAt = fmt.Sprintf("%02d-%02d-%02d %02d:%02d:%02d",
			t.Day(), t.Month(), t.Year(),
			t.Hour(), t.Minute(), t.Second())
		analitica.EmpleadoNombre, _ = GetNombreEmpleado(analitica.EmpleadoId)
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return analitica, err
	}
	return analitica, nil
}

func GetEstadisticasAnaliticas() (analiticas []util.AnaliticaHistorial_JSON, err error) {
	rows, err := db.Query(`SELECT id, leucocitos, hematies, plaquetas, glucosa, hierro FROM estadisticas_analiticas`) // check err
	if err == nil {
		var a util.AnaliticaHistorial_JSON
		defer rows.Close()
		for rows.Next() {
			var identificadorAnalitica string
			rows.Scan(&identificadorAnalitica, &a.Leucocitos, &a.Hematies, &a.Plaquetas, &a.Glucosa, &a.Hierro)
			a.EmpleadoNombre, _ = GetNombreEmpleado(a.EmpleadoId)
			a.Tags, _ = GetTagsIdEstadisticaAnalitica(identificadorAnalitica)
			analiticas = append(analiticas, a)
		}
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return analiticas, err
	}
	return analiticas, nil
}
