package dkg

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
)

// Handle deal message
func (dkg *DKG) HandleDkg(msg *types.Message) error {
	switch msg.Type {
	case "deal":
		err := dkg.HandleDeal(msg.Payload)
		if err != nil {
			util.LogError("DEAL", "HandleDeal err: ", err)
		}
		return err
	case "deal_resp":
		err := dkg.HandleDealResp(msg.Payload)
		if err != nil {
			util.LogError("DEAL", "HandleDealResp err: ", err)
		}
		return err
	case "justification":
		err := dkg.HandleJustification(msg.Payload)
		if err != nil {
			util.LogError("DEAL", "HandleJustification err: ", err)
		}
		return err
	case "secret_commits":
		err := dkg.HandleSecretCommits(msg.Payload)
		if err != nil {
			util.LogError("DEAL", "HandleSecretCommits err: ", err)
		}
		return err
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

func (dkg *DKG) HandleWorker(msg *types.Message) error {
	switch msg.Type {
	/// -------------------- Proof -----------------------
	case "upload_cluster_proof":
		hash, err := dkg.HandleUploadClusterProof(msg.Payload, msg.MsgID, msg.OrgId)
		if msg.OrgId != "" && msg.MsgID != "" {
			n := dkg.GetNode(msg.OrgId)
			if n == nil {
				return fmt.Errorf("node not found: %s", msg.OrgId)
			}
			errStr := ""
			if err != nil {
				errStr = err.Error()
			}
			if err := dkg.SendToNode(context.Background(), n, "worker", &types.Message{
				MsgID:   msg.MsgID,
				Type:    "upload_cluster_proof_reply",
				Payload: hash,
				Error:   errStr,
			}); err != nil {
				return errors.New("send to node: " + err.Error())
			}
		}

		return err
	case "sign_cluster_proof":
		err := dkg.HandleSignClusterProof(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			util.LogError("WORKER", "HandleSignClusterProof err: ", err)
		}
		return err
	case "sign_cluster_proof_reply":
		err := dkg.HandleSignClusterProofReply(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			util.LogError("WORKER", "HandleSignClusterProofReply err: ", err)
		}
		return err
	/// -------------------- Reencrypt -----------------------
	case "reencrypt_secret_remote_request":
		key, err := dkg.SendEncryptedSecretRequest(msg.Payload, msg.MsgID, msg.OrgId)
		if msg.OrgId != "" && msg.MsgID != "" {
			n := dkg.GetNode(msg.OrgId)
			if n == nil {
				return fmt.Errorf("node not found: %s", msg.OrgId)
			}

			errStr := ""
			var keyBt []byte
			if err != nil {
				errStr = err.Error()
			} else {
				keyBt, _ = json.Marshal(key)
			}

			if err := dkg.SendToNode(context.Background(), n, "worker", &types.Message{
				MsgID:   msg.MsgID,
				Type:    "reencrypt_secret_remote_reply",
				Payload: keyBt,
				Error:   errStr,
			}); err != nil {
				return errors.New("send to node: " + err.Error())
			}
		}
		return err
	case "reencrypt_secret_request":
		err := dkg.HandleProcessReencrypt(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			util.LogError("secret", "HandleReencryptSecretRequest err: ", err)
		}
		return err
	case "reencrypted_secret_reply":
		err := dkg.HandleReencryptedShare(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			util.LogError("secret", "HandleReencryptSecretRequest err: ", err)
		}
		return err
	/// -------------------- Reencrypt -----------------------
	case "work_launch_request":
		key, err := dkg.HandleWorkLaunchRequest(msg.Payload, msg.MsgID, msg.OrgId)
		if msg.OrgId != "" && msg.MsgID != "" {
			n := dkg.GetNode(msg.OrgId)
			if n == nil {
				return fmt.Errorf("node not found: %s", msg.OrgId)
			}

			errStr := ""
			var keyBt []byte
			if err != nil {
				errStr = err.Error()
			} else {
				keyBt, _ = json.Marshal(key)
			}

			if err := dkg.SendToNode(context.Background(), n, "worker", &types.Message{
				MsgID:   msg.MsgID,
				Type:    "work_launch_reply",
				Payload: keyBt,
				Error:   errStr,
			}); err != nil {
				return errors.New("send to node: " + err.Error())
			}
		}
		return err
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}
