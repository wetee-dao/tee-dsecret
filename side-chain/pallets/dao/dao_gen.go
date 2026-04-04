package dao

import (
	"errors"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types/codec"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/side-chain/pallets/base"
)

// ExecCall 解析 Tx.contract，再按 Method 与字符串参数列表调用对应变更。
func (d DaoMutation) ExecCall(call *model.ContractCall) error {
	if len(d.api.GetCaller()) == 0 {
		return errors.New("dao: missing caller")
	}

	dao := NewDAO(d.api)
	m := DaoMutation{DAO: *dao}

	method := call.Method
	args := call.Args
	switch method {
	case "init":
		if err := base.RequireArgLen(args, 3, method); err != nil {
			return err
		}
		members, err := base.DecodeScaleArgBytes[[]Member](args[0])
		if err != nil {
			return fmt.Errorf("init: members: %w", err)
		}
		pub, err := base.DecodeScaleArgBytes[bool](args[1])
		if err != nil {
			return fmt.Errorf("init: publicJoin: %w", err)
		}
		sudo, err := base.DecodeScaleArgBytes[[]byte](args[2])
		if err != nil {
			return fmt.Errorf("init: sudoAccount: %w", err)
		}
		var dt *TrackData
		if len(args) >= 4 {
			t, err := base.DecodeScaleArgBytes[TrackData](args[3])
			if err != nil {
				return fmt.Errorf("init: defaultTrack: %w", err)
			}
			dt = &t
		}
		return m.Init(members, pub, sudo, dt)
	case "public_join":
		return m.PublicJoin()
	case "join":
		if err := base.RequireArgLen(args, 2, method); err != nil {
			return err
		}
		newUser, err := base.DecodeScaleArgBytes[[]byte](args[0])
		if err != nil {
			return fmt.Errorf("join: newUser: %w", err)
		}
		bal, err := base.DecodeScaleArgBytes[[]byte](args[1])
		if err != nil {
			return fmt.Errorf("join: balance: %w", err)
		}
		return m.Join(newUser, bal)
	case "leave":
		return m.Leave()
	case "leave_with_burn":
		return m.LeaveWithBurn()
	case "submit_proposal":
		if err := base.RequireArgLen(args, 2, method); err != nil {
			return err
		}
		call, err := base.DecodeScaleArgBytes[CallContent](args[0])
		if err != nil {
			return fmt.Errorf("submit_proposal: call: %w", err)
		}
		tid, err := base.DecodeScaleArgBytes[uint32](args[1])
		if err != nil {
			return fmt.Errorf("submit_proposal: trackId: %w", err)
		}
		return m.SubmitProposal(call, tid)
	case "deposit_proposal":
		if err := base.RequireArgLen(args, 2, method); err != nil {
			return err
		}
		pid, err := base.DecodeScaleArgBytes[uint32](args[0])
		if err != nil {
			return fmt.Errorf("deposit_proposal: proposalId: %w", err)
		}
		amt, err := base.DecodeScaleArgBytes[[]byte](args[1])
		if err != nil {
			return fmt.Errorf("deposit_proposal: amount: %w", err)
		}
		return m.DepositProposal(pid, amt)
	case "submit_vote":
		if err := base.RequireArgLen(args, 3, method); err != nil {
			return err
		}
		pid, err := base.DecodeScaleArgBytes[uint32](args[0])
		if err != nil {
			return fmt.Errorf("submit_vote: proposalId: %w", err)
		}
		yes, err := base.DecodeScaleArgBytes[bool](args[1])
		if err != nil {
			return fmt.Errorf("submit_vote: opinionYes: %w", err)
		}
		lock, err := base.DecodeScaleArgBytes[[]byte](args[2])
		if err != nil {
			return fmt.Errorf("submit_vote: lockAmount: %w", err)
		}
		return m.SubmitVote(pid, yes, lock)
	case "cancel_vote":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return err
		}
		vid, err := base.DecodeScaleArgBytes[uint64](args[0])
		if err != nil {
			return fmt.Errorf("cancel_vote: voteId: %w", err)
		}
		return m.CancelVote(vid)
	case "unlock":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return err
		}
		vid, err := base.DecodeScaleArgBytes[uint64](args[0])
		if err != nil {
			return fmt.Errorf("unlock: voteId: %w", err)
		}
		return m.Unlock(vid)
	case "exec_proposal":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return err
		}
		pid, err := base.DecodeScaleArgBytes[uint32](args[0])
		if err != nil {
			return fmt.Errorf("exec_proposal: proposalId: %w", err)
		}
		return m.ExecProposal(pid)
	case "cancel_proposal":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return err
		}
		pid, err := base.DecodeScaleArgBytes[uint32](args[0])
		if err != nil {
			return fmt.Errorf("cancel_proposal: proposalId: %w", err)
		}
		return m.CancelProposal(pid)
	case "transfer":
		if err := base.RequireArgLen(args, 2, method); err != nil {
			return err
		}
		to, err := base.DecodeScaleArgBytes[[]byte](args[0])
		if err != nil {
			return fmt.Errorf("transfer: to: %w", err)
		}
		val, err := base.DecodeScaleArgBytes[[]byte](args[1])
		if err != nil {
			return fmt.Errorf("transfer: value: %w", err)
		}
		return m.Transfer(to, val)
	case "approve":
		if err := base.RequireArgLen(args, 2, method); err != nil {
			return err
		}
		sp, err := base.DecodeScaleArgBytes[[]byte](args[0])
		if err != nil {
			return fmt.Errorf("approve: spender: %w", err)
		}
		val, err := base.DecodeScaleArgBytes[[]byte](args[1])
		if err != nil {
			return fmt.Errorf("approve: value: %w", err)
		}
		return m.Approve(sp, val)
	case "transfer_from":
		if err := base.RequireArgLen(args, 3, method); err != nil {
			return err
		}
		from, err := base.DecodeScaleArgBytes[[]byte](args[0])
		if err != nil {
			return fmt.Errorf("transfer_from: from: %w", err)
		}
		to, err := base.DecodeScaleArgBytes[[]byte](args[1])
		if err != nil {
			return fmt.Errorf("transfer_from: to: %w", err)
		}
		val, err := base.DecodeScaleArgBytes[[]byte](args[2])
		if err != nil {
			return fmt.Errorf("transfer_from: value: %w", err)
		}
		return m.TransferFrom(from, to, val)
	case "spend":
		if err := base.RequireArgLen(args, 3, method); err != nil {
			return err
		}
		to, err := base.DecodeScaleArgBytes[[]byte](args[0])
		if err != nil {
			return fmt.Errorf("spend: to: %w", err)
		}
		amt, err := base.DecodeScaleArgBytes[[]byte](args[1])
		if err != nil {
			return fmt.Errorf("spend: amount: %w", err)
		}
		tid, err := base.DecodeScaleArgBytes[uint32](args[2])
		if err != nil {
			return fmt.Errorf("spend: trackId: %w", err)
		}
		return m.Spend(to, amt, tid)
	case "payout":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return err
		}
		sid, err := base.DecodeScaleArgBytes[uint64](args[0])
		if err != nil {
			return fmt.Errorf("payout: spendId: %w", err)
		}
		return m.Payout(sid)
	case "set_public_join":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return err
		}
		pj, err := base.DecodeScaleArgBytes[bool](args[0])
		if err != nil {
			return fmt.Errorf("set_public_join: %w", err)
		}
		return m.SetPublicJoin(pj)
	case "add_track":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return err
		}
		td, err := base.DecodeScaleArgBytes[TrackData](args[0])
		if err != nil {
			return fmt.Errorf("add_track: %w", err)
		}
		return m.AddTrack(td)
	case "set_default_track":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return err
		}
		tid, err := base.DecodeScaleArgBytes[uint32](args[0])
		if err != nil {
			return fmt.Errorf("set_default_track: trackId: %w", err)
		}
		return m.SetDefaultTrack(tid)
	default:
		return fmt.Errorf("dao: unknown method %q", method)
	}
}

// ExecuteQuery 按 method 将字符串参数解析为合约查询实参，返回 SCALE 编码结果。
func (q DaoQuery) ExecQuery(call *model.ContractCall) ([]byte, error) {
	if call.Method == "" {
		return nil, errors.New("dao: empty query method")
	}

	args := call.Args
	method := call.Method

	switch method {
	case "members":
		members, err := q.Members()
		if err != nil {
			return nil, err
		}

		return codec.Encode(members)
	case "public_join":
		v, err := q.PublicJoin()
		if err != nil {
			return nil, err
		}
		return codec.Encode(v)
	case "total_supply":
		b, err := q.TotalSupply()
		if err != nil {
			return nil, err
		}
		return codec.Encode(b)
	case "balance_of":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return nil, err
		}
		owner, err := base.DecodeScaleArgBytes[[]byte](args[0])
		if err != nil {
			return nil, fmt.Errorf("balance_of: owner: %w", err)
		}
		b, err := q.BalanceOf(owner)
		if err != nil {
			return nil, err
		}
		return codec.Encode(b)
	case "lock_balance_of":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return nil, err
		}
		owner, err := base.DecodeScaleArgBytes[[]byte](args[0])
		if err != nil {
			return nil, fmt.Errorf("lock_balance_of: owner: %w", err)
		}
		b, err := q.LockBalanceOf(owner)
		if err != nil {
			return nil, err
		}
		return codec.Encode(b)
	case "allowance":
		if err := base.RequireArgLen(args, 2, method); err != nil {
			return nil, err
		}
		owner, err := base.DecodeScaleArgBytes[[]byte](args[0])
		if err != nil {
			return nil, fmt.Errorf("allowance: owner: %w", err)
		}
		spender, err := base.DecodeScaleArgBytes[[]byte](args[1])
		if err != nil {
			return nil, fmt.Errorf("allowance: spender: %w", err)
		}
		b, err := q.Allowance(owner, spender)
		if err != nil {
			return nil, err
		}
		return codec.Encode(b)
	case "tracks":
		tracks, err := q.Tracks()
		if err != nil {
			return nil, err
		}

		return codec.Encode(tracks)
	case "track":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return nil, err
		}
		id, err := base.DecodeScaleArgBytes[uint32](args[0])
		if err != nil {
			return nil, fmt.Errorf("track: id: %w", err)
		}
		t, err := q.Track(id)
		if err != nil {
			return nil, err
		}

		return codec.Encode(t)
	case "default_track":
		id, err := q.DefaultTrack()
		if err != nil {
			return nil, err
		}
		return codec.Encode(id)
	case "proposal":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return nil, err
		}
		id, err := base.DecodeScaleArgBytes[uint32](args[0])
		if err != nil {
			return nil, fmt.Errorf("proposal: id: %w", err)
		}
		p, err := q.Proposal(id)
		if err != nil {
			return nil, err
		}

		return codec.Encode(p)
	case "proposals":
		ps, err := q.Proposals()
		if err != nil {
			return nil, err
		}
		return codec.Encode(ps)
	case "proposal_status":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return nil, err
		}
		id, err := base.DecodeScaleArgBytes[uint32](args[0])
		if err != nil {
			return nil, fmt.Errorf("proposal_status: id: %w", err)
		}
		st, err := q.ProposalStatus(id)
		if err != nil {
			return nil, err
		}
		return codec.Encode(st)
	case "vote":
		if err := base.RequireArgLen(args, 1, method); err != nil {
			return nil, err
		}
		id, err := base.DecodeScaleArgBytes[uint64](args[0])
		if err != nil {
			return nil, fmt.Errorf("vote: id: %w", err)
		}
		v, err := q.Vote(id)
		if err != nil {
			return nil, err
		}

		return codec.Encode(v)
	case "votes":
		vs, err := q.Votes()
		if err != nil {
			return nil, err
		}
		return codec.Encode(vs)
	default:
		return nil, fmt.Errorf("dao: unknown query method %q", method)
	}
}
