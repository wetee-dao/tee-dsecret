package dkg

import (
	"context"
	"fmt"

	types "wetee.app/dsecret/type"
)

func (i *DKG) GetSecretApi(ctx context.Context, rdrPk types.PubKey, sid string) (xncCmt []byte, encScrt [][]byte, err error) {
	req := &types.ReencryptSecretRequest{
		SecretId: string(sid),
		RdrPk:    &rdrPk,
	}

	// send request
	rawXncCmt, err := i.SendEncryptedSecretRequest(ctx, req)
	if err != nil {
		return nil, nil, fmt.Errorf("send encrypted secret request: %w", err)
	}

	// marshal xncCmt
	xncCmt, err = rawXncCmt.MarshalBinary()
	if err != nil {
		return nil, nil, fmt.Errorf("marshal xncCmt: %w", err)
	}

	// get secret
	scrt, err := i.GetSecretData(ctx, string(sid))
	if err != nil {
		return nil, nil, fmt.Errorf("encrypted secret for %s not found", string(sid))
	}

	return xncCmt, scrt.EncScrt, nil
}
