package sidechain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/dkg"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

func (app *SideChain) CheckEpochFromValidator() []byte {
	if chains.MainChain == nil {
		util.LogWithRed("CheckEpochFromValidator", "error chains.MainChain is nil")
		return nil
	}

	if app.dkg == nil {
		util.LogWithRed("CheckEpochFromValidator", "error app.dkg is nil")
		return nil
	}

	// Query epoch from main chain
	epoch, epochSolt, lastEpochBlock, now, _, err := chains.MainChain.GetEpoch()
	if err != nil {
		return nil
	}

	// STEP3 update new epoch to side chain (sync epoch form main chain)
	if app.GetEpoch() < epoch {
		util.LogWithYellow("NewEpoch", "P3")
		validators, err := chains.MainChain.GetValidatorList()
		if err != nil {
			return nil
		}

		sideValidators := make([]*model.SideValidator, 0, len(validators))
		for _, v := range validators {
			sideValidators = append(sideValidators, &model.SideValidator{
				Pubkey: v.ValidatorId.PublicKey,
				Power:  1,
			})
		}

		return GetTxBytes(&model.Tx{
			Payload: &model.Tx_EpochEnd{
				EpochEnd: &model.EpochEnd{
					Epoch:      epoch,
					Validators: sideValidators,
					DkgPub:     app.dkg.DkgPubKey.PublicKey,
				},
			},
		})
	}

	// Check if sync tx is submiting
	if IsSyncRuning() {
		util.LogWithYellow("CheckEpochFromValidator", "Sync is running, please wait...")
		return nil
	}

	// Query local epoch
	epochStatus := app.GetEpochStatus()

	// STEP1 check new epoch
	if now-lastEpochBlock >= epochSolt-1 || epoch == 0 {
		if time.Now().Unix()-int64(epochStatus) > 120 {
			validators, err := chains.MainChain.GetNextEpochValidatorList()
			if err != nil {
				util.LogWithYellow("GetNextEpochValidatorList error:", err)
				return nil
			}

			util.LogWithYellow("NewEpoch", "P1 at epoch", epoch)
			err = app.dkg.TryEpochConsensus(model.ConsensusMsg{
				Validators: validators,
				Epoch:      epoch + 1,
			}, app.newEpochSucceded, app.newEpochFail)
			if err == nil {
				return GetTxBytes(&model.Tx{
					Payload: &model.Tx_EpochStart{
						EpochStart: time.Now().Unix(), // start epoch, stop submit main chain tx
					},
				})
			}
		}
		return nil
	}

	return nil
}

// STEP2
func (app *SideChain) newEpochSucceded(signer *dkg.DssSigner, nodeId uint64) {
	util.LogWithYellow("NewEpoch", "P2 submit tx to main chain")
	if chains.MainChain == nil {
		util.LogWithRed("NewEpoch CheckEpochFromValidator", "error chains.MainChain is nil")
		return
	}

	// submit new epoch to main chain
	call, _ := chains.MainChain.TxCallOfSetNextEpoch(nodeId, signer.AccountID())

	client := chains.MainChain.GetClient()
	err := client.SignAndSubmit(signer, *call, false, 0)
	if err != nil {
		util.LogWithRed("NewEpoch client.SignAndSubmit", err.Error())
	}
}

// Callback when DKG consensus failed
func (app *SideChain) newEpochFail(err error) {
	util.LogWithYellow("NewEpoch", "P2 Error", err.Error())
}

// Callback when DKG consensus success
func (app *SideChain) GetEpochStatus() int64 {
	bt, err := model.GetKey(GLOABL_STATE, "epoch_status")
	if err != nil {
		return 0
	}

	return util.BytesToInt64(bt)
}

// SetEpochStatus last epoch timestamp
func (app *SideChain) SetEpochStatus(status int64) error {
	return model.SetKey(GLOABL_STATE, "epoch_status", util.Int64ToBytes(status))
}

// GetEpoch get last epoch
func (app *SideChain) GetEpoch() uint32 {
	bt, err := model.GetKey(GLOABL_STATE, "epoch")
	if err != nil {
		return 0
	}

	bytesBuffer := bytes.NewBuffer(bt)
	var x uint32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return x
}

// SetEpoch set new epoch
func (app *SideChain) SetEpoch(epoch *model.EpochEnd, txn *model.Txn) error {
	// Save new epoch to DKG
	if app.dkg != nil {
		app.dkg.ToNewEpoch()
	}

	// Save new epoch to DB
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, epoch.Epoch)
	txn.SetKey(GLOABL_STATE, "epoch", bytesBuffer.Bytes())

	// Save DKG pub key
	txn.SetKey(GLOABL_STATE, "dkg_pub_key", epoch.DkgPub)

	// Delete old epoch validators
	err := txn.DeletekeysByPrefix([]byte("G_validator"))
	if err != nil {
		return err
	}

	// Save new epoch validators
	for i, v := range epoch.Validators {
		if err = model.TxnSetProtoMessage(txn, []byte("G_validator"+fmt.Sprint(i)), v); err != nil {
			return err
		}
	}

	return nil
}

// Get validators form local db
func (app *SideChain) GetValidators() ([]*model.SideValidator, map[string]*model.PubKey, error) {
	var err error
	list, _, err := model.GetProtoMessageList[model.SideValidator](GLOABL_STATE, "validator")
	if err != nil {
		return nil, nil, err
	}

	validatorMap := map[string]*model.PubKey{}
	for _, v := range list {
		pub := model.PubKeyFromByte(v.Pubkey)
		validatorMap[pub.SS58()] = pub
	}

	return list, validatorMap, nil
}

// Init validator to db From init chain
func (app *SideChain) initValidators(vs []abci.ValidatorUpdate) error {
	tx := model.DBINS.NewTransaction()

	var err error
	for i, v := range vs {
		if err = model.TxnSetProtoMessage(tx, []byte("G_validator"+fmt.Sprint(i)), &model.SideValidator{
			Pubkey: v.GetPubKeyBytes(),
			Power:  v.GetPower(),
		}); err != nil {
			return err
		}
	}

	tx.Commit()
	return nil
}

// Calc validator updates
func (app *SideChain) calcValidatorUpdates(epoch *model.EpochEnd) {
	oldValidators, _, _ := app.GetValidators()
	newValidators := epoch.Validators

	for _, newv := range newValidators {
		isIn := false
		for i, oldv := range oldValidators {
			// update validator
			if bytes.Equal(oldv.Pubkey, newv.Pubkey) {
				oldValidators[i].Power = newv.Power
				isIn = true
			}
		}

		// add new validator
		if !isIn {
			oldValidators = append(oldValidators, newv)
		}
	}

	for i, oldv := range oldValidators {
		isIn := false
		for _, newv := range newValidators {
			if bytes.Equal(oldv.Pubkey, newv.Pubkey) {
				isIn = true
			}
		}

		// delete old validator
		if !isIn {
			oldValidators[i].Power = 0
		}
	}

	ongoing := make([]abci.ValidatorUpdate, 0, len(oldValidators))
	for _, v := range oldValidators {
		ongoing = append(ongoing, abci.ValidatorUpdate{
			PubKeyType:  "ed25519",
			PubKeyBytes: v.Pubkey,
			Power:       v.Power,
		})
	}

	app.onGoingValidators = ongoing
}
