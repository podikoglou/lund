package lund

type State struct {
	Servers []*Server
}

func (s *State) GetAliveServers() []*Server {
	var servers []*Server

	for _, server := range s.Servers {
		if server.Alive.Load() {
			servers = append(servers, server)
		}
	}

	return servers
}
