package main

import (
	"io"
	"log"
)

func HandleError(err error, msg string) {
	if err == io.EOF {
		log.Println(colorRed, "Exiting...")
	} else {
		log.Println(colorRed, msg, err)
	}
}

func HandleInstallError(err error, msg string) bool {
	if err != nil {
		log.Println(colorRed, msg)
		return true
	}
	return false
}
