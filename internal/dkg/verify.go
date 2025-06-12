package dkg

import (
	"encoding/hex"
	"errors"

	"github.com/vedhavyas/go-subkey/v2"
	"wetee.app/dsecret/internal/chain"
	"wetee.app/dsecret/internal/tee"
	types "wetee.app/dsecret/type"
	gtypes "wetee.app/dsecret/type/pallet/types"
)

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (d *DKG) VerifyWorker(reportData *types.TeeParam) ([]byte, error) {
	// 解码地址
	_, signer, err := subkey.SS58Decode(reportData.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	report, err := tee.VerifyReport(reportData)
	if err != nil {
		return nil, errors.New("verify cluster report: " + err.Error())
	}

	// 校验 worker 代码版本
	codeHash, codeSigner, err := chain.GetWorkerCode(chain.ChainIns.ChainClient)
	if err != nil {
		return nil, errors.New("GetWorkerCode error:" + err.Error())
	}
	if len(codeHash) > 0 || len(codeSigner) > 0 {
		if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
			return nil, errors.New("worker code hash error")
		}

		if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
			return nil, errors.New("worker signer error")
		}
	}

	return signer, nil
}

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (d *DKG) VerifyDsecret(reportData *types.TeeParam) ([]byte, error) {
	// 解码地址
	_, signer, err := subkey.SS58Decode(reportData.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	report, err := tee.VerifyReport(reportData)
	if err != nil {
		return nil, errors.New("verify cluster report: " + err.Error())
	}

	// 校验 worker 代码版本
	codeHash, codeSigner, err := chain.GetDsecretCode(chain.ChainIns.ChainClient)
	if err != nil {
		return nil, errors.New("GetWorkerCode error:" + err.Error())
	}
	if len(codeHash) > 0 || len(codeSigner) > 0 {
		if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
			return nil, errors.New("worker code hash error")
		}

		if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
			return nil, errors.New("worker signer error")
		}
	}

	return signer, nil
}

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (d *DKG) VerifyWorkLibos(wid gtypes.WorkId, reportData *types.TeeParam) ([]byte, error) {
	// 解码地址
	_, signer, err := subkey.SS58Decode(reportData.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	report, err := tee.VerifyReport(reportData)
	if err != nil {
		return nil, errors.New("verify cluster report: " + err.Error())
	}

	// 校验 worker 代码版本
	codeHash, codeSigner, err := chain.GetWorkCode(chain.ChainIns.ChainClient, wid)
	if err != nil {
		return nil, errors.New("GetWorkerCode error:" + err.Error())
	}
	if len(codeHash) > 0 || len(codeSigner) > 0 {
		if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
			return nil, errors.New("worker code hash error")
		}

		if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
			return nil, errors.New("worker signer error")
		}
	}

	return signer, nil
}
