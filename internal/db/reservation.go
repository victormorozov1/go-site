package db

import (
	"database/sql"
	"fmt"
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

func GetReservationsByUserId(db *sql.DB, userId int) ([]*Reservation, error) {
	f := func(rows *sql.Rows) (*Reservation, error) {
		var reservation = Reservation{}

		err := rows.Scan(&reservation.Id, &reservation.User_id, &reservation.Table_id,
			&reservation.Start_time, &reservation.End_time)

		return &reservation, err
	}
	return getFromDB(db, "select * from "+reservationsTable+" where user_id="+strconv.Itoa(userId), f)
}

func (r Reservation) String() string {
	return fmt.Sprintf("Reservation(id=%d user_id=%d table_id=%d start_time=%d end_time=%d",
		r.Id, r.User_id, r.Table_id, r.Start_time, r.End_time)
}
