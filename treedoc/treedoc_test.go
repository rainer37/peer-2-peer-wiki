package treedoc

import (
  "testing"
  "fmt"
)

func TestTemp(t *testing.T) {
  empty := &posId{path{}, ""}
  left := &posId{path{false}, ""}
  if empty.before(left) {
    fmt.Println("This is bad!")
  }
}



func TestTreedocInsert(t *testing.T) {
  // tree := &Treedoc{&[]node{node{"c",posId{path{}, "dC"},false}}, nil, nil}
  // tree.Insert("b", 1, "dB")
  // tree.Insert("a", 1, "dA")
  // tree.Insert("d", 4, "dD")
  // tree.Insert("e", 5, "dE")
  // tree.Insert("f", 6, "dF")
  // fmt.Println(tree.Contents())

  t2 := &Treedoc{&[]node{node{"c",posId{path{}, "dC"},false}}, nil, nil}
  fmt.Println(t2.Contents())
  t2.Insert("e", 10, "dE")
  t2.insertNode(&node{"x",posId{path{true}, "dC"},false})
  fmt.Println(t2.Contents())
  t2.Insert("a", 1, "dA")
  fmt.Println(t2.Contents())
  t2.Insert("d", 3, "dD")
  fmt.Println(t2.traverse())
  fmt.Println(t2.Contents())
  t2.Insert("f", 5, "dF")
  fmt.Println(t2.Contents())
  t2.Insert("b", 2, "dB")
  fmt.Println(t2.Contents())

  t2.Delete(3, "")
  fmt.Println(t2.Contents())
}













var paths = []*path {
  &path{},
  &path{false},
  &path{true},
  &path{false,false},
  &path{false,true},
  &path{true,false},
  &path{true,true},
}

func TestPathEquals(t *testing.T) {
  for i,p := range paths {
    if !p.equals(p) {
      t.Error("Expecting p1: %v to equal p2: %v", p, p)
    }
    q := paths[(i+1) % len(paths)]
    if p.equals(q) {
      t.Error("p1: %v should not equal p2: %v", p, q)
    }
  }
}

func TestPathPrefix(t *testing.T) {
  root := paths[0]
  p1 := paths[1]
  p2 := paths[2]

  if !root.prefix(p1) {
    t.Errorf("Expecting root to be prefix of %v.", p1)
  }
  if !root.prefix(p2) {
    t.Errorf("Expecting root to be prefix of %v.", p2)
  }



  // left
  if !p1.prefix(paths[3]) {
    t.Errorf("Expecting %v to be prefix of %v.", p1, paths[3])
  }
  if !p1.prefix(paths[4]) {
    t.Errorf("Expecting %v to be prefix of %v.", p1, paths[4])
  }
  if !p1.prefix(p1) {
    t.Errorf("Expecting %v to be prefix of %v.", p1, p1)
  }

  if p1.prefix(root) {
    t.Errorf("Not expecting %v to be prefix of %v.", p1, root)
  }
  if p1.prefix(paths[5]) {
    t.Errorf("Not expecting %v to be prefix of %v.", p1, paths[5])
  }
  if p1.prefix(paths[6]) {
    t.Errorf("Not expecting %v to be prefix of %v.", p1, paths[6])
  }

  // right
  if !p2.prefix(paths[5]) {
    t.Errorf("Expecting %v to be prefix of %v.", p2, paths[5])
  }
  if !p2.prefix(paths[6]) {
    t.Errorf("Expecting %v to be prefix of %v.", p2, paths[6])
  }
  if !p2.prefix(p2) {
    t.Errorf("Expecting %v to be prefix of %v.", p2, p2)
  }

  if p2.prefix(root) {
    t.Errorf("Not expecting %v to be prefix of %v.", p2, root)
  }
  if p2.prefix(paths[3]) {
    t.Errorf("Not expecting %v to be prefix of %v.", p2, paths[3])
  }
  if p2.prefix(paths[4]) {
    t.Errorf("Not expecting %v to be prefix of %v.", p2, paths[4])
  }
}

func TestPathCommonPrefix(t *testing.T) {
  root := paths[0]
  p1 := paths[1]
  p2 := paths[2]

  if root.commonPrefix(p1) != 0 {
    t.Errorf("Expecting common prefix of root to be 0.")
  }
  if root.commonPrefix(p2) != 0 {
    t.Errorf("Expecting common prefix of root to be 0.")
  }
  if p1.commonPrefix(p2) != 0 {
    t.Errorf("Expecting common prefix of %v and %v to be 0.", p1, p2)
  }
  if p2.commonPrefix(p1) != 0 {
    t.Errorf("Expecting common prefix of %v and %v to be 0.", p2, p1)
  }

  p3 := paths[3]
  p4 := paths[4]
  if p1.commonPrefix(p3) != 1 {
    t.Errorf("Expecting common prefix of %v and %v to 1.", p1, p3)
  }
  if p1.commonPrefix(p4) != 1 {
    t.Errorf("Expecting common prefix of %v and %v to 1.", p1, p4)
  }
  if p3.commonPrefix(p1) != 1 {
    t.Errorf("Expecting common prefix of %v and %v to 1.", p1, p3)
  }
  if p4.commonPrefix(p1) != 1 {
    t.Errorf("Expecting common prefix of %v and %v to 1.", p1, p4)
  }

  p5 := paths[5]
  p6 := paths[6]
  if p1.commonPrefix(p5) != 0 {
    t.Errorf("Expecting common prefix of %v and %v to 0.", p1, p5)
  }
  if p1.commonPrefix(p6) != 0 {
    t.Errorf("Expecting common prefix of %v and %v to 0.", p1, p6)
  }
}


// posId
func TestPosIdBefore(t *testing.T) {
  empty := &posId{path{}, ""}
  left := &posId{path{false}, ""}
  right := &posId{path{true}, ""}
  ids := [][2]*posId{
    // [2]*posId{&posId{path{false}, ""}, &posId{path{}, ""}},
    // [2]*posId{&posId{path{}, ""}, &posId{path{true}, ""}},
    [2]*posId{&posId{path{false}, ""}, &posId{path{true}, ""}},
    [2]*posId{&posId{path{false, true}, ""}, &posId{path{true, false}, ""}},
    [2]*posId{&posId{path{false, true, false}, ""}, &posId{path{false, true, true}, ""}},
    [2]*posId{&posId{path{false, false, false}, ""}, &posId{path{false, true, false}, ""}},
    [2]*posId{&posId{path{false}, "dA"}, &posId{path{false}, ""}},
    [2]*posId{&posId{path{false}, ""}, &posId{path{false}, "dB"}},
    [2]*posId{&posId{path{false}, "dA"}, &posId{path{false}, "dB"}},
  }

  if empty.before(empty) {
    t.Errorf("Expected %v to be ordered before %v.", empty, empty)
  }
  if !left.before(empty) {
    t.Errorf("Expected %v to be ordered before %v.", left, empty)
  }
  if !empty.before(right) {
    t.Errorf("Expected %v to be ordered before %v.", empty, right)
  }

  for _,tcase := range ids {
    if !tcase[0].before(tcase[1]) {
      t.Errorf("Expected %v to be ordered before %v.", tcase[0], tcase[1])
    }
  }

  for _,tcase := range ids {
    if tcase[1].before(tcase[0]) {
      t.Errorf("Expected %v to be ordered before %v.", tcase[1], tcase[0])
    }
  }
}


// node
func TestNodeParent(t *testing.T) {

}
func TestNodeAncestor(t *testing.T) {

}
func TestNodeMiniSibling(t *testing.T) {

}
