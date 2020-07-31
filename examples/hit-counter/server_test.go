package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/jordanorelli/tea"
)

type testStartServer struct {
	Server *httptest.Server `tea:"save"`
}

func (test *testStartServer) Run(t *testing.T) {
	handler := new(server)
	test.Server = httptest.NewServer(handler)
	t.Logf("started a test server on: %s", test.Server.URL)
}

func (test *testStartServer) After(t *testing.T) {
	t.Logf("closing a test server on: %s", test.Server.URL)
	test.Server.Close()
}

type testRequest struct {
	Server *httptest.Server `tea:"load"`

	path   string
	expect int
}

func (test *testRequest) Run(t *testing.T) {
	client := test.Server.Client()

	res, err := client.Get(test.Server.URL + test.path)
	if err != nil {
		t.Fatalf("request to %s failed: %v", test.path, err)
	}
	defer res.Body.Close()

	var body response
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("response at %s was not json: %v", test.path, err)
	}

	if body.Hits != test.expect {
		t.Errorf("expected a count of %d but saw %d", test.expect, body.Hits)
	}
}

func TestServer(t *testing.T) {
	type list []testRequest

	runSeries := func(node *tea.Tree, tests list) *tea.Tree {
		for i := 0; i < len(tests); i++ {
			node = node.Child(&tests[i])
		}
		return node
	}

	runParallel := func(node *tea.Tree, tests list) {
		for i := 0; i < len(tests); i++ {
			node.Child(&tests[i])
		}
	}

	root := tea.New(&testStartServer{})

	runSeries(root, list{
		{path: "/users/alice", expect: 1},
		{path: "/users/alice", expect: 2},
		{path: "/users/alice", expect: 3},
		{path: "/users/alice", expect: 4},
	})

	runSeries(root, list{
		{path: "/users/bob", expect: 1},
		{path: "/users/alice", expect: 1},
		{path: "/users/alice", expect: 2},
		{path: "/users/alice", expect: 3},
		{path: "/users/bob", expect: 2},
	})

	runParallel(root, list{
		{path: "/users/alice", expect: 1},
		{path: "/users/alice", expect: 1},
		{path: "/users/alice", expect: 1},
		{path: "/users/alice", expect: 1},
	})

	tea.Run(t, root)
}
