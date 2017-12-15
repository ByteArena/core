package types

import (
	commontypes "github.com/bytearena/core/common/types"
)

type WatcherMap struct {
	*commontypes.SyncMap
}

func NewWatcherMap() *WatcherMap {
	return &WatcherMap{
		commontypes.NewSyncMap(),
	}
}

func (wmap *WatcherMap) Get(id string) *Watcher {
	if res, ok := (wmap.GetGeneric(id)).(*Watcher); ok {
		return res
	}

	return nil
}
