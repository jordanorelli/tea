package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type response struct {
	OK   bool `json:"ok"`
	Hits int  `json:"hits"`
}

// server implements a clustered http hit-counter server. Each path is given a
// g-counter, and every node in the cluster acts as a read-write replica.
type server struct {
	sync.Mutex

	// my own id
	id int

	// a mapping of peer servers id -> addr
	peers map[int]string

	// distributed counts
	counters map[string]gcounter
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/join":
		s.join(w, r)
	case "/sync":
		s.sync(w, r)
	default:
		s.countHit(w, r)
	}
}

func (s *server) countHit(w http.ResponseWriter, r *http.Request) {
	hits := s.hit(r.URL.Path)

	fmt.Printf("% 8d %s\n", hits, r.URL.Path)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{OK: true, Hits: hits})
}

func (s *server) join(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Addr string `json:"addr"`
	}
}

func (s *server) sync(w http.ResponseWriter, r *http.Request) {
}

func (s *server) hit(path string) int {
	s.Lock()
	defer s.Unlock()

	if s.counters == nil {
		s.counters = map[string]gcounter{
			gcounter{s.id, map[int]int{id: 1}},
		}
		return 1
	}

	c, ok := s.counters[path]
	if !ok {
		s.counters[path] = gcounter{s.id, map[int]int{id: 1}}
		return 1
	}

	c.incr()
	return c.total()
}
