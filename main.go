package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	database "internal/db"
)

func printAllUsers(db *sql.DB) {
	println("All users:")
	allUsers, err := database.GetUsers(db)

	if err != nil {
		panic(err)
	}

	for _, user := range allUsers {
		println(user.String())
	}
}

func printReservationsByUser(db *sql.DB, userId int) {
	reservations, err := database.GetReservationsByUserId(db, 2)

	if err != nil {
		panic(err)
	}

	for _, r := range reservations {
		println(r.String())
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

	u := database.User{7, "Vova", "shelly", "sergeevich", "Brawler",
		"89124883837", "AMERICA@SOSYET.PUTINTOP228.ru", "/jkdfhjk", nil}
	u.SaveToDB(db)
	printAllUsers(db)

}
