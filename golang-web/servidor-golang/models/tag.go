package models

import (
	"fmt"

	util "../utils"
)

func GetTagsList() (tags []util.Tag_JSON, err error) {
	row, err := db.Query(`SELECT id, nombre FROM tags ORDER BY nombre`) // check err
	if err == nil {
		defer row.Close()
		for row.Next() {
			var t util.Tag_JSON
			row.Scan(&t.Id, &t.Nombre)
			tags = append(tags, t)
		}
		return tags, err
	} else {
		fmt.Println(err)
		util.PrintErrorLog(err)
		return nil, err
	}
}
