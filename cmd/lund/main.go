package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "lund",
		Usage: "A pretty simple load balancer",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "hostname",
				Aliases: []string{"host"},
				Value:   "127.0.0.1",
				Usage:   "Host to run the server on",
				EnvVars: []string{"HOSTNAME"},
			},
			&cli.IntFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   8080,
				Usage:   "Port to run the server on",
				EnvVars: []string{"PORT"},
			},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "helo")
	})

	// construct address to listen on
	addr := fmt.Sprintf("%s:%d", c.String("hostname"), c.Int("port"))

	log.Println("Starting to listen on", addr)

	return http.ListenAndServe(addr, r)
}
