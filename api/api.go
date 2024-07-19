package api

import (
	"context"
	"fmt"

	"wetee.app/dsecret/dkg"
	"wetee.app/dsecret/types"
)

type API struct {
	DKG *dkg.DKG
}

func (i *API) SetSecret(ctx context.Context, sid string, scrt []byte) error {
	i.DKG.SetSecret(ctx, scrt)

	return nil
}

func (i *API) ReencryptSecret(ctx context.Context, rdrPk types.PubKey, sid string) (xncCmt []byte, encScrt [][]byte, err error) {
	req := &types.ReencryptSecretRequest{
		SecretId: string(sid),
		RdrPk:    &rdrPk,
	}

	rawXncCmt, err := i.DKG.SendEncryptedSecretRequest(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("send encrypted secret request: %w", err)
	}

	xncCmt, err = rawXncCmt.MarshalBinary()
	if err != nil {
		return nil, nil, fmt.Errorf("marshal xncCmt: %w", err)
	}

	scrt, err := i.DKG.GetSecret(ctx, string(sid))
	if err != nil {
		return nil, nil, fmt.Errorf("encrypted secret for %s not found", string(sid))
	}

	return xncCmt, scrt.EncScrt, nil
}
