package api

import (
	"reflect"
	"testing"
	"time"

	"github.com/maxkaiser11/wowcombatparser/internal/combatlog"
)

func Test_toEncounterView(t *testing.T) {
	type args struct {
		encounter *combatlog.Encounter
	}
	tests := []struct {
		name      string
		encounter *combatlog.Encounter
		want      EncounterView
	}{
		// TODO: Add test cases.
		{name: "sorts players by damage", encounter: makeEncounter(t, 10*time.Second, []testPlayer{
			{guid: "p1", name: "Bob", damage: 500},
			{guid: "p2", name: "Alice", damage: 1000},
		}), want: EncounterView{
			Name:     "Test Boss",
			Duration: 10 * time.Second,
			Outcome:  "wipe",
			PlayerStats: []PlayerStat{
				{Name: "Alice", Damage: 1000, DPS: 100},
				{Name: "Bob", Damage: 500, DPS: 50},
			},
		},
		},
		{
			name: "zero duration gives zero DPS",
			encounter: makeEncounter(t, 0, []testPlayer{
				{guid: "p1", name: "Bob", damage: 500},
			}),
			want: EncounterView{
				Name:     "Test Boss",
				Duration: 0,
				Outcome:  "wipe",
				PlayerStats: []PlayerStat{
					{Name: "Bob", Damage: 500, DPS: 0},
				},
			},
		},
		{
			name:      "no players",
			encounter: makeEncounter(t, 10*time.Second, nil),
			want: EncounterView{
				Name:        "Test Boss",
				Duration:    10 * time.Second,
				Outcome:     "wipe",
				PlayerStats: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toEncounterView(tt.encounter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toEncounterView() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testPlayer struct {
	guid   string
	name   string
	damage int
}

func makeEncounter(t *testing.T, duration time.Duration, players []testPlayer) *combatlog.Encounter {
	t.Helper()

	enc := combatlog.NewEncounter("Test Boss", time.Time{})
	enc.Duration = duration

	for _, p := range players {
		enc.AddDamage(combatlog.DamageEvent{
			SourceGUID: p.guid,
			SourceName: p.name,
			Amount:     p.damage,
		})
	}
	return enc
}
