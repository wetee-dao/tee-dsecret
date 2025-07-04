// / Copyright (c) 2022 Sourcenetwork Developers. All rights reserved.
// / copy from https://github.com/sourcenetwork/orbis-go

package proxy_reenc

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v4/share"
	"go.dedis.ch/kyber/v4/suites"
	"go.dedis.ch/kyber/v4/util/random"
)

func TestReencryptAndVerify(t *testing.T) {
	var (
		n     = 5 //参与者的数量
		th    = 3 // 门限值
		suite = suites.MustFind("Ed25519")

		s       = suite.Scalar().Pick(suite.RandomStream())
		priPoly = share.NewPriPoly(suite, th, s, suite.RandomStream())
		pubPoly = priPoly.Commit(nil)

		dkgPk = pubPoly.Commit() // DKG 公钥

		rdrSk = suite.Scalar().Pick(suite.RandomStream())
		rdrPk = suite.Point().Mul(rdrSk, nil)
	)

	var pubShares []*share.PubShare

	// Generate a random secret.
	scrt := make([]byte, 32)
	random.Bytes(scrt, random.New())

	// 1. Encrypt the secret under the DKG public key.
	encCmt, encScrt := EncryptSecret(suite, dkgPk, scrt)

	for idx := range n {
		dkgSki := priPoly.Eval(uint32(idx)).V
		dkgCmt := pubPoly.Eval(uint32(idx)).V

		// 2. Re-encrypt the key under the reader's public key.
		xncSki, chlgi, proofi, err := reencrypt(suite, dkgSki, rdrPk, encCmt)
		require.NoErrorf(t, err, "failed to reencrypt for share %d", idx)

		// 3. Verify the re-encryption from other nodes.
		err = verify(suite, rdrPk, encCmt, xncSki, chlgi, proofi, dkgCmt)
		require.NoErrorf(t, err, "failed to verify reencryption for share %d", idx)

		pubShare := &share.PubShare{I: uint32(idx), V: xncSki}
		pubShares = append(pubShares, pubShare)
	}

	// 4 - Recover re-encrypted commmitment using Lagrange interpolation.
	// ski * (xG + rG) => rsG + xsG
	xncCmt, err := share.RecoverCommit(suite, pubShares, th, n)
	require.NoErrorf(t, err, "failed to recover commit")

	// 5 - Decode encrypted key with re-encrypted commitment and reader's privatekey.
	scrtHat, err := DecryptSecret(suite, encScrt, dkgPk, xncCmt, rdrSk)
	require.NoErrorf(t, err, "failed to decode key")
	require.Equal(t, scrt, scrtHat)
}
