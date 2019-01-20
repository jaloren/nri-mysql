package main

import (
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

func combineErrs(prefix string, errs []error) error {
	var msg string
	msg += prefix + ": "
	for _, e := range errs {
		msg += e.Error() + ": "
	}
	return errors.New(msg)
}

type semaphores struct {
	tokens                        []*token
	errs                          []error
	reservation                   int64
	signal                        int64
	spin_waits                    int64
	spin_rounds                   int64
	spin_os_waits                 int64
	rwshared_spin_waits           int64
	rwshared_spin_rounds          int64
	rwshared_spin_os_waits        int64
	rwexcl_spin_waits             int64
	rwexcl_spin_rounds            int64
	rwexcl_spin_os_waits          int64
	spin_rounds_per_wait          float64
	rwshared_spin_rounds_per_wait float64
	rwexcl_spin_rounds_per_wait   float64
}

func (s *semaphores) setCounts() {
	var (
		err                           error
		re                            = regexp.MustCompile(`.*(reservation|signal) count[\s]*([0-9]+)$`)
		isReservationSet, isSignalSet bool
	)
	for _, t := range s.tokens {
		if t.kind != PARAGRAPH {
			continue
		}
		line := t.literal
		matches := re.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		last := matches[len(matches)-1]
		switch matches[1] {
		case "reservation":
			s.reservation, err = strconv.ParseInt(last, 0, 64)
			if err != nil {
				s.errs = append(s.errs,
					errors.Wrapf(err, "failed to get reservation count from line %q", line))
			}
			isReservationSet = true
		case "signal":
			s.signal, err = strconv.ParseInt(last, 0, 64)
			if err != nil {
				s.errs = append(s.errs,
					errors.Wrapf(err, "failed to get signal count from line %q", line))
			}
			isSignalSet = true
		}
		if isSignalSet && isReservationSet {
			return
		}
	}
	s.errs = append(s.errs, errors.New("no line from the status output contained reservation or signal count"))
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
	lock_stats := &semaphores{tokens: tokens}
	lock_stats.setCounts()
	if len(lock_stats.errs) > 0 {
		log.Error("%s", combineErrs("failed to parse semaphore statistics", lock_stats.errs))
	}

	fmt.Printf("%d %d",lock_stats.reservation, lock_stats.signal)
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
