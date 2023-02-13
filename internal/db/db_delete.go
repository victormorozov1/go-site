package db

import (
	"database/sql"
	"strconv"
)

func Delete(db *sql.DB, tableName, condition string) {
	request := "DELETE FROM " + tableName + " WHERE " + condition
	println(request)
	db.Exec(request)
}

func DeleteBy(db *sql.DB, tableName, columnName, value string) {
	Delete(db, tableName, columnName+"="+value)
}

func DeleteById(db *sql.DB, tableName string, id int) {
	DeleteBy(db, tableName, "id", strconv.Itoa(id))
}
