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
	Print(io.Writer, []Test)
}

type PrettyPrinter struct{}

func (p PrettyPrinter) Print(w io.Writer, tests []Test) {
	for _, t := range tests {
		fmt.Fprintf(w, "%s\nOutput\n======\n%s\n\n\n", t.Input, t.Output)
	}
}

type JsonPrinter struct{}

func (j JsonPrinter) Print(w io.Writer, tests []Test) {
	out, err := json.Marshal(tests)
	HandleError(err, "nem tudtam létrehozni a json objektumot")
	fmt.Fprintln(w, string(out))
}

func Pio(c *cli.Context) {
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

const (
	TimeLimitExceed = iota
	WrongAnswer
	Accepted
	End
)

type TesterOptions struct {
	w         io.Writer
	binary    string
	timelimit int
	verbose   bool
	t         Test
}

func Tester(to TesterOptions) int {
	cmd := exec.Command("./" + to.binary)

	var out bytes.Buffer

	cmd.Stdin = strings.NewReader(to.t.Input)
	cmd.Stdout = &out

	err := cmd.Start()
	HandleError(err, "CANNOT START EXECUTABLE")

	l := make(chan bool, 1)
	go func(cmd *exec.Cmd, t chan bool) {
		cmd.Wait()
		t <- true
	}(cmd, l)

	select {
	case <-l:
		o, _ := ioutil.ReadAll(&out)
		if to.verbose {
			fmt.Fprintf(to.w, "\nANSWER\n")
			fmt.Fprintf(to.w, "======\n\n%s", to.t.Output)
			fmt.Fprintf(to.w, "\n\nYOUR ANSWER\n")
			fmt.Fprintf(to.w, "===========\n\n%s\n", o)
		}
		if RemoveWhitespaces(string(o)) == RemoveWhitespaces(to.t.Output) {
			return Accepted
		} else {
			return WrongAnswer

		}
	case <-time.After(time.Duration(int(time.Millisecond) * to.timelimit)):
		return TimeLimitExceed
		cmd.Process.Kill()
	}

	return End
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
			Action: Pio,
		},
		{
			Name:  "init",
			Usage: "egy fájlt inicializál a " + Default + " alapján",
			Action: func(c *cli.Context) {
				ValidateArgs(len(c.Args()), 1)

				file, err := os.Create(c.Args()[0])
				HandleError(err, "Hiba nem tudom megnyitni a fájlt")
				defer file.Close()

				d, err := os.Open(Default)
				HandleError(err, "Hiba nem tudom megnyitni az alapértelmezett fájlt")
				defer d.Close()

				dcontent, err := ioutil.ReadAll(d)
				HandleError(err, "Hiba nem tudom beolvasni az alapértelmezett fájlt")

				_, err = file.Write(dcontent)
				HandleError(err, "Hiba nem tudom az alapértelmezett fájl tartalmát beleírni a fájlba")
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
				cli.BoolFlag{
					Name:  "verbose",
					Usage: "Bőbeszédű output",
				},
			},
			Action: func(c *cli.Context) {
				ValidateArgs(len(c.Args()), 2)

				file, err := os.Open(c.Args()[0])
				HandleError(err, "Hiba nem tudom megnyitni a json-t!")

				content, err := ioutil.ReadAll(file)
				HandleError(err, "Hiba nem tudom beolvasni a json-t!")

				var tests []Test
				err = json.Unmarshal(content, &tests)
				HandleError(err, "Nem sikerült dekódolni a json fájlt!")

				to := TesterOptions{w: os.Stdout, binary: c.Args()[1], timelimit: c.Int("timelimit"), verbose: c.Bool("verbose")}
				for i, t := range tests {
					to.t = t
					fmt.Printf("TEST #%d\n", i+1)

					output := Tester(to)
					switch output {
					case WrongAnswer:
						fmt.Println("[WA]")
					case TimeLimitExceed:
						fmt.Println("[TE]")
					case Accepted:
						fmt.Println("[AC]")
					case End:
						fmt.Println("[??]")
					}
				}

			},
		},
	}

	app.Author = "Noszály Áron"
	app.Email = "noszalyaron4@gmail.com"
	app.Version = "v0.2.0"

	app.Action = func(c *cli.Context) {
		fmt.Println("don't know what to do, for help execute command \"codeforces help\" ;)")
	}

	app.Run(os.Args)
}
