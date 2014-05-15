// Copyright (c) 2014 Eric Robert. All rights reserved.

package shard

import (
	"fmt"
	"hash/fnv"
)

type Table struct {
	DefaultURL string
	ShardedURL string
	Shards     int
	URLs       map[string]string
}

func (t *Table) GetURL(id []byte) (url string, i int) {
	h := fnv.New32()
	h.Write(id)
	k := h.Sum32() % uint32(t.Shards)

	i = int(k)
	if s, ok := t.URLs[fmt.Sprintf("%d", i)]; ok {
		url = s
	} else {
		if t.ShardedURL != "" {
			url = fmt.Sprintf(t.ShardedURL, i)
		} else {
			url = t.DefaultURL
		}
	}

	return
}
