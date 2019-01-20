package main

import (
	"bytes"
	"regexp"
)

var (
	nameRegex = regexp.MustCompile(`^[A-Z]+[\s]+[A-Z\s]*$`)
)

const (
	UPPER TokenType = iota
	DASHED
	PARAGRAPH

	EOI = "END OF INNODB MONITOR OUTPUT"
)

type TokenType int

type token struct {
	kind    TokenType
	literal string
}

func isDashedLine(line string) bool{
	for _, char := range line {
		if char != '-' {
			return false
		}
	}
	return true
}

func isUpperLine(line string) bool{
	return nameRegex.MatchString(line)
}

func isEmptyLine(line []byte) bool{
	noSpace := bytes.TrimSpace(line)
	return len(noSpace) == 0
}