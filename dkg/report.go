package dkg

import (
	"encoding/json"
	"errors"

	"wetee.app/dsecret/tee"
	types "wetee.app/dsecret/type"
)

func (r *DKG) GetReport(hash string) (*types.TeeParam, *types.TeeReport, error) {
	// 获取数据
	secretData, err := r.GetData(hash)
	if err != nil {
		return nil, nil, err
	}

	teeParam := &types.TeeParam{}
	err = json.Unmarshal(secretData, teeParam)
	if err != nil {
		return nil, nil, errors.New("unmarshal tee param: " + err.Error())
	}

	report, err := tee.VerifyReport(teeParam)
	if err != nil {
		return nil, nil, errors.New("verify cluster report: " + err.Error())
	}

	return teeParam, report, nil
}
