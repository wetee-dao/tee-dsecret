package dkg

import (
	"errors"

	"github.com/vedhavyas/go-subkey/v2"
	gtypes "github.com/wetee-dao/tee-dsecret/chains/pallets/generated/types"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (d *DKG) VerifyWorker(reportData *model.TeeParam) ([]byte, error) {
	// 解码地址
	_, signer, err := subkey.SS58Decode(reportData.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	// TODO
	// report, err := tee.VerifyReport(reportData)
	// if err != nil {
	// 	return nil, errors.New("verify cluster report: " + err.Error())
	// }

	// // 校验 worker 代码版本
	// codeHash, codeSigner, err := chain.GetWorkerCode(chain.MainChain.ChainClient)
	// if err != nil {
	// 	return nil, errors.New("GetWorkerCode error:" + err.Error())
	// }
	// if len(codeHash) > 0 || len(codeSigner) > 0 {
	// 	if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
	// 		return nil, errors.New("worker code hash error")
	// 	}

	// 	if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
	// 		return nil, errors.New("worker signer error")
	// 	}
	// }

	return signer, nil
}

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (d *DKG) VerifyDsecret(reportData *model.TeeParam) ([]byte, error) {
	// 解码地址
	_, signer, err := subkey.SS58Decode(reportData.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	// TODO
	// report, err := tee.VerifyReport(reportData)
	// if err != nil {
	// 	return nil, errors.New("verify cluster report: " + err.Error())
	// }

	// // 校验 worker 代码版本
	// codeHash, codeSigner, err := chain.GetDsecretCode(chain.MainChain.ChainClient)
	// if err != nil {
	// 	return nil, errors.New("GetWorkerCode error:" + err.Error())
	// }
	// if len(codeHash) > 0 || len(codeSigner) > 0 {
	// 	if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
	// 		return nil, errors.New("worker code hash error")
	// 	}

	// 	if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
	// 		return nil, errors.New("worker signer error")
	// 	}
	// }

	return signer, nil
}

// VerifyWorker 函数验证工人报告并返回签名者或错误
func (d *DKG) VerifyWorkLibos(wid gtypes.WorkId, reportData *model.TeeParam) ([]byte, error) {
	// 解码地址
	_, signer, err := subkey.SS58Decode(reportData.Address)
	if err != nil {
		return nil, errors.New("SS58 decode: " + err.Error())
	}

	// TODO
	// report, err := tee.VerifyReport(reportData)
	// if err != nil {
	// 	return nil, errors.New("verify cluster report: " + err.Error())
	// }

	// 校验 worker 代码版本
	// codeHash, codeSigner, err := chain.GetWorkCode(chain.MainChain.ChainClient, wid)
	// if err != nil {
	// 	return nil, errors.New("GetWorkerCode error:" + err.Error())
	// }
	// if len(codeHash) > 0 || len(codeSigner) > 0 {
	// 	if hex.EncodeToString(codeHash) != hex.EncodeToString(report.CodeSignature) {
	// 		return nil, errors.New("worker code hash error")
	// 	}

	// 	if hex.EncodeToString(codeSigner) != hex.EncodeToString(report.CodeSigner) {
	// 		return nil, errors.New("worker signer error")
	// 	}
	// }

	return signer, nil
}
