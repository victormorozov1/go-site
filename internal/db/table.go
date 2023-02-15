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

func (table *Table) SaveToDB(db *sql.DB, tableTableName string) error {
	return SaveToDB(db, tableTableName, map[string]string{ // тут бы сделать interface{}
		"id":                  strconv.Itoa(table.Id),
		"description":         table.Description,
		"technical_equipment": table.TechnicalEquipment,
		"position":            XYPositionToText(table.X, table.Y),
		"hide":                strconv.FormatBool(table.Hide),
	})
}

func scanTablesFromDBRows(rows *sql.Rows) (*Table, error) {
	var newTable = Table{}
	var textPosition string
	err := rows.Scan(&newTable.Id, &newTable.Description, &newTable.TechnicalEquipment, &textPosition, newTable.Hide)

	newTable.X, newTable.Y, err = textPositionToXY(textPosition)
	return &newTable, err
}

func GetTables(db *sql.DB, tablesTableName, criteria string) ([]*Table, error) {
	if criteria != "" {
		criteria = " WHERE " + criteria
	}
	return getFromDB(db, "select * from "+tablesTableName+criteria, scanTablesFromDBRows)
}

func GetAllTables(db *sql.DB, tablesTableName string) ([]*Table, error) {
	return GetTables(db, tablesTableName, "")
}

func GetTablesBy(db *sql.DB, tablesTableName, columnName, columnValue string) ([]*Table, error) {
	return GetTables(db, tablesTableName, columnName+" = "+columnValue)
}

func GetTableById(db *sql.DB, tablesTableName string, id int) (*Table, error) {
	tables, err := GetTablesBy(db, tablesTableName, "id", strconv.Itoa(id))
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

func (table *Table) LoadReservations(db *sql.DB, reservationsTableName string) error {
	reservations, err := GetReservationsBy(db, reservationsTableName, "table_id", strconv.Itoa((table.Id)))
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
