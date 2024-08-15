package dkg

import (
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
	/// -------------------- Reencrypt -----------------------
	case "reencrypt_secret_request":
		err := dkg.HandleProcessReencrypt(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			util.LogError("DEAL", "HandleReencryptSecretRequest err: ", err)
		}
		return err
	case "reencrypted_secret_reply":
		err := dkg.HandleReencryptedShare(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			util.LogError("DEAL", "HandleReencryptSecretRequest err: ", err)
		}
		return err
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}

func (dkg *DKG) HandleWorker(msg *types.Message) error {
	switch msg.Type {
	case "upload_cluster_proof":
		err := dkg.HandleUploadClusterProof(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			util.LogError("DEAL", "HandleUploadClusterProof err: ", err)
		}
		return err
	case "sign_cluster_proof":
		err := dkg.HandleSignClusterProof(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			util.LogError("DEAL", "HandleSignClusterProof err: ", err)
		}
		return err
	case "sign_cluster_proof_reply":
		err := dkg.HandleSignClusterProofReply(msg.Payload, msg.MsgID, msg.OrgId)
		if err != nil {
			util.LogError("DEAL", "HandleSignClusterProofReply err: ", err)
		}
		return err
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}
