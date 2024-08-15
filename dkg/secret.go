package dkg

import (
	"context"
	"encoding/json"
	"fmt"

	proxy_reenc "wetee.app/dsecret/dkg/proxy-reenc"
	"wetee.app/dsecret/store"
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

func (r *DKG) SetSecret(ctx context.Context, scrt []byte) (string, error) {
	dkgPub := r.DkgPubKey

	// 加密秘密
	encCmt, encScrt := proxy_reenc.EncryptSecret(r.Suite, dkgPub, scrt)
	rawEncCmt, err := encCmt.MarshalBinary()
	if err != nil {
		return "", fmt.Errorf("marshal encCmt: %s", err)
	}

	// 转换秘文
	rawEncScrt := make([][]byte, len(encScrt))
	for i, encScrti := range encScrt {
		rawEncScrti, err := encScrti.MarshalBinary()
		if err != nil {
			return "", fmt.Errorf("marshal encScrt: %s", err)
		}
		rawEncScrt[i] = rawEncScrti
	}

	// 保存
	secret := &types.Secret{
		EncCmt:  rawEncCmt,
		EncScrt: rawEncScrt,
	}
	payload, err := json.Marshal(secret)
	if err != nil {
		return "", fmt.Errorf("marshal secret: %w", err)
	}

	cid, err := types.CidFromBytes(payload)
	if err != nil {
		return "", fmt.Errorf("cid from bytes: %w", err)
	}

	storeMsgID := string(cid.String())
	err = store.SetKey("secret", storeMsgID, payload)
	if err != nil {
		return "", fmt.Errorf("set secret: %w", err)
	}

	return storeMsgID, nil
}

func (r *DKG) GetSecretData(ctx context.Context, storeMsgID string) (*types.Secret, error) {
	buf, err := store.GetKey("secret", storeMsgID)
	if err != nil {
		return nil, fmt.Errorf("get secret: %w", err)
	}

	s := new(types.Secret)
	err = json.Unmarshal(buf, s)
	if err != nil {
		return nil, fmt.Errorf("unmarshal encrypted secret: %w", err)
	}

	return s, nil
}
