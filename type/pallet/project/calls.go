package project

import (
	types1 "github.com/centrifuge/go-substrate-rpc-client/v4/types"
	types "github.com/wetee-dao/go-sdk/pallet/types"
)

// 申请加入团队
func MakeProjectJoinRequestCall(daoId0 uint64, projectId1 uint64, who2 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsProjectJoinRequest:           true,
			AsProjectJoinRequestDaoId0:     daoId0,
			AsProjectJoinRequestProjectId1: projectId1,
			AsProjectJoinRequestWho2:       who2,
		},
	}
}

// 创建项目
func MakeCreateProjectCall(daoId0 uint64, name1 []byte, description2 []byte, creator3 [32]byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsCreateProject:             true,
			AsCreateProjectDaoId0:       daoId0,
			AsCreateProjectName1:        name1,
			AsCreateProjectDescription2: description2,
			AsCreateProjectCreator3:     creator3,
		},
	}
}

// 为项目申请资金
func MakeApplyProjectFundsCall(daoId0 uint64, projectId1 uint64, amount2 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsApplyProjectFunds:           true,
			AsApplyProjectFundsDaoId0:     daoId0,
			AsApplyProjectFundsProjectId1: projectId1,
			AsApplyProjectFundsAmount2:    amount2,
		},
	}
}

// 创建任务
func MakeCreateTaskCall(daoId0 uint64, projectId1 uint64, name2 []byte, description3 []byte, point4 uint16, priority5 byte, maxAssignee6 types.OptionTByte, skills7 types.OptionTByteSlice, assignees8 types.OptionTByteArray32Slice, reviewers9 types.OptionTByteArray32Slice, amount10 types1.U128) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsCreateTask:             true,
			AsCreateTaskDaoId0:       daoId0,
			AsCreateTaskProjectId1:   projectId1,
			AsCreateTaskName2:        name2,
			AsCreateTaskDescription3: description3,
			AsCreateTaskPoint4:       point4,
			AsCreateTaskPriority5:    priority5,
			AsCreateTaskMaxAssignee6: maxAssignee6,
			AsCreateTaskSkills7:      skills7,
			AsCreateTaskAssignees8:   assignees8,
			AsCreateTaskReviewers9:   reviewers9,
			AsCreateTaskAmount10:     amount10,
		},
	}
}

// 加入任务
func MakeJoinTaskCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsJoinTask:           true,
			AsJoinTaskDaoId0:     daoId0,
			AsJoinTaskProjectId1: projectId1,
			AsJoinTaskTaskId2:    taskId2,
		},
	}
}

// 离开项目
func MakeLeaveTaskCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsLeaveTask:           true,
			AsLeaveTaskDaoId0:     daoId0,
			AsLeaveTaskProjectId1: projectId1,
			AsLeaveTaskTaskId2:    taskId2,
		},
	}
}

// 加入项目审核团队
func MakeJoinTaskReviewCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsJoinTaskReview:           true,
			AsJoinTaskReviewDaoId0:     daoId0,
			AsJoinTaskReviewProjectId1: projectId1,
			AsJoinTaskReviewTaskId2:    taskId2,
		},
	}
}

// 离开任务审核
func MakeLeaveTaskReviewCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsLeaveTaskReview:           true,
			AsLeaveTaskReviewDaoId0:     daoId0,
			AsLeaveTaskReviewProjectId1: projectId1,
			AsLeaveTaskReviewTaskId2:    taskId2,
		},
	}
}

// 开始任务
func MakeStartTaskCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsStartTask:           true,
			AsStartTaskDaoId0:     daoId0,
			AsStartTaskProjectId1: projectId1,
			AsStartTaskTaskId2:    taskId2,
		},
	}
}

// 申请审核
func MakeRequestReviewCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsRequestReview:           true,
			AsRequestReviewDaoId0:     daoId0,
			AsRequestReviewProjectId1: projectId1,
			AsRequestReviewTaskId2:    taskId2,
		},
	}
}

// 完成任务
func MakeTaskDoneCall(daoId0 uint64, projectId1 uint64, taskId2 uint64) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsTaskDone:           true,
			AsTaskDoneDaoId0:     daoId0,
			AsTaskDoneProjectId1: projectId1,
			AsTaskDoneTaskId2:    taskId2,
		},
	}
}

// 发送审核报告
func MakeMakeReviewCall(daoId0 uint64, projectId1 uint64, taskId2 uint64, opinion3 types.ReviewOpinion, meta4 []byte) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsMakeReview:           true,
			AsMakeReviewDaoId0:     daoId0,
			AsMakeReviewProjectId1: projectId1,
			AsMakeReviewTaskId2:    taskId2,
			AsMakeReviewOpinion3:   opinion3,
			AsMakeReviewMeta4:      meta4,
		},
	}
}

// 创建非DAO项目
func MakeCreateProxyProjectCall(name0 []byte, description1 []byte, deposit2 types1.UCompact) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsCreateProxyProject:             true,
			AsCreateProxyProjectName0:        name0,
			AsCreateProxyProjectDescription1: description1,
			AsCreateProxyProjectDeposit2:     deposit2,
		},
	}
}
func MakeProxyCallCall(projectId0 uint64, call1 types.RuntimeCall) types.RuntimeCall {
	return types.RuntimeCall{
		IsProject: true,
		AsProjectField0: &types.WeteeProjectPalletCall{
			IsProxyCall:           true,
			AsProxyCallProjectId0: projectId0,
			AsProxyCallCall1:      &call1,
		},
	}
}
