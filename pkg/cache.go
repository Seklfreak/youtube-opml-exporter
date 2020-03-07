package pkg

import (
	"sync"
	"time"

	"google.golang.org/api/youtube/v3"
)

const expiry = 24 * time.Hour

type item struct {
	Updated time.Time
	Subs    []*youtube.Subscription
}

var (
	subsCache     = make(map[string]*item)
	subsCacheLock sync.RWMutex
)

func readSubsCache(key string) []*youtube.Subscription {
	subsCacheLock.RLock()
	defer subsCacheLock.RUnlock()

	item := subsCache[key]
	if item == nil {
		return nil
	}

	if !item.Updated.After(time.Now().Add(-expiry)) {
		return nil
	}

	subsCopy := make([]*youtube.Subscription, len(item.Subs))
	copy(subsCopy, item.Subs)
	return subsCopy
}

func storeSubsCache(key string, subs []*youtube.Subscription) {
	subsCacheLock.Lock()
	defer subsCacheLock.Unlock()

	if subs == nil {
		subsCache[key] = nil
		return
	}

	subsCopy := make([]*youtube.Subscription, len(subs))
	copy(subsCopy, subs)
	subsCache[key] = &item{
		Updated: time.Now(),
		Subs:    subsCopy,
	}
}
