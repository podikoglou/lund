package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/podikoglou/lund/internal/lund"
	"github.com/podikoglou/lund/internal/lund/discovery"
	"github.com/urfave/cli/v2"
	"github.com/valyala/fasthttp"
)

func main() {
	app := &cli.App{
		Name:  "lund",
		Usage: "A pretty simple load balancer",
		Flags: []cli.Flag{
			// server
			&cli.StringFlag{
				Name:     "hostname",
				Aliases:  []string{"host"},
				Value:    "127.0.0.1",
				Usage:    "Host to run the server on",
				Category: "Server:",
				EnvVars:  []string{"HOSTNAME"},
			},
			&cli.IntFlag{
				Name:     "port",
				Aliases:  []string{"p"},
				Value:    8080,
				Usage:    "Port to run the server on",
				Category: "Server:",
				EnvVars:  []string{"PORT"},
			},

			// health checking
			&cli.DurationFlag{
				Name:     "health-check-interval",
				Usage:    "Health Check Interval",
				Value:    time.Second * 5,
				Category: "Health Checking:",
				EnvVars:  []string{"HEALTH_CHECK_INTERVAL"},
			},
			&cli.DurationFlag{
				Name:     "health-check-write-timeout",
				Usage:    "Write timeout for health checking",
				Value:    time.Millisecond * 300,
				Category: "Health Checking:",
				EnvVars:  []string{"HEALTH_CHECK_WRITE_TIMEOUT"},
			},
			&cli.DurationFlag{
				Name:     "health-check-read-timeout",
				Usage:    "Read timeout for the health checking",
				Value:    time.Millisecond * 300,
				Category: "Health Checking:",
				EnvVars:  []string{"HEALTH_CHECK_READ_TIMEOUT"},
			},
			&cli.DurationFlag{
				Name:     "health-check-dns-cache-duration",
				Usage:    "How often to clear the DNS cache for the Health Check component",
				Value:    time.Minute,
				Category: "Health Checking:",
				EnvVars:  []string{"HEALTH_CHECK_DNS_CACHE_DURATION"},
			},
			&cli.IntFlag{
				Name:     "health-check-concurrency",
				Usage:    "How many Health Checks can run at the same time",
				Value:    4,
				Category: "Health Checking:",
				EnvVars:  []string{"HEALTH_CHECK_CONCURRENCY"},
			},

			&cli.DurationFlag{
				Name:     "proxy-write-timeout",
				Usage:    "Write timeout for the reverse proxy",
				Value:    time.Millisecond * 1000,
				Category: "Reverse Proxy:",
				EnvVars:  []string{"PROXY_CHECK_WRITE_TIMEOUT"},
			},
			&cli.DurationFlag{
				Name:     "proxy-read-timeout",
				Usage:    "Read timeout for the reverse proxy",
				Value:    time.Millisecond * 1000,
				Category: "Reverse Proxy:",
				EnvVars:  []string{"PROXY_CHECK_READ_TIMEOUT"},
			},
			&cli.DurationFlag{
				Name:     "proxy-dns-cache-duration",
				Usage:    "How often to clear the DNS cache for the reverse proxy",
				Value:    time.Millisecond * 300,
				Category: "Reverse Proxy:",
				EnvVars:  []string{"PROXY_DNS_CACHE_DURATION"},
			},
			&cli.IntFlag{
				Name:     "proxy-concurrency",
				Usage:    "How many requests can run at the same time by the reverse proxy http client (for each server)",
				Value:    4,
				Category: "Reverse Proxy:",
				EnvVars:  []string{"PROXY_CONCURRENCY"},
			},

			// discovery
			&cli.StringFlag{
				Name:     "discovery-strategy",
				Usage:    "Discovery Mode (possible values: docker, manual)",
				Required: true,
				Category: "Discovery:",
				EnvVars:  []string{"DISCOVERY_STRATEGY"},
				Action: func(c *cli.Context, val string) error {
					if val != "docker" && val != "manual" {
						return errors.New("Invalid discovery-strategy value (possible values: docker, manual)")
					}

					if val == "manual" && !c.IsSet("discovery-servers") {
						return errors.New("You must set the discovery-servers flag")
					}

					return nil
				},
			},

			&cli.StringSliceFlag{
				Name:     "discovery-servers",
				Usage:    "List of servers (URLs) -- use commas for separation (only applicable with manual discovery mode)",
				Category: "Discovery:",
				EnvVars:  []string{"DISCOVERY_SERVERS"},
				Action: func(c *cli.Context, val []string) error {
					// try to parse every url given, and if an error is found, return it
					for _, v := range val {
						url, err := url.Parse(v)

						if err != nil {
							return err
						}

						if url.Scheme != "http" && url.Scheme != "https" {
							return errors.New("Your servers must be URLs starting with http:// or https://")
						}
					}

					return nil
				},
			},

			&cli.StringFlag{
				Name:     "discovery-docker-sock",
				Usage:    "Docker Socket Path",
				Value:    "/var/run/docker.sock",
				Category: "Discovery:",
				EnvVars:  []string{"DISCOVERY_DOCKER_SOCK"},
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
	// create global state
	state := lund.State{}

	// initialize discovery strategy
	var strategy discovery.DiscoveryStrategy

	switch c.String("discovery-strategy") {
	case "docker":
		log.Panic("The docker discovery strategy hasn't been implemented yet")
		break

	case "manual":
		strategy = discovery.NewManualDiscoveryStrategy(c.StringSlice("discovery-servers"))
		break
	}

	// perform initial discovery
	log.Println("Performing Discovery...")

	servers := strategy.Discover()
	state.Servers = servers

	// load proxy options, and create HTTP clients for each server loaded
	// (not sure if we should do this here)
	proxyOpts := lund.ProxyOptions{
		WriteTimeout:     c.Duration("proxy-write-timeout"),
		ReadTimeout:      c.Duration("proxy-write-timeout"),
		DNSCacheDuration: c.Duration("proxy-dns-cache-duration"),
		Concurrency:      c.Int("proxy-concurrency"),
	}

	for _, server := range state.Servers {
		client := lund.CreateHTTPClient(&proxyOpts)

		server.Client = client
	}

	log.Println("Discovered", len(servers), "Servers")

	// start performing health checks
	go lund.HealthCheckLoop(&state, lund.HealthCheckOptions{
		Interval:         c.Duration("health-check-interval"),
		WriteTimeout:     c.Duration("health-check-write-timeout"),
		ReadTimeout:      c.Duration("health-check-read-timeout"),
		DNSCacheDuration: c.Duration("health-check-dns-cache-duration"),
		Concurrency:      c.Int("health-check-concurrency"),
	})

	// create server
	srv := fasthttp.Server{
		Handler: lund.MakeRequestHandler(&state),
	}

	// construct address to listen on
	addr := fmt.Sprintf("%s:%d", c.String("hostname"), c.Int("port"))

	log.Println("Starting to listen on", addr)

	return srv.ListenAndServe(addr)
}
