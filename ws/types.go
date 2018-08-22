package ws

import (
	"sync"
	"net"
	"time"
)

//Opts options to start server socket
type Opts struct {
	//Address address to listen clients
	Address  string
	//Time_out timeout for pong handling
	Time_out time.Duration
	clients  map[net.Conn]*connection
	lck      sync.Mutex
}

type connection struct {
	con__lock sync.Mutex
	con       net.Conn
	sig__kil  chan bool
}

//Socket_data data structure type for writing into clients
type Socket_data struct {
	//Event event name
	Event string   `bson:"event" json:"event"`
	//Rooms there could be rooms in the websocket side so you may define those rooms here
	Rooms []string `bson:"rooms" json:"rooms"`
	//Data actual data you want to send to websocket's clients
	Data  string   `bson:"data" json:"data"`
}
