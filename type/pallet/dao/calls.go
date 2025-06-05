package dao

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "github.com/wetee-dao/go-sdk/pallet/types"
)

// Create a DAO
// 从一个通证池,创建一个组织
func MakeCreateDaoCall(name0 []byte, desc1 []byte, purpose2 []byte, metaData3 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsDao: true,
		AsDaoField0: &types.WeteeDaoPalletCall{
			IsCreateDao:          true,
			AsCreateDaoName0:     name0,
			AsCreateDaoDesc1:     desc1,
			AsCreateDaoPurpose2:  purpose2,
			AsCreateDaoMetaData3: metaData3,
		},
	}
}

// init asset for dao
// 初始化 DAO 通证
func MakeInitAssetCall(daoId0 uint64, symbol1 []byte, decimals2 byte, initToken3 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsDao: true,
		AsDaoField0: &types.WeteeDaoPalletCall{
			IsInitAsset:           true,
			AsInitAssetDaoId0:     daoId0,
			AsInitAssetSymbol1:    symbol1,
			AsInitAssetDecimals2:  decimals2,
			AsInitAssetInitToken3: initToken3,
		},
	}
}

// update dao info
// 更新组织信息
func MakeUpdateDaoCall(daoId0 uint64, name1 types.OptionTByteSlice, desc2 types.OptionTByteSlice, purpose3 types.OptionTByteSlice, metaData4 types.OptionTByteSlice, status5 types.OptionTStatus) types.RuntimeCall {
	return types.RuntimeCall{
		IsDao: true,
		AsDaoField0: &types.WeteeDaoPalletCall{
			IsUpdateDao:          true,
			AsUpdateDaoDaoId0:    daoId0,
			AsUpdateDaoName1:     name1,
			AsUpdateDaoDesc2:     desc2,
			AsUpdateDaoPurpose3:  purpose3,
			AsUpdateDaoMetaData4: metaData4,
			AsUpdateDaoStatus5:   status5,
		},
	}
}

// create task
// 创建任务
func MakeCreateRoadmapTaskCall(daoId0 uint64, roadmapId1 uint32, name2 []byte, priority3 byte, tags4 types.OptionTByteSlice) types.RuntimeCall {
	return types.RuntimeCall{
		IsDao: true,
		AsDaoField0: &types.WeteeDaoPalletCall{
			IsCreateRoadmapTask:           true,
			AsCreateRoadmapTaskDaoId0:     daoId0,
			AsCreateRoadmapTaskRoadmapId1: roadmapId1,
			AsCreateRoadmapTaskName2:      name2,
			AsCreateRoadmapTaskPriority3:  priority3,
			AsCreateRoadmapTaskTags4:      tags4,
		},
	}
}

// update task
// 更新任务
func MakeUpdateRoadmapTaskCall(daoId0 uint64, roadmapId1 uint32, taskId2 uint64, priority3 byte, status4 byte, tags5 types.OptionTByteSlice) types.RuntimeCall {
	return types.RuntimeCall{
		IsDao: true,
		AsDaoField0: &types.WeteeDaoPalletCall{
			IsUpdateRoadmapTask:           true,
			AsUpdateRoadmapTaskDaoId0:     daoId0,
			AsUpdateRoadmapTaskRoadmapId1: roadmapId1,
			AsUpdateRoadmapTaskTaskId2:    taskId2,
			AsUpdateRoadmapTaskPriority3:  priority3,
			AsUpdateRoadmapTaskStatus4:    status4,
			AsUpdateRoadmapTaskTags5:      tags5,
		},
	}
}

// create app
// 创建APP
func MakeCreateAppCall(name0 []byte, desc1 []byte, icon2 []byte, url3 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsDao: true,
		AsDaoField0: &types.WeteeDaoPalletCall{
			IsCreateApp:      true,
			AsCreateAppName0: name0,
			AsCreateAppDesc1: desc1,
			AsCreateAppIcon2: icon2,
			AsCreateAppUrl3:  url3,
		},
	}
}

// update app status
// 更新APP状态
func MakeUpdateAppStatusCall(appId0 uint64, status1 types.Status) types.RuntimeCall {
	return types.RuntimeCall{
		IsDao: true,
		AsDaoField0: &types.WeteeDaoPalletCall{
			IsUpdateAppStatus:        true,
			AsUpdateAppStatusAppId0:  appId0,
			AsUpdateAppStatusStatus1: status1,
		},
	}
}

// org integrate app
// 组织集成应用
func MakeOrgIntegrateAppCall(daoId0 uint64, appId1 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsDao: true,
		AsDaoField0: &types.WeteeDaoPalletCall{
			IsOrgIntegrateApp:       true,
			AsOrgIntegrateAppDaoId0: daoId0,
			AsOrgIntegrateAppAppId1: appId1,
		},
	}
}

// 更新APP状态
func MakeUpdateOrgAppStatusCall(daoId0 uint64, appId1 uint64, status2 types.Status) types.RuntimeCall {
	return types.RuntimeCall{
		IsDao: true,
		AsDaoField0: &types.WeteeDaoPalletCall{
			IsUpdateOrgAppStatus:        true,
			AsUpdateOrgAppStatusDaoId0:  daoId0,
			AsUpdateOrgAppStatusAppId1:  appId1,
			AsUpdateOrgAppStatusStatus2: status2,
		},
	}
}
