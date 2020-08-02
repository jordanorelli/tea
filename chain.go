package tea

// xchain is a chain of xnodes. an xchain is an execution plan for executing a
// sequence of tests. somwhat ironically the nodes are actually in a slice.
type xchain struct {
	xnodes []*xnode
}
