package paxos

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
)

type Acceptor struct {
	lis            net.Listener
	id             int
	minNumber      int
	accepted       bool
	acceptedNumber int
	acceptedValue  any

	learners []int
}

func (a *Acceptor) Prepare(args *Message, reply *Acknowledge) error {
	log.Printf("Acceptor prepare: Receive message %v", args)
	if args.Number > a.minNumber {
		a.minNumber = args.Number
		reply.Ok = true
		reply.Accepted = false
		if a.accepted {
			reply.Accepted = true
			reply.Number = a.acceptedNumber
			reply.Value = a.acceptedValue
		}
	} else {
		reply.Ok = false
	}
	return nil
}

func (a *Acceptor) Accept(args *Message, reply *Acknowledge) error {
	log.Printf("Acceptor accept: Receive message %v", args)
	if args.Number >= a.minNumber {
		a.accepted = true
		a.minNumber = args.Number
		a.acceptedNumber = args.Number
		a.acceptedValue = args.Value
		reply.Ok = true
		length := len(a.learners)

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(length)

		for _, learner := range a.learners {
			go func(learner int) {
				addr := fmt.Sprintf("localhost:%d", learner)
				args.From = a.id
				args.To = learner
				resp := new(Acknowledge)
				ok := call(addr, "Learner.Learn", args, resp)
				waitGroup.Done()
				if !ok {
					return
				}
			}(learner)
		}
		waitGroup.Wait()
	} else {
		reply.Ok = false
	}
	return nil
}

func (a *Acceptor) server() {
	rpcs := rpc.NewServer()
	rpcs.Register(a)
	addr := fmt.Sprintf(":%d", a.id)
	l, e := net.Listen("tcp", addr)
	if e != nil {
		panic(e)
	}
	a.lis = l
	go func() {
		for {
			conn, err := a.lis.Accept()
			if err != nil {
				continue
			}
			go rpcs.ServeConn(conn)
		}
	}()
}

func (a *Acceptor) Close() {
	a.lis.Close()
}

func newAcceptor(id int, learners []int) *Acceptor {
	acceptor := &Acceptor{
		id:       id,
		learners: learners,
	}
	acceptor.server()
	return acceptor
}
