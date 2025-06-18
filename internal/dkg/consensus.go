package dkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.dedis.ch/kyber/v4"
	pedersen "go.dedis.ch/kyber/v4/share/dkg/pedersen"
	"go.dedis.ch/kyber/v4/sign/schnorr"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/util"
)

func (dkg *DKG) TryConsensus(msg model.ConsensusMsg) {
	if !dkg.inConsensus {
		if dkg.DkgKeyShare == nil && msg.Epoch > 1 {
			util.LogError("DKG Consensus", "Node is not old validator, cannot start consensus")
			return
		}

		if msg.Epoch > 1 {
			msg.ShareCommits = *util.DeepCopy(dkg.DkgKeyShare.Commits)
			msg.ConsensusNodeNum = len(dkg.DkgNodes)
			msg.OldValidators = *util.DeepCopy(dkg.DkgNodes)
		} else {
			msg.Epoch = 1
			msg.ShareCommits = model.KyberPoints{Public: []kyber.Point{}}
			msg.ConsensusNodeNum = 0
		}

		bt, _ := json.Marshal(msg)

		dkg.TryRun(&model.Message{
			Type:    "consensus",
			Payload: bt,
		})
	}
}

func (dkg *DKG) startConsensus(msg model.ConsensusMsg) error {
	if dkg.inConsensus {
		// util.LogError("DKG Consensus going +++++++++++++++++++++++++++++++++++++")
		return errors.New("DKG Consensus going")
	}

	dkg.inConsensus = true
	if msg.Epoch <= 1 {
		dkg.addConsensusTimeout()
		dkg.log.Error("StartConsensus", "epoch", msg.Epoch)
		return dkg.initConsensus(msg)
	}

	if dkg.Epoch >= msg.Epoch {
		// util.LogError("DKG Epoch is not need to update -----------------------------")
		return errors.New("DKG Epoch is not need to update")
	}

	dkg.addConsensusTimeout()
	dkg.log.Error("ReConsensus", "epoch", msg.Epoch)
	return dkg.reConsensus(msg)
}

func (dkg *DKG) addConsensusTimeout() {
	if dkg.failConsensusTimer != nil {
		dkg.failConsensusTimer.Stop()
	}
	dkg.failConsensusTimer = time.AfterFunc(time.Second*30, func() {
		dkg.stopConsensus(false)
	})
}

// Init Consensus
func (dkg *DKG) initConsensus(msg model.ConsensusMsg) error {
	// if flag.Lookup("test.v") == nil {
	// 	go dkg.HandleSecretSave()
	// }
	dkg.DkgNodes = msg.Validators
	dkg.NewNodes = msg.Validators
	dkg.Threshold = len(msg.Validators) * 2 / 3

	// 如果已经初始化，则直接返回
	if dkg.status == 1 {
		return nil
	}

	// dkg 节点列表
	nodes := make([]pedersen.Node, 0, len(dkg.DkgNodes))
	for i, p := range dkg.DkgNodes {
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
		return fmt.Errorf("failed to initialize DKG protocol: %w", err)
	}

	// 等待节点连接
	for {
		if dkg.connectLen()+1 < len(dkg.DkgNodes) {
			time.Sleep(time.Second * 10)
			fmt.Println("Number of nodes:", dkg.connectLen(), " len(dkg.DkgNodes) ", len(dkg.DkgNodes))
			fmt.Println("The number of nodes is insufficient, please wait for more nodes to join")
		} else {
			break
		}
	}

	// 获取当前节点的协议
	deal, err := dkg.DistKeyGenerator.Deals()
	if err != nil {
		return fmt.Errorf("failed to generate key shares: %w", err)
	}

	// 开启节点共识
	for _, node := range dkg.DkgNodes {
		err = dkg.sendDealMessage(&node.P2pId, &model.ConsensusMsg{
			DealBundle:   &model.DealBundle{DealBundle: deal},
			ShareCommits: model.KyberPoints{},
			Validators:   dkg.DkgNodes,
			Epoch:        0,
		})
		if err != nil {
			fmt.Println("Send error:", err)
		}
	}

	dkg.NewEoch = 1

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
	dkg.DkgNodes = msg.OldValidators
	dkg.Threshold = len(msg.OldValidators) * 2 / 3
	dkg.NewNodes = msg.Validators

	// new DKG 节点列表
	newNodes := make([]pedersen.Node, 0, len(msg.Validators))
	for i, p := range msg.Validators {
		newNodes = append(newNodes, pedersen.Node{
			Index:  uint32(i),
			Public: p.ValidatorId.Point(),
		})
	}

	// 获取旧节点列表
	oldNodes := make([]pedersen.Node, 0, len(dkg.DkgNodes))
	for i, p := range dkg.DkgNodes {
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
		fmt.Println("unable to create DistKeyGenerator", err.Error())
		return err
	}

	priShare := dkg.DkgKeyShare

	// 重置 DKG Key
	dkg.deals = map[string]*model.DealBundle{}
	dkg.responses = map[string]*pedersen.ResponseBundle{}
	dkg.justifs = []*pedersen.JustificationBundle{}

	// TODO
	dkg.DkgKeyShare = nil
	dkg.NewEoch = msg.Epoch

	// old node issue deals
	if priShare != nil {
		// 获取当前节点的协议
		deal, err := dkg.DistKeyGenerator.Deals()
		if err != nil {
			return fmt.Errorf("failed to generate key shares: %w", err)
		}

		// 开启节点共识
		for _, node := range dkg.NewNodes {
			msg.DealBundle = &model.DealBundle{DealBundle: deal}
			err = dkg.sendDealMessage(&node.P2pId, &msg)
			if err != nil {
				fmt.Println("Send error:", err)
			}
		}
	}

	// // 等待节点完成重组
	// for {
	// 	if dkg.DkgKeyShare.PriShare != nil {
	// 		break
	// 	}
	// 	time.Sleep(time.Second)
	// }

	// if dkg.log != nil {
	// 	dkg.log.Info("DKG uccessfully reshare <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	// }

	// dkg.deals = map[string]*model.DealBundle{}
	// dkg.responses = map[string]*pedersen.ResponseBundle{}
	// dkg.justifs = []*pedersen.JustificationBundle{}
	// dkg.DkgNodes = dkg.NewNodes
	return nil
}

func (dkg *DKG) stopConsensus(isok bool) {
	if dkg.failConsensusTimer != nil {
		dkg.failConsensusTimer.Stop()
	}
	dkg.deals = map[string]*model.DealBundle{}
	dkg.responses = map[string]*pedersen.ResponseBundle{}
	dkg.justifs = []*pedersen.JustificationBundle{}
	if isok {
		dkg.DkgNodes = dkg.NewNodes
		dkg.Epoch = dkg.NewEoch
		dkg.Threshold = len(dkg.NewNodes) * 2 / 3
		dkg.NewNodes = []*model.Validator{}
		util.LogWithGreen("DKG consensus", "successfully <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<", dkg.Epoch)
	} else {
		util.LogWithRed("DKG consensus", "failed <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
	}

	dkg.saveStore()
	dkg.inConsensus = false
}
