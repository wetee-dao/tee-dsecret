package dkg

import (
	"context"
	"encoding/json"
	"fmt"

	proxy_reenc "wetee.app/dsecret/dkg/proxy-reenc"
	"wetee.app/dsecret/store"
	types "wetee.app/dsecret/type"
)

func (r *DKG) SetSecretEnv(ctx context.Context, env types.Env) (*types.SecretEnvWithHash, error) {
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
	storeMsgID := cid.String()
	go r.SetData([]types.Kvs{
		{K: storeMsgID, V: payload},
		{K: storeMsgID + "_pub", V: pub},
	})

	return &types.SecretEnvWithHash{
		Hash:   storeMsgID,
		Secret: lenEnv,
	}, nil
}

func (r *DKG) GetSecretPubEnvData(ctx context.Context, storeMsgID string) (*types.SecretEnvWithHash, error) {
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
