package dkg

import (
	"encoding/json"
	"fmt"

	rabin "go.dedis.ch/kyber/v3/share/dkg/rabin"
	"wetee.app/dsecret/types"
)

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
	message, err := types.ProtocolToDeal(dkg.Suite, pmessage)
	if err != nil {
		return err
	}

	// 判断是否是当前节点需要处理的消息
	fmt.Println("HandleDeal message.Index: ", message.Index, "dkg.ID()", dkg.ID())
	if dkg.ID() == -1 || uint32(dkg.ID()) == message.Index {
		return nil
	}

	// 处理密钥份额。
	resp, err := dkg.DistKeyGenerator.ProcessDeal(message)
	if err != nil {
		return fmt.Errorf("HandleDeal error: %w", err)
	}
	fmt.Printf("HandleDeal : %+v\n", resp)

	// bt, err := json.Marshal(resp)
	// if err != nil {
	// 	return err
	// }

	// return dkg.Peer.Send(context.Background(), "response", bt)
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

	// 判断是否是当前节点需要处理的消息
	if dkg.ID() == -1 || uint32(dkg.ID()) == message.Index {
		return nil
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
		return nil
	}

	fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	// sc, err := dkg.DistKeyGenerator.SecretCommits()
	// if err != nil {
	// 	return fmt.Errorf("Generate secret commit: %w", err)
	// }

	// psc, err := types.SecretCommitsToProtocol(sc)
	// if err != nil {
	// 	return fmt.Errorf("SecretCommitsToProtocol : %w", err)
	// }

	// bt, err := json.Marshal(psc)
	// if err != nil {
	// 	return fmt.Errorf("HandleDealResp json.Marshal: %w", err)
	// }

	// return dkg.Peer.Send(context.Background(), "secret_commits", bt)
	return nil
}

// HandleDealMessage 处理密钥份额消息。
func (dkg *DKG) HandleJustification(data []byte) error {
	message := &rabin.Justification{}
	err := json.Unmarshal(data, message)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 判断是否是当前节点需要处理的消息
	if dkg.ID() == -1 || uint32(dkg.ID()) == message.Index {
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

	// 判断是否是当前节点需要处理的消息
	if dkg.ID() == -1 || uint32(dkg.ID()) != sc.Index {
		return nil
	}

	// 处理秘密提交
	_, err = dkg.DistKeyGenerator.ProcessSecretCommits(sc)
	if err != nil {
		return fmt.Errorf("ProcessSecretCommits: %w", err)
	}

	return nil
}

func (dkg *DKG) HandleMessage(msg *types.Message) error {
	switch msg.Type {
	case "deal":
		return dkg.HandleDeal(msg.Payload)
	case "deal_resp":
		return dkg.HandleDealResp(msg.Payload)
	case "justification":
		return dkg.HandleJustification(msg.Payload)
	case "secret_commits":
		return dkg.HandleSecretCommits(msg.Payload)
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}
