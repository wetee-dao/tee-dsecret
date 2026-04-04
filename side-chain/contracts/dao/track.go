package dao

import "github.com/wetee-dao/ink.go/util"

func (d DaoQuery) Tracks() ([]TrackData, error) {
	return d.tracks.List(d.api.GetTxn())
}

func (d DaoQuery) Track(id uint32) (util.Option[TrackData], error) {
	track, err := d.tracks.Get(d.api.GetTxn(), id)
	if err != nil {
		return util.NewNone[TrackData](), err
	}
	if track == nil {
		return util.NewNone[TrackData](), ErrNoTrack
	}
	return util.NewSome(*track), nil
}

func (d DaoQuery) DefaultTrack() (*uint32, error) {
	return d.defaultTrack.Get(d.api.GetTxn(), singletonKey)
}

func (d DaoMutation) AddTrack(track TrackData) error {
	if err := d.ensureGov(); err != nil {
		return err
	}
	next, err := d.nextTrackID()
	if err != nil {
		return err
	}
	if err := d.tracks.Set(d.api.GetTxn(), next, track); err != nil {
		return err
	}
	return d.nextTrackIDStore.Set(d.api.GetTxn(), singletonKey, next+1)
}

func (d DaoMutation) SetDefaultTrack(trackID uint32) error {
	if err := d.ensureGov(); err != nil {
		return err
	}
	track, err := d.tracks.Get(d.api.GetTxn(), trackID)
	if err != nil {
		return err
	}
	if track == nil {
		return ErrNoTrack
	}
	return d.defaultTrack.Set(d.api.GetTxn(), singletonKey, trackID)
}
