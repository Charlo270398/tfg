package models

import (
	"fmt"
	"strconv"

	util "../utils"
)

func InsertClinica(clinica util.Clinica_JSON) (clinicaId int, err error) {
	//INSERT
	res, err := db.Exec(`INSERT INTO clinicas (nombre, direccion, telefono) VALUES (?,?,?)`, clinica.Nombre, clinica.Direccion, clinica.Telefono)
	if err == nil {
		clinicaId, _ := res.LastInsertId()
		return int(clinicaId), nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return -1, err
	}
}

func EditClinicaData(clinica util.Clinica_JSON) (edited bool, err error) {
	//UPDATE
	_, err = db.Exec(`UPDATE clinicas set nombre = ?, direccion = ?, telefono = ? where id = ?`, clinica.Nombre, clinica.Direccion, clinica.Telefono, clinica.Id)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return false, nil
}

func DeleteClinica(clinica_id int) bool {
	_, err := db.Exec(`DELETE FROM clinicas where id = ` + strconv.Itoa(clinica_id))
	if err == nil {
		return true
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return false
}

func GetClinicaList() (clinicaList []util.Clinica, err error) {
	rows, err := db.Query("SELECT id, nombre, direccion, telefono FROM clinicas")
	if err == nil {
		defer rows.Close()
		var clinicas []util.Clinica
		for rows.Next() {
			var c util.Clinica
			rows.Scan(&c.Id, &c.Nombre, &c.Direccion, &c.Telefono)
			clinicas = append(clinicas, c)
		}
		return clinicas, err
	} else {
		fmt.Println(err)
		return nil, err
	}
}

func GetClinicaPagination(page int) []util.Clinica {
	firstRow := strconv.Itoa(page * 10)
	lastRow := strconv.Itoa((page * 10) + 10)
	rolEnfermero := strconv.Itoa(Rol_enfermero.Id)
	rolMedico := strconv.Itoa(Rol_medico.Id)
	rolAdmin := strconv.Itoa(Rol_administradorC.Id)
	rows, err := db.Query("SELECT id, nombre, direccion, telefono FROM clinicas LIMIT " + firstRow + "," + lastRow)
	if err == nil {
		defer rows.Close()
		var clinicas []util.Clinica
		for rows.Next() {
			var c util.Clinica
			rows.Scan(&c.Id, &c.Nombre, &c.Direccion, &c.Telefono)
			clinicaId := strconv.Itoa(c.Id)
			rowsUsarioClinica, _ := db.Query("SELECT count(*) from usuarios_clinicas where clinica_id = " + clinicaId + " and rol_id = " + rolEnfermero)
			rowsUsarioClinica.Next()
			rowsUsarioClinica.Scan(&c.NumeroEnfermeros)
			rowsUsarioClinica, _ = db.Query("SELECT count(*) from usuarios_clinicas where clinica_id = " + clinicaId + " and rol_id = " + rolMedico)
			rowsUsarioClinica.Next()
			rowsUsarioClinica.Scan(&c.NumeroMedicos)
			rowsUsarioClinica, _ = db.Query("SELECT count(*) from usuarios_clinicas where clinica_id = " + clinicaId + " and rol_id = " + rolAdmin)
			rowsUsarioClinica.Next()
			rowsUsarioClinica.Scan(&c.NumeroAdministradores)
			clinicas = append(clinicas, c)
		}
		fmt.Println(clinicas)
		return clinicas
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return nil
	}
}

func InsertarUserClinica(clinica_id int, usuario_id int, rol_id int) (result bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT INTO usuarios_clinicas (usuario_id, clinica_id, rol_id) VALUES (?,?,?)`, usuario_id, clinica_id, rol_id)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
}

func GetClinicaByAdmin(user_id string) (c util.Clinica, err error) {
	rolEnfermero := strconv.Itoa(Rol_enfermero.Id)
	rolMedico := strconv.Itoa(Rol_medico.Id)
	rolAdmin := strconv.Itoa(Rol_administradorC.Id)
	rows, err := db.Query("SELECT id, nombre, direccion, telefono FROM clinicas c, usuarios_clinicas u " +
		"WHERE c.id = u.clinica_id and usuario_id = " + user_id + " and rol_id = " + rolAdmin)
	if err == nil {
		defer rows.Close()
		var c util.Clinica
		rows.Next()
		rows.Scan(&c.Id, &c.Nombre, &c.Direccion, &c.Telefono)
		clinicaId := strconv.Itoa(c.Id)
		rowsUsarioClinica, _ := db.Query("SELECT count(*) from usuarios_clinicas where clinica_id = " + clinicaId + " and rol_id = " + rolEnfermero)
		rowsUsarioClinica.Next()
		rowsUsarioClinica.Scan(&c.NumeroEnfermeros)
		rowsUsarioClinica, _ = db.Query("SELECT count(*) from usuarios_clinicas where clinica_id = " + clinicaId + " and rol_id = " + rolMedico)
		rowsUsarioClinica.Next()
		rowsUsarioClinica.Scan(&c.NumeroMedicos)
		rowsUsarioClinica, _ = db.Query("SELECT count(*) from usuarios_clinicas where clinica_id = " + clinicaId + " and rol_id = " + rolAdmin)
		rowsUsarioClinica.Next()
		rowsUsarioClinica.Scan(&c.NumeroAdministradores)
		return c, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return c, err
	}
}
