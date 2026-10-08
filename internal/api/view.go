package api

import (
	"sort"
	"time"

	"github.com/maxkaiser11/wowcombatparser/internal/combatlog"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type PlayerStat struct {
	Name   string `json:"name"`
	Damage int    `json:"damage"`
	DPS    string `json:"dps"`
}

type EncounterView struct {
	Name        string        `json:"name"`
	Duration    time.Duration `json:"duration"`
	Outcome     string        `json:"outcome"`
	PlayerStats []PlayerStat  `json:"player_stats"`
}

func toEncounterView(encounter *combatlog.Encounter) EncounterView {

	var playerStats []PlayerStat

	p := message.NewPrinter(language.English)

	fightDuration := encounter.Duration.Seconds()

	for guid, total := range encounter.DamageByPlayer {
		name := encounter.Names[guid]
		var dps float64
		if fightDuration > 0 {
			dps = float64(total) / fightDuration
		}

		playerStats = append(playerStats, PlayerStat{
			Name:   name,
			Damage: total,
			DPS:    p.Sprintf("%.0f", dps), // this is a string so it can be formatted to 100,000 (with commas)
		})
	}

	sort.Slice(playerStats, func(i, j int) bool {
		return playerStats[i].Damage > playerStats[j].Damage
	})

	return EncounterView{
		Name:        encounter.Name,
		Duration:    encounter.Duration.Round(time.Second),
		Outcome:     encounter.Outcome(),
		PlayerStats: playerStats,
	}
}
