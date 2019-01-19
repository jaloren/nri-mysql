package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

const (
	dash = '-'
)

var (
	nameRegex = regexp.MustCompile(`^[A-Z]+[\s]+[A-Z\s]*$`)
)

type output struct {
	sections []section
	scanner  *bufio.Scanner
}

type section struct {
	header, body string
}

func NewOutput(reader io.Reader) *output {
	return &output{
		scanner: bufio.NewScanner(reader),
	}
}

func (o *output) skipScan() bool{
	o.skipEmptyLines()
	return o.scanner.Scan()
}

func(o *output) skipEmptyLines(){
	line := o.scanner.Bytes()
	noSpace := bytes.TrimSpace(line)
	if len(noSpace) == 0 {
		if ! o.scanner.Scan() {
			return
		}
		o.skipEmptyLines()
	}
}

func (o *output) isDashedLine() bool{
	line := o.scanner.Text()
	for _, char := range  line {
		if char != dash {
			return false
		}
	}
	return true
}

func (o *output) getSectionHdr() (string, bool) {
	if !o.isDashedLine() {
		return "", false
	}
	o.skipScan()
	hdr := o.scanner.Text()
	o.skipScan()
	if !o.isDashedLine() {
		return "", false
	}
	if hdr == "" {
		return "", false
	}
	return hdr, true
}

func (o *output) parse() () {
	for o.skipScan() {
		section := &section{}
		hdr, ok := o.getSectionHdr()
		if ok {
			section.header = hdr
			fmt.Println(hdr)
		}
	}


	//scanner.Scan()
	//sectionHdr := scanner.Text()
	//if re.MatchString(sectionHdr){
	//	scanner.Scan()
	//	return sectionHdr, true
	//}
	//return "", false
}

func main() {
	file, err := os.Open("engine-output.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	output := NewOutput(file)
	output.parse()

}
