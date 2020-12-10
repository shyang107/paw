package main

import (
	"fmt"

	"github.com/shyang107/paw/treeprint"
)

func exPrintTree2() {
	tree := treeprint.New()

	// create a new branch in the root
	one := tree.AddBranch("one")

	// add some nodes
	one.AddNode("subnode1").AddNode("subnode2")

	// create a new sub-branch
	one.AddBranch("two").
		AddNode("subnode1").AddNode("subnode2"). // add some nodes
		AddBranch("three").                      // add a new sub-branch
		AddNode("subnode1").AddNode("subnode2")  // add some nodes too

	// add one more node that should surround the inner branch
	one.AddNode("subnode3")

	// add a new node to the root
	tree.AddNode("outernode")

	fmt.Println(tree.String())
}

func exPrintTree1() {
	data := []treeprint.Org{
		// {"A001", "Dept1", "0 -----th top"},
		{"A001", "Dept1", "0"},
		{"A011", "Dept2", "0"},
		{"A002", "subDept1", "A001"},
		{"A005", "subDept2", "A001"},
		{"A003", "sub_subDept1", "A002"},
		{"A006", "gran_subDept", "A003"},
		{"A004", "sub_subDept2", "A002"},
		{"A012", "subDept1", "A011"},
	}

	treeprint.PrintOrgTree("ORG", data, "0", 3)
}
