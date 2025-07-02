package dkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
	"go.dedis.ch/kyber/v4"
	pedersen "go.dedis.ch/kyber/v4/share/dkg/pedersen"
	"go.dedis.ch/kyber/v4/sign/schnorr"
)

const StartEpoch = 1

func (dkg *DKG) TryEpochConsensus(msg model.ConsensusMsg, callback func(types.AccountID, [64]byte), fail func(error)) error {
	if dkg.ConsensusIsbusy() {
		util.LogError("DKG Consensus", "in consensus")
		return errors.New("in consensus")
	}

	if dkg.DkgKeyShare == nil && msg.Epoch > StartEpoch {
		util.LogError("DKG Consensus", "msg.Epoch", msg.Epoch, "| Node is not old validator, cannot start consensus")
		return errors.New("node is not old validator, cannot start consensus")
	}

	if msg.Epoch > StartEpoch {
		msg.ShareCommits = *util.DeepCopy(dkg.DkgKeyShare.Commits)
		msg.ConsensusNodeNum = len(dkg.Nodes)
		msg.OldValidators = *util.DeepCopy(dkg.Nodes)
	} else {
		msg.Epoch = StartEpoch
		msg.ShareCommits = model.KyberPoints{Public: []kyber.Point{}}
		msg.ConsensusNodeNum = 0
	}

	// create new side chain key
	shares, pub, err := NewSr25519(len(msg.Validators), len(msg.Validators)*2/3)
	if err != nil {
		util.LogError("DKG Consensus", "NewSr25519 error:", err)
		return err
	}

	// generate new validators sidekey
	newShares := map[string][]byte{}
	dkgSigner := dkg.Signer.GetPublic().SS58()
	validator := new(model.Validator)
	for i, v := range msg.Validators {
		newShares[v.P2pId.SS58()] = shares[i]
		if v.ValidatorId.SS58() == dkgSigner {
			validator = v
		}
	}

	if validator == nil {
		util.LogError("DKG Consensus", "Currunt Node is not in new validators")
		return err
	}

	msg.SideChainPub = pub
	msg.Sponsor = *validator

	// Must set nil
	dkg.NewSideKeyShares = newShares
	dkg.consensusSuccededBack = callback
	dkg.consensusFailBack = fail

	bt, _ := json.Marshal(msg)
	dkg.DkgOutHandler(&model.Message{
		Type:    "consensus",
		Payload: bt,
	})

	return nil
}

// start dkg consensus
func (dkg *DKG) startConsensus(msg model.ConsensusMsg) error {
	if dkg.ConsensusIsbusy() {
		return errors.New("DKG Consensus going")
	}

	if dkg.Epoch >= msg.Epoch {
		return errors.New("DKG Epoch is not need to update")
	}

	dkg.setConsensusBusy()
	dkg.addConsensusTimeout()

	if msg.Epoch <= StartEpoch {
		util.LogWithGreen("InitConsensus Epoch ======> ", msg.Epoch)
		return dkg.initConsensus(msg)
	}

	util.LogWithGreen("ReConsensus Epoch ======> ", msg.Epoch)
	return dkg.reConsensus(msg)
}

// Init Consensus
func (dkg *DKG) initConsensus(msg model.ConsensusMsg) error {
	// if flag.Lookup("test.v") == nil {
	// 	go dkg.HandleSecretSave()
	// }
	dkg.Nodes = msg.Validators
	dkg.NewNodes = msg.Validators
	dkg.Threshold = len(msg.Validators) * 2 / 3

	// 如果已经初始化，则直接返回
	if dkg.status == 1 {
		return nil
	}

	// dkg 节点列表
	nodes := make([]pedersen.Node, 0, len(dkg.Nodes))
	for i, p := range dkg.Nodes {
		nodes = append(nodes, pedersen.Node{
			Index:  uint32(i),
			Public: p.ValidatorId.Point(),
		})
	}
	signer := schnorr.NewScheme(dkg.Suite)

	// 初始化协议配置
	conf := pedersen.Config{
		Suite:     dkg.Suite,
		NewNodes:  nodes,
		Threshold: dkg.Threshold,
		Auth:      signer,
		FastSync:  true,
		Longterm:  dkg.Signer.Scalar(),
		Nonce:     epochToNonce(0),
		Log:       dkg.log,
	}

	// initialize dealer
	var err error
	dkg.DistKeyGenerator, err = pedersen.NewDistKeyHandler(&conf)
	if err != nil {
		dkg.finishDkgConsensusStep(false, "pedersen.NewDistKeyHandler")
		return fmt.Errorf("failed to initialize DKG protocol: %w", err)
	}

	// 等待节点连接
	if dkg.connectLen() < len(dkg.Nodes)-1 {
		util.LogError("DKG Consensus", "exapect node count:", len(dkg.Nodes)-1, ", got:", dkg.connectLen())
		dkg.finishDkgConsensusStep(false, "dkg.connectLen()+1 < len(dkg.Nodes)")
		return fmt.Errorf("waiting for nodes to connect")
	}

	// 获取当前节点的协议
	deal, err := dkg.DistKeyGenerator.Deals()
	if err != nil {
		dkg.finishDkgConsensusStep(false, "dkg.DistKeyGenerator.Deals")
		return fmt.Errorf("failed to generate key shares: %w", err)
	}

	// 开启节点共识
	for _, node := range dkg.Nodes {
		nodeShare := []byte{}
		if key, ok := dkg.NewSideKeyShares[node.P2pId.SS58()]; ok {
			nodeShare = key
		}
		err = dkg.sendDealMessage(&node.P2pId, &model.ConsensusMsg{
			DealBundle:        &model.DealBundle{DealBundle: deal},
			ShareCommits:      model.KyberPoints{},
			Validators:        dkg.Nodes,
			Epoch:             msg.Epoch,
			SideChainPub:      msg.SideChainPub,
			NodeNewEpochShare: nodeShare,
			Sponsor:           msg.Sponsor,
		})
		if err != nil {
			fmt.Println("Send error:", err)
		}
	}

	// for {
	// 	if dkg.DkgKeyShare.PriShare != nil {
	// 		break
	// 	}
	// 	time.Sleep(time.Second)
	// }

	// dkg.deals = map[string]*model.DealBundle{}
	// dkg.responses = map[string]*pedersen.ResponseBundle{}
	// dkg.justifs = []*pedersen.JustificationBundle{}

	// if dkg.log != nil {
	// 	dkg.log.Info("DKG uccessfully init <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	// }

	return nil
}

// Re-consensus DKG
func (dkg *DKG) reConsensus(msg model.ConsensusMsg) error {
	// old
	dkg.Threshold = len(msg.OldValidators) * 2 / 3
	dkg.Nodes = msg.OldValidators
	// new
	dkg.NewNodes = msg.Validators
	dkg.NewEoch = msg.Epoch

	// new DKG 节点列表
	newNodes := make([]pedersen.Node, 0, len(msg.Validators))
	for i, p := range msg.Validators {
		newNodes = append(newNodes, pedersen.Node{
			Index:  uint32(i),
			Public: p.ValidatorId.Point(),
		})
	}

	// 获取旧节点列表
	oldNodes := make([]pedersen.Node, 0, len(dkg.Nodes))
	for i, p := range dkg.Nodes {
		oldNodes = append(oldNodes, pedersen.Node{
			Index:  uint32(i),
			Public: p.ValidatorId.Point(),
		})
	}

	newThreshold := len(msg.Validators) * 2 / 3

	// 初始化协议配置
	conf := pedersen.Config{
		OldNodes:     oldNodes,
		OldThreshold: dkg.Threshold,
		Threshold:    newThreshold,
		NewNodes:     newNodes,
		Nonce:        epochToNonce(msg.Epoch),
		Suite:        dkg.Suite,
		Auth:         schnorr.NewScheme(dkg.Suite),
		FastSync:     true,
		Longterm:     dkg.Signer.Scalar(),
		Log:          dkg.log,
	}

	if dkg.DkgKeyShare != nil {
		priv := dkg.DkgKeyShare
		conf.Share = &pedersen.DistKeyShare{
			Commits: priv.Commits.Public,
			Share:   priv.PriShare.PriShare,
		}
	} else {
		conf.PublicCoeffs = msg.ShareCommits.Public
	}

	var err error
	dkg.DistKeyGenerator, err = pedersen.NewDistKeyHandler(&conf)
	if err != nil {
		dkg.finishDkgConsensusStep(false, "pedersen.NewDistKeyHandler(&conf)")
		fmt.Println("unable to create DistKeyGenerator", err.Error())
		return err
	}

	priShare := dkg.DkgKeyShare

	// 重置 DKG Key
	dkg.deals = map[string]*model.DealBundle{}
	dkg.responses = map[string]*pedersen.ResponseBundle{}
	dkg.justifs = []*pedersen.JustificationBundle{}

	// old node not issue deals
	if priShare == nil {
		util.LogWithCyan("DKG", "old node not issue deals")
		return nil
	}

	// 获取当前节点的协议
	deal, err := dkg.DistKeyGenerator.Deals()
	if err != nil {
		dkg.finishDkgConsensusStep(false, "dkg.DistKeyGenerator.Deals()")
		return fmt.Errorf("failed to generate key shares: %w", err)
	}

	// 开启节点共识
	for _, node := range dkg.NewNodes {
		nodeShare := []byte{}
		if key, ok := dkg.NewSideKeyShares[node.P2pId.SS58()]; ok {
			nodeShare = key
		}
		msg.DealBundle = &model.DealBundle{DealBundle: deal}
		msg.NodeNewEpochShare = nodeShare

		err = dkg.sendDealMessage(&node.P2pId, &msg)
		if err != nil {
			fmt.Println("Send error:", err)
		}
	}

	return nil
}

// set consensus time out
func (dkg *DKG) addConsensusTimeout() {
	if dkg.failConsensusTimer != nil {
		dkg.failConsensusTimer.Stop()
	}
	dkg.failConsensusTimer = time.AfterFunc(time.Second*30, func() {
		dkg.finishDkgConsensusStep(false, "timeout")
	})
}

// stop consensus
func (dkg *DKG) finishDkgConsensusStep(isok bool, tag string) {
	if dkg.failConsensusTimer != nil {
		dkg.failConsensusTimer.Stop()
	}

	dkg.deals = map[string]*model.DealBundle{}
	dkg.responses = map[string]*pedersen.ResponseBundle{}
	dkg.justifs = []*pedersen.JustificationBundle{}
	if !isok {
		util.LogWithRed("DKG dkg consensus", "failed <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< New Epoch", dkg.NewEoch, "error", tag)
		dkg.setConsensusFree()
		if dkg.consensusFailBack != nil {
			dkg.consensusFailBack(errors.New("DKG dkg consensus failed"))
		}
		return
	}

	// cancel new epoch if not sucsess
	dkg.failConsensusTimer = time.AfterFunc(time.Second*80, func() {
		dkg.cancelNewEpoch()
	})

	if dkg.NewEoch > StartEpoch {
		dkg.SendSideKeyToSponsor()
		dkg.saveStore()
	} else {
		if dkg.NewEochSponsor.ValidatorId.SS58() == dkg.Signer.GetPublic().SS58() && dkg.consensusSuccededBack != nil {
			dkg.consensusSuccededBack(dkg.NewSideKeyPub, [64]byte{})
		}
	}
}

// to next epoch
func (dkg *DKG) ToNewEpoch() {
	if dkg.failConsensusTimer != nil {
		dkg.failConsensusTimer.Stop()
	}

	defer dkg.setConsensusFree()

	if dkg.NewDkgKeyShare == nil {
		util.LogWithRed("DKG consensus ToNewEpoch", "failed <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< New Epoch", dkg.NewEoch)
		if dkg.consensusFailBack != nil {
			dkg.consensusFailBack(errors.New("DKG consensus failed"))
		}
		return
	}

	dkg.Nodes = dkg.NewNodes
	dkg.Epoch = dkg.NewEoch
	dkg.Threshold = len(dkg.NewNodes) * 2 / 3
	dkg.DkgPubKey = dkg.NewDkgPubKey
	dkg.DkgKeyShare = dkg.NewDkgKeyShare
	dkg.SideKeyPub = dkg.NewSideKeyPub
	dkg.SideKeyShare = dkg.NewSideKeyShare

	dkg.NewNodes = nil
	dkg.NewEoch = 0
	dkg.NewDkgPubKey = nil
	dkg.NewDkgKeyShare = nil
	dkg.NewSideKeyPub = types.AccountID{}
	dkg.NewSideKeyShare = nil

	dkg.NewEochSponsor = nil
	dkg.NewSideKeyShares = nil
	dkg.NewEochOldShares = nil

	util.LogWithGreen("DKG consensus", "successfully <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<< New Epoch", dkg.Epoch)
	dkg.saveStore()
}

func (dkg *DKG) cancelNewEpoch() {
	dkg.NewNodes = nil
	dkg.NewEoch = 0
	dkg.NewDkgPubKey = nil
	dkg.NewDkgKeyShare = nil
	dkg.NewSideKeyPub = types.AccountID{}
	dkg.NewSideKeyShare = nil

	dkg.NewEochSponsor = nil
	dkg.NewSideKeyShares = nil
	dkg.NewEochOldShares = nil
	dkg.setConsensusFree()
}

func (dkg *DKG) ConsensusIsbusy() bool {
	return time.Now().Unix()-dkg.lastConsensusTime < 90
}

func (dkg *DKG) setConsensusBusy() {
	dkg.lastConsensusTime = time.Now().Unix()
}

func (dkg *DKG) setConsensusFree() {
	dkg.lastConsensusTime = 0
}
