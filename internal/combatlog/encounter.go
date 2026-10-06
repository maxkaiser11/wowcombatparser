package combatlog

import (
	"fmt"
	"strconv"
	"time"
)

// Encounter is one boss pull, from ENCOUNTER_START to ENCOUNTER_END.
type Encounter struct {
	Name           string
	StartTime      time.Time
	EndTime        time.Time
	Success        bool
	Duration       time.Duration
	DamageByPlayer map[string]int    // player GUID -> total damage
	Names          map[string]string // player GUID -> player name
}

func NewEncounter(name string, start time.Time) *Encounter {
	return &Encounter{
		Name:           name,
		StartTime:      start,
		DamageByPlayer: make(map[string]int),
		Names:          make(map[string]string),
	}
}

// Finish fills in the end time, result and duration from an ENCOUNTER_END line.
func (e *Encounter) Finish(ts time.Time, fields []string) error {
	if len(fields) < 7 {
		return fmt.Errorf("ENCOUNTER_END line has %d fields, need at least 7", len(fields))
	}
	ms, err := strconv.Atoi(fields[6])
	if err != nil {
		return fmt.Errorf("parsing fight length %q: %w", fields[6], err)
	}
	e.EndTime = ts
	e.Success = fields[5] == "1"
	e.Duration = time.Duration(ms) * time.Millisecond
	return nil
}

// AddDamage adds one damage event to the player's total.
func (e *Encounter) AddDamage(ev DamageEvent) {
	e.DamageByPlayer[ev.SourceGUID] += ev.Amount
	e.Names[ev.SourceGUID] = ev.SourceName
}

// Outcome returns "kill" or "wipe".
func (e *Encounter) Outcome() string {
	if e.Success {
		return "kill"
	}
	return "wipe"
}
