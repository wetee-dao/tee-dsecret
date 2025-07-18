package model

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"github.com/edgelesssys/ego/attestation"
	"github.com/edgelesssys/ego/attestation/tcbstatus"
	"github.com/edgelesssys/ego/enclave"
	"github.com/vedhavyas/go-subkey/v2/ed25519"
	chain "github.com/wetee-dao/ink.go"

	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

func IssueReport(pk *chain.Signer, data []byte) (*TeeParam, error) {
	timestamp := time.Now().Unix()

	var buf bytes.Buffer
	buf.Write(util.Int64ToBytes(timestamp))
	buf.Write(pk.PublicKey)
	if len(data) > 0 {
		buf.Write(data)
	}
	sig, err := pk.Sign(buf.Bytes())
	if err != nil {
		return nil, err
	}

	report, err := enclave.GetRemoteReport(sig)
	if err != nil {
		return nil, err
	}

	return &TeeParam{
		Time:    timestamp,
		Address: pk.PublicKey,
		Report:  report,
		Data:    data,
	}, nil
}

func VerifyReport(reportData *TeeParam) (*TeeReport, error) {
	// TODO SEV/TDX not support
	if reportData.TeeType != 0 {
		return &TeeReport{
			CodeSignature: []byte{},
			CodeSigner:    []byte{},
			CodeProductID: []byte{},
		}, nil
	}

	var reportBytes, msgBytes, timestamp = reportData.Report, reportData.Data, reportData.Time

	// decode address
	signer := reportData.Address

	report, err := enclave.VerifyRemoteReport(reportBytes)
	if err == attestation.ErrTCBLevelInvalid {
		fmt.Printf("Warning: TCB level is invalid: %v\n%v\n", report.TCBStatus, tcbstatus.Explain(report.TCBStatus))
	} else if err != nil {
		return nil, err
	}

	pubkey, err := ed25519.Scheme{}.FromPublicKey(signer)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	buf.Write(util.Int64ToBytes(timestamp))
	buf.Write(signer)
	if len(msgBytes) > 0 {
		buf.Write(msgBytes)
	}

	sig := report.Data

	if !pubkey.Verify(buf.Bytes(), sig) {
		return nil, errors.New("invalid sgx report")
	}

	// if report.Debug {
	// 	return nil, errors.New("debug mode is not allowed")
	// }

	return &TeeReport{
		TeeType:       reportData.TeeType,
		CodeSigner:    report.SignerID,
		CodeSignature: report.UniqueID,
		CodeProductID: report.ProductID,
	}, nil
}
