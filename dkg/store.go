package dkg

import (
	"context"
	"encoding/json"

	"wetee.app/dsecret/store"
	types "wetee.app/dsecret/type"
)

func (r *DKG) SetData(datas []types.Kvs) error {
	// for _, data := range datas {
	// 	err := store.SetKey("secret", data.K, data.V)
	// 	if err != nil {
	// 		return fmt.Errorf("set key: %w", err)
	// 	}
	// }
	bt, _ := json.Marshal(datas)
	return r.Peer.Pub(context.Background(), "secret", bt)
}

func (r *DKG) GetData(k string) ([]byte, error) {
	return store.GetKey("secret", k)
}
