package regex_parser

import (
	"fmt"
)

type DFA struct {
	complete_states StateMap
	startState *State
}

func BlankDFA() *DFA {
	return &DFA{make(StateMap,0),nil}
}

/*
	Given an @nfa return a DFA
*/
func NewDFA(nfa NFA) *DFA {
	dfa := BlankDFA()
	dfa.convert(nfa)
	return dfa
}

/*
	Converts a given @nfa into a DFA using the epsilon-closure algorithm.
	Using this algorithm we eliminate non-deterministic behavior caused by
	epsilon edges by creating new states from unions of old states.  Each
	state in our new DFA will be complete, meaning that for each state,
	every transition in our alphabet is accounted for with an edge.
*/
func (this *DFA) convert(nfa NFA) {
	// Compute epsilon closure of nfa.startState and push into queue
	this.startState = NewStateUnion(nfa.startState.epsilonClosure())
	queue := []*State{this.startState}

	for len(queue) > 0 {
		// Pop a state from the queue and add it to this.complete_states
		var poppedState *State
		poppedState, queue = queue[0], queue[1:]
		this.complete_states[poppedState.label] = poppedState;

		for transition, _ := range nfa.alphabet {
			// Find the epsilon closure of the states reachable by applying a transition to poppedState.
			reachables := epsilonClosure(findNewTransitions(poppedState.compositionStates, transition))

			// If that state is unseen then add it to the queue this.complete_states
			newState := NewStateUnion(reachables)
			if state, ok := this.complete_states[newState.label]; ok {
				newState = state
			} else {
				queue = append(queue, newState)
				this.complete_states[newState.label] = newState
			}

			// Attach poppedState to newState
			edge := &Edge{transition, poppedState, newState}
			poppedState.addOutEdge(edge)
		}
	}
}

/*
	Starting at this.startState, we apply each character in @str
	as a transition.  After every character is applied, if our 
	state is an accept state we have matched.
*/
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
	return currentState.isAccept()
}

/* Print States for debugging purposes*/
func (this DFA) printStates() {
	for _, state := range this.complete_states {
		fmt.Print(printState(state))
	}
}
