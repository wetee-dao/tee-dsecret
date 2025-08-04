package sidechain

import (
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v4/share"
	"go.dedis.ch/kyber/v4/suites"

	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	proxy_reenc "github.com/wetee-dao/tee-dsecret/pkg/proxy-reenc"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

// PreRecerve is the channel to receive DecryptShares
var preRecerve map[uint64]chan *model.DecryptSharesResp = make(map[uint64]chan *model.DecryptSharesResp)

// Recive msg from p2p
func (s *SideChain) revSecret(m any) error {
	mbox := m.(*model.SecretBox)
	switch msg := mbox.Payload.(type) {
	case *model.SecretBox_Req:
		return s.HandleDecryptSecret(msg.Req, mbox.From)
	case *model.SecretBox_SharesResp:
		return s.VerifyDecryptSecret(msg.SharesResp)
	default:
		return fmt.Errorf("unknown secret message type")
	}
}

// BroadcastDecryptSecret broadcast decrypt secret request to all nodes
func (s *SideChain) BroadcastDecryptSecret(req *model.PodStart) (*model.DecryptResp, error) {
	suite := suites.MustFind("Ed25519")
	validators, err := chains.MainChain.GetValidatorList()
	if err != nil {
		return nil, fmt.Errorf("get validator list: %w", err)
	}
	threshold := len(validators)*2/3 + 1
	validatorP2Pkeys := make([]*model.PubKey, 0, len(validators))
	for _, v := range validators {
		validatorP2Pkeys = append(validatorP2Pkeys, &v.P2pId)
	}

	// 初始化重新加密回复
	preRecerve[req.Id] = make(chan *model.DecryptSharesResp, len(validators))

	// send decrypt secret request to all nodes
	dshares := make([]*model.DecryptSharesResp, 0, threshold)
	err = s.p2p.Send(model.SendToNodes(validatorP2Pkeys), &model.SecretBox{
		Payload: &model.SecretBox_Req{
			Req: req,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("send decrypt secret: %w", err)
	}

	// 收集至少达到阈值数量的节点响应
	for range threshold {
		select {
		case d := <-preRecerve[req.Id]:
			dshares = append(dshares, d)
		case <-time.After(30 * time.Second):
			util.LogError("BroadcastDecryptSecret", "Timeout receiving from channel")
			return nil, fmt.Errorf("timeout receiving from channel")
		}
	}

	// shares 收集解密后的份额
	shares := make(map[uint64][]*share.PubShare)
	for _, d := range dshares {
		for index, s := range d.Shares {
			if _, ok := shares[uint64(index)]; !ok {
				shares[uint64(index)] = make([]*share.PubShare, 0, threshold)
			}

			reply, err := DecodeDecryptShare(s, suite)
			if err != nil {
				return nil, fmt.Errorf("decode decrypt share: %w", err)
			}

			shares[uint64(index)] = append(shares[uint64(index)], &reply.Share)
		}
	}

	nameSpace := types.H160(req.NameSpace)
	secrets, err := s.GetSecrets(nameSpace, req.Indexs)
	if err != nil {
		return nil, fmt.Errorf("get secret: %w", err)
	}

	encodeSecret := make(map[uint64]*model.Secret)
	dkgPubKey, err := GetDkgPubkey()
	if err != nil {
		return nil, fmt.Errorf("get dkg pubkey: %w", err)
	}

	for index, shares := range shares {
		// 从收集的响应中恢复重加密承诺
		xncCmt, err := proxy_reenc.Recover(suite, shares, threshold, len(validators))
		if err != nil {
			return nil, fmt.Errorf("recover reencrypt reply: %s", err)
		}

		bt, err := xncCmt.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("marshal xnc cmt: %s", err)
		}

		encodeSecret[index] = &model.Secret{
			EncScrt: secrets[index].RawEncScrt,
			XncCmt:  bt,
		}
	}

	return &model.DecryptResp{
		DkgKey: dkgPubKey.ToBytes(),
		Lists:  encodeSecret,
	}, nil
	// scrtHat, err := proxy_reenc.DecryptSecret(suite, encScrt, dkg.DkgPubKey.Point(), xncCmt, rdrSk.Scalar())
	// if err != nil {
	// 	return fmt.Errorf("decrypt secret: %s", err)
	// }
}

// HandleDecryptSecret 处理解密请求
func (s *SideChain) HandleDecryptSecret(req *model.PodStart, from string) error {
	dkg := s.dkg

	// 获取重新加密所需的公钥和密文
	clientPubKey := model.PubKeyFromByte(req.PubKey)

	nameSpace := types.H160(req.NameSpace)
	secrets, err := s.GetSecrets(nameSpace, req.Indexs)
	if err != nil {
		return fmt.Errorf("get secret: %w", err)
	}

	// 获取本节点的份额，并进行重新加密操作
	dkgShare := dkg.Share()

	// 重加密所有的内容
	shares := make(map[uint64]*model.DecryptShare)
	for index, secret := range secrets {
		reply, err := proxy_reenc.Reencrypt(dkgShare, secret, *clientPubKey)
		if err != nil {
			return fmt.Errorf("reencrypt: %w", err)
		}

		// 编码重新加密的密文份额响应
		eshare, err := EncodeDecryptShare(reply, req.Id)
		if err != nil {
			return fmt.Errorf("encode decrypt share: %w", err)
		}

		// 构建重新加密的密文份额响应
		shares[index] = eshare
	}

	// 发送重新加密的密文份额响应
	formPubKey, err := model.PubKeyFromHex(from)
	if err != nil {
		return fmt.Errorf("pubkey from hex: %w", err)
	}
	err = s.p2p.Send(model.SendToNode(formPubKey), &model.SecretBox{
		Payload: &model.SecretBox_SharesResp{
			SharesResp: &model.DecryptSharesResp{
				Req:    req,
				Shares: shares,
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "P2P Send error")
	}

	return nil
}

// RevAndVerifyDecryptSecret 验证解密后的秘密
func (s *SideChain) VerifyDecryptSecret(shares *model.DecryptSharesResp) error {
	suite := suites.MustFind("Ed25519")
	req := shares.Req

	// 获取分布式密钥的份额和多项式承诺
	commits, err := GetDkgCommits()
	if err != nil {
		return fmt.Errorf("get dkg commits: %w", err)
	}
	poly := share.NewPubPoly(suite, nil, commits.Public)

	// 解析程序的空间
	nameSpace := types.H160(req.NameSpace)
	secrets, err := s.GetSecrets(nameSpace, req.Indexs)
	if err != nil {
		return fmt.Errorf("get secret: %w", err)
	}

	// 解析客户端的公钥
	clientPubKey := model.PubKeyFromByte(req.PubKey)

	// 验证所有的重新加密回复
	for index, share := range shares.Shares {
		reply, err := DecodeDecryptShare(share, suite)
		if err != nil {
			return fmt.Errorf("decode decrypt share: %w", err)
		}

		// 验证重新加密的回复
		secret := secrets[index]
		err = proxy_reenc.Verify(poly, secret, *clientPubKey, reply)
		if err != nil {
			return fmt.Errorf("verify reencrypt reply: %s", err)
		}
	}

	preRecerve[req.Id] <- shares
	return nil
}
