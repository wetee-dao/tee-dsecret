package gov

import (
	"errors"
	"math/big"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

var (
	ErrMustCallByGov         = errors.New("must call by gov")
	ErrPublicJoinNotAllowed  = errors.New("public join not allowed")
	ErrMemberExisted         = errors.New("member existed")
	ErrMemberNotExisted      = errors.New("member not existed")
	ErrMemberBalanceNotZero  = errors.New("member balance not zero")
	ErrLowBalance            = errors.New("low balance")
	ErrTransferDisabled      = errors.New("transfer disabled")
	ErrInsufficientAllowance = errors.New("insufficient allowance")
	ErrNoTrack               = errors.New("no track")
	ErrInvalidProposal       = errors.New("invalid proposal")
	ErrInvalidProposalStatus = errors.New("invalid proposal status")
	ErrInvalidProposalCaller = errors.New("invalid proposal caller")
	ErrInvalidDepositTime    = errors.New("invalid deposit time")
	ErrInvalidDeposit        = errors.New("invalid deposit")
	ErrPropNotOngoing        = errors.New("proposal not ongoing")
	ErrInvalidVoteTime       = errors.New("invalid vote time")
	ErrInvalidVote           = errors.New("invalid vote")
	ErrInvalidVoteUser       = errors.New("invalid vote user")
	ErrInvalidVoteStatus     = errors.New("invalid vote status")
	ErrVoteAlreadyUnlocked   = errors.New("vote already unlocked")
	ErrInvalidVoteUnlockTime = errors.New("invalid vote unlock time")
	ErrProposalNotConfirmed  = errors.New("proposal not confirmed")
	ErrProposalInDecision    = errors.New("proposal in decision")
	ErrSpendNotFound         = errors.New("spend not found")
	ErrSpendAlreadyExecuted  = errors.New("spend already executed")
)

const (
	ProposalPending    string = "pending"
	ProposalOngoing    string = "ongoing"
	ProposalConfirming string = "confirming"
	ProposalApproved   string = "approved"
	ProposalRejected   string = "rejected"
	ProposalCanceled   string = "canceled"
)

type Gov struct {
	api                 model.ContractApi
	members             *model.StoreMapping[model.UniAddr, model.Amount] //keyPfx:member_
	publicJoin          *model.StoreValue[bool]                          //key:public_join
	transferEnabled     *model.StoreValue[bool]                          //key:transfer
	totalIssuance       *model.StoreValue[model.Amount]                  //key:issuance
	defaultTrack        *model.StoreValue[uint32]                        //key:default_track
	nextProposalIDStore *model.StoreValue[uint32]                        //key:next_proposal
	nextVoteIDStore     *model.StoreValue[uint64]                        //key:next_vote
	nextSpendIDStore    *model.StoreValue[uint64]                        //key:next_spend
	nextTrackIDStore    *model.StoreValue[uint32]                        //key:next_track
	memberLocks         *model.StoreMapping[model.UniAddr, model.Amount] //keyPfx:member_lock_
	allowances          *model.StoreMapping[model.UniAddr, model.Amount] //keyPfx:allowance_
	tracks              *model.StoreMapping[uint32, TrackData]           //keyPfx:track_
	proposals           *model.StoreMapping[uint32, Proposal]            //keyPfx:proposal_
	votes               *model.StoreMapping[uint64, Vote]                //keyPfx:vote_
	voteUnlocks         *model.StoreMapping[uint64, bool]                //keyPfx:vote_unlock_
	spends              *model.StoreMapping[uint64, Spend]               //keyPfx:spend_
}

type GovQuery struct {
	Gov
}

type GovMutation struct {
	Gov
}

func (d GovMutation) Init() error {
	total := types.NewU256(*big.NewInt(0))
	publicJoin := true
	defaultTrack := TrackData{
		Name:               "default",
		PreparePeriod:      0,
		MaxDeciding:        100,
		ConfirmPeriod:      1,
		DecisionPeriod:     100,
		MinEnactmentPeriod: 0,
		DecisionDeposit:    model.Amount{Int: big.NewInt(1)},
		MaxBalance:         model.Amount{Int: big.NewInt(1_000_000)},
	}

	if err := d.publicJoin.Set(d.api.GetTxn(), publicJoin); err != nil {
		return err
	}

	if err := d.transferEnabled.Set(d.api.GetTxn(), false); err != nil {
		return err
	}

	if err := d.totalIssuance.Set(d.api.GetTxn(), total); err != nil {
		return err
	}

	if err := d.nextProposalIDStore.Set(d.api.GetTxn(), 0); err != nil {
		return err
	}

	if err := d.nextVoteIDStore.Set(d.api.GetTxn(), 0); err != nil {
		return err
	}

	if err := d.nextSpendIDStore.Set(d.api.GetTxn(), 0); err != nil {
		return err
	}

	if err := d.nextTrackIDStore.Set(d.api.GetTxn(), 0); err != nil {
		return err
	}

	if err := d.tracks.Set(d.api.GetTxn(), 0, defaultTrack); err != nil {
		return err
	}

	if err := d.defaultTrack.Set(d.api.GetTxn(), 0); err != nil {
		return err
	}

	if err := d.nextTrackIDStore.Set(d.api.GetTxn(), 1); err != nil {
		return err
	}

	return nil
}
