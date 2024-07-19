package lund

import (
	"errors"
	"sync/atomic"
)

type State struct {
	Servers []*Server

	// TODO: Reset this periodically! this will definitely overflow at some point
	Counter uint32
}

func (s *State) GetServersMap() map[string]*Server {
	servers := make(map[string]*Server)

	for _, server := range s.Servers {
		servers[server.URL] = server
	}

	return servers
}

// NOTE: cache?
func (s *State) GetAliveServers() []*Server {
	var servers []*Server

	for _, server := range s.Servers {
		if server.Alive.Load() {
			servers = append(servers, server)
		}
	}

	return servers
}

func (s *State) GetNextServer() (*Server, error) {
	servers := s.GetAliveServers()

	// if there's no alive servers, return an error
	if len(servers) == 0 {
		return nil, errors.New("There are no available servers at the moment")
	}

	counter := atomic.AddUint32(&s.Counter, 1)

	idx := int(counter % uint32(len(servers)))

	return servers[idx], nil
}
