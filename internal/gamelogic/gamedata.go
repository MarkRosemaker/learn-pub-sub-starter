package gamelogic

import (
	"fmt"
	"strings"
)

type Player struct {
	Username string
	Units    map[int]Unit
}

type UnitRank string

const (
	RankInfantry  = "infantry"
	RankCavalry   = "cavalry"
	RankArtillery = "artillery"
)

type Unit struct {
	ID       int
	Rank     UnitRank
	Location Location
}

func (u Unit) String() string {
	return fmt.Sprintf("%s (ID: %d)", u.Rank, u.ID)
}

type ArmyMove struct {
	Player     Player
	Units      []Unit
	ToLocation Location
}

func (m ArmyMove) String() string {
	units := make([]string, len(m.Units))
	for i, u := range m.Units {
		units[i] = u.String()
	}

	return fmt.Sprintf("%s -> %s", strings.Join(units, ", "), m.ToLocation)
}

type RecognitionOfWar struct {
	Attacker Player
	Defender Player
}

type Location string

func getAllRanks() map[UnitRank]struct{} {
	return map[UnitRank]struct{}{
		RankInfantry:  {},
		RankCavalry:   {},
		RankArtillery: {},
	}
}

func getAllLocations() map[Location]struct{} {
	return map[Location]struct{}{
		"americas":   {},
		"europe":     {},
		"africa":     {},
		"asia":       {},
		"australia":  {},
		"antarctica": {},
	}
}
