package dao

import (
	"errors"
	"math/big"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const singletonKey = "state"

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

type ProposalState string

const (
	ProposalPending    ProposalState = "pending"
	ProposalOngoing    ProposalState = "ongoing"
	ProposalConfirming ProposalState = "confirming"
	ProposalApproved   ProposalState = "approved"
	ProposalRejected   ProposalState = "rejected"
	ProposalCanceled   ProposalState = "canceled"
)

type DAO struct {
	api model.ContractApi

	members             *model.StoreMapping[[]byte, Member]    //keyPfx:member_
	publicJoin          *model.StoreMapping[string, bool]      //keyPfx:public_join_
	sudoAccount         *model.StoreMapping[string, []byte]    //keyPfx:sudo_
	transferEnabled     *model.StoreMapping[string, bool]      //keyPfx:transfer_
	totalIssuance       *model.StoreMapping[string, []byte]    //keyPfx:issuance_
	defaultTrack        *model.StoreMapping[string, uint32]    //keyPfx:default_track_
	nextProposalIDStore *model.StoreMapping[string, uint32]    //keyPfx:next_proposal_
	nextVoteIDStore     *model.StoreMapping[string, uint64]    //keyPfx:next_vote_
	nextSpendIDStore    *model.StoreMapping[string, uint64]    //keyPfx:next_spend_
	nextTrackIDStore    *model.StoreMapping[string, uint32]    //keyPfx:next_track_
	memberLocks         *model.StoreMapping[[]byte, []byte]    //keyPfx:member_lock_
	allowances          *model.StoreMapping[string, []byte]    //keyPfx:allowance_
	tracks              *model.StoreMapping[uint32, TrackData] //keyPfx:track_
	proposals           *model.StoreMapping[uint32, Proposal]  //keyPfx:proposal_
	votes               *model.StoreMapping[uint64, Vote]      //keyPfx:vote_
	voteUnlocks         *model.StoreMapping[uint64, bool]      //keyPfx:vote_unlock_
	spends              *model.StoreMapping[uint64, Spend]     //keyPfx:spend_
}

type DaoQuery struct {
	DAO
}

type DaoMutation struct {
	DAO
}

func (d DaoMutation) Init(initialMembers []Member, publicJoin bool, sudoAccount []byte, defaultTrack *TrackData) error {

	total := big.NewInt(0)
	for i := range initialMembers {
		member := &initialMembers[i]
		if len(member.Account) == 0 {
			continue
		}
		if err := d.members.Set(d.api.GetTxn(), member.Account, *member); err != nil {
			return err
		}
		total.Add(total, decodeAmount(member.Balance))
	}

	if err := d.publicJoin.Set(d.api.GetTxn(), singletonKey, publicJoin); err != nil {
		return err
	}

	if len(sudoAccount) > 0 {
		if err := d.sudoAccount.Set(d.api.GetTxn(), singletonKey, cloneBytes(sudoAccount)); err != nil {
			return err
		}
	}

	if err := d.transferEnabled.Set(d.api.GetTxn(), singletonKey, false); err != nil {
		return err
	}

	if err := d.totalIssuance.Set(d.api.GetTxn(), singletonKey, encodeAmount(total)); err != nil {
		return err
	}

	if err := d.nextProposalIDStore.Set(d.api.GetTxn(), singletonKey, 0); err != nil {
		return err
	}

	if err := d.nextVoteIDStore.Set(d.api.GetTxn(), singletonKey, 0); err != nil {
		return err
	}

	if err := d.nextSpendIDStore.Set(d.api.GetTxn(), singletonKey, 0); err != nil {
		return err
	}

	if err := d.nextTrackIDStore.Set(d.api.GetTxn(), singletonKey, 0); err != nil {
		return err
	}

	if defaultTrack != nil {
		if err := d.tracks.Set(d.api.GetTxn(), 0, *defaultTrack); err != nil {
			return err
		}
		if err := d.defaultTrack.Set(d.api.GetTxn(), singletonKey, 0); err != nil {
			return err
		}
		if err := d.nextTrackIDStore.Set(d.api.GetTxn(), singletonKey, 1); err != nil {
			return err
		}
	}

	return nil
}
