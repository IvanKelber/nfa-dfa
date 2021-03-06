package regex_parser

import (
	"fmt"
)

type NFA struct {
	input string
	size int
	startState, currentState *State
	states []*State
	alphabet ByteSet
}

func BlankNFA() *NFA {
	startState := NewState(createLabel(0))
	return &NFA{"", 1, startState, startState, []*State{startState}, make(ByteSet,0)}
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
			this.alphabet[char] = true
			break
		}
	}
	this.currentState.setAccept(true)
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
	state := NewState(createLabel(this.size));
	NewEdge(next, this.currentState, state)
	this.currentState = state
	this.states = append(this.states, state)
	this.size++
}

/* 
   The optional operation updates the NFA by adding another state reachable
   from the current state when given either a transition of @next or EPSILON
*/
func (this *NFA) optional(next byte) {
	state := NewState(createLabel(this.size));
	NewEdge(next, this.currentState, state)
	NewEdge(EPSILON, this.currentState, state)
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
	state := NewState(createLabel(this.size))
	NewEdge(next, this.currentState, state)
	NewEdge(EPSILON, this.currentState, state)
	NewEdge(next, state, state)
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

/* Print States for debugging purposes*/
func (this NFA) printStates() {
	for _, state := range this.states {
		fmt.Print(printState(state))
	}
}