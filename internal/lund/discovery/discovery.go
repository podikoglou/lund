package discovery

import (
	"log"
	"time"

	"github.com/podikoglou/lund/internal/lund"
)

type DiscoveryOptions struct {
	Interval  time.Duration
	ProxyOpts *lund.ProxyOptions
	Strategy  DiscoveryStrategy
}

func DiscoveryLoop(state *lund.State, opts DiscoveryOptions) {
	for {

		time.Sleep(opts.Interval)

		// perform discovery
		oldServersMap := state.GetServersMap()
		newServers := opts.Strategy.Discover()

		discovered := 0

		for _, s := range newServers {
			// check if the server is already in there
			_, exists := oldServersMap[s.URL]

			if !exists {
				// if it isn't, create an http client for it and update state.Servers

				s.Client = lund.CreateHTTPClient(opts.ProxyOpts)

				state.Servers = append(state.Servers, s)

				discovered++
			}
		}

		if discovered > 0 {
			log.Println("Discovered", discovered, "New Servers")
		}

	}
}
