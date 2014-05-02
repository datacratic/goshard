// Copyright (c) 2014 Eric Robert. All rights reserved.

package shard

import (
	"fmt"
	"hash/fnv"
)

type Table struct {
	DefaultUrl string
	ShardedUrl string
	Shards     int
	Url        map[int]string
}

func (t *Table) GetUrl(id []byte) (url string, i int) {
	h := fnv.New32()
	h.Write(id)
	k := h.Sum32() % uint32(t.Shards)

	i = int(k)
	if s, ok := t.Url[i]; ok {
		url = s
	} else {
		if t.ShardedUrl != "" {
			url = fmt.Sprintf(t.ShardedUrl, i)
		} else {
			url = t.DefaultUrl
		}
	}

	return
}
