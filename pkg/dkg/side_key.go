package dkg

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/hashicorp/vault/shamir"
	"github.com/vedhavyas/go-subkey/v2"
	"github.com/vedhavyas/go-subkey/v2/sr25519"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func NewSr25519Split(nodes int, threshold int) ([][]byte, types.AccountID, error) {
	// 生成64字节sr25519私钥种子（示例）
	kyr, err := sr25519.Scheme{}.Generate()
	if err != nil {
		return nil, types.AccountID{}, err
	}

	return SplitSr25519(kyr, nodes, threshold)
}

func SplitSr25519(key subkey.KeyPair, nodes, threshold int) ([][]byte, types.AccountID, error) {
	// 拆分私钥
	shares, err := shamir.Split(key.Seed(), nodes, threshold)
	if err != nil {
		return nil, types.AccountID{}, err
	}

	var pub [32]byte
	copy(pub[:], key.Public())

	return shares, types.AccountID(pub), nil
}

func CombineSr25519(shares [][]byte) (subkey.KeyPair, error) {
	oldkey, err := shamir.Combine(shares)
	if err != nil {
		return nil, err
	}

	scheme := sr25519.Scheme{}

	return scheme.FromSeed(oldkey)
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

	if dkg.NewOldSharesCache == nil {
		dkg.NewOldSharesCache = make(map[string][]byte)
	}

	dkg.NewOldSharesCache[OrgId] = kmessage.OldShare
	if len(dkg.NewOldSharesCache) <= dkg.Threshold {
		return nil
	}

	shares := make([][]byte, 0, len(dkg.NewOldSharesCache))
	for v := range dkg.NewOldSharesCache {
		shares = append(shares, dkg.NewOldSharesCache[v])
	}

	kyr, err := CombineSr25519(shares)
	if err != nil {
		dkg.consensusFailBack(errors.New("SideKeyRebuild CombineSr25519 error:" + err.Error()))
		return err
	}

	sig, err := kyr.Sign(dkg.NewSideKeyPub.ToBytes())
	if err != nil {
		dkg.consensusFailBack(errors.New("SideKeyRebuild sign:" + err.Error()))
		return err
	}

	var btsig [64]byte
	copy(btsig[:], sig)

	if !bytes.Equal(dkg.SideKeyPub[:], kyr.Public()) {
		dkg.consensusFailBack(errors.New("SideKeyRebuild old key is not equal to privkey public"))
		return nil
	}

	dkg.consensusSuccededBack(dkg.NewSideKeyPub, btsig)
	dkg.NewOldSharesCache = nil
	return nil
}

func Sr25519Verify(pub types.AccountID, msg []byte, sig [64]byte) bool {
	pubkey, err := sr25519.Scheme{}.FromPublicKey(pub.ToBytes())
	if err != nil {
		return false
	}
	return pubkey.Verify(msg, sig[:])
}
