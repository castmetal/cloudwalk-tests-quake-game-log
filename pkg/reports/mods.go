package reports

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
)

// KillMods Kill mods provided by quake game - examples: MOD_SHOTGUN, MOD_RAILGUN, MOD_GAUNTLET
type KillMods uint

const (
	MOD_SHOTGUN KillMods = iota
	MOD_RAILGUN
	MOD_GAUNTLET
)

const DEFAULT_KILLMOD_ERROR = "invalid_killmod"

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
		return []byte{}, fmt.Errorf(DEFAULT_KILLMOD_ERROR)
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

type KillModMeans map[KillMods]int32

func (k KillModMeans) MarshalJSON() ([]byte, error) {
	m := make(map[string]int32, 0)
	for mod, deaths := range k {
		modStr := mod.GetStrModByType()
		if modStr == "" {
			return []byte{}, fmt.Errorf(DEFAULT_KILLMOD_ERROR)
		}

		m[modStr] = deaths
	}

	return json.Marshal(&m)
}

func (k KillModMeans) AddDeath(mod KillMods) error {
	modType := mod.GetStrModByType()
	if modType == "" {
		return fmt.Errorf(DEFAULT_KILLMOD_ERROR)
	}

	addData := k[mod]

	atomic.AddInt32(&addData, 1)

	k[mod] = addData

	return nil
}

func (k KillModMeans) RemoveDeath(mod KillMods) error {
	modType := mod.GetStrModByType()
	if modType == "" {
		return fmt.Errorf(DEFAULT_KILLMOD_ERROR)
	}

	removeData := k[mod]

	atomic.AddInt32(&removeData, -1)

	k[mod] = removeData

	return nil
}
