package reports

import "encoding/json"

// DeathsReport - Report all deathmods containing all deaths for each game by killmod
type DeathsReport map[string]Game

type Game struct {
	KillsByMeans KillModMeans `json:"kills_by_means"`
}

type GameStr struct {
	KillsByMeans map[string]int `json:"kills_by_means"`
}

func (r *DeathsReport) UnmarshalJSON(d []byte) error {
	type GamesByStr map[string]GameStr
	var gamesStr GamesByStr
	var kMod KillMods

	deathsReport := make(map[string]Game, 0)

	if err := json.Unmarshal(d, &gamesStr); err != nil {
		return err
	}

	for game, kills := range gamesStr {
		killsByMeans := make(map[KillMods]int, 0)

		for item, deaths := range kills.KillsByMeans {
			killMean := kMod.GetModByString(item)
			killsByMeans[killMean] = deaths
		}

		deathsReport[game] = Game{
			KillsByMeans: killsByMeans,
		}
	}

	*r = deathsReport

	return nil
}
