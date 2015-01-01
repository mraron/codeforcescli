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
			Name: "problem",
			//ShortName: "pio",
			Usage: "egy problém inputjának és outputjának megszerzése",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "json",
					Usage: "turn it on to have json output",
				},
				cli.BoolFlag{
					Name:  "prettyjson",
					Usage: "szépen formázott json",
				},
			},
			Action: Problem,
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
		{
			Name:  "manager",
			Usage: "egyedi teszt esetek hozzáadása egy teszt esetek tartalmazó json fájlhoz (webes felületen!)",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Value: "8080",
					Usage: "milyen porton figyeljen a szerver",
				},
			},
			Action: Manager,
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
