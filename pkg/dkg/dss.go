package dkg

import (
	"crypto/ed25519"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"go.dedis.ch/kyber/v4"
	"go.dedis.ch/kyber/v4/sign/dss"
)

func NewDssSigner(dkg *DKG) *DssSigner {
	return &DssSigner{
		dkg:  dkg,
		sigs: make([]*dss.PartialSig, 0),
	}
}

type DssSigner struct {
	dkg  *DKG
	sigs []*dss.PartialSig
}

func (d *DssSigner) SetSigs(btsigs [][]byte) {
	sigs := make([]*dss.PartialSig, 0, len(btsigs))
	for _, bt := range btsigs {
		sig := &model.PartialSigWrap{}
		err := json.Unmarshal(bt, sig)
		if err != nil {
			fmt.Println(err)
			continue
		}
		sigs = append(sigs, &dss.PartialSig{
			Partial:   sig.Partial.PriShare,
			SessionID: sig.SessionID,
			Signature: sig.Signature,
		})
	}
	d.sigs = sigs
}

func (d *DssSigner) Public() []byte {
	if d.dkg.DkgPubKey != nil {
		return d.dkg.DkgPubKey.Byte()
	}
	return d.dkg.NewDkgPubKey.Byte()
}

func (d *DssSigner) AccountID() types.AccountID {
	if d.dkg.DkgPubKey != nil {
		return d.dkg.DkgPubKey.AccountID()
	}
	return d.dkg.NewDkgPubKey.AccountID()
}

func (d *DssSigner) Sign(msg []byte) ([]byte, error) {
	pubs, long, random, threshold := d.PubList()
	dss, err := dss.NewDSS(
		d.dkg.Suite,
		d.dkg.Signer.Scalar(),
		pubs,
		long,
		random,
		msg,
		threshold,
	)

	if dss == nil || err != nil {
		fmt.Println(err)
		return nil, errors.New("dss.NewDSS failed")
	}

	for _, sig := range d.sigs {
		err = dss.ProcessPartialSig(sig)
		if err != nil {
			return nil, err
		}
	}

	return dss.Signature()
}

func (d *DssSigner) PartialSign(msg []byte) ([]byte, error) {
	pubs, long, random, threshold := d.PubList()
	dss, err := dss.NewDSS(
		d.dkg.Suite,
		d.dkg.Signer.Scalar(),
		pubs,
		long,
		random,
		msg,
		threshold,
	)

	if dss == nil || err != nil {
		return nil, errors.New("dss.NewDSS failed")
	}

	sig, err := dss.PartialSig()
	if err != nil {
		return nil, errors.New("dss.PartialSig failed")
	}
	sigWrap := model.PartialSigWrap{
		Partial:   model.PriShare{PriShare: sig.Partial},
		SessionID: sig.SessionID,
		Signature: sig.Signature,
	}

	return json.Marshal(sigWrap)
}

func (d *DssSigner) PubList() ([]kyber.Point, *model.DistKeyShare, *model.DistKeyShare, int) {
	pubs := make([]kyber.Point, 0, len(d.dkg.Nodes))
	if len(d.dkg.NewNodes) > 0 {
		for _, k := range d.dkg.NewNodes {
			pubs = append(pubs, k.ValidatorId.Point())
		}
		return pubs, d.dkg.NewDkgKeyShare, d.dkg.NewDkgKeyShare, len(d.dkg.NewNodes) * 2 / 3
	}

	for _, k := range d.dkg.Nodes {
		pubs = append(pubs, k.ValidatorId.Point())
	}
	return pubs, d.dkg.DkgKeyShare, d.dkg.DkgKeyShare, d.dkg.Threshold
}

func (d *DssSigner) Verify(msg []byte, signature []byte) bool {
	return ed25519.Verify(d.dkg.DkgPubKey.Ed25519PublicKey(), msg, signature)
}

func (d *DssSigner) SignType() uint8 {
	return 1
}
