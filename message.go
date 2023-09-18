package paxos

import (
	"log"
	"net/rpc"
)

type Message struct {
	Number int
	Value  any
	From   int
	To     int
}

type Acknowledge struct {
	Ok     bool
	Accepted bool
	Number int
	Value  any
}

func call(srv, name string, args *Message, reply *Acknowledge) bool {
	c, err := rpc.Dial("tcp", srv)
	if err != nil {
		log.Println("rpc.Dial error")
		return false
	}
	defer c.Close()

	err = c.Call(name, args, reply)
	if err != nil {
		log.Println("c.Call error")
		return false
	}
	return true
}
