package qr

import (
	"context"
	"net/http"
	"net/url"
	"os"

	"github.com/corioders/terrain-story.api/foundation"
	"github.com/corioders/terrain-story.api/model/gamesCodeModel"
	"github.com/dimfeld/httptreemux"
)

type codesMapT map[string]string

type Controller struct {
	qrCodes codesMapT
}

func NewController(qrConfig foundation.QrConfig) (*Controller, error) {
	gamesCodeBytes, err := os.ReadFile(qrConfig.GamesCodeJsonPath)
	if err != nil {
		return nil, err
	}

	terrainGames, err := gamesCodeModel.Unmarshal(gamesCodeBytes)
	if err != nil {
		return nil, err
	}

	codesMap := codesMapT{}
	for _, terrainGame := range terrainGames {
		noAddons := len(terrainGame.Addons) == 0
		if noAddons {
			for _, code := range terrainGame.Codes {
				codesMap[code.Uuid] = code.To
			}
		} else {
			for _, addon := range terrainGame.Addons {
				for _, code := range terrainGame.Codes {
					// We need to unescape addon.Add as it is automatically escaped by httptreemux.ContextParams(ctx).
					unescapedAdd, err := url.PathUnescape(addon.Add.Uuid)
					if err != nil {
						return nil, err
					}

					codesMap[code.Uuid+unescapedAdd] = code.To + addon.Add.To
				}
			}
		}
	}

	return &Controller{qrCodes: codesMap}, nil
}

// RedirectHandler expects to be mounted with path that has one url param "uuid"
func (c *Controller) RedirectHandler(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	uuid := httptreemux.ContextParams(ctx)["uuid"]
	redirectionUrl, ok := c.qrCodes[uuid]
	if !ok {
		rw.WriteHeader(http.StatusNotFound)
		return nil
	}

	http.Redirect(rw, r, redirectionUrl, http.StatusFound)
	return nil
}
