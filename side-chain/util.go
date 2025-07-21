package sidechain

import (
	"bytes"
	"fmt"
	"time"

	"github.com/cometbft/cometbft/crypto"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

const (
	CodeTypeOK              uint32 = 0
	CodeTypeEncodingError   uint32 = 1
	CodeTypeInvalidTxFormat uint32 = 2
	CodeTypeBanned          uint32 = 3
	CodeInvalidTEE          uint32 = 4
	CodeInvalidNode         uint32 = 5
)

const (
	GLOABL_STATE = "G"
)

func LogWithTime(a ...any) {
	dim := "\033[2m"
	reset := "\033[0m"
	tag := dim + "> " + time.Now().Format("01/02 15:04:05") + reset
	a = append([]any{tag}, a...)
	fmt.Println(a...)
}

func (s *SideChain) ProposerAddressToNodeKey(proposer []byte) *model.PubKey {
	for _, node := range s.dkg.Nodes {
		if bytes.Equal(proposer, crypto.AddressHash(node.ValidatorId.PublicKey)) {
			return &node.P2pId
		}
	}
	panic("ProposerAddress not in DKG")
}
