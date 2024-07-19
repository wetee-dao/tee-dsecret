package dkg

import (
	"fmt"

	"wetee.app/dsecret/types"
)

func (dkg *DKG) HandleMessage(msg *types.Message) error {
	switch msg.Type {
	case "deal":
		err := dkg.HandleDeal(msg.Payload)
		if err != nil {
			fmt.Println("HandleDeal err: ", err)
		}
		return err
	case "deal_resp":
		err := dkg.HandleDealResp(msg.Payload)
		if err != nil {
			fmt.Println("HandleDealResp err: ", err)
		}
		return err
	case "justification":
		err := dkg.HandleJustification(msg.Payload)
		if err != nil {
			fmt.Println("HandleJustification err: ", err)
		}
		return err
	case "secret_commits":
		err := dkg.HandleSecretCommits(msg.Payload)
		if err != nil {
			fmt.Println("HandleSecretCommits err: ", err)
		}
		return err
	case "reencrypt_secret_request":
		_, err := dkg.HandleProcessReencrypt(msg.Payload, msg.MsgID)
		if err != nil {
			fmt.Println("HandleReencryptSecretRequest err: ", err)
		}
		return err
	case "reencrypted_secret_share":
		_, err := dkg.HandleProcessReencrypt(msg.Payload, msg.MsgID)
		if err != nil {
			fmt.Println("HandleReencryptSecretRequest err: ", err)
		}
		return err
	default:
		return fmt.Errorf("unknown message type: %s", msg.Type)
	}
}
