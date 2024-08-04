package api

import (
	"context"
	"fmt"

	"wetee.app/dsecret/dkg"
	types "wetee.app/dsecret/type"
)

type API struct {
	DKG *dkg.DKG
}

func (i *API) SetSecret(ctx context.Context, scrt []byte) (string, error) {
	return i.DKG.SetSecret(ctx, scrt)
}

func (i *API) GetSecret(ctx context.Context, rdrPk types.PubKey, sid string) (xncCmt []byte, encScrt [][]byte, err error) {
	req := &types.ReencryptSecretRequest{
		SecretId: string(sid),
		RdrPk:    &rdrPk,
	}

	// send request
	rawXncCmt, err := i.DKG.SendEncryptedSecretRequest(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("send encrypted secret request: %w", err)
	}

	// marshal xncCmt
	xncCmt, err = rawXncCmt.MarshalBinary()
	if err != nil {
		return nil, nil, fmt.Errorf("marshal xncCmt: %w", err)
	}

	// get secret
	scrt, err := i.DKG.GetSecret(ctx, string(sid))
	if err != nil {
		return nil, nil, fmt.Errorf("encrypted secret for %s not found", string(sid))
	}

	return xncCmt, scrt.EncScrt, nil
}
