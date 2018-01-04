package mysql

import (
	"log"
)

func Debug() {
	log.Println(lastQuery)
	return
}

func ShowErrors() {
	showErrors = true
}

func HideErrors() {
	showErrors = false
}

func printErrors(err error) {
	if err != nil && showErrors == true {
		log.Println(err)
	}
}
