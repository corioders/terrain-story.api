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
	Add  struct {
		Uuid string `json:"uuid"`
		To   string `json:"to"`
	} `json:"add"`
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

	for terrainGameKey := range terrainGamesJson {
		addons := terrainGamesJson[terrainGameKey].Addons
		if len(addons) == 0 {
			continue
		}

		for addonI := 0; addonI < len(addons); addonI++ {
			addons[addonI].Add.Uuid = url.PathEscape(addons[addonI].Add.Uuid)
		}
	}

	return terrainGamesJson, nil
}
