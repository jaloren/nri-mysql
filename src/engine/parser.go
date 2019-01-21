package main

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/pkg/errors"
	"strings"
)

func combineErrs(prefix string, errs []error) error {
	var msg string
	msg += prefix + ": "
	for idx, e := range errs {
		msg += e.Error()
		if idx != (len(errs) - 1) {
			msg += ": "
		}
	}
	return errors.New(msg)
}


type innodb_status struct {
	sections   map[string][]*token
	lock_stats semaphores
}

func (i *innodb_status) parse_semaphores() {
	tokens, ok := i.sections["semaphores"]
	if !ok {
		log.Warn("semaphore statistics are missing from engine innodb status")
		return
	}
	lock_stats := newSemaphores(tokens)
	lock_stats.setMetrics()
	if len(lock_stats.errs) > 0 {
		log.Error("%s", combineErrs("failed to parse semaphore statistics", lock_stats.errs))
	}

	spew.Dump(lock_stats)
}

func parse(tokens []*token) {
	engine_status := &innodb_status{
		sections: getSections(tokens),
	}
	engine_status.parse_semaphores()
}

func isHdr(idx int, t *token, tokens []*token) bool {
	if len(tokens) < (idx+1) || idx == 0 {
		return false
	}
	return tokens[idx-1].kind == DASHED && t.kind == UPPER && tokens[idx+1].kind == DASHED
}

func getHdrIndices(tokens []*token) []int {
	var hdrIndices []int
	for idx, t := range tokens {
		if isHdr(idx, t, tokens) {
			hdrIndices = append(hdrIndices, idx)
		}
	}
	return hdrIndices
}

func normalizeSectionName(token *token) string {
	lower := strings.ToLower(token.literal)
	underScores := strings.Replace(lower, " ", "_", -1)
	return strings.Replace(underScores, "/", "", -1)
}

func getSections(tokens []*token) map[string][]*token {
	sections := make(map[string][]*token)
	hdrIndices := getHdrIndices(tokens)
	last := hdrIndices[len(hdrIndices)-1]
	for i, hdrIdx := range hdrIndices {
		secName := normalizeSectionName(tokens[hdrIdx])
		if last == hdrIdx {
			sections[secName] = tokens[hdrIdx:]
			continue
		}
		end := hdrIndices[i+1] - 1
		sections[secName] = tokens[hdrIdx+1 : end]
	}
	return sections
}
