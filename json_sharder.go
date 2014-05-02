// Copyright (c) 2014 Eric Robert. All rights reserved.

package shard

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JsonSharder struct {
	Table
	Field string
}

func (h JsonSharder) Shard(content []byte) (url string, k int, err error) {
	item := make(map[string]interface{})

	decoder := json.NewDecoder(bytes.NewBuffer(content))
	err = decoder.Decode(&item)
	if err != nil {
		return
	}

	key, ok := item[h.Field]
	if !ok {
		err = fmt.Errorf("message body is missing a json object with field '%s'", h.Field)
		return
	}

	url, k = h.GetUrl([]byte(key.(string)))
	return
}
