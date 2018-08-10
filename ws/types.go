package ws

import (
	"sync"
	"net"
	"time"
)

type Opts struct {
	Address  string
	Time_out time.Duration
	clients  map[net.Conn]*connection
	lck      sync.Mutex
}

type connection struct {
	con__lock sync.Mutex
	con       net.Conn
	sig__kil  chan bool
	ticker    *time.Ticker
}

type Socket_data struct {
	Event string   `bson:"event" json:"event"`
	Rooms []string `bson:"rooms" json:"rooms"`
	Data  string   `bson:"data" json:"data"`
}
