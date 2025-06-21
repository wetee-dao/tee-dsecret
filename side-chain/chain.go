package sidechain

import (
	"context"

	abci "github.com/cometbft/cometbft/abci/types"
	"wetee.app/dsecret/chains"
	"wetee.app/dsecret/internal/dkg"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/util"

	"github.com/cometbft/cometbft/version"
)

const ApplicationVersion = 1
const CurseWordsLimitVE = 10

type SideChain struct {
	abci.BaseApplication

	dkg                 *dkg.DKG
	State               AppState
	onGoingBlock        *model.Txn
	onGoingValidators   []abci.ValidatorUpdate
	currProposerAddress []byte
}

func NewSideChain() (*SideChain, error) {
	state, err := loadAppState()
	if err != nil {
		return nil, err
	}

	return &SideChain{
		State: state,
	}, nil
}

// Info return application information
func (app *SideChain) Info(_ context.Context, info *abci.InfoRequest) (*abci.InfoResponse, error) {
	return &abci.InfoResponse{
		Version:          version.ABCIVersion,
		AppVersion:       ApplicationVersion,
		LastBlockHeight:  app.State.Height,
		LastBlockAppHash: app.State.Hash(),
	}, nil
}

func (app *SideChain) Query(ctx context.Context, query *abci.QueryRequest) (*abci.QueryResponse, error) {
	util.LogWithPurple("SideChain Query")
	resp := abci.QueryResponse{Key: query.Data}

	return &resp, nil
}

func (app *SideChain) CheckTx(_ context.Context, req *abci.CheckTxRequest) (*abci.CheckTxResponse, error) {
	util.LogWithPurple("SideChain", "CheckTx")

	// check req.Tx

	return &abci.CheckTxResponse{Code: CodeTypeOK}, nil
}

func (app *SideChain) InitChain(_ context.Context, req *abci.InitChainRequest) (*abci.InitChainResponse, error) {
	util.LogWithPurple("SideChain", "InitChain")
	app.saveValidators(req.Validators)
	appHash := app.State.Hash()

	// This parameter can also be set in the genesis file
	req.ConsensusParams.Feature.VoteExtensionsEnableHeight.Value = 1
	return &abci.InitChainResponse{ConsensusParams: req.ConsensusParams, AppHash: appHash}, nil
}

func (app *SideChain) PrepareProposal(_ context.Context, req *abci.PrepareProposalRequest) (*abci.PrepareProposalResponse, error) {
	util.LogWithPurple("SideChain", "PrepareProposal")

	app.CheckEpoch()

	finalProposal := make([][]byte, 0)
	for _, tx := range req.Txs {
		finalProposal = append(finalProposal, tx)
	}

	return &abci.PrepareProposalResponse{Txs: finalProposal}, nil
}

func (app *SideChain) ProcessProposal(_ context.Context, req *abci.ProcessProposalRequest) (*abci.ProcessProposalResponse, error) {
	util.LogWithPurple("SideChain", "ProcessProposal")

	status := app.ProcessTx(req.Txs, app.onGoingBlock)
	return &abci.ProcessProposalResponse{Status: status}, nil
}

func (app *SideChain) FinalizeBlock(_ context.Context, req *abci.FinalizeBlockRequest) (*abci.FinalizeBlockResponse, error) {
	util.LogWithPurple("SideChain", "FinalizeBlock")

	// Iterate over Tx in current block
	app.onGoingBlock = model.DBINS.NewTransaction()
	respTxs, err := app.FinalizeTx(req.Txs, app.onGoingBlock)
	if err != nil {
		return nil, err
	}

	// Sync validator updates to consensus
	var validatorUpdates []abci.ValidatorUpdate
	if app.onGoingValidators != nil {
		validatorUpdates = app.onGoingValidators
	}

	// save proposer of currut block
	app.currProposerAddress = req.ProposerAddress

	app.State.Height = req.Height
	response := &abci.FinalizeBlockResponse{
		TxResults:        respTxs,
		AppHash:          app.State.Hash(),
		ValidatorUpdates: validatorUpdates,
	}

	return response, nil
}

// Commit the application state
func (app *SideChain) Commit(_ context.Context, _ *abci.CommitRequest) (*abci.CommitResponse, error) {
	if err := app.onGoingBlock.Commit(); err != nil {
		return nil, err
	}

	app.onGoingValidators = nil
	err := saveAppState(&app.State)
	if err != nil {
		return nil, err
	}

	if app.currProposerAddress != nil {
		pub := model.PubKeyFromByte(app.currProposerAddress)
		if pub.SS58() == chains.MainChain.GetSignerAddress() {
			//TODO submit main chain transaction
		}
	}

	util.LogWithPurple("SideChain", "Commit")
	util.LogWithGreen("---------------------------------------------------------------")

	return &abci.CommitResponse{}, nil
}

func (app *SideChain) ExtendVote(_ context.Context, _ *abci.ExtendVoteRequest) (*abci.ExtendVoteResponse, error) {
	util.LogWithPurple("SideChain", "ExtendVote")

	return &abci.ExtendVoteResponse{VoteExtension: []byte("")}, nil
}

func (app *SideChain) VerifyVoteExtension(_ context.Context, req *abci.VerifyVoteExtensionRequest) (*abci.VerifyVoteExtensionResponse, error) {
	util.LogWithPurple("SideChain", "VerifyVoteExtension")

	// if len(curseWords) > CurseWordsLimitVE {
	// 	return &abci.VerifyVoteExtensionResponse{Status: abci.VERIFY_VOTE_EXTENSION_STATUS_REJECT}, nil
	// }

	return &abci.VerifyVoteExtensionResponse{Status: abci.VERIFY_VOTE_EXTENSION_STATUS_ACCEPT}, nil
}
