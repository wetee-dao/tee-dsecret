package types

import (
	"fmt"

	rabin "go.dedis.ch/kyber/v3/share/dkg/rabin"
	vss "go.dedis.ch/kyber/v3/share/vss/rabin"
	"go.dedis.ch/kyber/v3/suites"
)

func DealToProtocol(deal *rabin.Deal) (*Deal, error) {
	dkheyBytes, err := deal.Deal.DHKey.MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal dhkey: %w", err)
	}

	return &Deal{
		Index: deal.Index,
		Deal: &EncryptedDeal{
			DHKey:     dkheyBytes,
			Signature: deal.Deal.Signature,
			Nonce:     deal.Deal.Nonce,
			Cipher:    deal.Deal.Cipher,
		},
	}, nil
}

func ProtocolToDeal(suite suites.Suite, deal *Deal) (*rabin.Deal, error) {
	p := suite.Point()
	err := p.UnmarshalBinary(deal.Deal.DHKey)
	if err != nil {
		return nil, fmt.Errorf("unmarshal dhkey: %w", err)
	}
	return &rabin.Deal{
		Index: deal.Index,
		Deal: &vss.EncryptedDeal{
			DHKey:     p,
			Signature: deal.Deal.Signature,
			Nonce:     deal.Deal.Nonce,
			Cipher:    deal.Deal.Cipher,
		},
	}, nil
}

type Deal struct {
	Index uint32
	Deal  *EncryptedDeal
}

type EncryptedDeal struct {
	// Ephemeral Diffie Hellman key
	DHKey []byte
	// Signature of the DH key by the longterm key of the dealer
	Signature []byte
	// Nonce used for the encryption
	Nonce []byte
	// AEAD encryption of the deal marshalled by protobuf
	Cipher []byte
}
