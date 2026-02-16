package dao

import (
	"errors"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoAddTrack(caller []byte, p *DaoCallPayload, txn *model.Txn) error {
	st, err := loadDaoState(txn)
	if err != nil || st == nil {
		return errors.New("dao not initialized")
	}
	if !isSudo(caller, st) {
		return errors.New("must call by gov/sudo")
	}
	if p.Track == nil {
		return errors.New("track required")
	}
	state := newDaoStateState(txn)
	trackID := uint32(len(getDaoTracks(txn)))
	return model.SetMappingJson(state.Tracks, txn, trackID, p.Track)
}

func daoSetDefaultTrack(caller []byte, p *DaoCallPayload, txn *model.Txn) error {
	st, err := loadDaoState(txn)
	if err != nil || st == nil {
		return errors.New("dao not initialized")
	}
	if !isSudo(caller, st) {
		return errors.New("must call by gov/sudo")
	}
	st.DefaultTrack = &p.TrackId
	state := newDaoStateState(txn)
	return state.SetDefaultTrack(st.DefaultTrack)
}

// getDaoTracks 返回所有 track（按 index 0,1,2... 直到不存在）。
func getDaoTracks(txn *model.Txn) []*DaoTrackData {
	state := newDaoStateState(txn)
	var out []*DaoTrackData
	for i := uint32(0); ; i++ {
		t, _ := model.GetMappingJson[uint32, DaoTrackData](state.Tracks, txn, i)
		if t == nil {
			break
		}
		out = append(out, t)
	}
	return out
}
