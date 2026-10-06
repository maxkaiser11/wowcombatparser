package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/maxkaiser11/wowcombatlogger/internal/combatlog"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: wowparse <path-to-combat-log>")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var lineNum, parsed, skipped, failed int
	var encounters []*combatlog.Encounter
	var current *combatlog.Encounter // nil when we're not inside a boss fight

	for scanner.Scan() {
		lineNum++

		ts, fields, err := combatlog.SplitLine(scanner.Text())
		if err != nil {
			failed++
			log.Printf("line %d: %v", lineNum, err)
			continue
		}

		switch fields[0] {
		case "ENCOUNTER_START":
			if len(fields) < 3 {
				continue
			}
			current = combatlog.NewEncounter(strings.Trim(fields[2], `"`), ts)

		case "ENCOUNTER_END":
			if current == nil || len(fields) < 7 {
				continue
			}
			ms, err := strconv.Atoi(fields[6])
			if err != nil {
				failed++
				log.Printf("line %d: parsing fight length %q: %v", lineNum, fields[6], err)
				continue
			}
			current.EndTime = ts
			current.Success = fields[5] == "1"
			current.Duration = time.Duration(ms) * time.Millisecond
			encounters = append(encounters, current)
			current = nil

		default:
			if current == nil {
				skipped++
				continue // trash or downtime between bosses
			}
			ev, err := combatlog.ParseDamage(ts, fields)
			if errors.Is(err, combatlog.ErrNotDamage) {
				skipped++
				continue
			}
			if err != nil {
				failed++
				log.Printf("line %d: %v", lineNum, err)
				continue
			}
			if ev.SourceFlags&combatlog.FlagPlayer == 0 {
				continue
			}
			current.DamageByPlayer[ev.SourceGUID] += ev.Amount
			current.Names[ev.SourceGUID] = ev.SourceName
			parsed++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	fmt.Printf("parsed: %d, skipped: %d, failed: %d\n", parsed, skipped, failed)

	p := message.NewPrinter(language.English)

	for _, enc := range encounters {
		result := "wipe"
		if enc.Success {
			result = "kill"
		}
		fmt.Printf("\n%s (%s, %s)\n", enc.Name, enc.Duration.Round(time.Second), result)

		seconds := enc.Duration.Seconds()
		if seconds <= 0 {
			fmt.Println("  no duration recorded")
			continue
		}

		for guid, total := range enc.DamageByPlayer {
			dps := float64(total) / seconds
			p.Printf("%-25s %15d damage %12.0f DPS\n", enc.Names[guid], total, dps)
		}
	}
}
