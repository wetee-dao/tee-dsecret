package dao

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"sort"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func (d DAO) ensureGov() error {
	sudo, err := d.sudoAccount.Get(d.api.GetTxn(), singletonKey)
	if err != nil {
		return err
	}
	caller := d.api.GetCaller()
	if sudo == nil || !bytes.Equal(*sudo, caller) {
		return ErrMustCallByGov
	}
	return nil
}

func (d DAO) member(account []byte) (*Member, error) {
	return d.members.Get(d.api.GetTxn(), account)
}

func (d DAO) proposal(id uint32) (Proposal, error) {
	prop, err := d.proposals.Get(d.api.GetTxn(), id)
	if err != nil {
		return Proposal{}, err
	}
	if prop == nil {
		return Proposal{}, ErrInvalidProposal
	}
	return *prop, nil
}

func (d DAO) vote(id uint64) (Vote, error) {
	v, err := d.votes.Get(d.api.GetTxn(), id)
	if err != nil {
		return Vote{}, err
	}
	if v == nil {
		return Vote{}, ErrInvalidVote
	}
	return *v, nil
}

func (d DAO) track(id uint32) (TrackData, error) {
	t, err := d.tracks.Get(d.api.GetTxn(), id)
	if err != nil {
		return TrackData{}, err
	}
	if t == nil {
		return TrackData{}, ErrNoTrack
	}
	return *t, nil
}

func (d DAO) lockOf(account []byte) ([]byte, error) {
	lock, err := d.memberLocks.Get(d.api.GetTxn(), account)
	if err != nil {
		return nil, err
	}
	if lock == nil {
		return nil, nil
	}
	return cloneBytes(*lock), nil
}

func (d DAO) allowanceOf(owner, spender []byte) ([]byte, error) {
	val, err := d.allowances.Get(d.api.GetTxn(), allowanceKey(owner, spender))
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return cloneBytes(*val), nil
}

func (d DAO) bytesValue(store *model.StoreMapping[string, []byte]) ([]byte, error) {
	val, err := store.Get(d.api.GetTxn(), singletonKey)
	if err != nil {
		return nil, err
	}
	if val == nil {
		return nil, nil
	}
	return cloneBytes(*val), nil
}

func (d DAO) boolValue(store *model.StoreMapping[string, bool], fallback bool) (bool, error) {
	val, err := store.Get(d.api.GetTxn(), singletonKey)
	if err != nil {
		return false, err
	}
	if val == nil {
		return fallback, nil
	}
	return *val, nil
}

func (d DAO) nextProposalID() (uint32, error) {
	v, err := d.nextProposalIDStore.Get(d.api.GetTxn(), singletonKey)
	if err != nil {
		return 0, err
	}
	if v == nil {
		return 0, nil
	}
	return *v, nil
}

func (d DAO) nextVoteID() (uint64, error) {
	v, err := d.nextVoteIDStore.Get(d.api.GetTxn(), singletonKey)
	if err != nil {
		return 0, err
	}
	if v == nil {
		return 0, nil
	}
	return *v, nil
}

func (d DAO) nextSpendID() (uint64, error) {
	v, err := d.nextSpendIDStore.Get(d.api.GetTxn(), singletonKey)
	if err != nil {
		return 0, err
	}
	if v == nil {
		return 0, nil
	}
	return *v, nil
}

func (d DAO) nextTrackID() (uint32, error) {
	v, err := d.nextTrackIDStore.Get(d.api.GetTxn(), singletonKey)
	if err != nil {
		return 0, err
	}
	if v == nil {
		return 0, nil
	}
	return *v, nil
}

func (d DAO) calculateProposalStatus(prop Proposal) (bool, int64, TrackData, error) {
	track, err := d.track(prop.TrackID)
	if err != nil {
		return false, 0, TrackData{}, err
	}
	end := prop.Deposit.Block + int64(track.MaxDeciding)
	votes, err := d.votes.List(d.api.GetTxn())
	if err != nil {
		return false, 0, TrackData{}, err
	}

	totalSupply, err := d.bytesValue(d.totalIssuance)
	if err != nil {
		return false, 0, TrackData{}, err
	}

	yes := big.NewInt(0)
	support := big.NewInt(0)
	all := decodeAmount(totalSupply)

	sort.Slice(votes, func(i, j int) bool {
		return votes[i].VoteBlock < votes[j].VoteBlock
	})

	for _, vote := range votes {
		if vote.ProposalID != prop.ID || vote.Deleted {
			continue
		}
		pledge := decodeAmount(vote.Pledge)
		support.Add(support, pledge)
		if vote.OpinionYes {
			yes.Add(yes, pledge)
		}
	}

	if all.Sign() == 0 {
		return false, end, track, nil
	}
	if yes.Sign() > 0 && support.Cmp(decodeAmount(track.DecisionDeposit)) >= 0 {
		return true, end, track, nil
	}
	return false, end, track, nil
}

func (d DAO) calculateProposalEndBlock(id uint32) (int64, error) {
	height := d.api.GetHeight()
	prop, err := d.proposal(id)
	if err != nil {
		return 0, err
	}
	switch prop.Status.State {
	case ProposalRejected, ProposalApproved:
		return prop.Status.Block, nil
	case ProposalOngoing:
		confirmed, end, _, err := d.calculateProposalStatus(prop)
		if err != nil {
			return 0, err
		}
		if !confirmed && height > end {
			return end, nil
		}
	}
	return 0, ErrInvalidProposalStatus
}

func allowanceKey(owner, spender []byte) string {
	return hex.EncodeToString(owner) + ":" + hex.EncodeToString(spender)
}

func cloneCall(in *CallContent) *CallContent {
	if in == nil {
		return nil
	}
	return &CallContent{
		Contract:     cloneBytes(in.Contract),
		Selector:     cloneBytes(in.Selector),
		Input:        cloneBytes(in.Input),
		Amount:       cloneBytes(in.Amount),
		RefTimeLimit: in.RefTimeLimit,
		AllowReentry: in.AllowReentry,
	}
}

func cloneBytes(in []byte) []byte {
	if len(in) == 0 {
		return nil
	}
	out := make([]byte, len(in))
	copy(out, in)
	return out
}

func decodeAmount(b []byte) *big.Int {
	if len(b) == 0 {
		return big.NewInt(0)
	}
	return new(big.Int).SetBytes(b)
}

func encodeAmount(v *big.Int) []byte {
	if v == nil || v.Sign() == 0 {
		return nil
	}
	return v.Bytes()
}

func add(a, b []byte) []byte {
	x := decodeAmount(a)
	x.Add(x, decodeAmount(b))
	return encodeAmount(x)
}

func sub(a, b []byte) []byte {
	x := decodeAmount(a)
	x.Sub(x, decodeAmount(b))
	if x.Sign() < 0 {
		return nil
	}
	return encodeAmount(x)
}

func cmp(a, b []byte) int {
	return decodeAmount(a).Cmp(decodeAmount(b))
}

func isZero(a []byte) bool {
	return decodeAmount(a).Sign() == 0
}

func mustBytes(v []byte, err error) []byte {
	if err != nil {
		return nil
	}
	return v
}

func encodeUint64(v uint64) []byte {
	if v == 0 {
		return []byte{0}
	}
	out := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		out[i] = byte(v)
		v >>= 8
	}
	return out
}
