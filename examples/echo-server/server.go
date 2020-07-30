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

type server struct {
	sync.Mutex
	counters map[string]int
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hits := s.hit(r.URL.Path)

	fmt.Printf("% 8d %s\n", hits, r.URL.Path)

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{OK: true, Hits: hits})
}

func (s *server) hit(path string) int {
	s.Lock()
	defer s.Unlock()

	if s.counters == nil {
		s.counters = map[string]int{path: 1}
		return 1
	}

	s.counters[path]++
	return s.counters[path]
}
