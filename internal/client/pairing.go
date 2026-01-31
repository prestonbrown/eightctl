package client

import (
	"context"
	"fmt"
	"net/http"
)

type PairingActions struct{ c *Client }

func (c *Client) Pairing() *PairingActions { return &PairingActions{c: c} }

func (p *PairingActions) StartAutoPairing(ctx context.Context) (any, error) {
	id, err := p.c.EnsureDeviceID(ctx)
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("/devices/%s/auto-pairing/start", id)
	var res any
	err = p.c.do(ctx, http.MethodPost, path, nil, nil, &res)
	return res, err
}

func (p *PairingActions) GetPairingStatus(ctx context.Context, pairingID string, out any) error {
	id, err := p.c.EnsureDeviceID(ctx)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/devices/%s/auto-pairing/status/%s", id, pairingID)
	return p.c.do(ctx, http.MethodGet, path, nil, nil, out)
}
