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
func newLNode(test Test, parents ...*lnode) *lnode {
	if len(parents) == 0 {
		return rootLNode(test)
	}

	id := nextNodeID()
	node := lnode{
		id:      id,
		name:    parseName(test),
		test:    test,
		parents: parents,
	}

	xID := 0
	for _, parent := range parents {
		for i, _ := range parent.xnodes {
			node.xnodes = append(node.xnodes, xnode{
				id:     xID,
				lnode:  &node,
				parent: &parent.xnodes[i],
			})
			xID++
		}
	}

	for i, _ := range node.xnodes {
		xparent := node.xnodes[i]
		xparent.children = append(xparent.children, &node.xnodes[i])
	}
	return &node
}

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

// child adds an lnode as a child of the receiver. For every xnode in this
// receiver, the child lnode has a component xnode whose parent is the
// corresponding xnode in this lnode.
func (l *lnode) child(c *lnode, t Test) {
	panic("nuh")
	//c.parents = append(c.parents, l)
	//l.children = append(l.children, c)
	//for i, x := range l.xnodes {
	//	xchild := x.child(t)
	//	xchild.lnode = c
	//	xchild.id = i
	//	c.xnodes = append(c.xnodes, xchild)
	//}
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

// func newXNode(L *lnode

func (x *xnode) child(t Test) *xnode {
	panic("no")
	// child := &xnode{test: t, parent: x}
	// x.children = append(x.children, child)
	// return child
}

func (x *xnode) isOnlyTestInLNode() bool {
	return len(x.lnode.children) == 1
}

func (x *xnode) label() string {
	if x.isOnlyTestInLNode() {
		return x.lnode.name
	}
	switch {
	case len(x.lnode.children) < 10:
		return fmt.Sprintf("%s:%d", x.lnode.name, x.id)
	case len(x.lnode.children) < 100:
		return fmt.Sprintf("%s:%02d", x.lnode.name, x.id)
	case len(x.lnode.children) < 1000:
		return fmt.Sprintf("%s:%03d", x.lnode.name, x.id)
	default:
		return fmt.Sprintf("%s:%04d", x.lnode.name, x.id)
	}
}

// ancestry gives a slice of xnodes beginning at the root of the x graph and
// terminating at the receiver xnode.
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
