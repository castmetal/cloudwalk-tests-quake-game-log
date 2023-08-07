package reports

import (
	"encoding/json"
	"fmt"
)

// KillMods Kill mods provided by quake game - examples: MOD_SHOTGUN, MOD_RAILGUN, MOD_GAUNTLET
type KillMods uint

const (
	MOD_UNKNOWN KillMods = iota
	MOD_SHOTGUN
	MOD_GAUNTLET
	MOD_MACHINEGUN
	MOD_GRENADE
	MOD_GRENADE_SPLASH
	MOD_ROCKET
	MOD_ROCKET_SPLASH
	MOD_PLASMA
	MOD_PLASMA_SPLASH
	MOD_RAILGUN
	MOD_LIGHTNING
	MOD_BFG
	MOD_BFG_SPLASH
	MOD_WATER
	MOD_SLIME
	MOD_LAVA
	MOD_CRUSH
	MOD_TELEFRAG
	MOD_FALLING
	MOD_SUICIDE
	MOD_TARGET_LASER
	MOD_TRIGGER_HURT
	MOD_NAIL
	MOD_CHAINGUN
	MOD_PROXIMITY_MINE
	MOD_KAMIKAZE
	MOD_JUICED
	MOD_GRAPPLE
)

const DEFAULT_KILLMOD_ERROR = "invalid_killmod"

func (v KillMods) GetModByString(str string) KillMods {
	switch str {
	case "MOD_UNKNOWN":
		return MOD_UNKNOWN
	case "MOD_SHOTGUN":
		return MOD_SHOTGUN
	case "MOD_GAUNTLET":
		return MOD_GAUNTLET
	case "MOD_MACHINEGUN":
		return MOD_MACHINEGUN
	case "MOD_GRENADE":
		return MOD_GRENADE
	case "MOD_GRENADE_SPLASH":
		return MOD_GRENADE_SPLASH
	case "MOD_ROCKET":
		return MOD_ROCKET
	case "MOD_ROCKET_SPLASH":
		return MOD_ROCKET_SPLASH
	case "MOD_PLASMA":
		return MOD_PLASMA
	case "MOD_PLASMA_SPLASH":
		return MOD_PLASMA_SPLASH
	case "MOD_RAILGUN":
		return MOD_RAILGUN
	case "MOD_LIGHTNING":
		return MOD_LIGHTNING
	case "MOD_BFG":
		return MOD_BFG
	case "MOD_BFG_SPLASH":
		return MOD_BFG_SPLASH
	case "MOD_WATER":
		return MOD_WATER
	case "MOD_SLIME":
		return MOD_SLIME
	case "MOD_LAVA":
		return MOD_LAVA
	case "MOD_CRUSH":
		return MOD_CRUSH
	case "MOD_TELEFRAG":
		return MOD_TELEFRAG
	case "MOD_FALLING":
		return MOD_FALLING
	case "MOD_SUICIDE":
		return MOD_SUICIDE
	case "MOD_TARGET_LASER":
		return MOD_TARGET_LASER
	case "MOD_TRIGGER_HURT":
		return MOD_TRIGGER_HURT
	case "MOD_NAIL":
		return MOD_NAIL
	case "MOD_CHAINGUN":
		return MOD_CHAINGUN
	case "MOD_PROXIMITY_MINE":
		return MOD_PROXIMITY_MINE
	case "MOD_KAMIKAZE":
		return MOD_KAMIKAZE
	case "MOD_JUICED":
		return MOD_JUICED
	case "MOD_GRAPPLE":
		return MOD_GRAPPLE
	default:
		return 999
	}
}

func (v KillMods) GetModByInt(mod KillMods) KillMods {
	switch mod {
	case MOD_UNKNOWN:
		return MOD_UNKNOWN
	case MOD_SHOTGUN:
		return MOD_SHOTGUN
	case MOD_GAUNTLET:
		return MOD_GAUNTLET
	case MOD_MACHINEGUN:
		return MOD_MACHINEGUN
	case MOD_GRENADE:
		return MOD_GRENADE
	case MOD_GRENADE_SPLASH:
		return MOD_GRENADE_SPLASH
	case MOD_ROCKET:
		return MOD_ROCKET
	case MOD_ROCKET_SPLASH:
		return MOD_ROCKET_SPLASH
	case MOD_PLASMA:
		return MOD_PLASMA
	case MOD_PLASMA_SPLASH:
		return MOD_PLASMA_SPLASH
	case MOD_RAILGUN:
		return MOD_RAILGUN
	case MOD_LIGHTNING:
		return MOD_LIGHTNING
	case MOD_BFG:
		return MOD_BFG
	case MOD_BFG_SPLASH:
		return MOD_BFG_SPLASH
	case MOD_WATER:
		return MOD_WATER
	case MOD_SLIME:
		return MOD_SLIME
	case MOD_LAVA:
		return MOD_LAVA
	case MOD_CRUSH:
		return MOD_CRUSH
	case MOD_TELEFRAG:
		return MOD_TELEFRAG
	case MOD_FALLING:
		return MOD_FALLING
	case MOD_SUICIDE:
		return MOD_SUICIDE
	case MOD_TARGET_LASER:
		return MOD_TARGET_LASER
	case MOD_TRIGGER_HURT:
		return MOD_TRIGGER_HURT
	case MOD_NAIL:
		return MOD_NAIL
	case MOD_CHAINGUN:
		return MOD_CHAINGUN
	case MOD_PROXIMITY_MINE:
		return MOD_PROXIMITY_MINE
	case MOD_KAMIKAZE:
		return MOD_KAMIKAZE
	case MOD_JUICED:
		return MOD_JUICED
	case MOD_GRAPPLE:
		return MOD_GRAPPLE

	default:
		return 999
	}
}

func (v KillMods) GetStrModByType() string {
	switch v {
	case MOD_UNKNOWN:
		return "MOD_UNKNOWN"
	case MOD_SHOTGUN:
		return "MOD_SHOTGUN"
	case MOD_GAUNTLET:
		return "MOD_GAUNTLET"
	case MOD_MACHINEGUN:
		return "MOD_MACHINEGUN"
	case MOD_GRENADE:
		return "MOD_GRENADE"
	case MOD_GRENADE_SPLASH:
		return "MOD_GRENADE_SPLASH"
	case MOD_ROCKET:
		return "MOD_ROCKET"
	case MOD_ROCKET_SPLASH:
		return "MOD_ROCKET_SPLASH"
	case MOD_PLASMA:
		return "MOD_PLASMA"
	case MOD_PLASMA_SPLASH:
		return "MOD_PLASMA_SPLASH"
	case MOD_RAILGUN:
		return "MOD_RAILGUN"
	case MOD_LIGHTNING:
		return "MOD_LIGHTNING"
	case MOD_BFG:
		return "MOD_BFG"
	case MOD_BFG_SPLASH:
		return "MOD_BFG_SPLASH"
	case MOD_WATER:
		return "MOD_WATER"
	case MOD_SLIME:
		return "MOD_SLIME"
	case MOD_LAVA:
		return "MOD_LAVA"
	case MOD_CRUSH:
		return "MOD_CRUSH"
	case MOD_TELEFRAG:
		return "MOD_TELEFRAG"
	case MOD_FALLING:
		return "MOD_FALLING"
	case MOD_SUICIDE:
		return "MOD_SUICIDE"
	case MOD_TARGET_LASER:
		return "MOD_TARGET_LASER"
	case MOD_TRIGGER_HURT:
		return "MOD_TRIGGER_HURT"
	case MOD_NAIL:
		return "MOD_NAIL"
	case MOD_CHAINGUN:
		return "MOD_CHAINGUN"
	case MOD_PROXIMITY_MINE:
		return "MOD_PROXIMITY_MINE"
	case MOD_KAMIKAZE:
		return "MOD_KAMIKAZE"
	case MOD_JUICED:
		return "MOD_JUICED"
	case MOD_GRAPPLE:
		return "MOD_GRAPPLE"

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
