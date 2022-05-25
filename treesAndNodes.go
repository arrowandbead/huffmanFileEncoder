package main

type Node struct {
	lChild *Node
	rChild *Node
	parent *Node
	idNum  int
	isLeaf bool
	val    byte
	weight int
}

type Tree struct {
	root *Node
}

func solveTree(tree *Tree, counts map[byte]int) {

}

func encodeBytesWithTree(tree *Tree, bytes []byte) {

}
