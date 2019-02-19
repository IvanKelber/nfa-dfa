package regex_parser

import (
	"fmt"
)

type DFA struct {
	complete_states StateMap
	incomplete_states []*State //Treat this as a queue
	startState *State
}

func BlankDFA() *DFA {
	return &DFA{make(StateMap,0),make([]*State, 0),nil}
}

/*

*/
func NewDFA(nfa NFA) *DFA {
	dfa := BlankDFA()
	dfa.convert(nfa)
	return dfa
}

/*
	Converts a given @nfa into a @DFA.
*/
func (this *DFA) convert(nfa NFA) {
	this.incomplete_states = append(this.incomplete_states, nfa.startState);
	nfa.startState.compositionStates = nfa.startState.epsilonClosure()
	for len(this.incomplete_states) > 0 {
		var popped *State
		popped, this.incomplete_states = this.incomplete_states[0], this.incomplete_states[1:]
		// fmt.Printf("Printing popped state: %q\n", printState(popped))
		states := popped.compositionStates

		dfaState := NewStateUnion(states)
		if state, ok := this.complete_states[dfaState.label]; ok {
			// fmt.Printf("State %q already exists. %p vs %p\n", dfaState.label, state, dfaState)
			dfaState = state
		}
		// fmt.Printf("Printing DFA state before edges: %q\n", printState(dfaState))

		if this.startState == nil {
			this.startState = dfaState
		}
		this.complete_states[dfaState.label] = dfaState;
		for char, _ := range nfa.alphabet {
			transition := epsilonClosure(findNewTransitions(states, char))
			newState := NewStateUnion(transition)
			// fmt.Printf("\tPrinting newDFAstate: %q\n", printState(newState))
			// fmt.Printf("\t%q address: %p\n", newState.label, newState)

			if state, ok := this.complete_states[newState.label]; ok {
				newState = state
			} else {
				this.incomplete_states = append(this.incomplete_states, newState)
				this.complete_states[newState.label] = newState
			}

			edge := &Edge{char, dfaState, newState}
			dfaState.addOutEdge(edge)
		}
		// fmt.Printf("Printing DFA state after edges: %q\n", printState(dfaState))
	}
}

func (this DFA) match(str string) bool {
	currentState := this.startState;
	// this.printStates()
	for i := 0; i < len(str); {
		char := str[i]; 
		// outEdges := ""
		// for _, edges := range currentState.outEdges {
		// 	for _, edge := range edges {
		// 		outEdges += fmt.Sprintf(" %q: %q->%q ", edge.char, edge.source.label, edge.destination.label)
		// 	}
		// }
		// fmt.Printf("current state: %q, edges: %q, current char %q,  edges: ",currentState.label, outEdges, char)
		if edge, ok := currentState.outEdges[char]; ok {
			fmt.Printf("edge found.  traversing to %q\n", edge[0].destination.label)
			currentState = edge[0].destination
			i++
		} else if dot, ok := currentState.outEdges[DOT]; ok {
			currentState = dot[0].destination
			i++
		} else if epsilon, ok := currentState.outEdges[EPSILON]; ok {
			currentState = epsilon[0].destination;
		} else {
			this.printStates()
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
