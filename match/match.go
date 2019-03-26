package main

import (
	"fmt"
	"os"
    "github.com/ivankelber/regex_parser"
)

func match(pattern, str string) bool {
    dfa := regex_parser.NewDFA(pattern)
    return dfa.Match(str)
}

func main() {
	args := os.Args[1:]
	if len(args) != 2 {
		fmt.Println("./match <pattern> <string_to_match>")
		return
	}
	if match(args[0], args[1]) {
		fmt.Printf("Success! '%q' matched '%q'\n",args[0], args[1])
	} else {
		fmt.Printf("Failure! '%q' did not match '%q'\n",args[0], args[1])
	}
	
}
