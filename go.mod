module main

go 1.19

require internal/server v1.0.0

require internal/db v1.0.0

require internal/functions v1.0.0

replace internal/server => ./internal/server

replace internal/db => ./internal/db

replace internal/functions => ./internal/functions

require (
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	golang.org/x/exp v0.0.0-20230202163644-54bba9f4231b // indirect
)
