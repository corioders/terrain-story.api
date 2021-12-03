package main

import (
	"os"
	"path"

	"github.com/corioders/terrain-story.api/model/gamesCodeModel"
	"github.com/skip2/go-qrcode"
)

const rootCodesPath = "./codes"

func main() {
	gamesCodeBytes, err := os.ReadFile("../../data/gamesCode.jsonc")
	if err != nil {
		panic(err)
	}

	terrainGames, err := gamesCodeModel.Unmarshal(gamesCodeBytes)
	if err != nil {
		panic(err)
	}

	err = os.RemoveAll(rootCodesPath)
	if err != nil {
		panic(err)
	}

	err = mkdirAll(rootCodesPath)
	if err != nil {
		panic(err)
	}

	for terrainGameName, terrainGame := range terrainGames {
		rootFolder := path.Join(rootCodesPath, terrainGameName)
		err = mkdirAll(rootFolder)
		if err != nil {
			panic(err)
		}

		err = generateTerrainGame(rootFolder, terrainGame)
		if err != nil {
			panic(err)
		}
	}
}

func generateTerrainGame(rootFolder string, terrainGame gamesCodeModel.TerrainGameJson) error {
	noAddons := len(terrainGame.Addons) == 0

	if noAddons {
		for _, code := range terrainGame.Codes {
			err := generateCode(rootFolder, code)
			if err != nil {
				return err
			}
		}

		return nil
	}

	for _, addon := range terrainGame.Addons {
		addonFolder := path.Join(rootFolder, addon.Name)
		err := mkdirAll(addonFolder)
		if err != nil {
			return err
		}

		for _, code := range terrainGame.Codes {
			code.Uuid += addon.Add.Uuid
			err := generateCode(addonFolder, code)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func generateCode(rootFolder string, code gamesCodeModel.CodeJson) error {
	image, err := encodeQrCode("https://api.terrainstory.com/qr/", code.Uuid)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(rootFolder, code.Name)+".png", image, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func encodeQrCode(base, uuid string) ([]byte, error) {
	return qrcode.Encode(base+uuid, qrcode.Highest, 3000)
}
