package main

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/codegangsta/cli"
)

var tesztek []Test
var testfile string

func AddMenu(w io.Writer) {
	io.WriteString(w, "<a href=\"/\">Home</a>|<a href=\"/reload\">Reload test file</a>|<a href=\"/insert\">Insert</a>|<a href=\"/quit\">Quit</a><br><br>")
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	AddMenu(w)
	w.Header().Add("Content-Type", "text/html; charset=UTF-8")
	io.WriteString(w, "<b>TESZT ESETEK</b><br><table border=\"2\"><tr><td>ID</td><td>INPUT</td><td>YOUR OUTPUT</td><td>ANSWER</td><td>TIME</td><td>STATUS</td><td>SETTINGS</td></tr>")

	var l int = 0
	for i, t := range tesztek {
		io.WriteString(w, "<tr><td>"+strconv.Itoa(i)+"</td><td>"+EndlToBr(t.Input)+"</td><td>"+EndlToBr(t.Output)+"</td><td>"+EndlToBr(t.Answer)+"</td><td>"+strconv.Itoa(t.Time/1000000)+"ms</td><td bgcolor="+StatusToColor(t.Status)+">"+StatusToString(t.Status)+"</td><td><a href=\"/delete?id="+strconv.Itoa(i)+"\">DELETE</a> EDIT</td></tr>")
		l++
	}
	io.WriteString(w, "<tr><form method=\"post\" action=\"/add\"><td>"+strconv.Itoa(l)+"</td><td><textarea name=\"input\"></textarea></td><td></td><td><textarea name=\"answer\"></textarea></td><td><input type=\"submit\"></td></form></tr>")

	io.WriteString(w, "</table>")
}

func QuitHandler(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}

func ReloadHandler(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(testfile)
	HandleError(err, ".....")

	tesztek = JSON{}.Unmarshal(file)

	http.Redirect(w, r, "/", 302)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	HandleError(err, ".....")

	input := r.Form.Get("input")
	answer := r.Form.Get("answer")

	test := Test{Input: input, Answer: answer, Status: -1}
	tesztek = append(tesztek, test)

	http.Redirect(w, r, "/", 302)
}

func InsertHandler(w http.ResponseWriter, r *http.Request) {
	jsonout, err := os.Create(testfile)
	HandleError(err, "...")

	JSON{}.Marshal(jsonout, tesztek)

	http.Redirect(w, r, "/", 302)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if len(id) == 0 {
		return
	}
	id_int, err := strconv.Atoi(id)
	HandleError(err, "konvertálási hiba")

	tesztek = append(tesztek[:id_int], tesztek[id_int+1:]...)

	http.Redirect(w, r, "/", 302)
}

func Manager(c *cli.Context) {
	ValidateArgs(len(c.Args()), 1)
	port := c.String("port")

	testfile = c.Args()[0]
	file, err := os.Open(testfile)
	HandleError(err, ".....")

	tesztek = JSON{}.Unmarshal(file)
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/quit", QuitHandler)
	http.HandleFunc("/reload", ReloadHandler)
	http.HandleFunc("/add", AddHandler)
	http.HandleFunc("/insert", InsertHandler)
	http.HandleFunc("/delete", DeleteHandler)

	err = http.ListenAndServe(":"+port, nil)
	HandleError(err, "Nem tudtam elindítani a web szervert :(")
}
