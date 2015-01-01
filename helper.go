package main

import (
	"fmt"
	"os"
	"strings"
)

func StatusToColor(i int) string {
	switch i {
	case 0:
		return "blue"
	case 1:
		return "orange"
	case 2:
		return "green"
	}
	return "grey"
}

func StatusToString(i int) string {
	switch i {
	case 0:
		return "TE"
	case 1:
		return "WA"
	case 2:
		return "AC"
	}
	return "??"
}

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

func EndlToBr(s string) string {
	s = strings.Replace(s, "\n", "<br>", -1)
	return s
}

func RemoveWhitespaces(s string) string {
	s = strings.TrimSpace(s)
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
