package regex_parser

import (
	"testing"
	"fmt"
)

type MatchTest struct {
	pattern, str string
	expected bool
}

type ValidTest struct {
	pattern string
	expected bool
}


func TestValid(t *testing.T) {
	cases := []ValidTest {
		{"??",false},
		{"**a",false},
		{"",true},
		{"a.cd?d",true},
		{"a*", false},
		{"*.",true},
		{"?+",false},
	}
	
	for _, c := range cases {
		output, _ := isValid(c.pattern)
		if output != c.expected {
			t.Errorf("isValid(%q) == %t, expected %t", c.pattern, output, c.expected)
		}
	}
}

func TestNFA(t *testing.T) {
    cases := []MatchTest {
		{"a.c","abc",true},
		{"a.c","a.c",true},
		{"abc","def",false},
		{"a","",false},
		{"","",true},
		{"?a","", true},
		{"?a","a", true},
		{"?a","b", false},
		{"?.","a",true},
		{"?.a","aa",true},
		{"?abc","bc", true},
		{"?aab", "ab",true},
		{"hell?o w?orld", "hell world", true},
		{"test?s", "test", true},
		{"abc?de", "abcd", false},
		{"aa?bc","tt",false},
		{"aa?a?b?c?d?e","aa", true},
		{"*a","a", true},
		{"*a","", true},
		{"*a","aaaaaaaaaaa", true},
		{"*.", "Now this is the story all about how my life got flipped, turned upside down", true},
		{"*.lol","lmao",false},
		{"*.a","aa",true},
		{"*.ba","aa",false},
		{"abc*.","abc",true},
		{"*abcd*e","aaaaabcd", true},
		{"*abcd*e","aaaaabcdeeeeef", false},
		{"*a*b","", true},
		{"*a*b*c","c", true},
		{"*aab","ab",true},
		{"+a","",false},
		{"+a","a",true},
		{"+a","aaaaaa",true},
		{"+a+b+c+d","abcd",true},
		{"+a+a+a","aaa",true},
		{"+aab","ab",false},
		{"+.","",false},
		{"+.","And I'd like to take a minute just sit right there ", true},
		{"*a+b","ab",true},
		{"*a+b","a",false},
		{"*a+b","",false},
		{"*a+b","b",true},
		{"?a+b","b",true},
		{"?a*a","",true},
		{"?a*a","a",true},
		{"?a*a","aa",true},
		{"?a*b","ab",true},
		{"?a*b","aab",false},
		{"?a*b+c", "", false},
		{"?a*b+c","abc", true},
		{"*.+.","", false},
		{"*.+.","I'll tell you how I became the prince of a town called Bel-Air", true},
    }

	var nfa *NFA
	for _, c := range cases {
		nfa = NewNFA(c.pattern)
		dfa := NewDFA(*nfa)
		output := dfa.match(c.str)
		if output != c.expected {
			// nfa.printStates()
			// fmt.Println("===========")
			// dfa.printStates()
			t.Errorf("NFA[pattern == %q].match(%q) == %t, expected %t",
									c.pattern, c.str, output, c.expected)
		}
	}
	// nfa.printStates()

}

func TestDFA(t *testing.T) {
	// nfa := NewNFA("+b"); 
	// nfa.printStates()
	fmt.Println("====")
	// dfa := NewDFA(nfa)
	// dfa.printStates()
}

