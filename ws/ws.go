package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)
// it could be used as another functionality in the future so that's why commented out
//const __MIN_BETWEEN = int64(time.Millisecond * 10)

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

func Conn_client_num(opts *Opts) int {
	opts.lck.Lock()
	defer opts.lck.Unlock()
	return len(opts.clients)
}

//Broadcast writes the unstructured data into clients
func Broadcast(d []byte, opts *Opts) {
	opts.lck.Lock()
	defer opts.lck.Unlock()
	for c, con__struct := range opts.clients {
		go func() {
			go func() {
				c.Write(d)
				con__struct.stop__writing <- true
			}()

			select {
			case <-con__struct.stop__writing:
				return
			case <-time.After(time.Minute):
				con__struct.sig__kil <- true
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
		//var __last__msg__time = time.Now().UnixNano()
		//var __now = __last__msg__time + __MIN_BETWEEN + 4
		for {
			//if __now < __last__msg__time+__MIN_BETWEEN {
			//	conn.sig__kil <- true
			//	//TODO: logger needed
			//	fmt.Println("too often data!")
			//	return
			//}
			//__last__msg__time = __now
			_, err := conn.con.Read(buf)
			if nil != err {
				if !killSent {
					killSent = true
					conn.sig__kil <- true
				}
				return
			}
			//__now = time.Now().UnixNano()
			heartBeat <- true
		}
	}()
	for {
		select {
		case <-heartBeat:
		case <-conn.sig__kil:
			return
		case <-time.After(opts.Time_out):

			if !killSent {
				go func() {
					conn.sig__kil <- true
				}()
				killSent = true
			}

			return
		}
	}
}
