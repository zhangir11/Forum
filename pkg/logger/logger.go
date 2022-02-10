package logger

import (
	"log"
	"os"
)

//Logger ...
func Logger() *os.File {
	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Println(err.Error())
	}
	log.SetOutput(file)
	return file
}
