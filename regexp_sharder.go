// Copyright (c) 2014 Eric Robert. All rights reserved.

package shard

import (
	"fmt"
	"github.com/datacratic/goreports"
	"regexp"
	"sync"
)

type RegexpSharder struct {
	Table      Table
	Expression string
	re         *regexp.Regexp
	mutex      sync.Mutex
}

func (h *RegexpSharder) Shard(req *report.Request) (url string, k int, err error) {
	if h.re == nil {
		h.mutex.Lock()

		// double-checked locking is fine here
		if h.re == nil {
			h.re = regexp.MustCompile(h.Expression)
		}

		h.mutex.Unlock()
	}

	key := h.re.Find(req.Content)
	if len(key) == 0 {
		err = fmt.Errorf("HTTP message doesn't match regular expression '%s'", h.Expression)
		return
	}

	url, k = h.Table.GetURL(key)
	return
}
