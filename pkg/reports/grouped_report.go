package reports

// GroupedReport - Report all deaths by player containing all deaths for each player
type GroupedReport map[string]GroupedInformationReport

type GroupedInformationReport struct {
	TotalKills int32
	Players    []string
	Kills      map[string]int32
}
