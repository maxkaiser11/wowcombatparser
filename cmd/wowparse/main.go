package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/maxkaiser11/wowcombatparser/internal/combatlog"
)

var (
	killStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42")).MarginTop(1)  // green
	wipeStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196")).MarginTop(1) // red
	borderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("99"))                          // purple
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212")).Padding(0, 1)
	cellStyle   = lipgloss.NewStyle().Padding(0, 1)
	oddRowStyle = cellStyle.Foreground(lipgloss.Color("245")) // grey, for zebra stripes
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

	res, err := combatlog.Parse(file)
	if err != nil {
		log.Fatal(err)
	}

	for _, err := range res.Errors {
		log.Println(err)
	}
	fmt.Printf("parsed: %d, skipped: %d, failed: %d\n", res.Parsed, res.Skipped, res.Failed)

	p := message.NewPrinter(language.English)
	for _, enc := range res.Encounters {
		printEncounter(p, enc)
	}
}

// printEncounter prints one fight's header and DPS table.
func printEncounter(p *message.Printer, enc *combatlog.Encounter) {
	titleStyle := wipeStyle
	if enc.Success {
		titleStyle = killStyle
	}
	title := fmt.Sprintf("%s · %s · %s", enc.Name, enc.Duration.Round(time.Second), enc.Outcome())
	fmt.Println(titleStyle.Render(title))

	seconds := enc.Duration.Seconds()
	if seconds <= 0 {
		fmt.Println("  no duration recorded")
		return
	}

	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(borderStyle).
		Headers("Player", "Damage", "DPS").
		StyleFunc(func(row, col int) lipgloss.Style {
			var s lipgloss.Style
			switch {
			case row == table.HeaderRow:
				s = headerStyle
			case row%2 == 1:
				s = oddRowStyle
			default:
				s = cellStyle
			}
			if col > 0 {
				s = s.Align(lipgloss.Right) // numbers line up on the right
			}
			return s
		})

	for guid, total := range enc.DamageByPlayer {
		dps := float64(total) / seconds
		t.Row(enc.Names[guid], p.Sprintf("%d", total), p.Sprintf("%.0f", dps))
	}

	fmt.Println(t)
}
