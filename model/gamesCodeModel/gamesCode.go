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

type Background string

const DefaultBackground = NoneBackground
const (
	NoneBackground                Background = "none"
	NormalBackground              Background = "normal.png"
	CasualBackground              Background = "casual.png"
	QuizLekturowyCasualBackground Background = "quiz lekturowy/casual.png"
)

var AllBackgrounds = []Background{NormalBackground, CasualBackground, QuizLekturowyCasualBackground}

type AddonJson struct {
	Name       string     `json:"name"`
	Background Background `json:"background"`
	Add        struct {
		Uuid string `json:"uuid"`
		To   string `json:"to"`
	} `json:"add"`
}

func (a AddonJson) GetBackground() Background {
	if a.Background == "" {
		return DefaultBackground
	}
	return a.Background
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
