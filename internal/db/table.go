package db

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type Table struct {
	Id                              int
	Description, TechnicalEquipment string
	X, Y                            int
	Hide                            bool
	Users                           []*User
	Reservations                    []*Reservation
}

func (table Table) String() string {
	return fmt.Sprintf("Table#%d(%s, %s, pos=(%d, %d), Hidden=%v, Users=%v, Reservations=%v",
		table.Id, table.Description, table.TechnicalEquipment, table.X, table.Y, table.Hide, table.Users, table.Reservations)
}

func (table *Table) SaveToDB(db *Database) error {
	return SaveToDB(db.Connection, db.Tables.Tables.TableName, map[string]string{ // тут бы сделать interface{}
		db.Tables.Tables.Id:                 strconv.Itoa(table.Id),
		db.Tables.Tables.Description:        table.Description,
		db.Tables.Tables.TechnicalEquipment: table.TechnicalEquipment,
		db.Tables.Tables.Position:           XYPositionToText(table.X, table.Y),
		db.Tables.Tables.Hide:               strconv.FormatBool(table.Hide),
	})
}

func scanTablesFromDBRows(rows *sql.Rows) (*Table, error) {
	var newTable = Table{}
	var textPosition string
	err := rows.Scan(&newTable.Id, &newTable.Description, &newTable.TechnicalEquipment, &textPosition, newTable.Hide)

	newTable.X, newTable.Y, err = textPositionToXY(textPosition)
	return &newTable, err
}

func GetTables(db *Database, criteria string) ([]*Table, error) {
	if criteria != "" {
		criteria = " WHERE " + criteria
	}
	return getFromDB(db.Connection, "select * from "+db.Tables.Tables.TableName+criteria, scanTablesFromDBRows)
}

func GetAllTables(db *Database) ([]*Table, error) {
	return GetTables(db, "")
}

func GetTablesBy(db *Database, columnName, columnValue string) ([]*Table, error) {
	return GetTables(db, columnName+" = "+columnValue)
}

func GetTableById(db *Database, id int) (*Table, error) {
	tables, err := GetTablesBy(db, "id", strconv.Itoa(id))
	if err != nil {
		return nil, err
	}

	if len(tables) == 0 {
		return nil, errors.New("table not found")
	}

	if len(tables) > 1 {
		return nil, errors.New("Found many tables with id=" + strconv.Itoa(id))
	}

	return tables[0], nil
}

func (table *Table) LoadReservations(db *Database) error {
	reservations, err := GetReservationsBy(db, "table_id", strconv.Itoa((table.Id)))
	if err != nil {
		return err
	}
	table.Reservations = reservations
	return nil
}

//func (table *Table) LoadUsers(db *sql.DB, reservationsTableName string) error {
//	err := table.LoadReservations(db, reservationsTableName)
//	if err != nil {
//		return err
//	}
//
//}
