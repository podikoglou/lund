package lund

import (
	"log"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

type HealthCheckOptions struct {
	Interval         time.Duration
	WriteTimeout     time.Duration
	ReadTimeout      time.Duration
	DNSCacheDuration time.Duration
	Concurrency      int
}

// This function checks a server's health.
//
// Returns true if the server is up.
// Returns false if the server is down.
func CheckHealth(client *fasthttp.Client, url string) bool {
	// NOTE: should we reuse the request?
	// not the response though
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	// TODO: consider parsing the URI once and using setURI here
	// instead of setRequestURi
	req.SetRequestURI(url)

	// TODO: make option for using GET instead of HEAD
	req.Header.SetMethod(fasthttp.MethodHead)

	err := client.Do(req, resp)

	return err == nil
}

// This functions is typically ran in a goroutine and constantly sends
// health checks to the servers provided by the State.
func HealthCheckLoop(state *State, opt HealthCheckOptions) {
	// initialize fasthttp client
	client := &fasthttp.Client{
		ReadTimeout:                   opt.ReadTimeout,
		WriteTimeout:                  opt.WriteTimeout,
		NoDefaultUserAgentHeader:      false,
		DisableHeaderNamesNormalizing: true,
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      opt.Concurrency,
			DNSCacheDuration: opt.DNSCacheDuration,
		}).Dial,
	}

	for {
		// sleeps for a given interval
		time.Sleep(opt.Interval)

		var wg sync.WaitGroup

		for _, server := range state.Servers {
			wg.Add(1)

			go func() {
				defer wg.Done()

				log.Println("Performing Health Check on", server.URL)

				// perform health check
				alive := CheckHealth(client, server.URL)

				// update server state
				server.Alive.Store(alive)

				if alive {
					log.Println(server.URL, "is up")
				} else {
					log.Println(server.URL, "is down")
				}
			}()
		}
	}
}
