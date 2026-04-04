package dao

import (
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

type DAO struct {
	api model.ContractApi

	members             *model.StoreMapping[[]byte, Member]
	publicJoin          *model.StoreMapping[string, bool]
	sudoAccount         *model.StoreMapping[string, []byte]
	transferEnabled     *model.StoreMapping[string, bool]
	totalIssuance       *model.StoreMapping[string, []byte]
	defaultTrack        *model.StoreMapping[string, uint32]
	nextProposalIDStore *model.StoreMapping[string, uint32]
	nextVoteIDStore     *model.StoreMapping[string, uint64]
	nextSpendIDStore    *model.StoreMapping[string, uint64]
	nextTrackIDStore    *model.StoreMapping[string, uint32]
	memberLocks         *model.StoreMapping[[]byte, []byte]
	allowances          *model.StoreMapping[string, []byte]
	tracks              *model.StoreMapping[uint32, TrackData]
	proposals           *model.StoreMapping[uint32, Proposal]
	votes               *model.StoreMapping[uint64, Vote]
	voteUnlocks         *model.StoreMapping[uint64, bool]
	spends              *model.StoreMapping[uint64, Spend]
}

func NewDAO(api model.ContractApi) *DAO {
	return &DAO{
		api:                 api,
		members:             &model.StoreMapping[[]byte, Member]{Namespace: "dao", KeyPrefix: "member_"},
		publicJoin:          &model.StoreMapping[string, bool]{Namespace: "dao", KeyPrefix: "public_join_"},
		sudoAccount:         &model.StoreMapping[string, []byte]{Namespace: "dao", KeyPrefix: "sudo_"},
		transferEnabled:     &model.StoreMapping[string, bool]{Namespace: "dao", KeyPrefix: "transfer_"},
		totalIssuance:       &model.StoreMapping[string, []byte]{Namespace: "dao", KeyPrefix: "issuance_"},
		defaultTrack:        &model.StoreMapping[string, uint32]{Namespace: "dao", KeyPrefix: "default_track_"},
		nextProposalIDStore: &model.StoreMapping[string, uint32]{Namespace: "dao", KeyPrefix: "next_proposal_"},
		nextVoteIDStore:     &model.StoreMapping[string, uint64]{Namespace: "dao", KeyPrefix: "next_vote_"},
		nextSpendIDStore:    &model.StoreMapping[string, uint64]{Namespace: "dao", KeyPrefix: "next_spend_"},
		nextTrackIDStore:    &model.StoreMapping[string, uint32]{Namespace: "dao", KeyPrefix: "next_track_"},
		memberLocks:         &model.StoreMapping[[]byte, []byte]{Namespace: "dao", KeyPrefix: "member_lock_"},
		allowances:          &model.StoreMapping[string, []byte]{Namespace: "dao", KeyPrefix: "allowance_"},
		tracks:              &model.StoreMapping[uint32, TrackData]{Namespace: "dao", KeyPrefix: "track_"},
		proposals:           &model.StoreMapping[uint32, Proposal]{Namespace: "dao", KeyPrefix: "proposal_"},
		votes:               &model.StoreMapping[uint64, Vote]{Namespace: "dao", KeyPrefix: "vote_"},
		voteUnlocks:         &model.StoreMapping[uint64, bool]{Namespace: "dao", KeyPrefix: "vote_unlock_"},
		spends:              &model.StoreMapping[uint64, Spend]{Namespace: "dao", KeyPrefix: "spend_"},
	}
}

type DaoQuery struct {
	DAO
}

type DaoMutation struct {
	DAO
}
