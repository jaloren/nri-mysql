package main

import (
	"fmt"
	"log"
)

func parse(tokens []*token) {
	sections := getSections(tokens)

	for _, t := range sections["INDIVIDUAL BUFFER POOL INFO"] {
		if t.kind == DASHED_UPPPER {
			fmt.Println(t.literal)
		}
	}

}

func main() {
	tokens, err := lex("engine-output.txt")
	if err != nil {
		log.Fatal(err)
	}
	parse(tokens)
}
