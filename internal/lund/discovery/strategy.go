package discovery

import "github.com/podikoglou/lund/internal/lund"

type DiscoveryStrategy interface {
	// Synchronously discovers servers and returns their URLs
	DiscoverServers() []lund.Server
}
