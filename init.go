package main

import (
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"
)

func Init(c *cli.Context) {
	ValidateArgs(len(c.Args()), 1)
	CreateDefault(c.Args()[0], Default)
}

func CreateDefault(filename string, df string) {
	file, err := os.Create(filename)
	HandleError(err, "Hiba nem tudom megnyitni a fájlt")
	defer file.Close()

	d, err := os.Open(df)
	HandleError(err, "Hiba nem tudom megnyitni az alapértelmezett fájlt")
	defer d.Close()

	dcontent, err := ioutil.ReadAll(d)
	HandleError(err, "Hiba nem tudom beolvasni az alapértelmezett fájlt")

	_, err = file.Write(dcontent)
	HandleError(err, "Hiba nem tudom az alapértelmezett fájl tartalmát beleírni a fájlba")
}
