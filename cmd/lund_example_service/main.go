package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type PageData struct {
	Hostname string
	UnixTime int64
}

const tmpl = `
<h1>hello from {{ .Hostname }}</h1>
<p>unix time: {{ .UnixTime }}</p>
`

func main() {
	// get the hostname of the system (id of the container if running in docker)
	hostname, err := os.Hostname()

	if err != nil {
		hostname = "unknown"
	}

	// parse the template
	t := template.Must(template.New("page").Parse(tmpl))

	// handle requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Hostname: hostname,
			UnixTime: time.Now().Unix(),
		}

		err := t.Execute(w, data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// start the web server
	log.Fatal(http.ListenAndServe(":8081", nil))
}
