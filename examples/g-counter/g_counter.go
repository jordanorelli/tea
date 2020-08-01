package main

import (
	"encoding/json"
	"sort"
)

type gcounter struct {
	id     int
	counts map[int]int
}

func (c gcounter) incr() { c.counts[c.id]++ }

func (c gcounter) total() int {
	var n int
	for _, count := range c.counts {
		n += count
	}
	return n
}

func (c gcounter) merge(other gcounter) {
	for id, count := range other.counts {
		if c.[id] < count {
			c.[id] = count
		}
	}
}

type pair [2]int

func (c gcounter) MarshalJSON() ([]byte, error) {
	pairs := make([]pair, 0, len(c.counts))
	for id, count := range c.counts {
		pairs = append(pairs, pair{id, count})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i][0] < pairs[j][0] })
	return json.Marshal(pairs)
}

func (c *gcounter) UnmarshalJSON(b []byte) error {
	pairs := make([]pair, 0, len(c.counts))
	if err := json.Unmarshal(b, &pairs); err != nil {
		return err
	}
}
