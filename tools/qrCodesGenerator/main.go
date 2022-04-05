package main

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/png"
	"os"
	"path"

	"golang.org/x/image/draw"
	"golang.org/x/sync/errgroup"

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
			err := generateCode(rootFolder, code, gamesCodeModel.DefaultBackground)
			if err != nil {
				return err
			}
		}

		return nil
	}

	eg := errgroup.Group{}

	for _, addon := range terrainGame.Addons {
		localAddon := addon
		addonFolder := path.Join(rootFolder, localAddon.Name)
		err := mkdirAll(addonFolder)
		if err != nil {
			return err
		}

		for _, code := range terrainGame.Codes {
			// We need to redeclare this variable, in other case this happens https://dev.to/kkentzo/the-golang-for-loop-gotcha-1n35.
			localCode := code
			eg.Go(func() error {
				localCode.Uuid += localAddon.Add.Uuid
				err := generateCode(addonFolder, localCode, localAddon.GetBackground())
				if err != nil {
					return err
				}
				return nil
			})

		}
	}

	return eg.Wait()
}

//go:embed background/*
var backgroundFS embed.FS

type backgroundDetails struct {
	img      image.Image
	codeSize int
	pt       image.Point
}

var backgroundMap = map[gamesCodeModel.Background]*backgroundDetails{
	gamesCodeModel.NormalBackground: {
		codeSize: 630,
		pt:       image.Pt(910, 1455),
	},
	gamesCodeModel.CasualBackground: {
		codeSize: 550,
		pt:       image.Pt(940, 700),
	},
	gamesCodeModel.QuizLekturowyCasualBackground: {
		codeSize: 550,
		pt:       image.Pt(940, 700),
	},
}

func init() {
	for _, b := range gamesCodeModel.AllBackgrounds {
		backgroundPngBytes, err := backgroundFS.ReadFile(path.Join("background", string(b)))
		if err != nil {
			panic(err)
		}

		backgroundImage, err := png.Decode(bytes.NewReader(backgroundPngBytes))
		if err != nil {
			panic(err)
		}

		if backgroundMap[b] == nil {
			panic(fmt.Sprintf("You have not provided details for background: %s", b))
		}
		backgroundMap[b].img = backgroundImage
	}
}

func generateCode(rootFolder string, code gamesCodeModel.CodeJson, background gamesCodeModel.Background) error {
	size := 3000
	if background != gamesCodeModel.NoneBackground {
		size = backgroundMap[background].codeSize
	}

	codeImage, err := encodeQrCode("https://api.terrainstory.com/qr/", code.Uuid, size)
	if err != nil {
		return err
	}

	if background != gamesCodeModel.NoneBackground {
		backgroundDetails := backgroundMap[background]
		resultImage := image.NewNRGBA(backgroundDetails.img.Bounds())

		draw.Copy(resultImage, image.Pt(0, 0), backgroundDetails.img, backgroundDetails.img.Bounds(), draw.Over, nil)
		draw.Copy(resultImage, backgroundDetails.pt, codeImage, codeImage.Bounds(), draw.Over, nil)

		codeImage = resultImage
	}

	encoder := png.Encoder{CompressionLevel: png.BestCompression}
	buffer := bytes.Buffer{}
	err = encoder.Encode(&buffer, codeImage)
	if err != nil {
		return err
	}

	err = os.WriteFile(path.Join(rootFolder, code.Name)+".png", buffer.Bytes(), os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func encodeQrCode(base, uuid string, size int) (image.Image, error) {
	q, err := qrcode.New(base+uuid, qrcode.Highest)
	if err != nil {
		return nil, err
	}

	return q.Image(size), nil
}
