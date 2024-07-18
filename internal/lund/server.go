package lund

import (
	"sync/atomic"
)

type Server struct {
	URL   string
	Alive atomic.Bool
}
