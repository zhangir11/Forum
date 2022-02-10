package main

import (
	baselogic "forum/internal/BaseLogic"
	"forum/pkg/logger"
	"forum/pkg/sserver"
	"log"
)

func main() {
	logfile := logger.Logger()
	defer logfile.Close()
	server := sserver.CreateServer()

	if err := baselogic.Start(server); err != nil {
		log.Println("Failed launching server: ", err.Error())
	}
	DB := server.Database
	if err := DB.Close(); err != nil {
		log.Println(err.Error())
	}
}
