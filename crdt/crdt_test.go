package crdt

import (
  "testing"
  "fmt"
  //"reflect"
)


// var treeElem = []*node{
// //  &node{"A", "c", false},
//   &node{"dB", "b", false},
//   &node{"dA", "a", false},
//   &node{"A", "e", true},
//   &node{"A", "d", false},
//   &node{"A", "f", false},
//
//   //"c","b","a","e","d","f"
// }

// var treeElem = []*node{
// //  &node{"A", "c", false},
//   &node{"dB", "b", false},
//   &node{"dA", "a", false},
//   &node{"A", "g", true},
//   &node{"A", "e", false},
//   &node{"A", "d", false},
//   &node{"A", "f", false},
//   &node{"A", "h", false},
//   //"c","b","a","e","d","f"
// }


var treeElem = []*node{
//  &node{"A", "c", false},
  &node{"dB", "b", false},
  &node{"dA", "a", false},
  &node{"A", "c", false},
  &node{"A", "n", false},
  &node{"A", "o", false},
  &node{"A", "m", false},
  &node{"A", "p", false},
  //"c","b","a","e","d","f"
}


var paths = [][]bool {
  //[]bool{false},
  // []bool{false,false},
  // []bool{true},
  // []bool{true,false},
  // []bool{true,true},
}

// func TestTreedoc(t *testing.T) {
//   //fmt.Printf("treeElem is %v\n", treeElem)
//   //tree := &Treedoc{nil, []*node{treeElem[0]}, nil}
//   //tree := &Treedoc{nil, append(make([]*node,1), treeElem[0]), nil}
//   tree := NewTreedoc("dC", "g")
//
//   for _,e := range treeElem {
//       //fmt.Printf("Tree elem is %v\n", e)
//       tree = tree.Insert(e)
//   }
//   //var summary []nodeSummary
//   //x := tree.infix(&summary)
//   //fmt.Println(x)
//   fmt.Println(tree.String())
//   //tree.Print()
//   //fmt.Println(tree.infix())
//   fmt.Println(tree.size())
//
//
//   var prefix = []int{}
//   var path = []nodeSummary{}
//   fmt.Println(path)
//   tree.infixPath(prefix, &path)
//   fmt.Println(path)
//   tree.Delete(8)
//   fmt.Println(tree.Content())
// }

func TestTreedocInsertPos(t *testing.T) {
  s := NewTreedoc("dG", "g")
  s.insertPos(&node{"dB", "b", false},posId{[]int{0},"dB"})
  s.insertPos(&node{"dE", "e", false},posId{[]int{1},"dE"})
  s.insertPos(&node{"dA", "a", false},posId{[]int{0,0},"dA"})
  fmt.Println(s.String())
}

func TestTreedocInsert(t *testing.T) {
  s := NewTreedoc("dG", "g")
  fmt.Println(s.String())
  s.Insert("b", 1)
  fmt.Println(s.String())
  s.Insert("a", 1)
  //s.Insert("s",3)
  fmt.Println(s.String())
}
/*
func ExampleTreedoc() {
  tree := Treedoc{}
  for _,e := range treeElem {
    tree = *tree.Insert(e)
  }

  fmt.Println(tree.ToString())
  // Output: abcdef
}
*/
/*
func TestTreedocInfixPath(t *testing.T) {
  tree := Treedoc{}
  for _,e := range treeElem {
    tree = *tree.Insert(e)
  }

  var p = []bool{}
  tree.infixPath(&p)
  fmt.Println(p)



  var l = []bool{}
  var r = []bool{}
  tree.left.infixPath(&l) // prepend false
  tree.right.infixPath(&r) // prepend true
  fmt.Println(l)
  fmt.Println(r)

  // output_pos - # lnodes = offset
  // if output_pos > # lnodes then use r else use l
}


func ExampleTreedoc() {
  tree := Treedoc{}
  for _,e := range treeElem {
    tree = *tree.Insert(e)
  }

  // t = *t.Insert("a")
  // t = *t.Insert("b")
  // t = *t.Insert("c")
  // t = *t.Insert("d")
  // t = *t.Insert("e")
  // t = *t.Insert("f")
  fmt.Println(tree.ToString())
  // Output: abcdef
}

func TestTreedocHeight(t *testing.T) {
  tree := Treedoc{}
  for _,e := range treeElem {
    tree = *tree.Insert(e)
  }


  h := tree.left.height()
  fmt.Println(h)
}

func TestTreedocPath(t *testing.T) {
  tree := Treedoc{}
  for _,e := range treeElem {
    tree = *tree.Insert(e)
  }

  var p = []bool{}
  tree.Path("f", &p)
  fmt.Println(p)
}


func TestTreedocFind(t *testing.T) {

  // var d Treedoc
  //
  // fmt.Println(d)
  //
  // d = *d.Insert("c")
  // d = *d.Insert("b")
  // d = *d.Insert("a")
  // d = *d.Insert("e")
  // d = *d.Insert("d")
  // d = *d.Insert("f")
  //
  // fmt.Println(d)
  //
  // fmt.Println(tree.ToString())
  // fmt.Println(tree)
  // tree.Insert("b")
  // fmt.Println(tree.ToString())
  // fmt.Println(tree)
  //
  // tree.Insert("a")
  // fmt.Println(tree.ToString())
  // fmt.Println(tree)
  // tree.Insert("e")
  // fmt.Println(tree.ToString())
  // tree.Insert("d")
  // fmt.Println(tree.ToString())
  // tree.Insert("f")
  // fmt.Println(tree.ToString())

  tree := Treedoc{}
  for _,e := range treeElem {
    tree = *tree.Insert(e)
  }
  fmt.Println(tree)

  pos,err := tree.Find("c", 5)
  if err != nil {
    t.Error("Not expecting error.", err)
  } else if !reflect.DeepEqual(pos,[]bool{false,false})  {
    t.Error("Expecting posId to be [00], got %v", pos)
  }

  pos,err = tree.Find("f", 5)
  if err != nil {
    t.Error("Not expecting error.", err)
  } else if !reflect.DeepEqual(pos,[]bool{true,true})  {
    t.Error("Expecting posId to be [00], got %v", pos)
  }
}
*/



// Test the dynamic vector clock
func TestEquals(t *testing.T) {
  v := DynamicVectorClock{}
  w := DynamicVectorClock{}

  // check that nil maps are trivially equal
  if !v.Equals(w) {
    t.Error("Expected v to equal w.")
  }

  v.Append("a")
  if v.Equals(w) {
    t.Error("Expected v to not equal w after appending key 'a' to v.")
  }

  w.Append("a")
  if !v.Equals(w) {
    t.Error("Expected v to eqaul w after appending key 'a' to w.")
  }

  v.Increment("a")
  if v.Equals(w) {
    t.Error("Expected v to not equal w after incrementing key 'a' of v.")
  }

  w.Increment("a")
  if !v.Equals(w) {
    t.Error("Expected v to equal w after incrementing key 'a' of w.")
  }

  v.Append("b")
  if v.Equals(w) {
    t.Error("Expected v to not equal w after appending key 'b' to v.")
  }

  w.Append("b")
  if !v.Equals(w) {
    t.Error("Expected v to eqaul w after appending key 'b' to w.")
  }

  v.Increment("b")
  if v.Equals(w) {
    t.Error("Expected v to not equal w after incrementing key 'b' of v.")
  }

  w.Increment("b")
  if !v.Equals(w) {
    t.Error("Expected v to equal w after incrementing key 'b' of w.")
  }
}


func TestBefore(t *testing.T) {
  v := DynamicVectorClock{}
  w := DynamicVectorClock{}

  // v=[{}]; w=[{}]
  if v.Before(w) {
    t.Errorf("Not expecting v < w. Both v and w should be nil.")
  }

  // v=[{a,0}]; w=[{}]
  // w < v
  v.Append("a")
  if v.Before(w) {
    t.Errorf("Not expecting v < w. v should dominate w after appending key 'a'.")
  } else if !w.Before(v) {
    t.Errorf("Expected w < v after appending key 'a' to v.")
  }

  // v=[{a,1}]; w=[{}]
  // w < v
  v.Increment("a")
  if v.Before(w) {
    t.Errorf("Not expecting v < w. v should dominate w after incrementing key 'a'.")
  } else if !w.Before(v) {
    t.Errorf("Expected w < v after incrementing key 'a' to v.")
  }

  // v=[{a,1}]; w=[{b,0}]
  // v <> w
  w.Append("b")
  if v.Before(w) {
    t.Errorf("Not expecting v < w after appending key 'b' to w.")
  } else if w.Before(v) {
    t.Errorf("Not expecting w < v after appending key 'b' to w.")
  }

  // v=[{a,1},{b,0}]; w=[{b,0}]
  // w < v
  v.Append("b")
  if v.Before(w) {
    t.Errorf("Not expecting v < w after appending key 'b' to v.")
  } else if !w.Before(v) {
    t.Errorf("Expecting w < v after appending key 'b' to v.")
  }

  // v=[{a,1},{b,0}]; w=[{b,1}]
  // w <> v
  w.Increment("b")
  if v.Before(w) {
    t.Errorf("Not expecting v < w after incrementing key 'b' of w.")
  } else if w.Before(v) {
    t.Errorf("Not Expecting w < v after incrementing key 'b'of w.")
  }

  // v=[{a,1},{b,0}]; w=[{a,0},{b,1}]
  // w <> v
  w.Append("a")
  if v.Before(w) {
    t.Errorf("Not expecting v < w after appending key 'a' to w.")
  } else if w.Before(v) {
    t.Errorf("Not Expecting w < v after appending key 'a' to w.")
  }

  // v=[{a,1},{b,0}]; w=[{a,1},{b,1}]
  // v < w
  w.Increment("a")
  if !v.Before(w) {
    t.Errorf("Expected v < w after incrementing key 'a' of w.")
  } else if w.Before(v) {
    t.Errorf("Not expecting w < v after incrementing key 'a' of w.")
  }
}
