package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

const __MIN_BETWEEN = time.Millisecond * 50

//Start_listen starts listening clients; options for listening must be entered as parameter
func Start_listen(opts *Opts) error {
	//TODO: logger needed
	log.Println("ws listen address : ", opts.Address)
	accepter, err := net.Listen("tcp", opts.Address)
	if nil != err {
		return err
	}
	opts.clients = make(map[net.Conn]*connection)
	for {
		client, err := accepter.Accept()
		if nil != err {
			//TODO: logger needed
			log.Println("cannot accept client : ", err)
		}

		opts.lck.Lock()
		opts.clients[client] = new(connection)
		opts.clients[client].con = client
		opts.clients[client].sig__kil = make(chan bool)
		opts.clients[client].stop__writing = make(chan bool)
		opts.lck.Unlock()

		go pong__handler(opts.clients[client], opts)
		go func(client net.Conn) {
			<-opts.clients[client].sig__kil
			opts.lck.Lock()
			delete(opts.clients, client)
			opts.lck.Unlock()
			_ = client.Close()
			return
		}(client)
	}
	return nil
}

//Broadcast writes the unstructured data into clients
func Broadcast(d []byte, opts *Opts) {
	opts.lck.Lock()
	defer opts.lck.Unlock()
	for c := range opts.clients {
		go func() {
			go func() {
				c.Write(d)
				opts.clients[c].stop__writing <- true
			}()

			select {
			case <-opts.clients[c].stop__writing:
				return
			case <-time.After(time.Minute):
				opts.lck.Lock()
				opts.clients[c].sig__kil <- true
				opts.lck.Unlock()
				return
			}
		}()
	}
}

//BroadcastJSON writes the structured (json) data into clients
func BroadcastJSON(data *Socket_data, opts *Opts) {
	d, err := json.Marshal(data)
	if nil != err {
		//TODO: logger needed
		fmt.Errorf("cannot serialize json : %v", err)
	}
	Broadcast(d, opts)
}
func pong__handler(conn *connection, opts *Opts) {
	killSent := false
	var heartBeat = make(chan bool)
	go func() {
		var buf = make([]byte, 1024)
		var __last__msg__time = time.Now().UnixNano()
		var __now = __last__msg__time + int64(__MIN_BETWEEN) + 10
		for {
			if __now < __last__msg__time+int64(__MIN_BETWEEN) {
				conn.sig__kil <- true
				//TODO: logger needed
				fmt.Println("too often data!")
				return
			}
			__last__msg__time = __now
			_, err := conn.con.Read(buf)
			if nil != err {
				conn.con__lock.Lock()
				defer conn.con__lock.Unlock()
				if !killSent {
					killSent = true
					conn.sig__kil <- true
				}
				return
			}
			__now = time.Now().UnixNano()
			heartBeat <- true
		}
	}()
	for {
		select {
		case <-heartBeat:
		case <-conn.sig__kil:
			return
		case <-time.After(opts.Time_out):

			conn.con__lock.Lock()

			if !killSent {
				go func() {
					conn.sig__kil <- true
				}()
				killSent = true
			}
			conn.con__lock.Unlock()
			return
		}
	}
}
