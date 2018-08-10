package ws

import (
	"encoding/json"
	"net"
	"time"
	"log"
	"fmt"
)

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
func Broadcast(d []byte, opts *Opts) {
	opts.lck.Lock()
	defer opts.lck.Unlock()
	for c := range opts.clients {
		c.Write(d)
	}
}
func BroadcastJSON(data *Socket_data, opts *Opts) {
	d, err := json.Marshal(data)
	if nil != err {
		//TODO: logger needed
		fmt.Errorf("cannot serialize json : %v", err)
	}
	opts.lck.Lock()
	defer opts.lck.Unlock()
	for c := range opts.clients {
		c.Write(d)
	}
}
func pong__handler(conn *connection, opts *Opts) {
	conn.ticker = time.NewTicker(opts.Time_out)
	defer conn.ticker.Stop()
	killSent := false
	var heartBeat = make(chan bool)
	go func() {
		var buf = make([]byte, 1024)
		for {
			<-conn.ticker.C
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
		//this is here because ticker writes immediately it starts so this cause connection close immediately
		<-conn.ticker.C
		select {
		case <-heartBeat:
		case <-conn.sig__kil:
			return
		case <-conn.ticker.C:

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
