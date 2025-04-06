package dkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	pubsub "github.com/libp2p/go-libp2p-pubsub"

	"wetee.app/dsecret/internal/chain"
	"wetee.app/dsecret/internal/store"
	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
)

// HandleDkg 处理不同的DKG消息类型
// msg: 被处理的消息对象
// 返回：可能的错误
func (dkg *DKG) HandleDkg(msg *types.Message) error {
	// 根据消息类型执行相应的处理逻辑
	switch msg.Type {
	case "deal":
		// 处理交易消息
		err := dkg.HandleDeal(msg.OrgId, msg.Payload)
		if err != nil {
			util.LogError("DEAL <<<<<<<< ERROR", "HandleDeal:", err)
		}
		return err
	case "deal_resp":
		// 处理交易响应消息
		err := dkg.HandleDealResp(msg.OrgId, msg.Payload)
		if err != nil {
			util.LogError("DEAL <<<<<<<<<<<<<<<< ERROR", "HandleDealResp:", err)
		}
		return err
	// case "justification":
	// 	// 处理证明消息
	// 	err := dkg.HandleJustification(msg.OrgId, msg.Payload)
	// 	if err != nil {
	// 		util.LogError("DEAL <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< ERROR", "HandleJustification:", err)
	// 	}
	// 	return err
	default:
		// 如果消息类型未知，返回错误
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

// HandleWorker 处理来自worker的消息
// msg: 待处理的消息
func (dkg *DKG) HandleWorker(msg *types.Message) error {
	// 检查链的元数据
	err := chain.ChainIns.CheckMetadata()
	if err != nil {
		// 记录检查元数据时的错误
		util.LogError("CheckMetadata", err)
		return err
	}

	// 根据消息类型处理不同的逻辑
	switch msg.Type {
	/// -------------------- Proof -----------------------
	case "upload_cluster_proof":
		// 处理上传集群证明的消息
		hash, err := dkg.HandleUploadClusterProof(msg.Payload, msg.MsgID, msg.OrgId)
		if msg.OrgId != "" && msg.MsgID != "" {
			// 获取发送方节点
			n := dkg.GetNode(msg.OrgId)
			if n == nil {
				// 如果节点不存在，返回错误
				return fmt.Errorf("node not found: %s", msg.OrgId)
			}
			errStr := ""
			if err != nil {
				// 如果有错误，记录错误信息
				errStr = err.Error()
			}
			// 发送回复消息给节点
			if err := dkg.SendToNode(context.Background(), n, "worker", &types.Message{
				MsgID:   msg.MsgID,
				Type:    "upload_cluster_proof_reply",
				Payload: hash,
				Error:   errStr,
			}); err != nil {
				// 发送消息失败，返回错误
				return errors.New("send to node: " + err.Error())
			}
		}
		return err
	case "sign_cluster_proof":
		// 处理签名集群证明的消息
		err := dkg.HandleSignClusterProof(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			// 记录签名集群证明时的错误
			util.LogError("WORKER", "HandleSignClusterProof err: ", err)
		}
		return err
	case "sign_cluster_proof_reply":
		// 处理签名集群证明回复的消息
		err := dkg.HandleSignClusterProofReply(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			// 记录签名集群证明回复时的错误
			util.LogError("WORKER", "HandleSignClusterProofReply err: ", err)
		}
		return err
	/// -------------------- Reencrypt -----------------------
	case "reencrypt_secret_remote_request":
		// 发送加密的密钥请求
		key, err := dkg.SendEncryptedSecretRequest(msg.Payload, msg.MsgID, msg.OrgId)
		if msg.OrgId != "" && msg.MsgID != "" {
			// 获取发送方节点
			n := dkg.GetNode(msg.OrgId)
			if n == nil {
				// 如果节点不存在，返回错误
				return fmt.Errorf("node not found: %s", msg.OrgId)
			}

			errStr := ""
			var keyBt []byte
			if err != nil {
				// 如果有错误，记录错误信息
				errStr = err.Error()
			} else {
				// 将密钥转换为字节
				keyBt, _ = json.Marshal(key)
			}

			// 发送回复消息给节点
			if err := dkg.SendToNode(context.Background(), n, "worker", &types.Message{
				MsgID:   msg.MsgID,
				Type:    "reencrypt_secret_remote_reply",
				Payload: keyBt,
				Error:   errStr,
			}); err != nil {
				// 发送消息失败，返回错误
				return errors.New("send to node: " + err.Error())
			}
		}
		return err
	case "reencrypt_secret_request":
		// 处理重新加密密钥请求的消息
		err := dkg.HandleProcessReencrypt(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			// 记录处理重新加密密钥请求时的错误
			util.LogError("secret", "HandleReencryptSecretRequest err: ", err)
		}
		return err
	case "reencrypted_secret_reply":
		// 处理重新加密的密钥回复的消息
		err := dkg.HandleReencryptedShare(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			// 记录处理重新加密的密钥回复时的错误
			util.LogError("secret", "HandleReencryptSecretRequest err: ", err)
		}
		return err
	/// -------------------- Work Launch -----------------------
	case "work_launch_request":
		// 处理工作启动请求的消息
		key, err := dkg.HandleWorkLaunchRequest(msg.Payload, msg.MsgID, msg.OrgId)
		if msg.OrgId != "" && msg.MsgID != "" {
			// 获取发送方节点
			n := dkg.GetNode(msg.OrgId)
			if n == nil {
				// 如果节点不存在，返回错误
				return fmt.Errorf("node not found: %s", msg.OrgId)
			}

			errStr := ""
			keyBt := []byte{}
			if err != nil {
				// 如果有错误，记录错误信息
				errStr = err.Error()
			} else {
				// 将密钥转换为字节
				keyBt, _ = json.Marshal(key)
			}

			// 发送回复消息给节点
			if err := dkg.SendToNode(context.Background(), n, "worker", &types.Message{
				MsgID:   msg.MsgID,
				Type:    "work_launch_reply",
				Payload: keyBt,
				Error:   errStr,
			}); err != nil {
				// 发送消息失败，返回错误
				return errors.New("send to node: " + err.Error())
			}
		}
		return err
	default:
		// 默认情况下不执行任何操作，返回nil
		return nil
	}
}

// Handle pub msg data save
func (r *DKG) HandleSecretSave(ctx context.Context) {
	sub, err := r.Peer.Sub(ctx, "secret")
	if err != nil {
		return
	}

	// 使用缓冲通道异步处理消息
	msgCh := make(chan *pubsub.Message, 100)
	go func() {
		defer close(msgCh)
		for {
			msg, err := sub.Next(ctx)
			if err != nil {
				log.Printf("Error receiving message: %v", err)
				return
			}
			msgCh <- msg
		}
	}()

	for msg := range msgCh {
		// 解析消息
		var datas []types.Kvs
		err = json.Unmarshal(msg.Data, &datas)
		if err != nil {
			fmt.Println("Error unmarshalling message data: ", err)
			continue
		}

		for _, data := range datas {
			fmt.Println("-------------------------Save key: ", data.K)
			err := store.SetKey("secret", data.K, data.V)
			if err != nil {
				fmt.Println("Error setting key: ", err)
				continue
			}
		}
	}
}
