package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
)

func GetProblem(num int, sub string) ([]Test, error) {
	url := fmt.Sprintf(Url, num, sub)

	doc, err := goquery.NewDocument(url)

	if err != nil {
		return nil, errors.New("Failed to initalize goquery document: " + err.Error())
	}

	tests := make([]Test, 0)
	doc.Find(".input").Each(func(i int, s *goquery.Selection) {
		var t Test

		t.Input, _ = s.Find("pre").Html()
		t.Input = BrToEndl(t.Input)

		tests = append(tests, t)
	})

	l := 0
	doc.Find(".output").Each(func(i int, s *goquery.Selection) {

		tests[l].Answer, _ = s.Find("pre").Html()
		tests[l].Answer = BrToEndl(tests[l].Answer)
		tests[l].Status = -1
		tests[l].Time = 1
		l++
	})

	return tests, nil
}

type Printer interface {
	Print(io.Writer, []Test)
}

type PrettyPrinter struct{}

func (p PrettyPrinter) Print(w io.Writer, tests []Test) {
	for _, t := range tests {
		fmt.Fprintf(w, "%s\nAnswer\n======\n%s\n\n\n", t.Input, t.Answer)
	}
}

type JsonPrinter struct{}

func (j JsonPrinter) Print(w io.Writer, tests []Test) {
	out, err := json.Marshal(tests)
	HandleError(err, "nem tudtam létrehozni a json objektumot")
	fmt.Fprintln(w, string(out))
}

func Problem(c *cli.Context) {
	ValidateArgs(len(c.Args()), 2)

	problem, _ := strconv.Atoi(c.Args()[0])
	subproblem := c.Args()[1]

	HandleBoolError((problem < 1), "Hiba a probléma sorszáma nem lehet 1-nél kisebb!")

	tests, err := GetProblem(problem, subproblem)
	HandleError(err, "..")

	var printer Printer

	if c.Bool("json") {
		printer = JsonPrinter{}
	} else {
		printer = PrettyPrinter{}
	}
	printer.Print(os.Stdout, tests)
}
