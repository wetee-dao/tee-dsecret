package sidechain

import (
	"context"
	"crypto"
	"encoding/json"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	cryptoencoding "github.com/cometbft/cometbft/crypto/encoding"
	"wetee.app/dsecret/internal/dkg"
	"wetee.app/dsecret/internal/model"
	"wetee.app/dsecret/internal/util"

	"github.com/cometbft/cometbft/version"
)

const ApplicationVersion = 1
const CurseWordsLimitVE = 10

type SideChain struct {
	abci.BaseApplication
	valAddrToPubKeyMap map[string]crypto.PublicKey

	dkg          *dkg.DKG
	state        AppState
	onGoingBlock *model.Txn
}

func NewSideChain() (*SideChain, error) {
	state, err := loadState()
	if err != nil {
		return nil, err
	}
	return &SideChain{
		state:              state,
		valAddrToPubKeyMap: make(map[string]crypto.PublicKey),
	}, nil
}

// Info return application information
func (app *SideChain) Info(_ context.Context, info *abci.InfoRequest) (*abci.InfoResponse, error) {
	if len(app.valAddrToPubKeyMap) == 0 && app.state.Height > 0 {
		validators, err := app.GetValidators()
		if err != nil {
			return nil, err
		}
		for _, v := range validators {
			pubKey, err := cryptoencoding.PubKeyFromTypeAndBytes(v.PubKeyType, v.PubKeyBytes)
			if err != nil {
				return nil, fmt.Errorf("can't decode public key: %w", err)
			}

			app.valAddrToPubKeyMap[string(pubKey.Address())] = pubKey
		}
	}

	return &abci.InfoResponse{
		Version:         version.ABCIVersion,
		AppVersion:      ApplicationVersion,
		LastBlockHeight: app.state.Height,

		LastBlockAppHash: app.state.Hash(),
	}, nil
}

// Query the application state for specific information
func (app *SideChain) Query(ctx context.Context, query *abci.QueryRequest) (*abci.QueryResponse, error) {
	util.LogWithPurple("SideChain Query")

	resp := abci.QueryResponse{Key: query.Data}

	// Retrieve all message sent by the sender
	messages := map[string]string{}

	resultBytes, err := json.Marshal(messages)
	if err != nil {
		return nil, err
	}

	resp.Log = string(resultBytes)
	resp.Value = resultBytes

	return &resp, nil
}

// CheckTx handles validation of inbound transactions. If a transaction is not a valid message, or if a user
// does not exist in the database or if a user is banned it returns an error
func (app *SideChain) CheckTx(_ context.Context, req *abci.CheckTxRequest) (*abci.CheckTxResponse, error) {
	util.LogWithPurple("SideChain", "CheckTx")

	// check req.Tx

	return &abci.CheckTxResponse{Code: CodeTypeOK}, nil
}

// Consensus Connection
// InitChain initializes the blockchain with information sent from CometBFT such as validators or consensus parameters
func (app *SideChain) InitChain(_ context.Context, req *abci.InitChainRequest) (*abci.InitChainResponse, error) {
	util.LogWithPurple("SideChain", "InitChain")
	for _, v := range req.Validators {
		err := app.updateValidator(v)
		if err != nil {
			return nil, err
		}
	}
	appHash := app.state.Hash()

	// This parameter can also be set in the genesis file
	req.ConsensusParams.Feature.VoteExtensionsEnableHeight.Value = 1
	return &abci.InitChainResponse{ConsensusParams: req.ConsensusParams, AppHash: appHash}, nil
}

// PrepareProposal is used to prepare a proposal for the next block in the blockchain. The application can re-order, remove
// or add transactions
func (app *SideChain) PrepareProposal(_ context.Context, req *abci.PrepareProposalRequest) (*abci.PrepareProposalResponse, error) {
	util.LogWithPurple("SideChain", "PrepareProposal")

	app.CheckEpoch()

	finalProposal := make([][]byte, 0)
	for _, tx := range req.Txs {
		finalProposal = append(finalProposal, tx)
	}

	return &abci.PrepareProposalResponse{Txs: finalProposal}, nil
}

// ProcessProposal validates the proposed block and the transactions and return a status if it was accepted or rejected
func (app *SideChain) ProcessProposal(_ context.Context, req *abci.ProcessProposalRequest) (*abci.ProcessProposalResponse, error) {
	util.LogWithPurple("SideChain", "ProcessProposal")

	// for i, tx := range req.Txs {
	// }

	return &abci.ProcessProposalResponse{Status: abci.PROCESS_PROPOSAL_STATUS_ACCEPT}, nil
}

// FinalizeBlock Deliver the decided block to the Application
func (app *SideChain) FinalizeBlock(_ context.Context, req *abci.FinalizeBlockRequest) (*abci.FinalizeBlockResponse, error) {
	util.LogWithPurple("SideChain", "FinalizeBlock")

	// Iterate over Tx in current block
	app.onGoingBlock = model.DBINS.NewTransaction(true)

	respTxs := make([]*abci.ExecTxResult, len(req.Txs))

	for i, _ := range req.Txs {
		respTxs[i] = &abci.ExecTxResult{Code: abci.CodeTypeOK}
	}

	app.state.Height = req.Height
	response := &abci.FinalizeBlockResponse{
		TxResults: respTxs,
		AppHash:   app.state.Hash(),
		// ValidatorUpdates: []abci.ValidatorUpdate{},
	}

	return response, nil
}

// Commit the application state
func (app *SideChain) Commit(_ context.Context, _ *abci.CommitRequest) (*abci.CommitResponse, error) {
	if err := app.onGoingBlock.Commit(); err != nil {
		return nil, err
	}

	err := saveState(&app.state)
	if err != nil {
		return nil, err
	}

	util.LogWithPurple("SideChain", "Commit")
	util.LogWithGreen("---------------------------------------------------------------")

	return &abci.CommitResponse{}, nil
}

// ExtendVote returns curse words as vote extensions
func (app *SideChain) ExtendVote(_ context.Context, _ *abci.ExtendVoteRequest) (*abci.ExtendVoteResponse, error) {
	util.LogWithPurple("SideChain", "ExtendVote")

	return &abci.ExtendVoteResponse{VoteExtension: []byte("")}, nil
}

// VerifyVoteExtension verifies the vote extensions and ensure they include the curse words
// It will not be called for extensions generated by this validator
func (app *SideChain) VerifyVoteExtension(_ context.Context, req *abci.VerifyVoteExtensionRequest) (*abci.VerifyVoteExtensionResponse, error) {
	util.LogWithPurple("SideChain", "VerifyVoteExtension")

	// if len(curseWords) > CurseWordsLimitVE {
	// 	return &abci.VerifyVoteExtensionResponse{Status: abci.VERIFY_VOTE_EXTENSION_STATUS_REJECT}, nil
	// }

	return &abci.VerifyVoteExtensionResponse{Status: abci.VERIFY_VOTE_EXTENSION_STATUS_ACCEPT}, nil
}
