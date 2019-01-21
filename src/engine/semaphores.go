package main

import (
	"github.com/pkg/errors"
	"regexp"
	"strconv"
)

type semaphores struct {
	tokens                        []*token
	errs                          []error
	cntSet                        map[string]struct{}
	reservation                   int64
	signal                        int64
	spinWaits                     int64
	spinRounds                    int64
	spinOsWaits                   int64
	rwsharedSpinWaits           int64
	rwsharedSpinRounds          int64
	rwsharedSpinOsWaits        int64
	rwexclSpinWaits             int64
	rwexclSpinRounds            int64
	rwexclSpinOsWaits          int64
	spinRoundsPerWait          float64
	rwsharedSpinRoundsPerWait float64
	rwexclSpinRoundsPerWait   float64
}

func newSemaphores(tokens []*token) *semaphores{
	return &semaphores{
		tokens: tokens,
		cntSet: make(map[string]struct{}),
	}
}

func (s *semaphores) setMetrics() {
	for _, t := range s.tokens {
		if t.kind != PARAGRAPH {
			continue
		}

		key := "reservationCount"
	 	if _, ok := s.cntSet[key]; !ok {
			s.setReservation(t.literal, key)
		}

	 	key = "signalCount"
		if _, ok := s.cntSet[key]; !ok {
			s.setSignal(t.literal, key)
		}

		key = "spinMetrics"
		if _, ok := s.cntSet[key]; !ok {
			s.setSignal(t.literal, key)
		}
	}
	expectedKeys := []string{"reservationCount", "signalCount", "spinCounts"}
	for _, k := range expectedKeys {
		if _, ok := s.cntSet[k]; !ok {
			s.errs = append(s.errs, errors.Errorf("failed to set semaphore metrics %s", k))
		}
	}
}

func (s *semaphores) setReservation(line, key string) {
	var (
		err error
		re  = regexp.MustCompile(`.*reservation count ([0-9]+)$`)
	)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return
	}
	last := matches[len(matches)-1]
	s.reservation, err = strconv.ParseInt(last, 0, 64)
	if err != nil {
		s.errs = append(s.errs,
			errors.Wrapf(err, "failed to get reservation count from line %q", line))
		return
	}
	s.cntSet[key] = struct{}{}
}


func (s *semaphores) setSignal(line, key string) {
	var (
		err error
		re  = regexp.MustCompile(`.*signal count ([0-9]+)$`)
	)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return
	}
	last := matches[len(matches)-1]
	s.signal, err = strconv.ParseInt(last, 0, 64)
	if err != nil {
		s.errs = append(s.errs,
			errors.Wrapf(err, "failed to get signal count from line %q", line))
		return
	}
	s.cntSet[key] = struct{}{}
}

func (s *semaphores) setSpin(line, key string) {
	var (
		err error
		re  = regexp.MustCompile(`^Mutex spin waits ([0-9]+), rounds ([0-9]+), OS waits ([0-9]+)$`)
	)
	matches := re.FindStringSubmatch(line)
	if matches == nil {
		return
	}

	spinWaits := matches[1]
	rounds := matches[2]
	osWaits := matches[3]

	s.spinWaits, err = strconv.ParseInt(spinWaits, 0, 64)
	if err != nil {
		s.errs = append(s.errs,
			errors.Wrapf(err, "failed to get spin wait metric from line %q", line))
		return
	}

	s.spinRounds, err = strconv.ParseInt(rounds, 0, 64)
	if err != nil {
		s.errs = append(s.errs,
			errors.Wrapf(err, "failed to get spin rounds metric from line %q", line))
		return
	}

	s.spinOsWaits, err = strconv.ParseInt(osWaits, 0, 64)
	if err != nil {
		s.errs = append(s.errs,
			errors.Wrapf(err, "failed to get spin os wait metric from line %q", line))
		return
	}
	s.cntSet[key] = struct{}{}
}