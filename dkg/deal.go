package dkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	rabin "go.dedis.ch/kyber/v3/share/dkg/rabin"
	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
)

// SendDealMessage 发送交易信息到指定节点。
//
// 该函数将交易信息转换为协议格式，然后序列化并发送给指定的节点。
// 主要步骤包括：
// 1. 将交易信息转换为协议消息格式。
// 2. 将协议消息序列化为JSON字节切片。
// 3. 通过Peer发送序列化后的消息到目标节点。
//
// 参数:
// - ctx: 上下文，用于传递取消信号或超时信息。
// - node: 目标节点信息。
// - message: 待发送的交易信息。
//
// 返回值:
// - error: 如果转换、序列化或发送过程中发生错误，则返回相应的错误。
func (dkg *DKG) SendDealMessage(ctx context.Context, node *types.Node, message *rabin.Deal) error {
	// 将交易信息转换为协议消息格式
	pmessage, err := types.DealToProtocol(message)
	if err != nil {
		return err
	}

	// 将协议消息序列化为JSON字节切片
	bt, err := json.Marshal(pmessage)
	if err != nil {
		return err
	}

	// 通过Peer发送序列化后的消息到目标节点
	return dkg.Peer.Send(ctx, node, "dkg", &types.Message{
		Type:    "deal",
		Payload: bt,
	})
}

// HandleDeal 处理分发密钥生成协议中的交易消息。
// data 是接收到的交易数据。
func (dkg *DKG) HandleDeal(data []byte) error {
	// 加锁以确保线程安全。
	dkg.mu.Lock()
	defer dkg.mu.Unlock()

	// 初始化交易消息结构体。
	pmessage := &types.Deal{}
	// 解析接收到的数据到交易消息结构体。
	err := json.Unmarshal(data, pmessage)
	if err != nil {
		fmt.Println(err)
		return err
	}
	// 将协议消息转换为交易对象。
	deal, err := types.ProtocolToDeal(dkg.Suite, pmessage)
	if err != nil {
		return err
	}

	// 处理密钥份额。
	resp, err := dkg.DistKeyGenerator.ProcessDeal(deal)
	if err != nil {
		return fmt.Errorf("HandleDeal error: %w", err)
	}

	// 如果交易未被批准，则返回错误。
	if !resp.Response.Approved {
		return errors.New("deal rejected")
	}

	// 将响应对象序列化为字节切片。
	bt, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	// 发送 deal resp 到所有参与节点。
	for _, node := range dkg.DkgNodes {
		// 跳过自身。
		if node.PeerID() == dkg.Peer.ID() {
			continue
		}

		// 向节点发送交易响应。
		err = dkg.Peer.Send(context.Background(), node, "dkg", &types.Message{
			Type:    "deal_resp",
			Payload: bt,
		})
		if err != nil {
			util.LogError("DEAL", "Send deal_resp error", err)
		}
	}

	return nil
}

// HandleDealResp 处理交易响应消息。
// 该函数接收一个字节切片作为参数，预期其内容为JSON格式的交易响应。
// 它解析此响应，处理密钥份额，并相应地更新本地状态或与其他节点通信。
func (dkg *DKG) HandleDealResp(data []byte) error {
	// 使用互斥锁保证并发安全。
	dkg.mu.Lock()
	defer dkg.mu.Unlock()

	// 初始化一个交易响应对象，用于解析接收到的数据。
	message := &rabin.Response{}
	// 解析数据到交易响应对象。
	err := json.Unmarshal(data, message)
	if err != nil {
		// 如果解析失败，记录错误并返回。
		util.LogError("DEAL", err)
		return err
	}

	// 处理密钥份额。
	justification, err := dkg.DistKeyGenerator.ProcessResponse(message)
	if err != nil {
		// 如果处理过程中出现错误，返回错误。
		return fmt.Errorf("HandleDealResp ProcessResponse: %w", err)
	}

	// 检查是否生成了交易证明。
	// tJustification，证明 Deal 消息的无效性
	if justification != nil {
		// 如果生成了证明，记录并返回。
		util.LogError("DEAL", "Got justification during response process for ", message.Index, justification)
		return nil
	}

	// 检查是否已经通过认证。
	// 已经判断为有效了
	if !dkg.DistKeyGenerator.Certified() {
		// 如果未通过认证，返回。
		return nil
	}

	// 生成秘密提交。
	sc, err := dkg.DistKeyGenerator.SecretCommits()
	if err != nil {
		// 如果生成失败，返回错误。
		return fmt.Errorf("Generate secret commit: %w", err)
	}

	// 将秘密提交转换为协议格式。
	psc, err := types.SecretCommitsToProtocol(sc)
	if err != nil {
		// 如果转换失败，返回错误。
		return fmt.Errorf("SecretCommitsToProtocol : %w", err)
	}

	// 将协议格式的秘密提交序列化为JSON。
	bt, err := json.Marshal(psc)
	if err != nil {
		// 如果序列化失败，返回错误。
		return fmt.Errorf("HandleDealResp json.Marshal: %w", err)
	}

	// 向所有DKG节点广播秘密提交。
	for _, node := range dkg.DkgNodes {
		// 跳过自身。
		if node.PeerID() == dkg.Peer.ID() {
			continue
		}
		// 发送秘密提交给其他节点。
		err = dkg.Peer.Send(context.Background(), node, "dkg", &types.Message{
			Type:    "secret_commits",
			Payload: bt,
		})
		if err != nil {
			// 如果发送失败，记录错误。
			util.LogError("DEAL", "Send secret_commits error", err)
		}
	}

	return nil
}

// HandleJustification 处理合理性证明消息。
// 该函数接收一个字节切片数据作为输入，尝试解析并处理合理性证明。
// 如果消息与当前节点相关，则更新内部状态。
// 参数:
//   - data: 包含合理性证明信息的字节切片。
//
// 返回值:
//   - error: 如果解析或处理消息过程中发生错误，则返回该错误。
func (dkg *DKG) HandleJustification(data []byte) error {
	// 加锁以确保线程安全。
	dkg.mu.Lock()
	defer dkg.mu.Unlock()

	// 初始化一个合理性证明消息结构体。
	message := &rabin.Justification{}
	// 尝试将输入数据反序列化为合理性证明消息结构体。
	err := json.Unmarshal(data, message)
	if err != nil {
		// 如果反序列化失败，打印错误并返回。
		fmt.Println(err)
		return err
	}

	// 如果消息索引与当前节点ID不匹配，则忽略该消息。
	// 这里假设ID为-1表示无效或未初始化的节点。
	if dkg.ID() == -1 || uint32(dkg.ID()) != message.Index {
		return nil
	}

	// 调用分布式密钥生成器的ProcessJustification方法处理合理性证明消息。
	err = dkg.DistKeyGenerator.ProcessJustification(message)
	if err != nil {
		// 如果处理过程中发生错误，返回一个带有错误详情的新错误。
		return fmt.Errorf("ProcessJustification: %w", err)
	}

	// 如果一切顺利，返回nil表示没有错误。
	return nil
}

// HandleSecretCommits 处理秘密提交阶段的数据。
// 该函数接收一个字节数组作为输入，该数组包含一个序列化的秘密提交对象。
// 它的主要工作包括：
// 1. 将字节数组反序列化为协议对象。
// 2. 将协议对象转换为秘密提交对象。
// 3. 处理这些秘密提交，以进行分布式密钥生成的下一步。
// 4. 保存分发的密钥份额和公钥。
//
// 参数:
//
//	data []byte: 包含秘密提交信息的字节数组。
//
// 返回值:
//
//	error: 如果处理过程中发生错误，返回错误信息；否则返回nil。
func (dkg *DKG) HandleSecretCommits(data []byte) error {
	// 加锁以确保线程安全。
	dkg.mu.Lock()
	defer dkg.mu.Unlock()

	// 转换协议对象
	psc := &types.SecretCommitJson{}
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

	// 记录身份认证完成的日志
	util.LogOk("DEAL", "身份认证完成 ================================================")
	// 保存分发的密钥份额和公钥
	dkg.DkgKeyShare = types.DistKeyShare{
		Commits:  distkey.Commitments(),
		PriShare: distkey.PriShare(),
	}
	dkg.DkgPubKey = distkey.Public()
	// 保存数据到持久化存储
	dkg.store()

	return nil
}
