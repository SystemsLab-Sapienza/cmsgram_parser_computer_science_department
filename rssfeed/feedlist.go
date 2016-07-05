package rssfeed

import (
	"sync"

	"github.com/SlyMarbo/rss"
)

// Extends the type rss.Feed to add an ID field
type Feed struct {
	ID int
	*rss.Feed
}

// A list of feeds that can be safely accessed concurrently
type FeedList struct {
	sync.Mutex
	feeds []Feed
}

func (f *FeedList) Add(ID int, feed *rss.Feed) {
	f.Lock()
	f.feeds = append(f.feeds, Feed{ID, feed})
	f.Unlock()
}
