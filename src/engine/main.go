package main

import (
	"log"
)



func main() {
	tokens, err := lex("engine-output.txt")
	if err != nil {
		log.Fatal(err)
	}
	parse(tokens)
}
