package tea

import "fmt"

var lastID int

func nextNodeID() int {
	lastID++
	return lastID
}

// lnode is a node in the logical graph. Developers create a logical graph in
// which nodes may have more than one parent. Each test value written by a
// developer appears as one logical node in the logical graph. The public
// documentation refers to the logical graph as simply the test graph.
type lnode struct {
	id       int
	name     string
	xnodes   []xnode
	test     Test
	parents  []*lnode
	children []*lnode
}

// newLNode greates a new lnode with the provided test and the list of parents.
// If there are no parents, the lnode is a root node, and the provided test is
// run once. Otherwise, the provided test is used to create xnodes that are
// children of all of the provided parent nodes' xnodes.
func newLNode(test Test, sel Selection) *lnode {
	if len(sel.nodes) == 0 {
		return rootLNode(test)
	}

	node := lnode{
		id:      nextNodeID(),
		name:    parseName(test),
		test:    test,
		parents: make([]*lnode, len(sel.nodes)),
	}
	// not sure if this copy is necessary
	copy(node.parents, sel.nodes)

	xID := 0
	for _, parent := range node.parents {
		parent.children = append(parent.children, &node)
		for i, _ := range parent.xnodes {
			x := xnode{
				id:     xID,
				lnode:  &node,
				parent: &parent.xnodes[i],
			}
			node.xnodes = append(node.xnodes, x)
			xID++
		}
	}

	for i, x := range node.xnodes {
		x.parent.children = append(x.parent.children, &node.xnodes[i])
	}
	return &node
}

// rootLNode creates a root lnode. This case is a lot simpler so I split it out
// to keep newLNode a little more readable.
func rootLNode(test Test) *lnode {
	id := nextNodeID()
	node := lnode{
		id:   id,
		name: parseName(test),
		test: test,
	}
	node.xnodes = []xnode{{id: 0, lnode: &node}}
	return &node
}

// xnode is a node in the execution graph, representing one instance of a test
// to be executed. xnode is the unit test in tea. every xnode is either
// unparented or has one parent.
type xnode struct {
	id       int    // id within the parent lnode
	lnode    *lnode // corresponding node in the logical test graph
	parent   *xnode
	children []*xnode
}

func (x *xnode) isOnlyTestInLNode() bool {
	return len(x.lnode.xnodes) == 1
}

// label must be unique or some other shit will break, I'm using this as a way
// to globally identify xnodes, which may be very flawed and maybe I should
// have an actual global ID system.
func (x *xnode) label() string {
	if x.parent == nil {
		switch {
		case len(x.lnode.children) < 10:
			return fmt.Sprintf("%s.%d", x.lnode.name, x.id)
		case len(x.lnode.children) < 100:
			return fmt.Sprintf("%s.%02d", x.lnode.name, x.id)
		case len(x.lnode.children) < 1000:
			return fmt.Sprintf("%s.%03d", x.lnode.name, x.id)
		default:
			return fmt.Sprintf("%s.%04d", x.lnode.name, x.id)
		}
	} else {
		switch {
		case len(x.lnode.children) < 10:
			return fmt.Sprintf("%s.%d.%s", x.lnode.name, x.id, x.parent.lnode.name)
		case len(x.lnode.children) < 100:
			return fmt.Sprintf("%s.%02d.%s", x.lnode.name, x.id, x.parent.lnode.name)
		case len(x.lnode.children) < 1000:
			return fmt.Sprintf("%s.%03d.%s", x.lnode.name, x.id, x.parent.lnode.name)
		default:
			return fmt.Sprintf("%s.%04d.%s", x.lnode.name, x.id, x.parent.lnode.name)
		}
	}
}

// ancestry gives a slice of xnodes beginning at the root of the x graph and
// terminating at the receiver xnode. The ancestry list of a leaf node in the x
// graph is a single chain of tests.
func (x *xnode) ancestry() []*xnode {
	if x.parent == nil {
		return []*xnode{x}
	}
	return append(x.parent.ancestry(), x)
}

// descendents gives a slice of all xnodes whose ancestry includes the receiver
// xnode, in depth-first order.
func (x *xnode) descendents() []*xnode {
	if len(x.children) == 0 {
		return nil
	}

	descendents := make([]*xnode, 0, len(x.children))
	for _, c := range x.children {
		descendents = append(descendents, c)
		descendents = append(descendents, c.descendents()...)
	}
	return descendents
}

// leaves descends the x graph from the receiver xnode, returning a slice
// containing all of the leaves of the x graph having the receiver x as an
// ancestor.
func (x *xnode) leaves() []*xnode {
	if len(x.children) == 0 {
		return []*xnode{x}
	}

	var leaves []*xnode
	for _, child := range x.children {
		leaves = append(leaves, child.leaves()...)
	}

	return leaves
}
