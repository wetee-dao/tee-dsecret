package dkg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"
	proxy_reenc "wetee.app/dsecret/dkg/proxy-reenc"
	"wetee.app/dsecret/store"
	"wetee.app/dsecret/types"
)

// SendEncryptedSecretRequest sends a request to reencrypt a secret
// and waits for responses from all nodes.
func (d *DKG) SendEncryptedSecretRequest(ctx context.Context, req *types.ReencryptSecretRequest) (kyber.Point, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal reencrypt secret request: %w", err)
	}

	d.preRecerve[req.SecretId] = make(chan *share.PubShare)
	for _, n := range d.Nodes {
		if n.PeerID() == d.Peer.ID() {
			go func() {
				s, err := d.HandleProcessReencrypt(payload, string(d.Peer.ID()))
				if err != nil {
					fmt.Println("do process reencrypt: ", err)
					return
				}

				b, err := json.Marshal(s)
				if err != nil {
					fmt.Println("marshal reencrypted secret share: ", err)
					return
				}

				r, err := d.HandleReencryptedShare(b, string(d.Peer.ID()))
				if err != nil {
					fmt.Println("do process reencrypt:", err)
					return
				}
				d.preRecerve[req.SecretId] <- r
			}()

			continue
		}

		err = d.Peer.Send(context.Background(), n, "dkg", &types.Message{
			Type:    "reencrypt_secret_request",
			Payload: payload,
		})
		if err != nil {
			fmt.Println("send reencrypt secret request: ", err)
		}
	}

	psk := make([]*share.PubShare, 0, len(d.Nodes))
	for i := 0; i < d.Threshold; i++ {
		select {
		case data := <-d.preRecerve[req.SecretId]:
			fmt.Println("Received:", data)
		case <-time.After(15 * time.Second):
			fmt.Println("Timeout receiving from channel")
		}
	}

	xncCmt, err := proxy_reenc.Recover(d.Suite, psk, d.Threshold, len(d.Nodes))
	if err != nil {
		return nil, fmt.Errorf("recover reencrypt reply: %s", err)
	}

	return xncCmt, nil
}

// HandleProcessReencrypt processes a reencrypt request.
func (d *DKG) HandleProcessReencrypt(reqBt []byte, msgID string) (*types.ReencryptedSecretShare, error) {
	req := &types.ReencryptSecretRequest{}
	err := json.Unmarshal(reqBt, req)
	if err != nil {
		return nil, fmt.Errorf("unmarshal reencrypt secret request: %w", err)
	}

	rdrPk := req.RdrPk
	scrt, err := d.GetSecret(context.TODO(), req.SecretId)
	if err != nil {
		return nil, fmt.Errorf("get secret: %w", err)
	}

	// if r.DKG.State() != d.CERTIFIED.String() {
	// 	return nil, fmt.Errorf("dkg not certified yet: %s", r.DKG.State())
	// }

	share := d.Share()
	reply, err := proxy_reenc.Reencrypt(share, scrt, *rdrPk)
	if err != nil {
		return nil, fmt.Errorf("reencrypt: %w", err)
	}

	xncski, err := reply.Share.V.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal xncski: %w", err)
	}

	chlgi, err := reply.Challenge.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal chlgi: %w", err)
	}

	proofi, err := reply.Proof.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal proofi: %w", err)
	}

	resp := &types.ReencryptedSecretShare{
		OrgId:    string(d.Peer.ID()),
		SecretId: req.SecretId,
		Index:    int32(reply.Share.I),
		XncSki:   xncski,
		Chlgi:    chlgi,
		Proofi:   proofi,
	}

	if peer.ID(req.OrgId) != d.Peer.ID() {
		bt, err := json.Marshal(resp)
		for _, n := range d.Nodes {
			if n.PeerID() != peer.ID(req.OrgId) {
				continue
			}
			err = d.Peer.Send(context.Background(), n, "dkg", &types.Message{
				MsgID:   msgID,
				Type:    "reencrypted_secret_share",
				Payload: bt,
			})
			if err != nil {
				fmt.Println("send reencrypted secretshare: ", err)
			}
		}
	}

	return resp, nil
}

func (d *DKG) HandleReencryptedShare(reqBt []byte, msgID string) (*share.PubShare, error) {
	var resp types.ReencryptedSecretShare

	err := json.Unmarshal(reqBt, &resp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal reencrypt request: %s", err)
	}
	fmt.Printf("handling PRE response: secretid=%s from=%s \n", resp.SecretId, resp.OrgId)

	rdrPk := resp.RdrPk
	ste, err := types.SuiteForType(rdrPk.Type())
	if err != nil {
		return nil, fmt.Errorf("suite for type: %s", err)
	}

	reply := proxy_reenc.ReencryptReply{
		Share: share.PubShare{
			I: int(resp.Index),
			V: ste.Point().Base(),
		},
		Challenge: ste.Scalar(),
		Proof:     ste.Scalar(),
	}

	reply.Share.I = int(resp.Index)

	err = reply.Share.V.UnmarshalBinary(resp.XncSki)
	if err != nil {
		return nil, fmt.Errorf("unmarshal xncski: %s", err)
	}

	err = reply.Challenge.UnmarshalBinary(resp.Chlgi)
	if err != nil {
		return nil, fmt.Errorf("unmarshal chlgi: %s", err)
	}

	err = reply.Proof.UnmarshalBinary(resp.Proofi)
	if err != nil {
		return nil, fmt.Errorf("unmarshal proofi: %s", err)
	}

	distKeyShare := d.Share()
	poly := share.NewPubPoly(ste, nil, distKeyShare.Commits)

	scrt, err := d.GetSecret(context.TODO(), string(resp.SecretId))
	if err != nil {
		return nil, fmt.Errorf("getting secret: %w", err)
	}
	rawEncCmt := scrt.EncCmt

	encCmt := ste.Point().Base()
	err = encCmt.UnmarshalBinary(rawEncCmt)
	if err != nil {
		return nil, fmt.Errorf("unmarshal encrypted commitment: %s", err)
	}

	fmt.Printf("handling PRE response: verifying reencrypt reply share")
	err = proxy_reenc.Verify(*rdrPk, poly, encCmt, reply)
	if err != nil {
		return nil, fmt.Errorf("verify reencrypt reply: %s", err)
	}

	return &reply.Share, nil

	// reencryptMsgID := msgID
	// r.xncSki[reencryptMsgID] = append(r.xncSki[reencryptMsgID], &reply.Share)
	// xncSki := r.xncSki[reencryptMsgID]
	// if len(xncSki) < d.Threshold {
	// 	log.Printf("not enough shares to recover %d/%d", len(xncSki), d.Threshold)
	// 	return nil
	// }

	// log.Info("handling PRE response: recovering reencrypted commitment")
	// xncCmt, err := proxy_reenc.Recover(ste, xncSki, d.Threshold, len(d.Nodes))
	// if err != nil {
	// 	return fmt.Errorf("recover reencrypt reply: %s", err)
	// }

	// ch, ok := r.xncCmts[reencryptMsgID]
	// if !ok {
	// 	return fmt.Errorf("xncCmt channel for %s not found", reencryptMsgID)
	// }
	// log.Info("handling PRE response: returning reencrypted commitment")
	// ch <- xncCmt
	// log.Info("handling PRE response: done!")
	// return nil
}

var s []byte = nil

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

func (r *DKG) GetSecret(ctx context.Context, storeMsgID string) (*types.Secret, error) {
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
