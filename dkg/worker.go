package dkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"github.com/vedhavyas/go-subkey/v2"
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

	psk := make([][]byte, 0, len(dkg.DkgNodes))
	for i := 0; i < dkg.Threshold; i++ {
		select {
		case d := <-dkg.preRecerve[msgId]:
			data := d.([]byte)
			psk = append(psk, data)
		case <-time.After(30 * time.Second):
			fmt.Println("Timeout receiving from channel")
			return fmt.Errorf("timeout receiving from channel")
		}
	}

	fmt.Println(psk)

	dkg.mu.Lock()
	delete(dkg.preRecerve, msgId)
	dkg.mu.Unlock()

	return nil
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

	sig, err := siger.Sign(data)
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
	if _, ok := dkg.preRecerve[msgID]; !ok {
		return nil
	}

	dkg.mu.Lock()
	dkg.preRecerve[msgID] <- data
	dkg.mu.Unlock()
	return nil
}
