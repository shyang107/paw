package filetree

import (
	"github.com/shyang107/paw/treeprint"
)

type node struct {
	Root  *node
	Meta  treeprint.MetaValue
	Value treeprint.Value
	Nodes []*node
}

// // AddNode adds a new node to a branch.
// func (n *node) AddNode(v treeprint.Value) treeprint.Tree {

// }

// // AddMetaNode adds a new node with meta value provided to a branch.
// AddMetaNode(meta MetaValue, v Value) Tree
// // AddBranch adds a new branch node (a level deeper).
// AddBranch(v Value) Tree
// // AddMetaBranch adds a new branch node (a level deeper) with meta value provided.
// AddMetaBranch(meta MetaValue, v Value) Tree
// // Branch converts a leaf-node to a branch-node,
// // applying this on a branch-node does no effect.
// Branch() Tree
// // FindByMeta finds a node whose meta value matches the provided one by reflect.DeepEqual,
// // returns nil if not found.
// FindByMeta(meta MetaValue) Tree
// // FindByValue finds a node whose value matches the provided one by reflect.DeepEqual,
// // returns nil if not found.
// FindByValue(value Value) Tree
// //  returns the last node of a tree
// FindLastNode() Tree
// // String renders the tree or subtree as a string.
// String() string
// // Bytes renders the tree or subtree as byteslice.
// Bytes() []byte

// SetValue(value Value)
// SetMetaValue(meta MetaValue)
// func New() treeprint.Tree {
// 	return &node{Value: "."}
// }
