package app

import types "github.com/wetee-dao/go-sdk/pallet/types"

// App create
// 注册任务
func MakeCreateCall(name0 []byte, templateId1 types.OptionTU128, image2 []byte, signer3 []byte, signature4 []byte, meta5 []byte, port6 []types.Service, command7 types.Command, env8 []types.EnvInput, secretEnv9 types.OptionTByteSlice, cpu10 uint32, memory11 uint32, disk12 []types.Disk, sideContainer13 []types.Container, level14 byte, teeVersion15 types.TEEVersion) types.RuntimeCall {
	return types.RuntimeCall{
		IsApp: true,
		AsAppField0: &types.WeteeAppPalletCall{
			IsCreate:                true,
			AsCreateName0:           name0,
			AsCreateTemplateId1:     templateId1,
			AsCreateImage2:          image2,
			AsCreateSigner3:         signer3,
			AsCreateSignature4:      signature4,
			AsCreateMeta5:           meta5,
			AsCreatePort6:           port6,
			AsCreateCommand7:        command7,
			AsCreateEnv8:            env8,
			AsCreateSecretEnv9:      secretEnv9,
			AsCreateCpu10:           cpu10,
			AsCreateMemory11:        memory11,
			AsCreateDisk12:          disk12,
			AsCreateSideContainer13: sideContainer13,
			AsCreateLevel14:         level14,
			AsCreateTeeVersion15:    teeVersion15,
		},
	}
}

// App update
// 更新任务
func MakeUpdateCall(appId0 uint64, newName1 types.OptionTByteSlice, newImage2 types.OptionTByteSlice, newSigner3 types.OptionTByteSlice, newSignature4 types.OptionTByteSlice, newPort5 types.OptionTServiceSlice, newCommand6 types.OptionTCommand, newEnv7 []types.EnvInput, secretEnv8 types.OptionTByteSlice, withRestart9 bool) types.RuntimeCall {
	return types.RuntimeCall{
		IsApp: true,
		AsAppField0: &types.WeteeAppPalletCall{
			IsUpdate:              true,
			AsUpdateAppId0:        appId0,
			AsUpdateNewName1:      newName1,
			AsUpdateNewImage2:     newImage2,
			AsUpdateNewSigner3:    newSigner3,
			AsUpdateNewSignature4: newSignature4,
			AsUpdateNewPort5:      newPort5,
			AsUpdateNewCommand6:   newCommand6,
			AsUpdateNewEnv7:       newEnv7,
			AsUpdateSecretEnv8:    secretEnv8,
			AsUpdateWithRestart9:  withRestart9,
		},
	}
}

// App restart
// 更新任务
func MakeRestartCall(appId0 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsApp: true,
		AsAppField0: &types.WeteeAppPalletCall{
			IsRestart:       true,
			AsRestartAppId0: appId0,
		},
	}
}

// update price
// 更新价格
func MakeUpdatePriceCall(level0 byte, price1 types.Price) types.RuntimeCall {
	return types.RuntimeCall{
		IsApp: true,
		AsAppField0: &types.WeteeAppPalletCall{
			IsUpdatePrice:       true,
			AsUpdatePriceLevel0: level0,
			AsUpdatePricePrice1: price1,
		},
	}
}
