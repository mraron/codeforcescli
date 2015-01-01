package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/codegangsta/cli"
)

type Test struct {
	Input  string
	Output string
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

func Testcli(c *cli.Context) {
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

}
