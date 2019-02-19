package regex_parser

import(
	"fmt"
	"sort"
)

type DFA struct {
	complete_states StateMap
	incomplete_states []*State //Treat this as a queue
	startState *State
}

func NewDFA() *DFA {
	return &DFA{make(StateMap,0),make([]*State, 0),nil}
}

func (this *DFA) convert(nfa NFA) {
	//Calculate E(S0)
	this.incomplete_states = append(this.incomplete_states, nfa.startState);
	this.startState = nfa.startState
	nfa.startState.compositionStates = nfa.startState.epsilonClosure()
	for len(this.incomplete_states) > 0 {
		var popped *State
		popped, this.incomplete_states = this.incomplete_states[0], this.incomplete_states[1:]
		states := popped.compositionStates
		dfaState := createDFAState(states)
		this.complete_states[dfaState.label] = dfaState;
		for char, _ := range nfa.alphabet {
			transition := epsilonClosure(findNewTransitions(states, char))
			newState := createDFAState(transition)
			edge := &Edge{char, dfaState, newState}
			dfaState.addOutEdge(edge)
			if _, ok := this.complete_states[newState.label]; !ok {
				this.incomplete_states = append(this.incomplete_states, newState)
				this.complete_states[newState.label] = newState
			}
		}
	}
}

func createDFAState(states StateMap) *State {
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

/* Print States for debugging purposes*/
func (this DFA) printStates() {
	for _, state := range this.complete_states {
		outEdges := ""
		for _, edges := range state.outEdges {
			for _, edge := range edges {
				outEdges += fmt.Sprintf(" %q: %q->%q ", edge.char, edge.source.label, edge.destination.label)
			}
		}
		fmt.Printf("State: %q; accept: %t; outEdges: %v; compositionStates : %v\n", state.label, state.accept, outEdges, state.compositionStates)
	}
}

func (this DFA) match(str string) bool {
	currentState := this.startState;

	for i := 0; i < len(str); {
		char := str[i]; 
		if edge, ok := currentState.outEdges[char]; ok {
			currentState = edge[0].destination
			i++
		} else if dot, ok := currentState.outEdges[DOT]; ok {
			currentState = dot[0].destination
			i++
		} else if epsilon, ok := currentState.outEdges[EPSILON]; ok {
			currentState = epsilon[0].destination;
		} else {
			return false
		}
	}

	// Iterate through epsilon edges until successful 
	epsilon, ok := currentState.outEdges[EPSILON]
	for ok {
		if currentState.isAccept() {
			return true
		}
		currentState = epsilon[0].destination
		epsilon, ok = currentState.outEdges[EPSILON]
	}

	return currentState.isAccept()
}