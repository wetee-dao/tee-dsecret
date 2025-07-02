package sidechain

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/wetee-dao/tee-dsecret/pkg/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

func (app *SideChain) CheckEpochFromValidator() []byte {
	if chains.MainChain == nil {
		util.LogWithGreen("SideChain CheckEpochFromValidator", "error chains.MainChain is nil")
		return nil
	}

	if app.dkg == nil {
		util.LogWithGreen("SideChain CheckEpochFromValidator", "error app.dkg is nil")
		return nil
	}

	// query epoch
	epoch, epochSolt, lastEpochBlock, now, _, err := chains.MainChain.GetEpoch()
	if err != nil {
		return nil
	}

	// query local epoch
	epochStatus := app.GetEpochStatus()

	if now-lastEpochBlock >= epochSolt-1 || epoch == 0 {
		if time.Now().Unix()-int64(epochStatus) > 60 {
			validators, err := chains.MainChain.GetNextEpochValidatorList()
			if err != nil {
				return nil
			}

			util.LogWithYellow("SideChain NewEpoch", "P1 at epoch", epoch)
			err = app.dkg.TryEpochConsensus(model.ConsensusMsg{
				Validators: validators,
				Epoch:      epoch + 1,
			}, app.newEpochSucceded, app.newEpochFail)
			if err == nil {
				return GetTxBytes(&model.Tx{
					Payload: &model.Tx_EpochStatus{
						EpochStatus: time.Now().Unix(), // start epoch, stop submit main chain tx
					},
				})
			}
		}
		return nil
	}

	// Step4 update new epoch to side chain
	if app.GetEpoch() < epoch {
		util.LogWithYellow("SideChain NewEpoch", "P2")
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
			Payload: &model.Tx_Epoch{
				Epoch: &model.Epoch{
					Epoch:      epoch,
					Validators: sideValidators,
				},
			},
		})
	}
	return nil
}

func (app *SideChain) newEpochSucceded(pubkey types.AccountID, sig [64]byte) {
	util.LogWithYellow("SideChain NewEpoch", "P2 submit tx to main chain")
	if chains.MainChain == nil {
		util.LogWithYellow("SideChain CheckEpochFromValidator", "error chains.MainChain is nil")
		return
	}

	// fmt.Println(hex.EncodeToString(pubkey[:]))
	// fmt.Println(hex.EncodeToString(sig[:]))

	// submit new epoch to main chain
	err := chains.MainChain.SetNewEpoch(pubkey, sig)
	if err != nil {
		util.LogWithRed("SideChain next epoch main chain", "error %v", err)
	} else {
		util.LogWithYellow("SideChain next epoch main chain success")
	}

}

func (app *SideChain) newEpochFail(err error) {
	util.LogWithYellow("SideChain NewEpoch", "step2 Error", err.Error())
}

func (app *SideChain) GetEpochStatus() int64 {
	bt, err := model.GetKey("G", "epoch_status")
	if err != nil {
		return 0
	}

	return util.BytesToInt64(bt)
}

func (app *SideChain) SetEpochStatus(status int64) error {
	model.SetKey("G", "epoch_status", util.Int64ToBytes(status))

	return nil
}

func (app *SideChain) GetEpoch() uint32 {
	bt, err := model.GetKey("G", "epoch")
	if err != nil {
		return 0
	}

	bytesBuffer := bytes.NewBuffer(bt)
	var x uint32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return x
}

func (app *SideChain) SetEpoch(epoch *model.Epoch, txn *model.Txn) error {
	// Save new epoch to DKG
	app.dkg.ToNewEpoch()

	// Save new epoch to DB
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, epoch.Epoch)
	txn.SetKey("G", "epoch", bytesBuffer.Bytes())

	// Delete old epoch validators
	err := txn.DeletekeysByPrefix([]byte("G_validator"))
	if err != nil {
		return err
	}

	// save new epoch validators
	for i, v := range epoch.Validators {
		if err = model.TxnSetProtoMessage(txn, []byte("G_validator"+fmt.Sprint(i)), v); err != nil {
			return err
		}
	}

	return nil
}

func (app *SideChain) GetValidators() ([]*model.SideValidator, map[string]*model.PubKey, error) {
	var err error
	list, err := model.GetProtoMessageList[model.SideValidator]("G", "validator")
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

func (app *SideChain) calcValidatorUpdates(epoch *model.Epoch) {
	oldValidators, _, _ := app.GetValidators()
	newValidators := epoch.Validators

	for _, newv := range newValidators {
		isIn := false
		for i, oldv := range oldValidators {
			// 更新
			if bytes.Equal(oldv.Pubkey, newv.Pubkey) {
				oldValidators[i].Power = newv.Power
				isIn = true
			}
		}

		// 新增
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

		// 删除
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
