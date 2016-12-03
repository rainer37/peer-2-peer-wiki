package article

import (
  "fmt"
)

// // This is a standard binary tree except that each node can contain many sibling
// // nodes.
// type Treedoc struct {
//   MiniNodes []*Node
//   Left *Treedoc
//   Right *Treedoc
// }
//
// // A node in the treedoc. Nodes have a value, path, disambiguator (siteId) and
// // an indicator of whether the node is visible (deleted).
// type Node struct {
//   Value Atom
//   Site Disambiguator
//   Tombstone bool  // true if node has been deleted
//   Left *Treedoc
//   Right *Treedoc
// }
//
// // Nodes are identified in a treedoc by their path and their disambiguator (siteId)
// type PosId struct {
//   Dir Direction
//   Site Disambiguator
// }
//
// // Represents the smallest unit that can be modified atomically
// type Atom string
//
// // A treedoc is a binary tree so a path is a bitstring (represented as an array)
// // starting from the root where a 0 indicates a left branch and a 1 indicates a
// // right branch.
// type Path []PosId
//
// // A globally-unique identifier for the process making the action
// type Disambiguator string



var t Treedoc

// func ExampleInfix() {
//   nodeA := Node{"A", "dA", false, nil, nil}
//   nodeC := Node{"C", "dB", false, nil, nil}
//   t.MiniNodes = append(t.MiniNodes, &nodeA)
//
//   t1 := Treedoc{}
//   t1.MiniNodes = append(t1.MiniNodes, &nodeC)
//   t.Right = &t1
//   var paths []Path
//   var nodes []*Node
//   t.infix(Path{}, &paths, &nodes)
//   fmt.Println(t.Contents())
//   for _,path := range paths {
//     fmt.Println(path.String())
//   }
//   // Output:
// }

func ConstructTreeManually() {
  nodeA := Node{"A", "dA", false, nil, nil}
  nodeC := Node{"C", "dB", false, nil, nil}
  nodeE := Node{"E", "dA", false, nil, nil}
  nodeG := Node{"G", "dA", false, nil, nil}
  nodeH := Node{"H", "dB", false, nil, nil}
  nodeJ := Node{"J", "dB", false, nil, nil}
  nodeD := Node{"D", "dB", false, nil, nil}

  t1 := Treedoc{}
  t1.MiniNodes = append(t1.MiniNodes, &nodeA)

  t2 := Treedoc{}
  t2.MiniNodes = append(t2.MiniNodes, &nodeC)
  nodeB := Node{"B", "dA", false, nil, &t2}

  t3 := Treedoc{}
  t3.MiniNodes = append(t3.MiniNodes, &nodeE)

  t4 := Treedoc{}
  t4.MiniNodes = append(t4.MiniNodes, &nodeG)
  nodeF := Node{"F", "dA", false, nil, &t4}

  t5 := Treedoc{}
  t5.MiniNodes = append(t5.MiniNodes, &nodeH)
  nodeI := Node{"I", "dB", false, &t5, nil}

  t6 := Treedoc{}
  t6.MiniNodes = append(t6.MiniNodes, &nodeJ)


  t7 := Treedoc{}
  t7.MiniNodes = append(t7.MiniNodes, &nodeB, &nodeD)
  t7.Left = &t1
  t7.Right = &t3

    // root
  t.MiniNodes = append(t.MiniNodes, &nodeF, &nodeI)
  t.Left = &t7
  t.Right = &t6
}


func ExampleContents() {
  ConstructTreeManually()
  var paths []Path
  var nodes []*Node
  t.infix(Path{}, &paths, &nodes)
  fmt.Println(t.Contents())
  for _,path := range paths {
    fmt.Println(path.String())
    //fmt.Println(path)
  }
  // Output:

}


// func ExampleWalk() {
//   ConstructTreeManually()
//   pA := Path{PosId{Left, ""}, PosId{Left, ""}}
//   // pB := Path{PosId{Left, "dA"}}
//   // pC := Path{PosId{Left, "dA"}, PosId{Right, ""}}
//   // pD := Path{PosId{Left, "dB"}}
//   // pE := Path{PosId{Left, ""}, PosId{Right, ""}}
//   // pF := Path{PosId{Empty, "dA"}}
//   // pG := Path{PosId{Empty, "dA"}, PosId{Right, ""}}
//   // pH := Path{PosId{Empty, "dB"}, PosId{Left, ""}}
//   // pI := Path{PosId{Empty, "dB"}}
//   // pJ := Path{PosId{Right, ""}}
//   r := t.walk(pA)
//   fmt.Println((*r).MiniNodes[0].Value)
//   fmt.Println(t.traverse())
//
//   //Output:
// }

// func ExampleInsertNode() {
//
//   // pA := Path{PosId{Left, ""}, PosId{Left, ""}, PosId{Left, ""}}
//   // t.insertNode(pA, &Node{"X", "dA", false, nil, nil})
//   // r,_ := t.walk(Path{PosId{Left, ""}, PosId{Left, ""}})
//   // fmt.Printf("%p\n", r) //(*r).Left)
//
//
//   //t.insertNode(Path{PosId{Left, ""}, PosId{Left, ""}}, &Node{"X", "dA", false, nil, nil})
//   t.insertNode(Path{PosId{Empty, ""}}, &Node{"X", "d", false, nil, nil})
//   fmt.Println(t.Contents())
//   // Output:
// }
//
// func ExampleDeleteNode() {
//   p := Path{PosId{Left, ""}, PosId{Left, "dA"}}
//   fmt.Println("Path is",p)
//   t.deleteNode(p)
//   fmt.Println(t.Contents())
//   // Output:
// }
