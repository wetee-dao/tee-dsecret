// / Copyright (c) 2022 Sourcenetwork Developers. All rights reserved.
// / copy from https://github.com/sourcenetwork/orbis-go

package dkg

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/share"

	proxy_reenc "wetee.app/dsecret/dkg/proxy-reenc"
	types "wetee.app/dsecret/type"
)

// SendEncryptedSecretRequest sends a request to reencrypt a secret
// and waits for responses from all nodes.
func (d *DKG) SendEncryptedSecretRequest(ctx context.Context, req *types.ReencryptSecretRequest) (kyber.Point, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal reencrypt secret request: %w", err)
	}

	msgId := uuid.NewV4().String()
	d.mu.Lock()
	d.preRecerve[msgId] = make(chan interface{})
	d.mu.Unlock()

	for _, n := range d.DkgNodes {
		err = d.SendToNode(context.Background(), n, "worker", &types.Message{
			Type:    "reencrypt_secret_request",
			Payload: payload,
		})
		if err != nil {
			fmt.Println("send reencrypt secret request: ", err)
		}
	}

	psk := make([]*share.PubShare, 0, len(d.DkgNodes))
	for i := 0; i < d.Threshold; i++ {
		select {
		case d := <-d.preRecerve[msgId]:
			data := d.(*share.PubShare)
			fmt.Println("Received:", data)
			psk = append(psk, data)
		case <-time.After(30 * time.Second):
			fmt.Println("Timeout receiving from channel")
			return nil, fmt.Errorf("timeout receiving from channel")
		}
	}

	d.mu.Lock()
	delete(d.preRecerve, msgId)
	d.mu.Unlock()

	xncCmt, err := proxy_reenc.Recover(d.Suite, psk, d.Threshold, len(d.DkgNodes))
	if err != nil {
		return nil, fmt.Errorf("recover reencrypt reply: %s", err)
	}

	return xncCmt, nil
}

// HandleProcessReencrypt processes a reencrypt request.
func (d *DKG) HandleProcessReencrypt(reqBt []byte, msgID string, OrgId string) error {
	req := &types.ReencryptSecretRequest{}
	err := json.Unmarshal(reqBt, req)
	if err != nil {
		return fmt.Errorf("unmarshal reencrypt secret request: %w", err)
	}

	rdrPk := req.RdrPk
	scrt, err := d.GetSecretData(context.TODO(), req.SecretId)
	if err != nil {
		return fmt.Errorf("get secret: %w", err)
	}

	// if r.DKG.State() != d.CERTIFIED.String() {
	// 	return nil, fmt.Errorf("dkg not certified yet: %s", r.DKG.State())
	// }

	share := d.Share()
	reply, err := proxy_reenc.Reencrypt(share, scrt, *rdrPk)
	if err != nil {
		return fmt.Errorf("reencrypt: %w", err)
	}

	xncski, err := reply.Share.V.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal xncski: %w", err)
	}

	chlgi, err := reply.Challenge.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal chlgi: %w", err)
	}

	proofi, err := reply.Proof.MarshalBinary()
	if err != nil {
		return fmt.Errorf("marshal proofi: %w", err)
	}

	resp := &types.ReencryptedSecretShare{
		SecretId: req.SecretId,
		Index:    int32(reply.Share.I),
		XncSki:   xncski,
		Chlgi:    chlgi,
		Proofi:   proofi,
	}
	bt, _ := json.Marshal(resp)

	n := d.GetNode(OrgId)
	if n == nil {
		return fmt.Errorf("node not found: %s", OrgId)
	}
	err = d.SendToNode(context.Background(), n, "worker", &types.Message{
		MsgID:   msgID,
		Type:    "reencrypted_secret_reply",
		Payload: bt,
	})
	if err != nil {
		fmt.Println("send reencrypted secretshare: ", err)
	}

	return nil
}

func (d *DKG) HandleReencryptedShare(reqBt []byte, msgID string, OrgId string) error {
	var req types.ReencryptedSecretShare

	err := json.Unmarshal(reqBt, &req)
	if err != nil {
		return fmt.Errorf("unmarshal reencrypt request: %s", err)
	}
	fmt.Printf("handling PRE response: secretid=%s from=%s \n", req.SecretId, OrgId)

	rdrPk := req.RdrPk
	ste, err := types.SuiteForType(rdrPk.Type())
	if err != nil {
		return fmt.Errorf("suite for type: %s", err)
	}

	reply := proxy_reenc.ReencryptReply{
		Share: share.PubShare{
			I: int(req.Index),
			V: ste.Point().Base(),
		},
		Challenge: ste.Scalar(),
		Proof:     ste.Scalar(),
	}

	reply.Share.I = int(req.Index)

	err = reply.Share.V.UnmarshalBinary(req.XncSki)
	if err != nil {
		return fmt.Errorf("unmarshal xncski: %s", err)
	}

	err = reply.Challenge.UnmarshalBinary(req.Chlgi)
	if err != nil {
		return fmt.Errorf("unmarshal chlgi: %s", err)
	}

	err = reply.Proof.UnmarshalBinary(req.Proofi)
	if err != nil {
		return fmt.Errorf("unmarshal proofi: %s", err)
	}

	distKeyShare := d.Share()
	poly := share.NewPubPoly(ste, nil, distKeyShare.Commits)

	scrt, err := d.GetSecretData(context.TODO(), string(req.SecretId))
	if err != nil {
		return fmt.Errorf("getting secret: %w", err)
	}
	rawEncCmt := scrt.EncCmt

	encCmt := ste.Point().Base()
	err = encCmt.UnmarshalBinary(rawEncCmt)
	if err != nil {
		return fmt.Errorf("unmarshal encrypted commitment: %s", err)
	}

	fmt.Printf("handling PRE response: verifying reencrypt reply share")
	err = proxy_reenc.Verify(*rdrPk, poly, encCmt, reply)
	if err != nil {
		return fmt.Errorf("verify reencrypt reply: %s", err)
	}

	if _, ok := d.preRecerve[msgID]; !ok {
		return nil
	}

	d.preRecerve[msgID] <- &reply.Share

	return nil

}
