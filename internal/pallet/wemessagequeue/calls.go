package wemessagequeue

import types "github.com/wetee-dao/go-sdk/pallet/types"

// Remove a page which has no more messages remaining to be processed or is stale.
func MakeReapPageCall(messageOrigin0 types.MessageOrigin, pageIndex1 uint32) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeMessageQueue: true,
		AsWeMessageQueueField0: &types.WeteeMessageQueuePalletCall{
			IsReapPage:               true,
			AsReapPageMessageOrigin0: messageOrigin0,
			AsReapPagePageIndex1:     pageIndex1,
		},
	}
}

// Execute an overweight message.
//
// Temporary processing errors will be propagated whereas permanent errors are treated
// as success condition.
//
//   - `origin`: Must be `Signed`.
//   - `message_origin`: The origin from which the message to be executed arrived.
//   - `page`: The page in the queue in which the message to be executed is sitting.
//   - `index`: The index into the queue of the message to be executed.
//   - `weight_limit`: The maximum amount of weight allowed to be consumed in the execution
//     of the message.
//
// Benchmark complexity considerations: O(index + weight_limit).
func MakeExecuteOverweightCall(messageOrigin0 types.MessageOrigin, page1 uint32, index2 uint32, weightLimit3 types.Weight) types.RuntimeCall {
	return types.RuntimeCall{
		IsWeMessageQueue: true,
		AsWeMessageQueueField0: &types.WeteeMessageQueuePalletCall{
			IsExecuteOverweight:               true,
			AsExecuteOverweightMessageOrigin0: messageOrigin0,
			AsExecuteOverweightPage1:          page1,
			AsExecuteOverweightIndex2:         index2,
			AsExecuteOverweightWeightLimit3:   weightLimit3,
		},
	}
}
