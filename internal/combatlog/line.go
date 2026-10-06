package combatlog

import (
	"fmt"
	"strings"
	"time"
)

const TimeLayout = "1/2/2006 15:04:05.0000"

// SplitLine separates the timestamp from the comma-separated fields.
func SplitLine(line string) (time.Time, []string, error) {
	tstmp, rest, found := strings.Cut(line, "  ")
	if !found {
		return time.Time{}, nil, fmt.Errorf("no timestamp separator in line %q", line)
	}

	ts, err := time.Parse(TimeLayout, tstmp)
	if err != nil {
		return time.Time{}, nil, fmt.Errorf("parsing timestamp: %w", err)
	}

	rest = strings.TrimRight(rest, "\r") // the log has Windows line endings
	return ts, strings.Split(rest, ","), nil
}
