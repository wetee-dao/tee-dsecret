package dkg

import (
	"encoding/json"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/hashicorp/vault/shamir"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func NewSr25519(nodes int, threshold int) ([][]byte, types.AccountID, error) {
	// 生成64字节sr25519私钥种子（示例）
	kyr, err := sr25519.Scheme{}.Generate()
	if err != nil {
		return nil, types.AccountID{}, err
	}

	// 拆分私钥
	shares, err := shamir.Split(kyr.Seed(), nodes, threshold)
	if err != nil {
		return nil, types.AccountID{}, err
	}

	var pub [32]byte
	copy(pub[:], kyr.Public())

	return shares, types.AccountID(pub), nil
}

func (dkg *DKG) SaveSideKey(newMsg *model.ConsensusMsg) {
	if len(newMsg.NodeNewEpochShare) > 0 {
		dkg.NewEochSponsor = &newMsg.Sponsor
		dkg.NewSideKeyPub = newMsg.SideChainPub
		dkg.NewSideKeyShare = newMsg.NodeNewEpochShare
	}
	newMsg.NodeNewEpochShare = nil
}

func (dkg *DKG) SendSideKeyToSponsor() {
	if dkg.NewEochSponsor == nil {
		return
	}

	if len(dkg.SideKeyShare) == 0 {
		return
	}

	msg := model.SideKeyRebuildMsg{
		OldShare: dkg.SideKeyShare,
	}
	bt, _ := json.Marshal(msg)

	dkg.sendToNode(&dkg.NewEochSponsor.P2pId, "dkg", &model.Message{
		Type:    "consensus_side_key_rebuild",
		Payload: bt,
	})
}

func (dkg *DKG) SideKeyRebuild(OrgId string, data []byte) error {
	if !dkg.ConsensusIsbusy() {
		return nil
	}

	kmessage := &model.SideKeyRebuildMsg{}
	err := json.Unmarshal(data, kmessage)
	if err != nil {
		return err
	}

	if dkg.NewEochOldShares == nil {
		dkg.NewEochOldShares = make(map[string][]byte)
	}

	dkg.NewEochOldShares[OrgId] = kmessage.OldShare
	if len(dkg.NewEochOldShares) < dkg.Threshold {
		return nil
	}

	shares := make([][]byte, 0, len(dkg.NewEochOldShares))
	for v := range dkg.NewEochOldShares {
		shares = append(shares, dkg.NewEochOldShares[v])
	}

	if dkg.NewEoch <= StartEpoch {
		dkg.consensusSuccededBack(dkg.NewSideKeyPub, [64]byte{})
		return nil
	}

	oldkey, err := shamir.Combine(shares)
	if err != nil {
		return err
	}

	scheme := sr25519.Scheme{}
	kyr, err := scheme.FromSeed(oldkey)
	if err != nil {
		return err
	}

	sig, err := kyr.Sign(dkg.NewSideKeyPub.ToBytes())
	if err != nil {
		return err
	}

	var btsig [64]byte
	copy(btsig[:], sig)

	dkg.consensusSuccededBack(dkg.NewSideKeyPub, btsig)

	return nil
}

func Sr25519Verify(pub types.AccountID, msg []byte, sig [64]byte) bool {
	pubkey, err := sr25519.Scheme{}.FromPublicKey(pub.ToBytes())
	if err != nil {
		return false
	}
	return pubkey.Verify(msg, sig[:])
}
