package tea

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
	xnodes   []*xnode
	parents  []*lnode
	children []*lnode
}

// child adds an lnode as a child of the receiver. For every xnode in this
// lnode, the child lnode has a component xnode whose parent is the
// corresponding xnode in this lnode.
func (l *lnode) child(c *lnode, t Test) {
	c.parents = append(c.parents, l)
	l.children = append(l.children, c)
	for _, x := range l.xnodes {
		xchild := x.child(t)
		xchild.lnode = c
		c.xnodes = append(c.xnodes, xchild)
	}
}

// xnode is a node in the execution graph, representing one instance of a test
// to be executed. xnode is the unit test in tea. every xnode is either
// unparented or has one parent.
type xnode struct {
	lnode    *lnode
	test     Test
	parent   *xnode
	children []*xnode
}

func (x *xnode) child(t Test) *xnode {
	child := &xnode{test: t, parent: x}
	x.children = append(x.children, child)
	return child
}
