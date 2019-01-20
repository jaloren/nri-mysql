package main

import (
	"fmt"
	"log"
	"strings"
)

func build(tokens []*token) string{
	body := strings.Builder{}
	for idx, t := range tokens{
		if isHeader(idx,t, tokens) {
			return body.String()
		}
		if t.kind != PARAGRAPH {
			continue
		}
		body.WriteString(t.literal + "\n")
	}
	return body.String()
}

func parse(tokens []*token) {
	//sections := make(map[string]string)
	for idx, t := range tokens {
		if isHeader(idx, t, tokens){
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
