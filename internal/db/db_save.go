package db

import (
	"database/sql"
	"golang.org/x/exp/constraints"
)

func Join(sep string, arr []string) string {
	var s = ""
	for ind, el := range arr {
		if ind != 0 {
			s += sep
		}
		s += el
	}
	return s
}

func KeysValues[keyType constraints.Ordered, valueType interface{}](mp map[keyType]valueType, keys, values bool) ([]keyType, []valueType) {
	keysSlice := make([]keyType, len(mp))
	valuesSlice := make([]valueType, len(mp))
	i := 0
	for key, value := range mp {
		if keys {
			keysSlice[i] = key
		}
		if values {
			valuesSlice[i] = value
		}
		i++
	}
	return keysSlice, valuesSlice
}

func Keys[keyType constraints.Ordered, valueType interface{}](mp map[keyType]valueType) []keyType {
	keys, _ := KeysValues(mp, true, false)
	return keys
}

func Values(mp map[string]string) []string {
	_, values := KeysValues(mp, false, true)
	return values
}

func SaveToDB(db *sql.DB, tableName string, data map[string]string) error {
	keys, values := KeysValues(data, true, true)
	columns := "(`" + Join("`, `", keys) + "`)"
	columnsValues := "('" + Join("', '", values) + "'	)"

	request := "INSERT INTO `" + tableName + "` " + columns + " VALUES" + columnsValues
	println(request)

	insert, err := db.Query(request) // add reservations
	if err != nil {
		return err
	}
	return insert.Close()
}
