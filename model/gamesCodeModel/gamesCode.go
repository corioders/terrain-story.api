package gamesCodeModel

import (
	"net/url"

	"muzzammil.xyz/jsonc"
)

type CodeJson struct {
	Name string `json:"name"`
	Uuid string `json:"uuid"`
	To   string `json:"to"`
}

type AddonJson struct {
	Name string `json:"name"`
	Add  string `json:"add"`
}

type TerrainGameJson struct {
	Codes  []CodeJson  `json:"codes"`
	Addons []AddonJson `json:"addons"`
}

type TerrainGamesJson map[string]TerrainGameJson

func Unmarshal(data []byte) (TerrainGamesJson, error) {
	terrainGamesJson := TerrainGamesJson{}
	err := jsonc.Unmarshal(data, &terrainGamesJson)
	if err != nil {
		return nil, err
	}

	return terrainGamesJson, nil
}

func NormalizeAddonAdd(add string) string {
	return url.PathEscape(add)
}
