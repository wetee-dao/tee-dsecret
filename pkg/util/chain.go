package util

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/wetee-dao/tee-dsecret/chains/pallets/generated/types"
)

// work type to string
func GetWorkTypeStr(work types.WorkId) string {
	if work.Wtype.IsAPP {
		return "s"
	}

	if work.Wtype.IsTASK {
		return "t"
	}

	if work.Wtype.IsGPU {
		return "g"
	}

	return "unknown"
}

// string to work type
func GetWorkType(ty string) types.WorkType {
	if ty == "s" {
		return types.WorkType{IsAPP: true}
	}
	if ty == "t" {
		return types.WorkType{IsTASK: true}
	}
	if ty == "g" {
		return types.WorkType{IsGPU: true}
	}
	return types.WorkType{}
}

func GetWorkTypeFromWorkId(workId string) types.WorkId {
	ws := strings.Split(workId, "::")
	wty := GetWorkType(ws[0])
	num, _ := strconv.ParseUint(ws[1], 10, 64)
	return types.WorkId{
		Wtype: wty,
		Id:    num,
	}
}

func GetWorkIdFromWorkType(wtype types.WorkId) string {
	return GetWorkTypeStr(wtype) + "::" + fmt.Sprint(wtype.Id)
}
