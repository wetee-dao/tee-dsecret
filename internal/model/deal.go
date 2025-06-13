package model

import (
	"fmt"

	"go.dedis.ch/kyber/v4"
	pedersen "go.dedis.ch/kyber/v4/share/dkg/pedersen"
	"go.dedis.ch/kyber/v4/suites"
)

type Deal struct {
	DealerIndex uint32
	Deals       []pedersen.Deal
	// Public coefficients of the public polynomial used to create the shares
	Public [][]byte
	// SessionID of the current run
	SessionID []byte
	// Signature over the hash of the whole bundle
	Signature []byte
	Reshare   int
}

func DealToProtocol(deal *pedersen.DealBundle) (*Deal, error) {
	public := deal.Public
	points := make([][]byte, len(public))
	for i, p := range public {
		b, err := p.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("MarshalBinary public error: %w", err)
		}
		points[i] = b
	}

	return &Deal{
		DealerIndex: deal.DealerIndex,
		Deals:       deal.Deals,
		Public:      points,
		SessionID:   deal.SessionID,
		Signature:   deal.Signature,
	}, nil
}

func ProtocolToDeal(suite suites.Suite, deal *Deal) (*pedersen.DealBundle, error) {
	public := deal.Public
	points := make([]kyber.Point, len(public))
	for i, p := range public {
		point := suite.Point()
		err := point.UnmarshalBinary(p)
		if err != nil {
			return nil, fmt.Errorf("unmarshal commitment: %w", err)
		}

		points[i] = point
	}

	return &pedersen.DealBundle{
		DealerIndex: deal.DealerIndex,
		Deals:       deal.Deals,
		Public:      points,
		SessionID:   deal.SessionID,
		Signature:   deal.Signature,
	}, nil
}

type JustificationBundle struct {
	DealerIndex    uint32
	Justifications []Justification
	// SessionID of the current run
	SessionID []byte
	// Signature over the hash of the whole bundle
	Signature []byte
}

type Justification struct {
	ShareIndex uint32
	Share      []byte
}

func JustificationToProtocol(deal *pedersen.JustificationBundle) (*JustificationBundle, error) {
	js := deal.Justifications
	newJs := make([]Justification, len(js))
	for i, p := range js {
		b, err := p.Share.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("MarshalBinary public error: %w", err)
		}
		newJs[i] = Justification{
			ShareIndex: p.ShareIndex,
			Share:      b,
		}
	}

	return &JustificationBundle{
		DealerIndex:    deal.DealerIndex,
		Justifications: newJs,
		SessionID:      deal.SessionID,
		Signature:      deal.Signature,
	}, nil
}

func ProtocolToJustification(suite suites.Suite, deal *JustificationBundle) (*pedersen.JustificationBundle, error) {
	js := deal.Justifications
	newJs := make([]pedersen.Justification, len(js))
	for i, p := range js {
		point := suite.Scalar()
		err := point.UnmarshalBinary(p.Share)
		if err != nil {
			return nil, fmt.Errorf("unmarshal commitment: %w", err)
		}

		newJs[i] = pedersen.Justification{
			ShareIndex: p.ShareIndex,
			Share:      point,
		}
	}

	return &pedersen.JustificationBundle{
		DealerIndex:    deal.DealerIndex,
		Justifications: newJs,
		SessionID:      deal.SessionID,
		Signature:      deal.Signature,
	}, nil
}
