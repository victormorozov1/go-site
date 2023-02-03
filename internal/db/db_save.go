package db

import (
	"database/sql"
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

func keysValues(mp map[string]string, keys, values bool) map[string][]string {
	keysSlice := make([]string, len(mp))
	valuesSlice := make([]string, len(mp))
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
	return map[string][]string{
		"keys":   keysSlice,
		"values": valuesSlice,
	}
}

// func Keys[keyType constraints.Ordered](mp map[keyType]interface{}) []keyType {
func Keys(mp map[string]string) []string {
	return keysValues(mp, true, false)["keys"]
}

func Values(mp map[string]string) []string {
	return keysValues(mp, false, true)["values"]
}

func SaveToDB(db *sql.DB, tableName string, data map[string]string) error {
	keysValues := keysValues(data, true, true)
	names := "(`" + Join("`, `", keysValues["keys"]) + "`)"
	values := "('" + Join("', '", keysValues["values"]) + "'	)"

	request := "INSERT INTO `" + tableName + "` " + names + " VALUES" + values
	println(request)

	insert, err := db.Query(request) // add reservations
	if err != nil {
		return err
	}
	return insert.Close()
}

//func GetFromDB(db *sql.DB, tableName string, columns, []string) ([]interface{}, error) {
//	values := make([]interface{}, len(columns))
//
//	names := "`" + Join("`, `", columns) + "`"
//	request := "SELECT " + names + " FROM " + tableName
//	println(request)
//	result, err := db.Query(request)
//	if err != nil {
//		return values, err
//	}
//	defer result.Close()
//
//	for result.Next() {
//		err := result.Scan(&p.id, &p.model, &p.company, &p.price)
//		if err != nil {
//			fmt.Println(err)
//			continue
//		}
//		products = append(products, p)
//	}
//
//	return values, err
//}
