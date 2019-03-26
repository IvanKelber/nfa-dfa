package regex_parser

import (
	"fmt"
	"errors"
)

// Constants
const EPSILON = 0
const DOT = byte('.')

// typedefs
type ByteSet map[byte]bool
type StateMap map[string]*State

/* Set Union*/
func Union(states1, states2 map[string]*State) map[string]*State {
	for label, state := range states1 {
		states2[label] = state
	}
	return states2
}

/*
	Tests whether a @pattern is a valid.
	Returns an error if invalid
*/
func isValid(pattern string) (bool, error) {
	modifier_count := 0
	for i := range pattern {
		char := pattern[i]
		switch char {
		case '?':
			fallthrough
		case '*':
			fallthrough
		case '+':
			modifier_count++
			if modifier_count == 2 {
				return false, errors.New("Two consecutive modifiers")
			}
			if  i == len(pattern) - 1 {
				return false, errors.New("A modifier is the last character in the string")
			}
		default:
			modifier_count = 0
		}
	}
	return true, nil
}


/*
	Finds the union of states that are reached by applying transition
	@char to all of the states in @states.
*/
func findNewTransitions(states StateMap, char byte) StateMap {
	closed_states := make(StateMap,0)
	for _, state := range states {
		closed_states = Union(state.findAllOutStates(char), closed_states)
	}
	return closed_states
}

/*
	Finds the union of states that are reached by applying epsilon
	to all of the states in @states.
*/
func epsilonClosure(states StateMap) StateMap {
	closed_states := make(StateMap,0)
	for _, state := range states {
		closed_states = Union(state.epsilonClosure(), closed_states)
	}
	return closed_states
}

/* Utility function for numbering states */
func createLabel(size int) string {
	return fmt.Sprintf("S%d", size)
}

func printState(state *State) string {
	outEdges := ""
	for _, edges := range state.outEdges {
		for _, edge := range edges {
			outEdges += fmt.Sprintf(" %q: %q-%p => %q-%p ", edge.char, edge.source.label, edge.source, edge.destination.label, edge.destination)
		}
	}
	return fmt.Sprintf("Address %p; State: %q; accept: %t; outEdges: %v;\n",state, state.label, state.accept, outEdges)

}
