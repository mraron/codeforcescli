package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
)

type Test struct {
	Input  string
	Output string
}

const Url = "http://codeforces.com/problemset/problem/%d/%s"
const Codeforcesdir = "/home/aron/cpp/codeforces"
const Default = "/home/aron/.codeforces/default.cpp"

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

		tests[l].Output, _ = s.Find("pre").Html()
		tests[l].Output = BrToEndl(tests[l].Output)
		l++
	})

	return tests, nil
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

type Printer interface {
	Print(*io.Writer, []Test)
}

type PrettyPrinter struct{}

func (p PrettyPrinter) Print(w *io.Writer, tests []Test) {
	for _, t := range tests {
		fmt.Fprintf(w, "%s\nOutput\n======\n%s\n\n\n", t.Input, t.Output)
	}
}

type JsonPrinter struct{}

func (j JsonPrinter) Print(w *io.Writer, tests []Test) {
	out, err := json.Marshal(tests)
	HandleError(err, "nem tudtam létrehozni a json objektumot:", err.Error())
	fmt.Fprintln(string(out))
}

func main() {
	app := cli.NewApp()
	app.Name = "codeforces"
	app.Usage = "a simple tool for codeforces.com"

	app.Commands = []cli.Command{
		{
			Name:      "problemio",
			ShortName: "pio",
			Usage:     "egy problém inputjának és outputjának megszerzése",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "turn it on to have json output",
				},
			},
			Action: func(c *cli.Context) {
				ValidateArgs(len(c.Args()), 2)

				problem, _ := strconv.Atoi(c.Args()[0])
				subproblem := c.Args()[1]

				HandleErrorBool((problem < 1), "Hiba a probléma sorszáma nem lehet 1-nél kisebb!")

				tests, err := GetProblem(problem, subproblem)
				HandleError(err, err.Error())

				var printer Printer

				if c.Bool("json") {
					printer = JsonPrinter{}
				} else {
					printer = PrettyPrinter{}
				}

				printer.Print(os.Stdout, tests)
			},
		},
		{
			Name:  "init",
			Usage: "egy fájlt inicializál a " + Default + " alapján",
			Action: func(c *cli.Context) {
				ValidateArgs(len(c.Args()), 1)
				//if len(c.Args()) < 1 {
				//	fmt.Println("Hiba meg kell adni egy fájl nevet ahova inicializálni akarjuk az alap.cpp-t")
				//	return
				//}

				file, err := os.Create(c.Args()[0])
				if err != nil {
					fmt.Println("Hiba nem tudom megnyitni a fájlt")
					return
				}
				defer file.Close()

				d, err := os.Open(Default)
				if err != nil {
					fmt.Println("Hiba nem tudom megnyitni az alapértelmezett fájlt")
					return
				}
				defer d.Close()

				dcontent, err := ioutil.ReadAll(d)
				if err != nil {
					fmt.Println("Hiba nem tudom beolvasni az alapértelmezett fájlt")
					return
				}
				file.Write(dcontent)
			},
		},
		{
			Name:  "test",
			Usage: "egy json fájlt és egy binárist kell megadni ami segítségével leelenőrzi a megadott inputokra az outputokat",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "timelimit",
					Value: 1000,
					Usage: "time limit in milliseconds",
				},
			},
			Action: func(c *cli.Context) {
				ValidateArgs(len(c.Args()), 2)

				file, err := os.Open(c.Args()[0])
				if err != nil {
					fmt.Println("Hiba nem tudom megnyitni a json-t!")
					return
				}
				content, err := ioutil.ReadAll(file)
				if err != nil {
					fmt.Println("Hiba nem tudom beolvasni a json-t!")
				}

				var tests []Test
				err = json.Unmarshal(content, &tests)
				if err != nil {
					fmt.Println("Nem sikerült dekódolni a json fájlt!")
				}

				for i, t := range tests {
					fmt.Printf("TEST #%d\n", i+1)
					cmd := exec.Command("./" + c.Args()[1])

					cmd.Stdin = strings.NewReader(t.Input)
					var out bytes.Buffer
					cmd.Stdout = &out

					err := cmd.Start()
					if err != nil {
						fmt.Println("CANNOT START EXECUTABLE")
						continue
					}
					l := make(chan bool, 1)

					go func(cmd *exec.Cmd, t chan bool) {
						cmd.Wait()

						t <- true
					}(cmd, l)

					select {
					case <-l:

						o, _ := ioutil.ReadAll(&out)
						if RemoveWhitespaces(string(o)) == RemoveWhitespaces(t.Output) {
							fmt.Println("[AC]")
						} else {
							fmt.Println("[WA]")
						}
					case <-time.After(time.Duration(int(time.Millisecond) * c.Int("timelimit"))):
						fmt.Println("[TE]")
						o, _ := ioutil.ReadAll(&out)
						fmt.Println(string(o))
						continue

					}
				}

			},
		},
	}

	app.Author = "Noszály Áron"
	app.Email = "noszalyaron4@gmail.com"
	app.Version = "v0.1.0"

	app.Action = func(c *cli.Context) {
		fmt.Println("don't know what to do, for help execute command \"codeforces help\" ;)")
	}

	app.Run(os.Args)
}
