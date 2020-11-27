package main

import "github.com/shyang107/paw/treeprint"

func exPrintTree() {
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
