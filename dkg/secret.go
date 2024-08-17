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

func (r *DKG) SetSecret(ctx context.Context, env types.Env) (*types.SecretEnvWithHash, error) {
	scrt, err := json.Marshal(env)
	if err != nil {
		return nil, err
	}

	// 计算用于展示的谜文
	lenvs := make([]*types.LenValue, len(env.Envs))
	lfiles := make([]*types.LenValue, len(env.Files))
	for i, v := range env.Envs {
		lenvs[i] = &types.LenValue{
			K: v.K,
			V: len(v.V),
		}
	}
	for i, v := range env.Files {
		lfiles[i] = &types.LenValue{
			K: v.K,
			V: len(v.V),
		}
	}
	lenEnv := &types.SecretEnv{
		Envs:  lenvs,
		Files: lfiles,
	}
	pub, _ := json.Marshal(lenEnv)

	// 使用公钥
	dkgPub := r.DkgPubKey

	// 加密秘密
	encCmt, encScrt := proxy_reenc.EncryptSecret(r.Suite, dkgPub, scrt)
	rawEncCmt, err := encCmt.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal encCmt: %s", err)
	}

	// 转换秘文
	rawEncScrt := make([][]byte, len(encScrt))
	for i, encScrti := range encScrt {
		rawEncScrti, err := encScrti.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("marshal encScrt: %s", err)
		}
		rawEncScrt[i] = rawEncScrti
	}

	// 保存
	secret := &types.Secret{
		EncCmt:  rawEncCmt,
		EncScrt: rawEncScrt,
	}

	// 格式化数据保存
	payload, err := json.Marshal(secret)
	if err != nil {
		return nil, fmt.Errorf("marshal secret: %w", err)
	}

	// 获取唯一id
	cid, err := types.CidFromBytes(payload)
	if err != nil {
		return nil, fmt.Errorf("cid from bytes: %w", err)
	}

	// 保存数据
	storeMsgID := cid.KeyString()
	err = r.SetSecretData([]types.Kvs{
		{K: storeMsgID, V: payload},
		{K: storeMsgID + "_pub", V: pub},
	})
	if err != nil {
		return nil, fmt.Errorf("set secret: %w", err)
	}

	return &types.SecretEnvWithHash{
		Hash:   storeMsgID,
		Secret: lenEnv,
	}, nil
}

func (r *DKG) SetSecretData(datas []types.Kvs) error {
	// 保存数据
	for _, data := range datas {
		err := store.SetKey("secret", data.K, data.V)
		if err != nil {
			return fmt.Errorf("set key: %w", err)
		}
	}
	bt, _ := json.Marshal(datas)
	return r.Peer.Pub(context.Background(), "secret", bt)
}

func (r *DKG) HandleSecretSave(ctx context.Context) {
	sub, err := r.Peer.Sub(ctx, "secret")
	if err != nil {
		return
	}

	for {
		msg, err := sub.Next(ctx)
		if err != nil {
			fmt.Println("Error receiving message:", err)
			continue
		}

		// 解析消息
		var datas []types.Kvs
		err = json.Unmarshal(msg.Data, &datas)
		for _, data := range datas {
			err := store.SetKey("secret", data.K, data.V)
			if err != nil {
				fmt.Println("set key: %w", err)
				continue
			}
		}
	}
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

func (r *DKG) GetSecretPubData(ctx context.Context, storeMsgID string) (*types.SecretEnvWithHash, error) {
	buf, err := store.GetKey("secret", storeMsgID+"_pub")
	if err != nil {
		return nil, fmt.Errorf("get secret: %w", err)
	}

	s := new(types.SecretEnv)
	err = json.Unmarshal(buf, s)
	if err != nil {
		return nil, fmt.Errorf("unmarshal encrypted secret: %w", err)
	}

	return &types.SecretEnvWithHash{
		Hash:   storeMsgID,
		Secret: s,
	}, nil
}
