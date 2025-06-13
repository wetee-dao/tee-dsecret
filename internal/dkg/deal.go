package dkg

import (
	"context"
	"encoding/json"
	"fmt"

	pedersen "go.dedis.ch/kyber/v4/share/dkg/pedersen"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/util"
)

// SendDealMessage 发送交易信息到指定节点
//
// 该函数将交易信息转换为协议格式，然后序列化并发送给指定的节点
// 主要步骤包括：
// 1. 将交易信息转换为协议消息格式
// 2. 将协议消息序列化为JSON字节切片
// 3. 通过Peer发送序列化后的消息到目标节点
//
// 参数:
// - ctx: 上下文，用于传递取消信号或超时信息
// - node: 目标节点信息
// - message: 待发送的交易信息
//
// 返回值:
// - error: 如果转换、序列化或发送过程中发生错误，则返回相应的错误
func (dkg *DKG) SendDealMessage(ctx context.Context, node *model.Node, message *pedersen.DealBundle, reshare int) error {
	// 将交易信息转换为协议消息格式
	pmessage, err := model.DealToProtocol(message)
	if err != nil {
		return err
	}
	pmessage.Reshare = reshare

	// 将协议消息序列化为JSON字节切片
	bt, err := json.Marshal(pmessage)
	if err != nil {
		return err
	}

	// 通过Peer发送序列化后的消息到目标节点
	return dkg.SendToNode(ctx, node, "dkg", &model.Message{
		Type:    "deal",
		Payload: bt,
	})
}

// HandleDeal 处理分发密钥生成协议中的交易消息
// data 是接收到的交易数据
func (dkg *DKG) HandleDeal(OrgId string, data []byte) error {
	// 加锁以确保线程安全
	dkg.mu.Lock()
	defer dkg.mu.Unlock()

	// 初始化交易消息结构体
	pmessage := &model.Deal{}
	err := json.Unmarshal(data, pmessage)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 将协议消息转换为交易对象
	deal, err := model.ProtocolToDeal(dkg.Suite, pmessage)
	if err != nil {
		return err
	}

	dkg.deals[OrgId] = deal
	num := len(dkg.DkgNodes)
	if pmessage.Reshare > 0 {
		num = pmessage.Reshare
	}
	if len(dkg.deals) < num {
		// 如果交易数量小于阈值，则返回错误
		return nil
	}

	deals := make([]*pedersen.DealBundle, 0, len(dkg.deals))
	for _, d := range dkg.deals {
		deals = append(deals, d)
	}

	// 处理密钥份额
	resp, err := dkg.DistKeyGenerator.ProcessDeals(deals)
	if err != nil || resp == nil {
		return fmt.Errorf("ProcessDeals error: %w", err)
	}

	// 如果交易未被批准，则返回错误
	// all nodes in the new group should have reported an error
	errNum := 0
	var errorLog []any = []any{"ProcessDeals ===> "}
	for _, r := range resp.Responses {
		errorLog = append(errorLog, fmt.Sprint(r.DealerIndex)+"="+fmt.Sprint(r.Status))
		if r.Status != pedersen.Success {
			errNum++
		}
	}
	fmt.Println(errorLog...)
	if errNum > 1 {
		return fmt.Errorf("ProcessDeals error: %w", err)
	}

	// 将响应对象序列化为字节切片
	bt, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	// 发送 deal resp 到所有参与节点
	for _, node := range dkg.DkgNodes {
		// 向节点发送交易响应
		err = dkg.SendToNode(context.Background(), node, "dkg", &model.Message{
			Type:    "deal_resp",
			Payload: bt,
		})
		if err != nil {
			util.LogError("DEAL", "Send deal_resp error", err)
		}
	}

	return nil
}

// HandleDealResp 处理交易响应消息
// 该函数接收一个字节切片作为参数，预期其内容为JSON格式的交易响应
// 它解析此响应，处理密钥份额，并相应地更新本地状态或与其他节点通信
func (dkg *DKG) HandleDealResp(OrgId string, data []byte) error {
	// 使用互斥锁保证并发安全
	dkg.mu.Lock()
	defer dkg.mu.Unlock()

	// 初始化一个交易响应对象，用于解析接收到的数据
	message := &pedersen.ResponseBundle{}
	// 解析数据到交易响应对象
	err := json.Unmarshal(data, message)
	if err != nil {
		// 如果解析失败，记录错误并返回
		util.LogError("DEAL", err)
		return err
	}

	dkg.responses[OrgId] = message
	if len(dkg.responses) < len(dkg.DkgNodes) {
		// 如果交易数量小于阈值，则返回错误
		return nil
	}

	responses := make([]*pedersen.ResponseBundle, 0, len(dkg.responses))
	for _, d := range dkg.responses {
		responses = append(responses, d)
	}

	// 处理密钥份额
	res, justification, err := dkg.DistKeyGenerator.ProcessResponses(responses)
	if err != nil {
		// 如果处理过程中出现错误，返回错误
		return fmt.Errorf("ProcessResponse: %w", err)
	}

	// 检查是否生成了密钥份额
	if res != nil {
		dkg.DkgKeyShare = model.DistKeyShare{
			Commits:  res.Key.Commits,
			PriShare: res.Key.Share,
		}
		dkg.DkgPubKey = res.Key.Public()

		// 保存密钥份额
		dkg.saveStore()
		return nil
	}

	// Justification 为 nil
	if justification == nil {
		// reshare 可能在这里获取私钥
		res, err := dkg.DistKeyGenerator.ProcessJustifications(nil)
		if err == nil {
			dkg.DkgKeyShare = model.DistKeyShare{
				Commits:  res.Key.Commits,
				PriShare: res.Key.Share,
			}
			dkg.DkgPubKey = res.Key.Public()

			// 保存密钥份额
			dkg.saveStore()
			return nil
		}
	}

	// // 将交易信息转换为协议消息格式
	// pmessage, err := model.JustificationToProtocol(justification)
	// if err != nil {
	// 	return err
	// }

	// // 将协议格式的秘密提交序列化为JSON
	// bt, err := json.Marshal(pmessage)
	// if err != nil {
	// 	// 如果序列化失败，返回错误
	// 	return fmt.Errorf("HandleDealResp json.Marshal: %w", err)
	// }

	// // 向所有DKG节点广播秘密提交
	// for _, node := range dkg.DkgNodes {
	// 	// 发送秘密提交给其他节点
	// 	err = dkg.SendToNode(context.Background(), node, "dkg", &model.Message{
	// 		Type:    "justification",
	// 		Payload: bt,
	// 	})
	// 	if err != nil {
	// 		// 如果发送失败，记录错误
	// 		util.LogError("DEAL", "Send justification error", err)
	// 	}
	// }

	return fmt.Errorf("HandleDealResp not implemented")
}

// HandleJustification 处理合理性证明消息
// 该函数接收一个字节切片数据作为输入，尝试解析并处理合理性证明
// 如果消息与当前节点相关，则更新内部状态
// 参数:
//   - data: 包含合理性证明信息的字节切片
//
// 返回值:
//   - error: 如果解析或处理消息过程中发生错误，则返回该错误
// func (dkg *DKG) HandleJustification(OrgId string, data []byte) error {
// 	// 加锁以确保线程安全
// 	dkg.mu.Lock()
// 	defer dkg.mu.Unlock()

// 	// 初始化交易消息结构体
// 	pmessage := &model.JustificationBundle{}
// 	err := json.Unmarshal(data, pmessage)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}

// 	// 将协议消息转换为交易对象
// 	message, err := model.ProtocolToJustification(dkg.Suite, pmessage)
// 	if err != nil {
// 		return err
// 	}

// 	// 如果消息索引与当前节点ID不匹配，则忽略该消息
// 	// 这里假设ID为-1表示无效或未初始化的节点
// 	if dkg.ID() == -1 || uint32(dkg.ID()) != message.DealerIndex {
// 		return nil
// 	}

// 	dkg.justifs = append(dkg.justifs, message)
// 	if len(dkg.responses) < len(dkg.DkgNodes) {
// 		// 如果交易数量小于阈值，则返回错误
// 		return nil
// 	}

// 	// 调用分布式密钥生成器的 ProcessJustification 方法处理合理性证明消息
// 	res, err := dkg.DistKeyGenerator.ProcessJustifications(dkg.justifs)
// 	if err != nil || res == nil {
// 		// 如果处理过程中发生错误，返回一个带有错误详情的新错误
// 		return fmt.Errorf("ProcessJustification: %w", err)
// 	}

// 	dkg.results = res
// 	dkg.DkgKeyShare = model.DistKeyShare{
// 		Commits:  res.Key.Commits,
// 		PriShare: res.Key.Share,
// 	}
// 	dkg.DkgPubKey = res.Key.Public()

// 	// 保存密钥份额
// 	dkg.saveStore()

// 	// 如果一切顺利，返回nil表示没有错误
// 	return nil
// }
