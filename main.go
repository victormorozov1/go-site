package main

import "internal/server"

func main() {
	s := server.Server{"127.0.0.1", 8080}
	s.Start()
}
