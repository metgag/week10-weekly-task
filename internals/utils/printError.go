package utils

import (
	"log"
	"strings"
)

func PrintError(head string, rep int, err error) {
	log.Printf("%s %s %s", strings.Repeat("=", rep), head, strings.Repeat("=", rep))
	if err != nil {
		log.Println(err.Error())
	}
}
