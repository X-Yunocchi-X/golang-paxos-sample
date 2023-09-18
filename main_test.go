package paxos

import (
	"log"
	"testing"
)

func start(acceptorsIds, learnerIds []int) ([]*Acceptor, []*Learner) {
	acceptors := make([]*Acceptor, 0)
	for _, aid := range acceptorsIds {
		acceptors = append(acceptors, newAcceptor(aid, learnerIds))
	}

	learners := make([]*Learner, 0)
	for _, lid := range learnerIds {
		learners = append(learners, newLearner(lid, acceptorsIds))
	}

	return acceptors, learners
}

func cleanup(acceptors []*Acceptor, learners []*Learner) {
	for _, a := range acceptors {
		a.Close()
	}
	for _, l := range learners {
		l.Close()
	}
}

func TestSingleProposer(t *testing.T) {
	acceptorIds := []int{9901, 9902, 9903}
	learnerIds := []int{9006}
	acceptors, learners := start(acceptorIds, learnerIds)
	defer cleanup(acceptors, learners)

	proposer := &Proposer{
		id:        1,
		round:    0,
		acceptors: acceptorIds,
		totalAcceptor: len(acceptorIds),
	}

	var value string
	value = proposer.propose("Hello, world")
	log.Printf("Value: %s", value)
	if value != "Hello, world" {
		t.Errorf("Expected 'Hello, world', got %s", value)
	}

	learnValue := learners[0].chosen()
	if learnValue != "Hello, world" {
		t.Errorf("Expected 'Hello, world', got %v", learnValue)
	}
}

func TestTwoProposers(t *testing.T) {
	acceptorIds := []int{9001, 9002, 9003}
	learnerIds := []int{9006}
	acceptors, learners := start(acceptorIds, learnerIds)
	defer cleanup(acceptors, learners)

	proposer1 := &Proposer{
		id:        1,
		round:    0,
		acceptors: acceptorIds,
		totalAcceptor: len(acceptorIds),
	}

	proposer2 := &Proposer{
		id:        2,
		round:    0,
		acceptors: acceptorIds,
		totalAcceptor: len(acceptorIds),
	}

	value1 := proposer1.propose("Hello, world")
	// time.Sleep(1 * time.Second)
	value2 := proposer2.propose("Hello, world 2")

	if value1 != value2 {
		t.Errorf("Expected %v, got %v", value1, value2)
	}

	learnValue := learners[0].chosen()
	if learnValue != value1 {
		t.Errorf("Expected %v, got %v", value1, learnValue)
	}
}
