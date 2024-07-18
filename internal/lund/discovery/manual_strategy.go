package discovery

import (
	"sync/atomic"

	"github.com/podikoglou/lund/internal/lund"
)

type ManualDiscoveryStrategy struct {
	servers []string
}

func NewManualDiscoveryStrategy(servers []string) ManualDiscoveryStrategy {
	return ManualDiscoveryStrategy{
		servers: servers,
	}
}

func (s ManualDiscoveryStrategy) Discover() []lund.Server {
	var servers []lund.Server

	for _, url := range s.servers {

		// we *could* do a health check here, but let's just
		// leave it to the health check component so that
		// we don't repeat code.
		servers = append(servers, lund.Server{
			URL:   url,
			Alive: atomic.Bool{},
		})
	}

	return servers
}
