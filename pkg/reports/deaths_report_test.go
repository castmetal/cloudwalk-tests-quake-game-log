package reports

import (
	"encoding/json"
	"testing"
)

func TestDeathsReportInterface(t *testing.T) {

	testCases := []struct {
		description string
		reportInput string
		expected    int32
	}{
		{
			description: "Testing marshalling report input and if it's given the MOD_SHOTGUN with the same value",
			reportInput: `{"game-1":{"kills_by_means":{"MOD_GAUNTLET":1,"MOD_RAILGUN":2,"MOD_SHOTGUN":10}}}`,
			expected:    10,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			jsonData := tc.reportInput
			var deathsReport DeathsReport

			err := json.Unmarshal([]byte(jsonData), &deathsReport) // Testing expected input
			if err != nil {
				t.Errorf("unmarshal error = %v", err)
			}

			if deathsReport["game-1"].KillsByMeans[MOD_SHOTGUN] != tc.expected {
				t.Errorf("error testing MOD_SHOTGUN = %d, want %d", deathsReport["game-1"].KillsByMeans[MOD_SHOTGUN], tc.expected)

			}

			b, err := json.Marshal(&deathsReport) // Testing Marshalling
			if err != nil {
				t.Errorf("marshalling error = %v", err)
			}

			reportStr := string(b)
			if reportStr != tc.reportInput {
				t.Errorf("marshalling comparison error = %s, want %s", reportStr, tc.reportInput)
			}

		})
	}

}
