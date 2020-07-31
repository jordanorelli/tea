package tea

var lastID int

func nextNodeID() int {
	lastID++
	return lastID
}

type node struct {
	id       int
	test     Test
	name     string
	parents  []*node
	children []*node
}
