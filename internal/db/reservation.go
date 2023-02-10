package db

import (
	"database/sql"
	"fmt"
	"strconv"
)

type Reservation struct {
	Id, Table_id, User_id int
	Start_time, End_time  int // Потом сделать Time
	User                  *User
}

func (reservation *Reservation) SaveToDB(db *sql.DB, reservationTableName string) error {
	return SaveToDB(db, reservationTableName, map[string]string{
		"id":         strconv.Itoa(reservation.Id),
		"table_id":   strconv.Itoa(reservation.Table_id),
		"user_id":    strconv.Itoa(reservation.User_id),
		"start_time": strconv.Itoa(reservation.Start_time),
		"end_time":   strconv.Itoa(reservation.End_time),
	})
}

func scanReservation(rows *sql.Rows) (*Reservation, error) {
	var reservation = Reservation{}

	err := rows.Scan(&reservation.Id, &reservation.User_id, &reservation.Table_id,
		&reservation.Start_time, &reservation.End_time)

	return &reservation, err
}

func GetReservationsByUserId(db *sql.DB, reservationTableName string, userId int) ([]*Reservation, error) {
	return getFromDB(db, "select * from "+reservationTableName+" where user_id="+strconv.Itoa(userId), scanReservation)
}

func GetReservationById(db *sql.DB, reservationTableName string, id int) (*Reservation, error) {
	reservations, err := getFromDB(db, "select * from "+reservationTableName+" where id="+strconv.Itoa(id), scanReservation)
	if err != nil {
		return nil, err
	}
	if len(reservations) == 0 {
		return nil, &Error{"Reservation #" + strconv.Itoa(id) + " not found"}
	}
	if len(reservations) > 1 {
		panic("many reservations with id = " + strconv.Itoa(id))
	}
	return reservations[0], nil
}

func (reservation Reservation) String() string {
	s := fmt.Sprintf("Reservation(id=%d user_id=%d table_id=%d start_time=%d end_time=%d",
		reservation.Id, reservation.User_id, reservation.Table_id, reservation.Start_time, reservation.End_time)
	if reservation.User != nil {
		s += ", " + reservation.User.String()
	}
	s += ")"
	return s
}

func (reservation *Reservation) LoadUser(db *sql.DB, usersTableName, reservationTableName string) error {
	// Убрать цикличность
	userId := strconv.Itoa(reservation.User_id)
	users, err := GetUserBy(db, "id", userId, usersTableName, reservationTableName)
	if err != nil {
		return err
	}
	reservation.User = users[0] // Создать функцию getUserById и ипользовать здесь ее
	return nil
}
