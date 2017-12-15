package mq

import (
	commontypes "github.com/bytearena/core/common/types"
)

type SubscriptionCallback func(msg BrokerMessage)

type SubscriptionMap struct {
	*commontypes.SyncMap
}

func NewSubscriptionMap() *SubscriptionMap {
	return &SubscriptionMap{
		commontypes.NewSyncMap(),
	}
}

func (smap *SubscriptionMap) Get(id string) SubscriptionCallback {
	if res, ok := (smap.GetGeneric(id)).(SubscriptionCallback); ok {
		return res
	}

	return nil
}
