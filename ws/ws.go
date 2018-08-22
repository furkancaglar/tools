package ws

import (
	"encoding/json"
	"net"
	"time"
	"log"
	"fmt"
)

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
			var stop = make(chan bool)
			go func() {
				c.Write(d)
				stop <- true
			}()

			select {
			case <-stop:
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
	opts.lck.Lock()
	defer opts.lck.Unlock()
	for c := range opts.clients {
		go func() {
			var stop = make(chan bool)
			go func() {
				c.Write(d)
				stop <- true
			}()

			select {
			case <-stop:
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
func pong__handler(conn *connection, opts *Opts) {
	killSent := false
	var heartBeat = make(chan bool)
	go func() {
		var buf = make([]byte, 1024)
		for {
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
