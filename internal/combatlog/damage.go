package combatlog

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const FlagPlayer = 0x400

var ErrNotDamage = errors.New("not a damage event")

type DamageEvent struct {
	Timestamp   time.Time
	EventType   string
	SourceGUID  string
	SourceName  string
	SourceFlags uint64
	Amount      int
}

// amountIndex returns where the damage amount sits for each damage event type.
func amountIndex(eventType string) (int, bool) {
	switch eventType {
	case "SPELL_DAMAGE", "SPELL_PERIODIC_DAMAGE", "RANGE_DAMAGE":
		return 31, true
	case "SWING_DAMAGE":
		return 28, true
	default:
		return 0, false
	}
}

// ParseDamage turns the fields of one damage line into a DamageEvent.
func ParseDamage(ts time.Time, fields []string) (DamageEvent, error) {
	idx, ok := amountIndex(fields[0])
	if !ok {
		return DamageEvent{}, ErrNotDamage
	}
	if len(fields) <= idx {
		return DamageEvent{}, fmt.Errorf("%s line has %d fields, need at least %d", fields[0], len(fields), idx+1)
	}

	flags, err := strconv.ParseUint(fields[3], 0, 64)
	if err != nil {
		return DamageEvent{}, fmt.Errorf("parsing source flags %q: %w", fields[3], err)
	}

	amount, err := strconv.Atoi(fields[idx])
	if err != nil {
		return DamageEvent{}, fmt.Errorf("parsing amount %q: %w", fields[idx], err)
	}

	return DamageEvent{
		Timestamp:   ts,
		EventType:   fields[0],
		SourceGUID:  fields[1],
		SourceName:  strings.Trim(fields[2], `"`),
		SourceFlags: flags,
		Amount:      amount,
	}, nil
}
