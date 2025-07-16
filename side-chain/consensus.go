package sidechain

import (
	"context"
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/version"

	"github.com/wetee-dao/tee-dsecret/pkg/dkg"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	bftbrigde "github.com/wetee-dao/tee-dsecret/pkg/network/bft-brigde"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

const ApplicationVersion = 1

type SideChain struct {
	abci.BaseApplication

	dkg *dkg.DKG
	p2p *bftbrigde.BTFReactor

	txCh *model.PersistChan[*model.BlockPartialSign]

	state               AppState
	onGoingBlock        *model.Txn
	onGoingValidators   []abci.ValidatorUpdate
	currProposerAddress []byte
}

func NewSideChain(light bool) (*SideChain, error) {
	state, err := loadAppState()
	if err != nil {
		return nil, err
	}

	c := &SideChain{
		state: state,
	}

	if !light {
		txCh, err := model.NewPersistChan[*model.BlockPartialSign]("back_tx", 1000)
		if err != nil {
			return nil, err
		}
		c.txCh = txCh
	}

	return c, nil
}

func (app *SideChain) Info(_ context.Context, info *abci.InfoRequest) (*abci.InfoResponse, error) {
	return &abci.InfoResponse{
		Version:          version.ABCIVersion,
		AppVersion:       ApplicationVersion,
		LastBlockHeight:  app.state.Height,
		LastBlockAppHash: app.state.Hash(),
	}, nil
}

func (app *SideChain) Query(ctx context.Context, query *abci.QueryRequest) (*abci.QueryResponse, error) {
	util.LogWithGreen("Query")
	resp := abci.QueryResponse{Key: query.Data}

	return &resp, nil
}

func (app *SideChain) InitChain(_ context.Context, req *abci.InitChainRequest) (*abci.InitChainResponse, error) {
	util.LogWithGreen("InitChain")
	app.initValidators(req.Validators)
	appHash := app.state.Hash()

	// This parameter can also be set in the genesis file
	req.ConsensusParams.Feature.VoteExtensionsEnableHeight.Value = 1
	return &abci.InitChainResponse{ConsensusParams: req.ConsensusParams, AppHash: appHash}, nil
}

func (app *SideChain) CheckTx(_ context.Context, req *abci.CheckTxRequest) (*abci.CheckTxResponse, error) {
	fmt.Println()
	util.LogWithGreen("START BLOCK", "--------------------------------------------------------------")
	LogWithTime("🚀 CheckTx")

	return &abci.CheckTxResponse{Code: app.checkTx(req.Tx)}, nil
}

func (app *SideChain) PrepareProposal(_ context.Context, req *abci.PrepareProposalRequest) (*abci.PrepareProposalResponse, error) {
	LogWithTime("🎁 PrepareProposal")

	// Check if the current epoch is valid
	epochTx := app.CheckEpochFromValidator()
	finalProposal := make([][]byte, 0, len(req.Txs)+2)
	if len(epochTx) > 0 {
		finalProposal = append(finalProposal, epochTx)
	}

	epochStatus := app.GetEpochStatus()
	// Check if it is in the epoch transition phase
	if len(epochTx) == 0 && time.Now().Unix()-int64(epochStatus) > 120 {
		app.PrepareTx(req.Txs, &finalProposal, true)
	} else {
		app.PrepareTx(req.Txs, &finalProposal, false)
	}

	return &abci.PrepareProposalResponse{Txs: finalProposal}, nil
}

func (app *SideChain) ProcessProposal(_ context.Context, req *abci.ProcessProposalRequest) (*abci.ProcessProposalResponse, error) {
	LogWithTime("🌈 ProcessProposal")

	status := app.ProcessTx(req.Txs, app.onGoingBlock)
	return &abci.ProcessProposalResponse{Status: status}, nil
}

func (app *SideChain) FinalizeBlock(_ context.Context, req *abci.FinalizeBlockRequest) (*abci.FinalizeBlockResponse, error) {
	// Iterate over Tx in current block
	app.onGoingBlock = model.DBINS.NewTransaction()
	respTxs, err := app.FinalizeTx(req.Txs, app.onGoingBlock, req.Height, req.ProposerAddress)
	if err != nil {
		return nil, err
	}

	// Sync validator updates to consensus
	var validatorUpdates []abci.ValidatorUpdate
	if app.onGoingValidators != nil {
		validatorUpdates = app.onGoingValidators
		ss58 := []string{}
		for _, v := range app.onGoingValidators {
			ss58 = append(ss58, model.PubKeyFromByte(v.PubKeyBytes).SS58())
		}
		util.LogWithPurple("Validator updates", ss58)
	}

	// save proposer of currut block
	app.currProposerAddress = req.ProposerAddress

	app.state.Height = req.Height
	response := &abci.FinalizeBlockResponse{
		TxResults:        respTxs,
		AppHash:          app.state.Hash(),
		ValidatorUpdates: validatorUpdates,
	}

	LogWithTime("📦 Finalize Block =>", util.Green+" "+fmt.Sprint(req.Height)+" "+util.Reset)
	return response, nil
}

// Commit the application state
func (app *SideChain) Commit(_ context.Context, _ *abci.CommitRequest) (*abci.CommitResponse, error) {
	if err := app.onGoingBlock.Commit(); err != nil {
		return nil, err
	}

	app.onGoingValidators = nil
	err := saveAppState(&app.state)
	if err != nil {
		return nil, err
	}

	LogWithTime("💤 Commit")
	util.LogWithGreen("END BLOCK  ", "--------------------------------------------------------------")

	return &abci.CommitResponse{}, nil
}

func (app *SideChain) ExtendVote(_ context.Context, _ *abci.ExtendVoteRequest) (*abci.ExtendVoteResponse, error) {
	LogWithTime("💊 Issue TEE report")

	return &abci.ExtendVoteResponse{VoteExtension: []byte("")}, nil
}

func (app *SideChain) VerifyVoteExtension(_ context.Context, req *abci.VerifyVoteExtensionRequest) (*abci.VerifyVoteExtensionResponse, error) {
	LogWithTime("💊 Verify TEE report")

	// if len(curseWords) > CurseWordsLimitVE {
	// 	return &abci.VerifyVoteExtensionResponse{Status: abci.VERIFY_VOTE_EXTENSION_STATUS_REJECT}, nil
	// }

	return &abci.VerifyVoteExtensionResponse{Status: abci.VERIFY_VOTE_EXTENSION_STATUS_ACCEPT}, nil
}
