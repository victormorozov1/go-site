package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	database "internal/db"
	db2 "internal/db"
	"internal/server"
)

//func printAllUsers(db *sql.DB) {
//	println("All users:")
//	allUsers, err := database.GetUsers(db, "users2", "reservations2")
//
//	if err != nil {
//		panic(err)
//	}
//
//	for _, user := range allUsers {
//		println(user.String())
//	}
//}

func printReservationsByUser(db *sql.DB, userId int) {
	reservations, err := database.GetReservationsByUserId(db, "reservations2", 2)

	if err != nil {
		panic(err)
	}

	for _, r := range reservations {
		println(r.String())
	}
}

func mainProgram() {
	s := server.Server{
		Host:                  "127.0.0.1",
		Port:                  8080,
		UsersTableName:        "users",
		ReservationsTableName: "reservations2",
		CookieName:            "data",
		Routes: &server.Routes{
			MainPage:        "/main",
			AllUsersPage:    "/users",
			UserCabinet:     "/me",
			RegisterPage:    "/register",
			LoginPage:       "/login",
			TestPage:        "/test",
			ReservationPage: "/reservation",
		},
	}
	s.CountBaseTemplateData()

	s.Start()
}

func main() {

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/go_site")

	if err != nil {
		panic(err)
	}
	defer db.Close()

	r, err := db2.GetReservationById(db, "reservations2", 7)
	print(r)
	r.Delete(db, "reservations2")
}
