package main

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/jordanorelli/tea"
)

// testStartServer is a test that checks that the server can start. If this
// test passes, we retain a reference to the server for future tests to
// utilize.
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

// testHits sends a request to a hitcount server created in a previous test,
// checking that the number of hits returned matches what we expect.
type testHits struct {
	Server *httptest.Server `tea:"load"`

	path string
	hits int
}

func (test *testHits) Run(t *testing.T) {
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

	if body.Hits != test.hits {
		t.Errorf("expected a count of %d hits but saw %d instead", test.hits, body.Hits)
	}
}

func TestServer(t *testing.T) {
	// start with a root node that creates our test server
	root := tea.New(&testStartServer{})

	// add a child node: this test is run if the root test passes. If the root
	// test is failed, this test and all of its descendents are logged as
	// skipped.
	one := root.Child(&testHits{path: "/alice", hits: 1})

	// the effects of the first test create the initial state for the second test.
	two := one.Child(&testHits{path: "/alice", hits: 2})

	// since we have never visited /bob, we know that bob should only have one hit.
	two.Child(&testHits{path: "/bob", hits: 1})

	// but we could also run the exact same test off of the root, like so:
	root.Child(&testHits{path: "/bob", hits: 1})

	// since tests are values in tea, we can re-use the exact same test from
	// different initial states by saving the test as a variable.
	bob := &testHits{path: "/bob", hits: 1}

	// these two executions of the same test value are operating on different
	// program states. Since they are not in the same sequence, they have no
	// effect on one another, even though they're utilizing the same test
	// value.
	two.Child(bob)
	root.Child(bob)

	tea.Run(t, root)
}
