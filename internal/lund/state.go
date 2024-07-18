package lund

type State struct {
	Servers []struct {
		URL   string
		Alive bool
	}
}
