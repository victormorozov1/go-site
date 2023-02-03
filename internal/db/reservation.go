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
