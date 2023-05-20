package db

import "database/sql"

type Database struct {
	Name                         string
	Driver, Login, Password, Url string
	Tables                       Tables
	Connection                   *sql.DB
}

type Tables struct {
	Users        UsersTable
	Reservations ReservationsTable
	Tables       TablesTable
}

type ReservationsTable struct {
	TableName string
	Id        string
	UserId    string
	TableId   string
	StartTime string
	EndTime   string
}

type UsersTable struct {
	TableName      string
	Id             string
	Name           string
	Surname        string
	Patronymic     string
	Role           string
	JobTile        string
	Department     string
	Phone          string
	Email          string
	PhotoSrc       string
	HashedPassword string
}

type TablesTable struct {
	TableName          string
	Id                 string
	Description        string
	TechnicalEquipment string
	Position           string
	Hide               string
}
