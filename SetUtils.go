package regex_parser


type ByteSet struct {
	set map[byte]bool
}

func (this *ByteSet) add(b byte) {
	this.set[b] = true
}

func (this ByteSet) contains(b byte) bool {
	_, ok := this.set[b]
	return ok
}

type StateSet struct {
	set map[string]*State
}

func (this *StateSet) add(s *State) {
	this.set[s.label] = s
}

func (this StateSet) contains(s State) bool {
	_, ok := this.set[s.label]
	return ok
}

func Union(states1, states2 map[string]*State) map[string]*State {
	for label, state := range states1 {
		states2[label] = state
	}
	return states2
}
