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
	outEdges map[byte][]*Edge
	compositionStates StateMap
	accept bool
}

func NewState(label string) *State {
	return &State{label,make(map[byte][]*Edge,0),make(StateMap,0),false}
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

func (this State) findOutState(char byte) *State {
	if edge, ok := this.outEdges[char]; ok {
		return edge[0].destination
	}
	return nil
}

func (this State) findAllOutStates(char byte) StateMap {
	states := make(StateMap, 0)
	if edges, ok := this.outEdges[char]; ok {
		for _, edge := range edges {
			states[edge.destination.label] = edge.destination
		}
	}
	return states
}

func (this State) epsilonClosure() StateMap {
	states := make(StateMap, 0)
	if edges, ok := this.outEdges[EPSILON]; ok {
		fmt.Println("start state: ", this.label)
		for _, edge := range edges {
			states[edge.destination.label] = edge.destination
			fmt.Println("\trecursing into ", edge.destination.label)
			states = Union(states, edge.destination.epsilonClosure())
		}
	}
	states[this.label] = &this
	fmt.Println(this.label, " epsilon closure: ", states)
	return states
}

/* Utility function for numbering states */
func createLabel(size int) string {
	return fmt.Sprintf("S%d", size)
}