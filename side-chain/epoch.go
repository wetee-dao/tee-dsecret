package sidechain

import (
	"bytes"
	"encoding/binary"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"
	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/wetee-dao/tee-dsecret/chains"
	"github.com/wetee-dao/tee-dsecret/pkg/model"
	"github.com/wetee-dao/tee-dsecret/pkg/util"
)

func (app *SideChain) CheckEpochFromValidator() {
	if chains.MainChain == nil {
		util.LogWithGreen("SideChain CheckEpochFromValidator", "error chains.MainChain is nil")
		return
	}

	if app.dkg == nil {
		util.LogWithGreen("SideChain CheckEpochFromValidator", "error app.dkg is nil")
		return
	}

	epoch, lastEpochBlock, now, err := chains.MainChain.GetEpoch()
	if err != nil {
		return
	}

	// util.LogWithGreen("SideChain CheckEpoch", "epoch:", epoch, "lastEpochBlock:", lastEpochBlock, "now:", now)
	if epoch > app.dkg.Epoch || now-lastEpochBlock >= 72000 || epoch <= 1 {
		validators, err := chains.MainChain.GetValidatorList()
		if err != nil {
			return
		}

		if app.GetEpoch() < epoch {
			sideValidators := make([]*model.SideValidator, 0, len(validators))
			for _, v := range validators {
				sideValidators = append(sideValidators, &model.SideValidator{
					Pubkey: v.ValidatorId.PublicKey,
					Power:  1,
				})
			}

			// Submit new epoch
			go SubmitTx(&model.Tx{
				Payload: &model.Tx_Epoch{
					Epoch: &model.Epoch{
						Epoch:      epoch,
						Validators: sideValidators,
					},
				},
			})
		}

		// util.LogWithGreen("SideChain GetValidatorList", "validators:", validators)
		app.dkg.TryConsensus(model.ConsensusMsg{
			Validators: validators,
			Epoch:      epoch,
		})
	}
}

func (app *SideChain) GetEpoch() uint32 {
	bt, err := model.Get("G_epoch")
	if err != nil {
		return 0
	}

	bytesBuffer := bytes.NewBuffer(bt)
	var x uint32
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return x
}

func (app *SideChain) SetEpoch(epoch *model.Epoch, txn *model.Txn) error {
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, epoch.Epoch)
	txn.Set([]byte("G_epoch"), bytesBuffer.Bytes())

	err := txn.DeletekeysByPrefix([]byte("G_validator"))
	if err != nil {
		return err
	}

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

func (app *SideChain) saveValidators(vs []abcitypes.ValidatorUpdate) error {
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

	oldMap := make(map[string]int64, len(oldValidators)) // map[pubkey]power
	newMap := make(map[string]int64, len(newValidators))
	for _, v := range oldValidators {
		pub := model.PubKeyFromByte(v.Pubkey)
		oldMap[pub.SS58()] = v.Power
	}
	for _, v := range oldValidators {
		pub := model.PubKeyFromByte(v.Pubkey)
		newMap[pub.SS58()] = v.Power
	}

	vups := make([]abci.ValidatorUpdate, 0, len(oldValidators))

	// 找新增（包括新加和 power 改变的）
	for pubkey, newPower := range newMap {
		oldPower, exists := oldMap[pubkey]
		if !exists || oldPower != newPower {
			_, pubkey, _ := model.SS58Decode(pubkey)
			vups = append(vups, abci.ValidatorUpdate{
				PubKeyBytes: pubkey,
				PubKeyType:  "ed25519",
				Power:       newPower,
			})
		}
	}

	// 找删除
	for pubkey := range oldMap {
		_, exists := newMap[pubkey]
		if !exists {
			_, pubkey, _ := model.SS58Decode(pubkey)
			vups = append(vups, abci.ValidatorUpdate{
				PubKeyBytes: pubkey,
				PubKeyType:  "ed25519",
				Power:       0,
			})
		}
	}

	app.onGoingValidators = vups
}
