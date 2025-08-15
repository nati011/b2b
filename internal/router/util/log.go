package handler

import (
	"log"
	"os"
)

var logger = log.New(os.Stdout, "api: ", log.LstdFlags)

func logError(err error) {
	logger.Println(err)
}
