package models

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"

	util "../utils"
)

var p = &util.Params_argon2{
	Memory:      64 * 1024,
	Iterations:  1,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

func LoginUser(userId int, password []byte) bool {
	user, err := GetUserById(userId)
	if err != nil {
		util.PrintErrorLog(err)
		return false
	}
	if string(user.Password) == "" {
		util.PrintLog("No hay contraseña")
		return false
	}
	match, err := util.Argon2comparePasswordAndHash(password, string(user.Password))
	if err != nil {
		util.PrintErrorLog(err)
	}
	return match
}

func InsertUser(user util.User_JSON) (userId int, err error) {
	//ARGON2
	encodedHash, err := util.Argon2generateFromPassword(user.Password, p)
	if err != nil {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return -1, err
	}
	//INSERT
	createdAt := time.Now()
	res, err := db.Exec(`INSERT INTO usuarios (dni, nombre, apellidos, email, password, created_at, clave, clave_maestra) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, user.Identificacion,
		user.Nombre, user.Apellidos, user.Email, encodedHash, createdAt.Local(), user.Clave, user.ClaveMaestra)
	if err == nil {
		userId, _ := res.LastInsertId()
		return int(userId), nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return -1, nil
}

func InsertUserDniHash(user_id int, user_dni string) (inserted bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT INTO usuarios_dnihashes (usuario_id, dni_hash) VALUES (?, ?)`, user_id, user_dni)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return false, err
}

func CheckUserDniHash(user_dni string) (usuarioId int, err error) {
	//SHA 256, cogemos la primera mitad
	sha_256 := sha256.New()
	sha_256.Write([]byte(user_dni))
	hash := sha_256.Sum(nil)
	stringHash := fmt.Sprintf("%x", hash) //Pasamos a hexadecimal el hash
	var bdString string
	//INSERT
	row, err := db.Query(`SELECT dni_hash, usuario_id FROM usuarios_dnihashes where dni_hash = '` + stringHash + `'`) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&bdString, &usuarioId)
		if stringHash == bdString {
			return usuarioId, nil
		}
		return -1, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return -1, err
	}
}

func EditUserData(user util.User_JSON) (edited bool, err error) {
	//Editar los DATOS del usuario
	//UPDATE

	_, err = db.Exec(`UPDATE usuarios set dni = ?, nombre = ?, apellidos = ?, email = ? where dni = ?`, user.Identificacion,
		user.Nombre, user.Apellidos, user.Email, user.Identificacion)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return false, nil
}

func DeleteUser(user_id int) (deleted bool, err error) {
	_, err = db.Exec(`DELETE FROM usuarios where id = ` + strconv.Itoa(user_id))
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return false, err
}

func GetUsersList() (usersList []util.User, err error) {
	rows, err := db.Query("SELECT id, dni, nombre, apellidos, email, created_at FROM usuarios")
	if err == nil {
		defer rows.Close()
		var users []util.User
		for rows.Next() {
			var u util.User
			rows.Scan(&u.Id, &u.Identificacion, &u.Nombre, &u.Apellidos, &u.Email, &u.CreatedAt)
			users = append(users, u)
		}
		return users, err
	} else {
		fmt.Println(err)
		return nil, err
	}
}

func GetUserById(id int) (user util.User_JSON, err error) {
	row, err := db.Query(`SELECT id, dni, nombre, apellidos, email, created_at, clave, password FROM usuarios where id = ` + strconv.Itoa(id)) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&user.Id, &user.Identificacion, &user.Nombre, &user.Apellidos, &user.Email, &user.CreatedAt, &user.Clave, &user.Password)
		return user, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return user, err
	}
}

func GetUserByIdentificacion(Identificacion string) (user util.User_JSON, err error) {
	row, err := db.Query(`SELECT u.id, dni, nombre, apellidos, email, clave FROM usuarios_dnihashes ud, usuarios u where dni_hash = '` + Identificacion + `' and u.id = ud.usuario_id`) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&user.Id, &user.Identificacion, &user.Nombre, &user.Apellidos, &user.Email, &user.Clave)
		return user, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return user, err
	}
}

func GetUserByEmail(email string) (user util.User_JSON, err error) {
	row, err := db.Query(`SELECT id, dni, nombre, apellidos, email, password, created_at, clave FROM usuarios where email = '` + email + `'`) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&user.Id, &user.Identificacion, &user.Nombre, &user.Apellidos, &user.Email, &user.Password, &user.CreatedAt, &user.Clave)
		return user, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return user, err
	}
}

func GetUsersPagination(page int) []util.User_JSON {
	firstRow := strconv.Itoa(page * 10)
	lastRow := strconv.Itoa((page * 10) + 10)
	rows, err := db.Query("SELECT id, dni, nombre, apellidos, email, created_at, clave, clave_maestra FROM usuarios LIMIT " + firstRow + "," + lastRow)
	if err == nil {
		defer rows.Close()
		var users []util.User_JSON
		for rows.Next() {
			var u util.User_JSON
			rows.Scan(&u.Id, &u.Identificacion, &u.Nombre, &u.Apellidos, &u.Email, &u.CreatedAt, &u.Clave, &u.ClaveMaestra)
			users = append(users, u)
		}
		return users
	} else {
		fmt.Println(err)
		return nil
	}
}

//GESTION DE TOKEN DEL USUARIO

func ProveUserToken(user_id int, token string) (result bool, err error) {
	timeNow := time.Now().Local().UTC()
	var idToken int
	row, err := db.Query(`SELECT id, token, fecha_expiracion FROM usuarios_tokens where usuario_id = ` + strconv.Itoa(user_id))
	if err != nil {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
	defer row.Close()
	var tokenBD string
	var fechaExpiracionBD time.Time
	for row.Next() {
		row.Scan(&idToken, &tokenBD, &fechaExpiracionBD)
	}
	if tokenBD == token && timeNow.Before(fechaExpiracionBD) {
		//Reseteamos el tiempo
		timeNow = time.Now().Local().Add(time.Minute * time.Duration(30))
		_, err = db.Exec(`UPDATE usuarios_tokens SET fecha_expiracion = (?) WHERE id = (?)`, timeNow, idToken)
		return true, nil
	}
	return false, nil
}

func InsertUserToken(user_id int) (token string, err error) {
	//INSERT
	//30 minutos de tiempo
	timeNow := time.Now().Local().Add(time.Minute * time.Duration(30))
	//Generamos token
	token, err = util.GenerateRandomString(156)
	if err != nil {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return "", err
	}

	//Primero borramos el token ya existente
	_, err = db.Exec(`DELETE FROM usuarios_tokens where usuario_id = ` + strconv.Itoa(user_id))
	if err != nil {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return "", err
	}
	//Luego insertamos uno nuevo
	_, err = db.Exec(`INSERT INTO usuarios_tokens (usuario_id, token, fecha_expiracion) VALUES (?, ?, ?)`,
		user_id, token, timeNow)
	if err == nil {
		return token, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return "", err
	}
}

//PAIRKEYS

func InsertUserPairKeys(user_id int, pairKeys util.PairKeys) (result bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT INTO usuarios_pairkeys (usuario_id, public_key, private_key) VALUES (?, ?, ?)`, user_id,
		pairKeys.PublicKey, pairKeys.PrivateKey)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return false, nil
}

func InsertUserMasterPairKeys(user_id int, pairKeys util.PairKeys) (result bool, err error) {
	//INSERT
	_, err = db.Exec(`INSERT INTO usuarios_master_pairkeys (usuario_id, public_key, private_key) VALUES (?, ?, ?)`, user_id,
		pairKeys.PublicKey, pairKeys.PrivateKey)
	if err == nil {
		return true, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return false, nil
}

func GetUserPairKeys(user_id string) (result util.PairKeys, err error) {
	//GET
	row, err := db.Query(`SELECT public_key, private_key FROM usuarios_pairkeys WHERE usuario_id = ` + user_id)
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&result.PublicKey, &result.PrivateKey)
		return result, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return result, nil
}

func GetUserPairKeysByHistorialId(historial_id string) (result util.PairKeys, err error) {
	historial_idString, _ := strconv.Atoi(historial_id)
	historial, _ := GetHistorialById(historial_idString)
	user_id := strconv.Itoa(historial.PacienteId)
	row, err := db.Query(`SELECT public_key, private_key FROM usuarios_pairkeys WHERE usuario_id = ` + user_id)
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&result.PublicKey, &result.PrivateKey)
		return result, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return result, nil
}

func GetUserMasterPairKeys(user_id string) (result util.PairKeys, err error) {
	//GET
	row, err := db.Query(`SELECT public_key, private_key FROM usuarios_master_pairkeys WHERE usuario_id = ` + user_id)
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&result.PublicKey, &result.PrivateKey)
		return result, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return result, nil
}

func GetPublicMasterKey() (result util.PairKeys, err error) {
	//GET
	row, err := db.Query(`SELECT public_key FROM usuarios_master_pairkeys`)
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&result.PublicKey)
		return result, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return result, nil
}

func GetUserPublicKey(user_id string) (result util.PairKeys, err error) {
	//GET
	row, err := db.Query(`SELECT public_key FROM usuarios_pairkeys WHERE usuario_id = ` + user_id)
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&result.PublicKey)
		return result, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return result, nil
}

func GetUserPrivateKey(user_id string) (result util.PairKeys, err error) {
	//GET
	row, err := db.Query(`SELECT private_key FROM usuarios_pairkeys WHERE usuario_id = ` + user_id)
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&result.PrivateKey)
		return result, nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return result, nil
}

/*
func GetUsuariosPorNombreODni(nombreApellidos string, dni string) (user util.User, err error) {
	if nombreApellidos != "" {
		row, err := db.Query(`SELECT id, dni, nombre, apellidos, email, created_at FROM usuarios where id = ` + strconv.Itoa(id))
	} else {

	}
	row, err := db.Query(`SELECT id, dni, nombre, apellidos, email, created_at FROM usuarios where id = ` + strconv.Itoa(id)) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&user.Id, &user.Identificacion, &user.Nombre, &user.Apellidos, &user.Email, &user.CreatedAt)
		return user, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return user, err
	}
}
*/
