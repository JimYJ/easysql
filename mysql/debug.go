package mysql

import (
	"log"
)

func Debug() {
	log.Println(lastQuery)
	return
}

func DebugMode() {
	showErrors = true
}

func ReleaseMode() {
	showErrors = false
}

func printErrors(err error) {
	if err != nil && showErrors == true {
		log.Println(err)
	}
}
