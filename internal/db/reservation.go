package db

import (
	"database/sql"
	"strconv"
)

const reservationsTable = "reservations2"

type Reservation struct {
	Id, Table_id, User_id int
	Start_time, End_time  int // Потом сделать Time
}

func (reservation *Reservation) SaveToDB(db *sql.DB) error {
	return SaveToDB(db, reservationsTable, map[string]string{
		"id":         strconv.Itoa(reservation.Id),
		"table_id":   strconv.Itoa(reservation.Table_id),
		"user_id":    strconv.Itoa(reservation.User_id),
		"start_time": strconv.Itoa(reservation.Start_time),
		"end_time":   strconv.Itoa(reservation.End_time),
	})
}

func GetReservationsByUserId(db *sql.DB, userId int) {
	//result, err := db.Query("select * from users2")
	//
	//var all_users []*User
	//
	//if err != nil {
	//	return err, all_users
	//}
	//
	//for result.Next() {
	//	var new_user = User{}
	//
	//	err := result.Scan(&new_user.Id, &new_user.Name, &new_user.Surname, &new_user.Patronymic,
	//		&new_user.Role, &new_user.Phone, &new_user.Email, &new_user.Photo_src)
	//
	//	if err != nil {
	//		return err, all_users
	//	}
	//
	//	all_users = append(all_users, &new_user)
	//}
	//
	//return nil, all_users
}
