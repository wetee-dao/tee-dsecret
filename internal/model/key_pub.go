package model

import (
	gocrypto "crypto"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/p2p"
	"go.dedis.ch/kyber/v4"
	"go.dedis.ch/kyber/v4/suites"
)

type PubKey struct {
	ed25519.PublicKey
	suite suites.Suite
}

func (p *PubKey) Suite() suites.Suite {
	return p.suite
}

func (p *PubKey) Point() kyber.Point {
	buf := p.PublicKey
	point := p.suite.Point()
	point.UnmarshalBinary(buf)
	return point
}

func (p *PubKey) Std() (gocrypto.PublicKey, error) {
	return p.PublicKey, nil
}

func (p *PubKey) String() string {
	return hex.EncodeToString(p.PublicKey)
}

func (p *PubKey) SideChainNodeID() p2p.ID {
	return p2p.ID(hex.EncodeToString(crypto.AddressHash(p.PublicKey)))
}

func (p *PubKey) Byte() ([]byte, error) {
	return p.PublicKey, nil
}

func (p *PubKey) SS58() string {
	k, err := p.Std()
	if err != nil {
		return "error key"
	}

	return SS58Encode(k.(ed25519.PublicKey), 42)
}

func (p PubKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.PublicKey)
}

func (p *PubKey) UnmarshalJSON(bt []byte) error {
	var key ed25519.PublicKey = []byte{}
	err := json.Unmarshal(bt, &key)
	if err != nil {
		return err
	}

	p.PublicKey = key
	p.suite = suites.MustFind("Ed25519")
	return nil
}

func PubKeyFromStdPubKey(pubkey ed25519.PublicKey) (*PubKey, error) {
	return &PubKey{
		PublicKey: pubkey,
		suite:     suites.MustFind("Ed25519"),
	}, nil
}

func PubKeyFromPoint(suite suites.Suite, point kyber.Point) (*PubKey, error) {
	buf, err := point.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal point: %w", err)
	}

	var pk ed25519.PublicKey

	switch strings.ToLower(suite.String()) {
	case "ed25519":
		pk = buf
	default:
		return nil, fmt.Errorf("unknown suite: %v", suite)
	}

	return &PubKey{
		PublicKey: pk,
	}, nil
}

func PubKeyFromByte(pubkey []byte) *PubKey {
	return &PubKey{
		PublicKey: ed25519.PublicKey(pubkey),
		suite:     suites.MustFind("Ed25519"),
	}
}
