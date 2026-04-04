package dao

import "errors"

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

