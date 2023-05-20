package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type Reservation struct {
	Id, Table_id, User_id int
	Start_time, End_time  int // Потом сделать Time
	User                  *User
	Table                 *Table
	//добавить поле стола
}

func (reservation *Reservation) SaveToDB(db *Database) error {
	return SaveToDB(db.Connection, db.Tables.Reservations.TableName, map[string]string{
		db.Tables.Reservations.Id:        strconv.Itoa(reservation.Id),
		db.Tables.Reservations.TableId:   strconv.Itoa(reservation.Table_id),
		db.Tables.Reservations.UserId:    strconv.Itoa(reservation.User_id),
		db.Tables.Reservations.StartTime: strconv.Itoa(reservation.Start_time),
		db.Tables.Reservations.EndTime:   strconv.Itoa(reservation.End_time),
	})
}

func scanReservationFromDBRows(rows *sql.Rows) (*Reservation, error) {
	var reservation = Reservation{}

	err := rows.Scan(&reservation.Id, &reservation.User_id, &reservation.Table_id,
		&reservation.Start_time, &reservation.End_time)

	return &reservation, err
}

func GetReservations(db *Database, criteria string) ([]*Reservation, error) {
	if criteria != "" {
		criteria = "WHERE " + criteria
	}
	return getFromDB(db.Connection, "select * from "+db.Tables.Reservations.TableName+" "+criteria, scanReservationFromDBRows)
}

func GetAllReservations(db *Database) ([]*Reservation, error) {
	return GetReservations(db, "")
}

func GetReservationsBy(db *Database, columnName, columnValue string) ([]*Reservation, error) {
	return GetReservations(db, columnName+" = "+columnValue)
}

func GetReservationsByUserId(db *Database, userId int) ([]*Reservation, error) {
	return GetReservationsBy(db, db.Tables.Reservations.UserId, strconv.Itoa(userId))
}

func GetReservationById(db *Database, id int) (*Reservation, error) {
	reservations, err := getFromDB(db.Connection, "select * from "+db.Tables.Reservations.TableName+
		" where "+db.Tables.Reservations.Id+"="+strconv.Itoa(id),
		scanReservationFromDBRows)
	if err != nil {
		return nil, err
	}
	if len(reservations) == 0 {
		return nil, errors.New("Reservation #" + strconv.Itoa(id) + " not found")
	}
	if len(reservations) > 1 {
		panic("many reservations with id = " + strconv.Itoa(id))
	}
	return reservations[0], nil
}

func (reservation Reservation) String() string {
	s := fmt.Sprintf("Reservation(id=%d user_id=%d table_id=%d start_time=%d end_time=%d table=%s",
		reservation.Id, reservation.User_id, reservation.Table_id, reservation.Start_time, reservation.End_time, reservation.Table.String())
	if reservation.User != nil {
		s += ", " + reservation.User.String()
	}
	s += ")"
	return s
}

func (reservation *Reservation) LoadUser(db *Database) error {
	// Убрать цикличность
	user, err := GetUserById(db, reservation.User_id)
	if err != nil {
		return err
	}
	reservation.User = user
	return nil
}

func (reservation *Reservation) LoadTable(db *Database) error {
	table, err := GetTableById(db, reservation.Table_id)
	if err != nil {
		return err
	}
	reservation.Table = table
	return nil
}

func (reservation *Reservation) Delete(db *Database) {
	DeleteById(db.Connection, db.Tables.Reservations.TableName, reservation.Id)
}

func ReservationOverlaps(r1, r2 *Reservation) bool {
	if r1.Table_id != r2.Table_id {
		return false
	}

	between := func(a, b, c int) bool {
		return a >= b && a <= c
	}
	return between(r1.Start_time, r2.Start_time, r2.End_time) || between(r2.Start_time, r1.Start_time, r1.End_time)
}

func (reservation *Reservation) OverlapsOther(allReservations []*Reservation) (bool, int) {
	for _, r := range allReservations {
		if r.Id != reservation.Id {
			if ReservationOverlaps(reservation, r) {
				return true, r.Id
			}
		}
	}
	return false, -1
}

func (reservation *Reservation) Check(allReservations []*Reservation) error {
	if reservation.Start_time >= reservation.End_time {
		return errors.New("start_time >= end_time")
	}

	overlap, ind := reservation.OverlapsOther(allReservations)
	if overlap {
		return errors.New("Overlaps with reservation#" + strconv.Itoa(ind))
	}

	return nil
}
