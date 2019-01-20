package main

import (
	"bufio"
	"bytes"
	"os"
	"regexp"
	"strings"
)

var (
	nameRegex = regexp.MustCompile(`^[A-Z]+[A-Z\s/]*$`)
)

const (
	UPPER TokenType = iota
	DASHED
	PARAGRAPH
	DASHED_UPPPER

	EOM = "END OF INNODB MONITOR OUTPUT"
)

type TokenType int

type token struct {
	kind    TokenType
	literal string
}

func isDashedLine(line string) bool {
	for _, char := range line {
		if char != '-' {
			return false
		}
	}
	return true
}

func isUpperLine(line string) bool {
	return nameRegex.MatchString(line)
}

func isEmptyLine(line []byte) bool {
	noSpace := bytes.TrimSpace(line)
	return len(noSpace) == 0
}

func lex(filePath string) ([]*token, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tokens []*token
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if isEmptyLine(scanner.Bytes()) {
			continue
		}
		t := &token{literal: scanner.Text()}
		if t.literal == EOM {
			return tokens, nil
		}
		if isDashedLine(t.literal) {
			t.kind = DASHED
		} else if isUpperLine(t.literal) {
			t.kind = UPPER
		} else if strings.HasPrefix(t.literal, "-"){
			t.kind = DASHED_UPPPER
		} else {
			t.kind = PARAGRAPH
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}
