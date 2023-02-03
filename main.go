package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	database "internal/db"
)

func main() {
	//s := server.Server{"127.0.0.1", 8080}
	//s.Start()
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_site")

	if err != nil {
		panic(err)
	}

	user := database.User{4, "Victor223", "Morov", "Victorovich", "admin", "89228311", "kakaska@ya.ru", "-", nil}
	err = user.SaveToDB(db)

	if err != nil {
		panic(err)
	}

	reservation := database.Reservation{4, 2, 3, 4, 5}
	err = reservation.SaveToDB(db)
	if err != nil {
		panic(err)
	}

	//print(database.GetUsers(db))

}
