package main

import (
	"encoding/json"
	"testing"
)

func TestGCounterJSON(t *testing.T) {
	t.Run("marshal", func(t *testing.T) {
		type test struct {
			counter gcounter
			expect  string
		}

		tests := []test{
			{
				gcounter{1, map[int]int{1: 5, 3: 8}},
				`[[1,5],[3,8]]`,
			},
			{
				gcounter{8, map[int]int{3: 10, 7: 16, 12: 9, 2: 37}},
				`[[2,37],[3,10],[7,16],[12,9]]`,
			},
		}

		for _, test := range tests {
			b, err := json.Marshal(test.counter)
			if err != nil {
				t.Errorf("marshal failed: %s", err)
				continue
			}
			s := string(b)
			if s != test.expect {
				t.Errorf("expected json: %s received: %s", test.expect, s)
			}
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		type test struct {
			in     string
			expect gcounter
		}

		tests := []test{
			{
				`[[1,5],[3,8]]`,
				gcounter{1, map[int]int{1: 5, 3: 8}},
			},
			{
				`[[2,37],[3,10],[7,16],[12,9]]`,
				gcounter{8, map[int]int{3: 10, 7: 16, 12: 9, 2: 37}},
			},
		}

		for _, test := range tests {
			var c gcounter
			if err := json.Unmarshal(test.in, &c); err != nil {
				t.Errorf("unmarshal failed: %s", err)
				continue
			}

			for id, count := range c.counts {
				n := test.expect.counts[id]
				if n != count {
					t.Errorf("mismatched counts for id %d: expected %d, saw %d", id, count, n)
				}
			}
		}
	})
}
