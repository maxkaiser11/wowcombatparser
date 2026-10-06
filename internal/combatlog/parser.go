package combatlog

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Result is everything Parse found in a combat log.
type Result struct {
	Encounters []*Encounter
	Parsed     int     // player damage events counted
	Skipped    int     // lines that aren't relevant (other events, trash)
	Failed     int     // lines that couldn't be parsed
	Errors     []error // details for the failed lines
}

// Parse reads a combat log and splits it into boss encounters with damage totals.
func Parse(r io.Reader) (*Result, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	res := &Result{}
	var current *Encounter // nil when we're not inside a boss fight
	lineNum := 0

	fail := func(err error) {
		res.Failed++
		res.Errors = append(res.Errors, fmt.Errorf("line %d: %w", lineNum, err))
	}

	for scanner.Scan() {
		lineNum++

		ts, fields, err := SplitLine(scanner.Text())
		if err != nil {
			fail(err)
			continue
		}

		switch fields[0] {
		case "ENCOUNTER_START":
			if len(fields) < 3 {
				fail(fmt.Errorf("ENCOUNTER_START line has %d fields, need at least 3", len(fields)))
				continue
			}
			current = NewEncounter(strings.Trim(fields[2], `"`), ts)

		case "ENCOUNTER_END":
			if current == nil {
				continue
			}
			if err := current.Finish(ts, fields); err != nil {
				fail(err)
				continue
			}
			res.Encounters = append(res.Encounters, current)
			current = nil

		default:
			if current == nil {
				res.Skipped++
				continue // trash or downtime between bosses
			}
			ev, err := ParseDamage(ts, fields)
			if errors.Is(err, ErrNotDamage) {
				res.Skipped++
				continue
			}
			if err != nil {
				fail(err)
				continue
			}
			if ev.SourceFlags&FlagPlayer == 0 {
				continue
			}
			current.AddDamage(ev)
			res.Parsed++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading log: %w", err)
	}
	return res, nil
}
