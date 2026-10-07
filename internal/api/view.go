package api

import (
	"sort"
	"time"

	"github.com/maxkaiser11/wowcombatparser/internal/combatlog"
)

type PlayerStat struct {
	Name   string  `json:"name"`
	Damage int     `json:"damage"`
	DPS    float64 `json:"dps"`
}

type EncounterView struct {
	Name        string        `json:"name"`
	Duration    time.Duration `json:"duration"`
	Outcome     string        `json:"outcome"`
	PlayerStats []PlayerStat  `json:"player_stats"`
}

func toEncounterView(encounter *combatlog.Encounter) EncounterView {

	var playerStats []PlayerStat

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
			DPS:    dps,
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
