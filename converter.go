package main

import "io"
import "encoding/json"
import "io/ioutil"

type Marshaller interface {
	Marshal(io.Writer, []Test)
}

type Unmarshaller interface {
	Unmarshal(io.Reader) []Test
}

type JSON struct {
	Pretty bool
}

func (j JSON) Marshal(w io.Writer, tests []Test) {
	var data []byte
	var err error

	if j.Pretty {
		data, err = json.MarshalIndent(tests, "", "\t")
	} else {
		data, err = json.Marshal(tests)
	}
	HandleError(err, "...")

	_, err = w.Write(data)
	HandleError(err, "...")
}

func (j JSON) Unmarshal(r io.Reader) (t []Test) {
	all, err := ioutil.ReadAll(r)
	HandleError(err, "...")

	err = json.Unmarshal(all, &t)
	HandleError(err, "...")

	return
}

//TODO XML :)
