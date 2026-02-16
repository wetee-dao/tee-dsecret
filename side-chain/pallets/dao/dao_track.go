package dao

import (
	"errors"
	"strconv"

	"github.com/wetee-dao/tee-dsecret/pkg/model"
)

func daoAddTrack(caller []byte, m *model.DaoAddTrack, txn *model.Txn) error {
	state := newDaoStateState(txn)
	if len(state.Members()) == 0 {
		return errors.New("dao not initialized")
	}
	if !isSudo(caller, state) {
		return errors.New("must call by gov/sudo")
	}
	if m.GetTrack() == nil {
		return errors.New("track required")
	}
	trackID := uint32(len(getDaoTracks(txn)))
	return model.SetMappingJson(state.Tracks, txn, trackID, trackFromProto(m.Track))
}

func daoSetDefaultTrack(caller []byte, m *model.DaoSetDefaultTrack, txn *model.Txn) error {
	state := newDaoStateState(txn)
	if len(state.Members()) == 0 {
		return errors.New("dao not initialized")
	}
	if !isSudo(caller, state) {
		return errors.New("must call by gov/sudo")
	}
	trackId := m.GetTrackId()
	return state.SetDefaultTrack(&trackId)
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

// DefaultTrack 返回默认轨道 ID（对应 ink defalut_track()）。
func DefaultTrack() *uint32 {
	v, _ := model.GetJson[uint32]("dao", stateKeydaoStateDefaultTrack)
	return v
}

// TrackList 分页返回轨道列表（对应 ink track_list(page, size)）。
func TrackList(page, size uint32) []*DaoTrackData {
	all := getDaoTracksFromDB()
	if page < 1 {
		page = 1
	}
	start := (page - 1) * size
	if start >= uint32(len(all)) {
		return nil
	}
	end := start + size
	if end > uint32(len(all)) {
		end = uint32(len(all))
	}
	return all[start:end]
}

// Track 按 ID 返回轨道（对应 ink track(id)）。
func Track(id uint32) *DaoTrackData {
	key := stateKeydaoStateTracks + strconv.FormatUint(uint64(id), 10)
	t, _ := model.GetJson[DaoTrackData]("dao", key)
	return t
}

func getDaoTracksFromDB() []*DaoTrackData {
	list, _, _ := model.GetJsonList[DaoTrackData]("dao", stateKeydaoStateTracks)
	return list
}
