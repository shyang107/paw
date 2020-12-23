package main

import (
	"fmt"

	"github.com/xlab/treeprint"
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

// func exPrintTree1() {

// 	data := []treeprint.Org{
// 		// {"A001", "Dept1", "0 -----th top"},
// 		{OrgID: "A001", OrgName: "Dept1", ParentID: "0"},
// 		{OrgID: "A011", OrgName: "Dept2", ParentID: "0"},
// 		{OrgID: "A002", OrgName: "subDept1", ParentID: "A001"},
// 		{OrgID: "A005", OrgName: "subDept2", ParentID: "A001"},
// 		{OrgID: "A003", OrgName: "sub_subDept1", ParentID: "A002"},
// 		{OrgID: "A006", OrgName: "gran_subDept", ParentID: "A003"},
// 		{OrgID: "A004", OrgName: "sub_subDept2", ParentID: "A002"},
// 		{OrgID: "A012", OrgName: "subDept1", ParentID: "A011"},
// 	}

// 	treeprint.PrintOrgTree("ORG", data, "0", 3)
// }
