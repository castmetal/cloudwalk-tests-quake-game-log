package reports

import (
	"encoding/json"
	"fmt"
)

// KillMods Kill mods provided by quake game - examples: MOD_SHOTGUN, MOD_RAILGUN, MOD_GAUNTLET
type KillMods uint

const (
	MOD_SHOTGUN KillMods = iota
	MOD_RAILGUN
	MOD_GAUNTLET
)

func (v KillMods) GetModByString(str string) KillMods {
	switch str {
	case "MOD_SHOTGUN":
		return MOD_SHOTGUN
	case "MOD_RAILGUN":
		return MOD_RAILGUN
	case "MOD_GAUNTLET":
		return MOD_GAUNTLET

	default:
		return 999
	}
}

func (v KillMods) GetModByInt(mod KillMods) KillMods {
	switch mod {
	case MOD_SHOTGUN:
		return MOD_SHOTGUN
	case MOD_RAILGUN:
		return MOD_RAILGUN
	case MOD_GAUNTLET:
		return MOD_GAUNTLET

	default:
		return 999
	}
}

func (v KillMods) GetStrModByType() string {
	switch v {
	case MOD_SHOTGUN:
		return "MOD_SHOTGUN"
	case MOD_RAILGUN:
		return "MOD_RAILGUN"
	case MOD_GAUNTLET:
		return "MOD_GAUNTLET"

	default:
		return ""
	}
}

func (v KillMods) MarshalJSON() ([]byte, error) {
	modStr := v.GetStrModByType()

	if modStr == "" {
		return []byte{}, fmt.Errorf("killmod type is invalid")
	}

	return []byte(modStr), nil
}

func (v KillMods) UnmarshalJSON(data []byte) error {
	var typ struct {
		Type interface{} `json:"type"`
	}
	if err := json.Unmarshal(data, &typ); err != nil {
		return err
	}

	str, ok := typ.Type.(string)
	if ok {
		mod := v.GetModByString(str)
		v = mod
		return nil
	}

	return json.Unmarshal(data, &v)
}

type KillModMeans map[KillMods]int

func (k KillModMeans) MarshalJSON() ([]byte, error) {
	m := make(map[string]int, 0)
	for mod, deaths := range k {
		modStr := mod.GetStrModByType()
		if modStr == "" {
			return []byte{}, fmt.Errorf("invalid_killmod")
		}

		m[modStr] = deaths
	}

	return json.Marshal(&m)
}
