package dkg

import (
	"fmt"

	types "wetee.app/dsecret/type"
	"wetee.app/dsecret/util"
)

// Handle deal message
func (dkg *DKG) HandleMessage(msg *types.Message) error {
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
	case "reencrypt_secret_request":
		_, err := dkg.HandleProcessReencrypt(msg.Payload, msg.MsgID)
		if err != nil {
			util.LogError("DEAL", "HandleReencryptSecretRequest err: ", err)
		}
		return err
	case "reencrypted_secret_share":
		_, err := dkg.HandleProcessReencrypt(msg.Payload, msg.MsgID)
		if err != nil {
			util.LogError("DEAL", "HandleReencryptSecretRequest err: ", err)
		}
		return err
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}
