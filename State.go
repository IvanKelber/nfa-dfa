package regex_parser

import (
	"sort"
)

type Edge struct {
	char byte
	source, destination *State
}

func NewEdge(char byte, source, destination *State) *Edge {
	edge :=& Edge{char, source, destination}
	source.addOutEdge(edge)
	return edge
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

/*
	Creates a new state out of existing states
	found in %states.  
	The new state will be:
	 - labeled the concatenation of all state labels in %states
	 - an accept state if and only if one of the states in %states is an accept state
	 - will contain pointers to all of the states used to create it in compositionStates
*/
func NewStateUnion(states StateMap) *State {
	newLabel := ""
	accept := false
	keys := make([]string, 0)
	for key, _ := range states {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, label := range keys {
		accept = accept || states[label].isAccept()
		newLabel += label
	}
	state := NewState(newLabel)
	state.compositionStates = states
	state.setAccept(accept)
	return state
}

func (this *State) setAccept(accept bool) {
	this.accept = accept
}

func (this State) isAccept() bool {
	return this.accept
}

/*
	Appends an edge to this.outEdges[@edge.char], an array
*/
func (this *State) addOutEdge(edge *Edge) {
	if _, ok := this.outEdges[edge.char]; !ok {
		this.outEdges[edge.char] = make([]*Edge,0)
	}
	this.outEdges[edge.char] = append(this.outEdges[edge.char], edge)
}

/*
	Find all states reachable from @this by applying transition
	@char
*/
func (this State) findAllOutStates(char byte) StateMap {
	states := make(StateMap, 0)
	if edges, ok := this.outEdges[char]; ok {
		for _, edge := range edges {
			states[edge.destination.label] = edge.destination
		}
	}
	return states
}

/*
	Find all states reachable from @this by applying epsilon.
	Note that if a reachable state has an outgoing epsilon that
	we need to recursively apply this function.
*/

func (this State) epsilonClosure() StateMap {
	states := make(StateMap, 0)
	if edges, ok := this.outEdges[EPSILON]; ok {
		for _, edge := range edges {
			states[edge.destination.label] = edge.destination
			states = Union(states, edge.destination.epsilonClosure())
		}
	}
	states[this.label] = &this
	return states
}

