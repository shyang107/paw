package treeprint

import "fmt"

// Org use to PrintTree
type Org struct {
	OrgID    string // 節點編號
	OrgName  string // 節點（欲列印）名稱
	ParentID string // 父節點編號
}

// PrintOrgTree prints tree of Org[]
// 	`name string`: name this
// 	`tbl []Org`: data to print
// 	`parentID string`: from parentID to print
// 	`depth int`: depth to print form parent
func PrintOrgTree(name string, tbl []Org, parentID string, depth int) {
	tree := New()
	pot(tree, tbl, parentID, depth)
	if len(name) > 0 {
		tree.SetValue(name)
	}
	fmt.Println(tree.String())
}
func pot(tree Tree, tbl []Org, parentID string, depth int) {
	for _, r := range tbl {
		if r.ParentID == parentID {
			// one := tree.AddMetaBranch(r.OrgID, r.OrgName+"("+r.ParentID+")")
			one := tree.AddBranch(r.OrgName)
			if depth < 1 {
				return
			}
			pot(one, tbl, r.OrgID, depth-1)
		}
	}
}

func printTreeOrg(tbl []Org, parentID string, depth int, padding string) {
	for id, r := range tbl {
		if r.ParentID == parentID {
			// fmt.Print(padding + "├")
			fmt.Print(padding)
			for i := 0; i < depth; i++ {
				fmt.Print("--")
			}
			fmt.Print(id, r.OrgName, "\n")
			// printTree(tbl, r.OrgID, depth+1, "  "+padding)
			printTreeOrg(tbl, r.OrgID, depth+1, ""+padding)
		}
		// fmt.Print(id)
	}
}
