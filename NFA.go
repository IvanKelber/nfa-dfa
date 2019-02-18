package regex_parser

import (
	"fmt"
	"errors"
)
const EPSILON = 0
const DOT = byte('.')


type NFA struct {
	input string
	size int
	startState, currentState *State
	states []*State
	alphabet ByteSet
}

func BlankNFA() *NFA {
	startState := NewState(0)
	return &NFA{"", 1, startState, startState, []*State{startState}, ByteSet{make(map[byte]bool,0)}}
}

func NewNFA(input string) *NFA {
	nfa := BlankNFA()
	err := nfa.setPattern(input)
	if err != nil {
		fmt.Printf("Error setting pattern: %q.  NFA is blank.\n",err)
	}
	return nfa
}

/*
   Constructs an NFA associated with the input string in @this.
*/
func (this *NFA) construct() {
	// If the character is a modifier store the modifier
	// for use on the next character
	var modifier func(byte)
	for i := range this.input {
		char := this.input[i];
		switch char {
		case '?':
			modifier = this.optional
			break
		case '*':
			modifier = this.star
			break
		case '+':
			modifier = this.plus
			break
		default:
			if modifier == nil {
				modifier = this.concat
			} 
			modifier(char)
			modifier = nil
			this.alphabet.add(char)
			break
		}
	}
	this.currentState.setAccept(true)
}

func (this NFA) epsilonClosure(states map[string]*State) map[string]*State {
	closed_states := make(map[string]*State,0)
	for _, state := range states {
		closed_states = Union(state.epsilonClosure(), closed_states)
	}
	return closed_states
}

func (this NFA) findNewTransition(states map[string]*State, char byte) map[string]*State {
	closed_states := make(map[string]*State,0)
	for _, state := range states {
		closed_states = Union(state.findAllOutStates(char), closed_states)
	}
	return closed_states
}

func (this NFA) convertToDFA() *NFA {
	return nil
}


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
				return false, errors.New("Modifier is the last character in the string")
			}
		default:
			modifier_count = 0
		}
	}
	return true, nil
}

/*
	If the given pattern is not valid, then do not construct the NFA
*/
func (this *NFA) setPattern(pattern string) error {
	if ok, err := isValid(pattern); !ok {
		return err
	}
	this.input = pattern;
	this.construct()
	return nil
}

/* =========== NFA operations =========*/

/* 
   The concat operation updates the NFA by adding another state only reachable
   from the current state when given a transition of @next
 */
func (this *NFA) concat(next byte) {
	state := NewState(this.size);
	edge := &Edge{next, this.currentState, state}
	this.currentState.addOutEdge(edge)
	state.addInEdge(edge)
	this.currentState = state
	this.states = append(this.states, state)
	this.size++
}

/* 
   The optional operation updates the NFA by adding another state reachable
   from the current state when given either a transition of @next or EPSILON
*/
func (this *NFA) optional(next byte) {
	state := NewState(this.size);
	edge := &Edge{next, this.currentState, state}
	epsilon := &Edge{EPSILON, this.currentState, state}
	this.currentState.addOutEdge(edge)
	this.currentState.addOutEdge(epsilon)
	state.addInEdge(edge)
	state.addInEdge(epsilon)
	this.currentState = state
	this.states = append(this.states, state)
	this.size++
}

/* 
   The star operation updates the NFA by adding another state reachable
   from the current state when given either a transition of @next or EPSILON.
   In addition the new state will recursively direct towards itself given a
   transition of @next.
*/
func (this *NFA) star(next byte) {
	state := NewState(this.size)
	edge := &Edge{next, this.currentState, state}
	epsilon := &Edge{EPSILON, this.currentState, state}
	recursive := &Edge{next, state, state}
	this.currentState.addOutEdge(edge)
	this.currentState.addOutEdge(epsilon)
	state.addInEdge(edge)
	state.addInEdge(epsilon)
	state.addOutEdge(recursive)
	state.addInEdge(recursive)
	this.currentState = state
	this.states = append(this.states, state)
	this.size++
}

/*
   The plus operation updates the NFA by adding another state reachable
   from the current state only when given a transition of @next.
   In addition the new state will recursively direct towards itself given a
   transition of @next.
   Note that this is simply the composition of concat() and star()
*/
func (this *NFA) plus(next byte) {
	this.concat(next)
	this.star(next)
}

/*
*/
func (this NFA) match(str string) bool {
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

/* Print States for debugging purposes*/
func (this NFA) printStates() {
	for _, state := range this.states {
		fmt.Printf("State: %q; accept: %t; outEdges: %v; inEdges: %v\n", state.label, state.accept, state.outEdges, state.inEdges)
	}

	fmt.Printf("Epsilon Closure of %v: %v\n", this.states, this.epsilonClosure(this.states[3].epsilonClosure()))
	fmt.Printf("All outgoing Dot paths of %v: %v\n", this.states, this.findNewTransition(this.states[1].epsilonClosure(), DOT))

}