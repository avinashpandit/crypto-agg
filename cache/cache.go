package cache

import (
	"bytes"
	"context"
	"encoding/gob"
	"sync"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/avinashpandit/crypto-agg/exchange"
)

var once sync.Once
var cache *bigcache.BigCache

func InitCache() *bigcache.BigCache {
	once.Do(func() {
		config := bigcache.Config{
			// number of shards (must be a power of 2)
			Shards: 1024,

			// time after which entry can be evicted
			LifeWindow: 10 * time.Minute,

			// Interval between removing expired entries (clean up).
			// If set to <= 0 then no action is performed.
			// Setting to < 1 second is counterproductive — bigcache has a one second resolution.
			CleanWindow: 5 * time.Minute,

			// rps * lifeWindow, used only in initial memory allocation
			MaxEntriesInWindow: 1000 * 10 * 60,

			// max entry size in bytes, used only in initial memory allocation
			MaxEntrySize: 500,

			// prints information about additional memory allocation
			Verbose: true,

			// cache will not allocate more memory than this limit, value in MB
			// if value is reached then the oldest entries can be overridden for the new ones
			// 0 value means no size limit
			HardMaxCacheSize: 8192,

			// callback fired when the oldest entry is removed because of its expiration time or no space left
			// for the new entry, or because delete was called. A bitmask representing the reason will be returned.
			// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
			OnRemove: nil,

			// OnRemoveWithReason is a callback fired when the oldest entry is removed because of its expiration time or no space left
			// for the new entry, or because delete was called. A constant representing the reason will be passed through.
			// Default value is nil which means no callback and it prevents from unwrapping the oldest entry.
			// Ignored if OnRemove is specified.
			OnRemoveWithReason: nil,
		}

		cache, _ = bigcache.New(context.Background(), config)

	})

	return cache
}

func GetQuote(exchange1 string, pair string) (*exchange.Quote, *exchange.Quote, error) {
	// Get the value in the byte format it is stored in
	valueBytes, err := cache.Get(exchange1 + "_" + pair)
	if err != nil {
		return nil, nil, err
	}

	// Deserialize the bytes of the value
	value, err := deserialize(valueBytes)
	if err != nil {
		return nil, nil, err
	}

	//get array
	val := value.([]interface{})

	bid := val[0].(exchange.Quote)
	ask := val[1].(exchange.Quote)

	return &bid, &ask, nil

}

func SetQuote(exchange string, pair string, bid *exchange.Quote, ask *exchange.Quote) bool {
	//fist get and see if same Quote
	bid1, ask1, err := GetQuote(exchange, pair)
	if err == nil {
		if bid1.Rate == bid.Rate && ask1.Rate == ask.Rate {
			// rates are same so no need to update
			return false
		}
	}

	// Assert the key is of string type
	keyString := exchange + "_" + pair

	// Serialize the value into bytes
	value := []interface{}{bid, ask}
	valueBytes, err := serialize(value)
	if err != nil {
		return false
	}

	if cache.Set(keyString, valueBytes) != nil {
		return false
	}
	return true
}

func serialize(value interface{}) ([]byte, error) {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	gob.Register(value)
	gob.Register(exchange.Quote{})

	err := enc.Encode(&value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func deserialize(valueBytes []byte) (interface{}, error) {
	var value interface{}
	buf := bytes.NewBuffer(valueBytes)
	dec := gob.NewDecoder(buf)

	err := dec.Decode(&value)
	if err != nil {
		return nil, err
	}

	return value, nil
}
