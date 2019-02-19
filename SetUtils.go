package regex_parser


type ByteSet map[byte]bool
type StateMap map[string]*State

func Union(states1, states2 map[string]*State) map[string]*State {
	for label, state := range states1 {
		states2[label] = state
	}
	return states2
}
