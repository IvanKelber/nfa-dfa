package regex_parser

import (
	"fmt"
)

type Edge struct {
	char byte
	source, destination *State
}

type State struct {
	label string
	outEdges, inEdges map[byte][]*Edge
	accept bool
}

func NewState(size int) *State {
	return &State{createLabel(size),make(map[byte][]*Edge,0),make(map[byte][]*Edge,0),false}
}

func (this *State) setAccept(accept bool) {
	this.accept = accept
}

func (this State) isAccept() bool {
	return this.accept
}

func (this *State) addOutEdge(edge *Edge) {
	if _, ok := this.outEdges[edge.char]; !ok {
		this.outEdges[edge.char] = make([]*Edge,0)
	}
	this.outEdges[edge.char] = append(this.outEdges[edge.char], edge)
}

func (this *State) addInEdge(edge *Edge) {
	if _, ok := this.inEdges[edge.char]; !ok {
		this.inEdges[edge.char] = make([]*Edge,0)
	}
	this.inEdges[edge.char] = append(this.inEdges[edge.char], edge)
}

func (this State) findOutState(char byte) *State {
	if edge, ok := this.outEdges[char]; ok {
		return edge[0].destination
	}
	return nil
}

func (this State) findAllOutStates(char byte) map[string]*State {
	states := make(map[string]*State, 0)
	if edges, ok := this.outEdges[char]; ok {
		for _, edge := range edges {
			states[edge.destination.label] = edge.destination
		}
	}
	return states
}

func (this State) epsilonClosure() map[string]*State {
	states := this.findAllOutStates(EPSILON);
	states[this.label] = &this
	return states
}

/* Utility function for numbering states */
func createLabel(size int) string {
	return fmt.Sprintf("S%d", size)
}