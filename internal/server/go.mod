module server

go 1.19

require internal/functions v1.0.0 // indirect

replace internal/functions => ../../internal/functions

require internal/db v1.0.0

require golang.org/x/exp v0.0.0-20230202163644-54bba9f4231b // indirect

replace internal/db => ../../internal/db

require (
	github.com/gorilla/mux v1.8.0 // indirect
)