package db

import "database/sql"

func getFromDB[elementType interface{}](db *sql.DB, request string, scanElement func(rows *sql.Rows) (elementType, error)) ([]elementType, error) {
	println(request)
	result, err := db.Query(request)

	var all_users []elementType

	if err != nil {
		return all_users, err
	}

	for result.Next() {
		newElement, err := scanElement(result)

		if err != nil {
			return all_users, err
		}

		all_users = append(all_users, newElement)
	}

	return all_users, nil
}
