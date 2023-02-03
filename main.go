package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	database "internal/db"
)

func printAllUsers(db *sql.DB) {
	println("All users:")
	err, allUsers := database.GetUsers(db)

	if err != nil {
		panic(err)
	}

	for _, user := range allUsers {
		println(user.String())
	}
}

func main() {
	//s := server.Server{"127.0.0.1", 8080}
	//s.Start()
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_site")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	printAllUsers(db)

}
