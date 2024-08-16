package dkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	stypes "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	uuid "github.com/satori/go.uuid"
	"github.com/vedhavyas/go-subkey/v2"
	gtypes "github.com/wetee-dao/go-sdk/pallet/types"
	"github.com/wetee-dao/go-sdk/pallet/weteedsecret"
	"github.com/wetee-dao/go-sdk/pallet/weteeworker"
	"golang.org/x/crypto/blake2b"

	"wetee.app/dsecret/chain"
	"wetee.app/dsecret/tee"
	types "wetee.app/dsecret/type"
)

// HandleDeal 处理密钥份额消息
func (dkg *DKG) HandleUploadClusterProof(data []byte, msgID string, OrgId string) error {
	workerReport := &types.TeeParam{}
	err := json.Unmarshal(data, workerReport)
	if err != nil {
		return err
	}

	// 通过地址获取集群信息
	_, account, _ := subkey.SS58Decode("5C5NWbLEkbb6gb7prqZEDcJnNe2y4SmBwZHbaoxyLdqHm3v2") //(workerReport.Address)
	var account32 [32]byte
	copy(account32[:], account)
	cid, ok, err := weteeworker.GetK8sClusterAccountsLatest(chain.ChainIns.GetClient().Api.RPC.State, account32)
	if err != nil || !ok {
		return errors.New("get k8s cluster error")
	}

	msgId := uuid.NewV4().String()

	dkg.mu.Lock()
	dkg.preRecerve[msgId] = make(chan interface{})
	dkg.mu.Unlock()

	// 请求节点验证签名
	errNum := 0
	for _, node := range dkg.DkgNodes {
		err := dkg.SendToNode(context.Background(), node, "worker", &types.Message{
			MsgID:   msgId,
			Type:    "sign_cluster_proof",
			Payload: data,
		})
		if err != nil {
			errNum++
		}
	}

	if len(dkg.DkgNodes)-errNum < dkg.Threshold {
		return errors.New("not enough nodes")
	}

	pubs := make([][32]byte, 0, len(dkg.DkgNodes))
	sigs := make([]gtypes.MultiSignature, 0, len(dkg.DkgNodes))
	for i := 0; i <= dkg.Threshold; i++ {
		select {
		case d := <-dkg.preRecerve[msgId]:
			data := d.(*ReportSign)
			pubs = append(pubs, data.account)
			sigs = append(sigs, data.sig)
		case <-time.After(30 * time.Second):
			fmt.Println("Timeout receiving from channel")
			return fmt.Errorf("timeout receiving from channel")
		}
	}

	// 获取交易帐户
	s, err := dkg.Signer.ToSigner()
	if err != nil {
		return errors.New("signer to signer: " + err.Error())
	}

	dkg.mu.Lock()
	delete(dkg.preRecerve, msgId)
	dkg.mu.Unlock()

	ins := chain.ChainIns.GetClient()
	ins.CheckMetadata()

	// 提交证明
	hash := blake2b.Sum512(workerReport.Report)
	call := weteedsecret.MakeUploadClusterProofCall(cid, hash[:], pubs, sigs)
	return ins.SignAndSubmit(s, call, false)
}

func (dkg *DKG) HandleSignClusterProof(data []byte, msgID string, OrgId string) error {
	workerReport := &types.TeeParam{}
	err := json.Unmarshal(data, workerReport)
	if err != nil {
		return fmt.Errorf("unmarshal reencrypt secret reply: %w", err)
	}

	// decode address
	_, signer, err := subkey.SS58Decode(workerReport.Address)
	if err != nil {
		return errors.New("SS58 decode: " + err.Error())
	}

	_, err = tee.VerifyReport(workerReport.Report, workerReport.Data, signer, workerReport.Time)
	if err != nil {
		return errors.New("verify report: " + err.Error())
	}

	// TODO
	// 校验代码版本

	siger, err := dkg.Signer.ToSigner()
	if err != nil {
		return errors.New("signer to signer: " + err.Error())
	}

	hash := blake2b.Sum512(workerReport.Report)
	sig, err := siger.Sign(hash[:])
	if err != nil {
		return errors.New("sign: " + err.Error())
	}

	n := dkg.GetNode(OrgId)
	if n == nil {
		return fmt.Errorf("node not found: %s", OrgId)
	}
	if err := dkg.SendToNode(context.Background(), n, "worker", &types.Message{
		MsgID:   msgID,
		Type:    "sign_cluster_proof_reply",
		Payload: sig,
	}); err != nil {
		return errors.New("send to node: " + err.Error())
	}

	return nil
}

func (dkg *DKG) HandleSignClusterProofReply(data []byte, msgID string, OrgId string) error {
	account := dkg.GetNode(OrgId)
	if account == nil {
		return fmt.Errorf("node not found: %s", OrgId)
	}

	// 还原公钥
	pub, err := types.PublicKeyFromLibp2pHex(account.ID)
	if err != nil {
		return errors.New("public key from libp2p hex: " + err.Error())
	}

	// 计算 account32
	bt, err := pub.Raw()
	if err != nil {
		return errors.New("public key raw: " + err.Error())
	}
	var account32 [32]byte
	copy(account32[:], bt)

	// 如果已经满足签名需求，则直接返回
	if _, ok := dkg.preRecerve[msgID]; !ok {
		return nil
	}

	sig := stypes.NewSignature(data)
	dkg.preRecerve[msgID] <- &ReportSign{
		account: account32,
		sig: gtypes.MultiSignature{
			IsEd25519:       true,
			AsEd25519Field0: sig,
		},
	}

	return nil
}

type ReportSign struct {
	account [32]byte
	sig     gtypes.MultiSignature
}
