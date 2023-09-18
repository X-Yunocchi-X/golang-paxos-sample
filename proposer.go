package paxos

import (
	"fmt"
	"log"
)

type Proposer struct {
	id            int
	round         int
	acceptors     []int
	totalAcceptor int
}

func (p *Proposer) majority() int {
	return p.totalAcceptor/2 + 1
}

func (p *Proposer) proposalNumber() int {
	return p.round<<16 | p.id
}

func (p *Proposer) propose(v string) string {

	log.Printf("Proposer %d proposing value %v\n", p.id, v)

	p.round++
	number := p.proposalNumber()

	// Phase 1: Prepare
	prepareCount := 0
	for _, aid := range p.acceptors {
		args := &Message{
			Number: number,
			From:   p.id,
			To:     aid,
		}
		reply := new(Acknowledge)
		ok := call(fmt.Sprintf("localhost:%d", aid), "Acceptor.Prepare", args, reply)
		if !ok {
			log.Println("Acceptor.Prepare error")
			continue
		}
		if reply.Ok {
			prepareCount++
			if reply.Number > number {
				number = reply.Number
				v = reply.Value.(string)
			}
			if reply.Accepted {
				v = reply.Value.(string)
			}
		} else {
			log.Println("Acceptor.Prepare reply not ok")
			continue
		}

		if prepareCount == p.majority() {
			break
		}
	}

	log.Printf("Proposer: Phase 1: %d/%d/Propose: %s\n", prepareCount, p.majority(), v)

	// Phase 2: accept
	acceptCount := 0
	if prepareCount >= p.majority() {
		for _, aid := range p.acceptors {
			args := &Message{
				Value:  v,
				Number: number,
				From:   p.id,
				To:     aid,
			}
			reply := new(Acknowledge)
			ok := call(fmt.Sprintf("localhost:%d", aid), "Acceptor.Accept", args, reply)
			if !ok {
				log.Println("Acceptor.Accept error")
				continue
			}
			if reply.Ok {
				acceptCount++
			}
			if acceptCount >= p.majority() {
				log.Println("Proposer: Phase 2 success, value:", v)
				return v
			}
		}
	}
	return ""
}
