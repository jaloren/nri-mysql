package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func lex(filePath string) ([]*token, error){
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tokens []*token

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		if isEmptyLine(scanner.Bytes()){
			continue
		}
		t := &token{literal: scanner.Text()}
		if t.literal == EOI {
			return tokens, nil
		}
		if isDashedLine(t.literal) {
			t.kind = DASHED
		} else if isUpperLine(t.literal) {
			t.kind = UPPER
		} else {
			t.kind = PARAGRAPH
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

func main() {
	tokens, err := lex("engine-output.txt")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range tokens {
		if t.kind == UPPER {
			fmt.Println(t.literal)
		}
	}
}
