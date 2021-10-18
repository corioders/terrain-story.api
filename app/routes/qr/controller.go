package qr

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/corioders/terrain-story.api/foundation"
	"github.com/dimfeld/httptreemux"
)

type qrCodesMap map[string]string

type qrCodesJson []struct {
	Uuid string `json:"uuid"`
	To   string `json:"to"`
}

type Controller struct {
	qrCodes qrCodesMap
}

func NewController(qrConfig foundation.QrConfig) (*Controller, error) {
	qrCodesBytes, err := os.ReadFile(qrConfig.QrCodesJsonPath)
	if err != nil {
		return nil, err
	}

	qrCodesJ := qrCodesJson{}
	err = json.Unmarshal(qrCodesBytes, &qrCodesJ)
	if err != nil {
		return nil, err
	}

	qrCodesM := qrCodesMap{}
	for _, qrCode := range qrCodesJ {
		qrCodesM[qrCode.Uuid] = qrCode.To
	}

	return &Controller{qrCodes: qrCodesM}, nil
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
