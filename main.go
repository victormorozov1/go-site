package main

import (
	_ "github.com/go-sql-driver/mysql"
	database "internal/db"
	"internal/server"
)

func startServer() {
	s := server.Server{
		Host:       "127.0.0.1",
		Port:       8080,
		CookieName: "data",
		Roles: server.Roles{
			Admin: "admin",
			User:  "user",
		},
		Routes: &server.Routes{
			MainPage:                     "/main",
			AllUsersPage:                 "/users",
			UserCabinet:                  "/me",
			RegisterPage:                 "/register",
			LoginPage:                    "/login",
			TestPage:                     "/test",
			ReservationPage:              "/reservation",
			DeleteReservationAjaxHandler: "/delete_reservation",
			MapPage:                      "/map",
			WorkplacePage:                "/workplace",
		},
		DataBase: database.Database{
			Name:     "go_site",
			Driver:   "mysql",
			Login:    "root",
			Password: "",
			Url:      "127.0.0.1:3306",
			Tables: database.Tables{
				Users: database.UsersTable{
					TableName:      "users",
					Id:             "id",
					Name:           "name",
					Surname:        "surname",
					Patronymic:     "patronymic",
					Role:           "role",
					Phone:          "phone",
					Email:          "email",
					PhotoSrc:       "photo_src",
					HashedPassword: "hashed_password",
				},
				Reservations: database.ReservationsTable{
					TableName: "reservations2",
					Id:        "id",
					UserId:    "user_id",
					TableId:   "table_id",
					StartTime: "start_time",
					EndTime:   "end_time",
				},
				Tables: database.TablesTable{
					TableName:          "tables",
					Id:                 "id",
					Description:        "description",
					TechnicalEquipment: "technical_equipment",
					Position:           "position",
					Hide:               "hide",
				},
			},
		},
	}
	s.CountBaseTemplateData()

	s.Start()
}

func main() {
	startServer()
}
