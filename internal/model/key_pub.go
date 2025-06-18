package model

import (
	gocrypto "crypto"
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"

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

func (p *PubKey) Byte() []byte {
	return p.PublicKey
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

func PubKeyFromPoint(point kyber.Point) (*PubKey, error) {
	buf, err := point.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal point: %w", err)
	}

	var pk ed25519.PublicKey
	pk = buf

	return &PubKey{
		PublicKey: pk,
		suite:     suites.MustFind("Ed25519"),
	}, nil
}

func PubKeyFromByte(pubkey []byte) *PubKey {
	return &PubKey{
		PublicKey: ed25519.PublicKey(pubkey),
		suite:     suites.MustFind("Ed25519"),
	}
}

func PubKeyFromSS58(ss58 string) (*PubKey, error) {
	_, pubkey, err := SS58Decode(ss58)
	if err != nil {
		return nil, fmt.Errorf("decode ss58: %w", err)
	}

	return PubKeyFromByte(pubkey), nil
}
