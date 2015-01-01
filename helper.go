package main

import (
	"fmt"
	"os"
	"strings"
)

func ValidateArgs(length, minimum int) {
	if minimum > length {
		fmt.Println("Túl kevés argumentum a megadott parancs futattásához, lásd codeforces help!")
		os.Exit(1)
	}
}

// br-t \n-né konvertál
func BrToEndl(s string) string {
	s = strings.Replace(s, "<br/>", "\n", -1)
	s = strings.Replace(s, "<br />", "\n", -1)
	s = strings.Replace(s, "<br>", "\n", -1)
	return s
}

func RemoveWhitespaces(s string) string {
	s = strings.Replace(s, " ", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\t", "", -1)
	return s
}

func HandleError(e error, message string) {
	if e != nil {
		panic(message)
	}
}

func HandleBoolError(b bool, message string) {
	if b {
		panic(message)
	}
}
