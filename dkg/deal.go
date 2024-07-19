package dkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	rabin "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"wetee.app/dsecret/types"
)

// SendDealMessage 发送Deal消息
func (dkg *DKG) SendDealMessage(ctx context.Context, node *types.Node, message *rabin.Deal) error {
	pmessage, err := types.DealToProtocol(message)
	if err != nil {
		return err
	}

	bt, err := json.Marshal(pmessage)
	if err != nil {
		return err
	}

	return dkg.Peer.Send(ctx, node, "dkg", &types.Message{
		Type:    "deal",
		Payload: bt,
	})
}

// HandleDeal 处理密钥份额消息。
func (dkg *DKG) HandleDeal(data []byte) error {
	dkg.mu.Lock()
	defer dkg.mu.Unlock()

	pmessage := &types.Deal{}
	err := json.Unmarshal(data, pmessage)
	if err != nil {
		fmt.Println(err)
		return err
	}
	deal, err := types.ProtocolToDeal(dkg.Suite, pmessage)
	if err != nil {
		return err
	}

	// 处理密钥份额。
	resp, err := dkg.DistKeyGenerator.ProcessDeal(deal)
	if err != nil {
		return fmt.Errorf("HandleDeal error: %w", err)
	}

	if !resp.Response.Approved {
		return errors.New("deal rejected")
	}

	bt, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	// 发送 deal resp
	for _, node := range dkg.Nodes {
		if node.PeerID() == dkg.Peer.ID() {
			continue
		}
		err = dkg.Peer.Send(context.Background(), node, "dkg", &types.Message{
			Type:    "deal_resp",
			Payload: bt,
		})
		if err != nil {
			fmt.Println("Send deal_resp error", err)
		}
	}

	return nil
}

// HandleDealMessage 处理密钥份额消息。
func (dkg *DKG) HandleDealResp(data []byte) error {
	dkg.mu.Lock()
	defer dkg.mu.Unlock()

	message := &rabin.Response{}
	err := json.Unmarshal(data, message)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 处理密钥份额。
	justification, err := dkg.DistKeyGenerator.ProcessResponse(message)
	if err != nil {
		return fmt.Errorf("HandleDealResp ProcessResponse: %w", err)
	}

	// tJustification，证明 Deal 消息的无效性
	if justification != nil {
		fmt.Println("Got justification during response process for ", message.Index, justification)
		return nil
	}

	// 已经判断为有效了
	if !dkg.DistKeyGenerator.Certified() {
		fmt.Println("DistKeyGenerator not certified")
		return nil
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	sc, err := dkg.DistKeyGenerator.SecretCommits()
	if err != nil {
		return fmt.Errorf("Generate secret commit: %w", err)
	}

	psc, err := types.SecretCommitsToProtocol(sc)
	if err != nil {
		return fmt.Errorf("SecretCommitsToProtocol : %w", err)
	}

	bt, err := json.Marshal(psc)
	if err != nil {
		return fmt.Errorf("HandleDealResp json.Marshal: %w", err)
	}

	// 发送 deal resp
	for _, node := range dkg.Nodes {
		if node.PeerID() == dkg.Peer.ID() {
			continue
		}
		err = dkg.Peer.Send(context.Background(), node, "dkg", &types.Message{
			Type:    "secret_commits",
			Payload: bt,
		})
		if err != nil {
			fmt.Println("Send secret_commits error", err)
		}
	}

	return nil
}

// HandleDealMessage 处理密钥份额消息。
func (dkg *DKG) HandleJustification(data []byte) error {
	dkg.mu.Lock()
	defer dkg.mu.Unlock()

	message := &rabin.Justification{}
	err := json.Unmarshal(data, message)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 判断是否是当前节点需要处理的消息
	if dkg.ID() == -1 || uint32(dkg.ID()) != message.Index {
		return nil
	}

	// 处理密钥份额。
	err = dkg.DistKeyGenerator.ProcessJustification(message)
	if err != nil {
		return fmt.Errorf("ProcessJustification: %w", err)
	}

	return nil
}

func (dkg *DKG) HandleSecretCommits(data []byte) error {
	dkg.mu.Lock()
	defer dkg.mu.Unlock()

	// 转换协议对象。
	psc := &types.SecretCommits{}
	err := json.Unmarshal(data, psc)
	if err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}
	sc, err := types.SecretCommitsFromProtocol(dkg.Suite, psc)
	if err != nil {
		return fmt.Errorf("ProtocolToSecretCommits: %w", err)
	}

	// 处理秘密提交
	_, err = dkg.DistKeyGenerator.ProcessSecretCommits(sc)
	if err != nil {
		return fmt.Errorf("ProcessSecretCommits: %w", err)
	}

	// interpolate shared public key
	distkey, err := dkg.DistKeyGenerator.DistKeyShare()
	if err != nil {
		return fmt.Errorf("rabin dkg dist key share: %w", err)
	}

	fmt.Println("================================================", distkey)

	dkg.DkgKeyShare = types.DistKeyShare{
		Commits:  distkey.Commitments(),
		PriShare: distkey.PriShare(),
	}
	dkg.DkgPubKey = distkey.Public()

	return nil
}
