package paxos

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
)

type Learner struct {
	lis net.Listener
	id  int

	acceptedMessage map[int]Message
}

func (l *Learner) Learn(args *Message, reply *Acknowledge) error {
	log.Printf("Learner learn: Receive message %v", args)
	a := l.acceptedMessage[args.From]
	if a.Number < args.Number {
		l.acceptedMessage[args.From] = *args
		reply.Ok = true
	} else {
		reply.Ok = false
	}
	return nil
}

func (l *Learner) chosen() any {
	acceptedCount := make(map[int]int)
	for _, m := range l.acceptedMessage {
		acceptedCount[m.Number]++
		if acceptedCount[m.Number] > len(l.acceptedMessage)/2 {
			return m.Value
		}
	}
	return nil
}

func (l *Learner) server() {
	rpcs := rpc.NewServer()
	rpcs.Register(l)
	addr := fmt.Sprintf(":%d", l.id)
	lis, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	l.lis = lis
	go func() {
		for {
			conn, err := l.lis.Accept()
			if err == nil {
				go rpcs.ServeConn(conn)
			} else {
				break
			}
		}
	}()
}

func (l *Learner) Close() {
	l.lis.Close()
}

func newLearner(id int, acceptorIds []int) *Learner {
	l := new(Learner)
	l.id = id
	l.acceptedMessage = make(map[int]Message)

	for _, aid := range acceptorIds {
		l.acceptedMessage[aid] = Message{
			Number: 0,
			Value:  nil,
		}
	}
	l.server()
	return l
}
