package db

import (
	"database/sql"
	"internal/functions"
)

func SaveToDB(db *sql.DB, tableName string, data map[string]string) error {
	keys, values := functions.KeysValues(data, true, true)
	columns := "(`" + functions.Join("`, `", keys) + "`)"
	columnsValues := "('" + functions.Join("', '", values) + "'	)"

	request := "INSERT INTO `" + tableName + "` " + columns + " VALUES" + columnsValues
	println(request)

	insert, err := db.Query(request) // add reservations
	if err != nil {
		return err
	}
	return insert.Close()
}
