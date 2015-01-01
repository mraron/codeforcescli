package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

const Url = "http://codeforces.com/problemset/problem/%d/%s"
const Codeforcesdir = "/home/aron/cpp/codeforces"
const Default = "/home/aron/.codeforces/default.cpp"

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
			Name:   "init",
			Usage:  "egy fájlt inicializál a " + Default + " alapján",
			Action: Init,
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
			Action: Testcli,
		},
	}

	app.Author = "Noszály Áron"
	app.Email = "noszalyaron4@gmail.com"
	app.Version = "v0.3.0"

	app.Action = func(c *cli.Context) {
		fmt.Println("don't know what to do, for help execute command \"codeforces help\" ;)")
	}

	app.Run(os.Args)
}
