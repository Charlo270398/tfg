package models

import (
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

func LoginUser(email string, password []byte) bool {
	user, err := GetUserByEmail(email)
	if err != nil {
		util.PrintErrorLog(err)
		return false
	}
	if user.Password == "" {
		util.PrintLog("El usuario no existe")
		return false
	}
	match, err := util.Argon2comparePasswordAndHash(password, user.Password)
	if err != nil {
		util.PrintErrorLog(err)
	}
	return match
}

func InsertUser(user util.User_JSON) (userId int, err error) {
	//ARGON2
	encodedHash, err := util.Argon2generateFromPassword(user.Password, p)
	util.PrintLog(encodedHash)
	if err != nil {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return -1, err
	}
	//INSERT
	createdAt := time.Now()
	res, err := db.Exec(`INSERT INTO usuarios (dni, nombre, apellidos, email, password, created_at) VALUES (?, ?, ?, ?, ?, ?)`, user.Identificacion,
		user.Nombre, user.Apellidos, user.Email, encodedHash, createdAt)
	if err == nil {
		userId, _ := res.LastInsertId()
		return int(userId), nil
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
	}
	return -1, nil
}

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

func GetUserById(id int) (user util.User, err error) {
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

func GetUserByEmail(email string) (user util.User, err error) {
	row, err := db.Query(`SELECT id, dni, nombre, apellidos, email, password, created_at FROM usuarios where email = '` + email + `'`) // check err
	if err == nil {
		defer row.Close()
		row.Next()
		row.Scan(&user.Id, &user.Identificacion, &user.Nombre, &user.Apellidos, &user.Email, &user.Password, &user.CreatedAt)
		return user, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return user, err
	}
}

func GetUsersPagination(page int) []util.User {
	firstRow := strconv.Itoa(page * 10)
	lastRow := strconv.Itoa((page * 10) + 10)
	rows, err := db.Query("SELECT id, dni, nombre, apellidos, email, created_at FROM usuarios LIMIT " + firstRow + "," + lastRow)
	if err == nil {
		defer rows.Close()
		var users []util.User
		for rows.Next() {
			var u util.User
			rows.Scan(&u.Id, &u.Identificacion, &u.Nombre, &u.Apellidos, &u.Email, &u.CreatedAt)
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
	row, err := db.Query(`SELECT token, fecha_expiracion FROM usuarios_tokens where usuario_id = ` + strconv.Itoa(user_id))
	if err != nil {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return false, err
	}
	defer row.Close()
	var tokenBD string
	var fechaExpiracionBD time.Time
	for row.Next() {
		row.Scan(&tokenBD, &fechaExpiracionBD)
	}
	if tokenBD == token && timeNow.Before(fechaExpiracionBD) {
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
	return "", nil
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
