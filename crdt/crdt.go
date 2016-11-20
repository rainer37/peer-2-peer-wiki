/* TODO Packaage description
 *
 */

package crdt


import (
  "fmt"
  "bytes"
  //"strings"
)

/* Treedoc structures */
type posId struct {
  path []int
  disambiguator string
}
// func (p *posId) prefix(q *posId) bool {
//   if len(*p) > len(*q) {
//     return false
//   }
//   for i,e := range *p {
//     if e != (*q)[i] {
//       return false
//     }
//   }
//   return true
// }

type path []bool

type node struct {
  disambiguator string
  atom string
  tombstone bool
}

type nodeSummary struct {
  posId posId
  node *node
}

type Treedoc struct {
  left *Treedoc
  node []*node
  right *Treedoc
}
func NewTreedoc(nodeId string, atom string) *Treedoc {
  n := &node{nodeId, atom, false}
  return &Treedoc{nil, []*node{n}, nil}
}
func (t *Treedoc) Content() []string {
  var prefix = []int{}
  var ns = []nodeSummary{}
  t.infixPath(prefix, &ns)

  var content []string
  for _,n := range ns {
    if !n.node.tombstone {
      content = append(content, n.node.atom)
    }
  }
  return content
}

func (t *Treedoc) Delete(pos int) {
  var prefix = []int{}
  var ns = []nodeSummary{}
  t.infixPath(prefix, &ns)

  ns[pos-1].node.tombstone = true
}

func (t *Treedoc) Insert(atom string, pos int) {
  var prefix = []int{}
  var ns = []nodeSummary{}
  t.infixPath(prefix, &ns)

  if len(ns) < 2 {
    if pos == 1 {
      t.insertPos(&node{"", atom, false}, posId{[]int{0}, ""})
    } else {
      t.insertPos(&node{"", atom, false}, posId{[]int{1}, ""})
    }
  } else {
    pathBefore := ns[pos-1].posId.path
    pathAfter := ns[pos].posId.path

    newPos := t.newPosId(pathBefore, pathAfter)
    fmt.Printf("Inserting \"%s\" in position %d at %v with pbefore = %v, pafter = %v\n", atom, pos, newPos,pathBefore,pathAfter)
    t.insertPos(&node{"", atom, false}, posId{newPos, ""})
  }
}

func (t *Treedoc) insertPos(n *node, p posId) {
  newNode := &Treedoc{nil, []*node{n}, nil}

  if len(p.path) > 1 {
    for _,e := range p.path[:len(p.path)-1] {
      if e == 0 {
        t = t.left
      } else {
        t = t.right
      }
    }
  }

  if p.path[len(p.path)-1] == 0 {
    t.left = newNode
  } else {
    t.right = newNode
  }
}

func (t *Treedoc) newPosId(pathBefore []int, pathAfter []int) []int {
  switch {
  // case //if exists:
  //   return newPosId(pathBefore, elem)
case len(pathBefore) == 0:
    return append(pathAfter, 0)
  case len(pathAfter) == 0:
    return append(pathBefore, 1)
  case pathBefore[0] == pathAfter[0]:
    return append(pathAfter, 0)
  default:
    return append(pathBefore, 1)
  }

}


// func (t *Treedoc) infix(ns *[]nodeSummary) *[]nodeSummary {
//   if t.left != nil {
//     ns = t.left.infix(ns)
//
//
//     //pid := posId{append(q[0].posId.path, false), t.node[0].disambiguator}
//     // pid := posId{[]bool{false},t.node[0].disambiguator}
//     // *ns = append(*ns, nodeSummary{pid, t.node[0].atom})
//     //append(t.left.size(), nodeSummary{posId{},t.node.atom})
//   } else {
//   //ns = []nodeSummary{{posId{[]bool{false},t.node[0].disambiguator},t.node[0].atom}}
//   for i,x := range *ns {
//     (*ns)[i].posId.path = append(x.posId.path, false)
//   }
//   *ns = append(*ns, nodeSummary{posId{[]bool{false},t.node[0].disambiguator},t.node[0].atom})
// }
//   if t.right != nil {
//     ns = t.right.infix(ns)
//
//     for i,x := range *ns {
//       (*ns)[i].posId.path = append(x.posId.path, true)
//     }
//     //pid := posId{append(q[0].posId.path, false), t.node[0].disambiguator}
//     pid := posId{[]bool{true},t.node[0].disambiguator}
//     *ns = append(*ns, nodeSummary{pid, t.node[0].atom})
//     //append(t.left.size(), nodeSummary{posId{},t.node.atom})
//   }
//
//
//   // if t.right != nil {
//   //   q := t.right.infix()
//   //   pid := posId{append(q[0].posId.path, true), t.node[0].disambiguator}
//   //   ns = append(q, nodeSummary{pid, t.node[0].atom})
//   // }
//
//
//
//   return ns
// }
/*
func (t *Treedoc) infixPath(prefix []bool, path *[][]bool) {
  //fmt.Printf("%s.infixPath(%v, %v<%p>)\n",t.String(),prefix, path, path)
  if t.left != nil {
    //fmt.Printf("left path is %v at %p\n", path, path)
    *prefix = append(*prefix, false)
    t.left.infixPath(prefix, path)
  }
  //fmt.Printf("Appending %v to %v\n", prefix, path)
  *path = append(*path, *prefix)
  //fmt.Printf("Path is now %v\n", path)
  if t.right != nil {
    //fmt.Printf("right path is %v at %p\n", path, path)
    *prefix = append(*prefix, true)

    t.right.infixPath(prefix, path)
    //fmt.Printf("Pointer of swtiching value %p", (*path)[3][2])
    //fmt.Printf("after right is %v\n", path)
  }
  //fmt.Printf("Path is %v\n", path)
}
*/

func (t *Treedoc) infixPath(prefix []int, path *[]nodeSummary) {
  // fmt.Printf("prefix: %v, path: %v\n", prefix, path)
  //fmt.Printf("%s.infixPath(%v, %v<%p>)\n",t.String(),prefix, path, path)
  if t.left != nil {

    t.left.infixPath(append(prefix, 0), path)
  }

  *path = append(*path, nodeSummary{posId{prefix,t.node[0].disambiguator}, t.node[0]})

  if t.right != nil {
    // fmt.Printf("Path is %v\n", path)

    t.right.infixPath(append(prefix, 1), path)
  }
}

// func (t *Treedoc) infixContent([]string)






// func (t *Treedoc) Insert(n *node) *Treedoc {
//   if t == nil {
//     return &Treedoc{nil, []*node{n}, nil}
//   }
//   if n.atom < t.node[0].atom {
//     t.left = t.left.Insert(n)
//     return t
//   }
//   t.right = t.right.Insert(n)
//   return t
// }






//
// func (t *Treedoc) infixPath(prefix []bool, path *[][]bool) {
//   if t.left != nil {
//     (*path) = append(*path, []bool{false})
//     t.left.infixPath(append(prefix, false), path)
//     (*path)[0] = append((*path)[0], false)
//
//     // for i,x := range *path {
//     //   (*path)[i] = append(x, false)
//     // }
//
//   }
//
//   if t.right != nil {
//     (*path) = append(*path, []bool{true})
//     t.right.infixPath(path)
//   }
// //  *path = append(*path, []bool{})
//
// }

// func (t *Treedoc) infixPath(path *[][]bool) {
//   if t.left != nil {
//     (*path) = append(*path, []bool{false})
//     t.left.infixPath(path)
//     (*path)[0] = append((*path)[0], false)
//
//     // for i,x := range *path {
//     //   (*path)[i] = append(x, false)
//     // }
//
//   }
//
//   if t.right != nil {
//     (*path) = append(*path, []bool{true})
//     t.right.infixPath(path)
//   }
// //  *path = append(*path, []bool{})
//
// }




/*
func (t *Treedoc) infix() []nodeSummary {
  var ns []nodeSummary

  if t.left != nil {
    ns = t.left.infix()

    for i,x := range ns {
      ns[i].posId.path = append(x.posId.path, false)
    }
    //pid := posId{append(q[0].posId.path, false), t.node[0].disambiguator}
    pid := posId{[]bool{false},t.node[0].disambiguator}
    ns = append(ns, nodeSummary{pid, t.node[0].atom})
    //append(t.left.size(), nodeSummary{posId{},t.node.atom})
  }
  //ns = []nodeSummary{{posId{[]bool{false},t.node[0].disambiguator},t.node[0].atom}}
  ns = append(ns, nodeSummary{posId{[]bool{false},t.node[0].disambiguator},t.node[0].atom})

  if t.right != nil {
    t.right.infix()

    for i,x := range ns {
      ns[i].posId.path = append(x.posId.path, true)
    }
    //pid := posId{append(q[0].posId.path, false), t.node[0].disambiguator}
    pid := posId{[]bool{true},t.node[0].disambiguator}
    ns = append(ns, nodeSummary{pid, t.node[0].atom})
    //append(t.left.size(), nodeSummary{posId{},t.node.atom})
  }


  // if t.right != nil {
  //   q := t.right.infix()
  //   pid := posId{append(q[0].posId.path, true), t.node[0].disambiguator}
  //   ns = append(q, nodeSummary{pid, t.node[0].atom})
  // }



  return ns
}
*/






/*

// works for the right tree
func (t *Treedoc) infix() []nodeSummary {
  ns := []nodeSummary{}

  if t.left != nil {
    q := t.left.infix()

    for i,x := range q {
      q[i].posId.path = append(x.posId.path, false)
    }
    //pid := posId{append(q[0].posId.path, false), t.node[0].disambiguator}
    pid := posId{[]bool{false},t.node[0].disambiguator}
    ns = append(q, nodeSummary{pid, t.node[0].atom})
    //append(t.left.size(), nodeSummary{posId{},t.node.atom})
  }
  ns = []nodeSummary{{posId{[]bool{false},t.node[0].disambiguator},t.node[0].atom}}

  if t.right != nil {
    q := t.right.infix()

    for i,x := range q {
      q[i].posId.path = append(x.posId.path, true)
    }
    //pid := posId{append(q[0].posId.path, false), t.node[0].disambiguator}
    pid := posId{[]bool{true},t.node[0].disambiguator}
    ns = append(q, nodeSummary{pid, t.node[0].atom})
    //append(t.left.size(), nodeSummary{posId{},t.node.atom})
  }


  // if t.right != nil {
  //   q := t.right.infix()
  //   pid := posId{append(q[0].posId.path, true), t.node[0].disambiguator}
  //   ns = append(q, nodeSummary{pid, t.node[0].atom})
  // }



  return ns
} */

// func (t *Treedoc) Print() {
//
//   for _,ns := range t.infix() {
//     path := strings.Join(ns.path,"")
//     fmt.Printf("%s {%s:%s}", ns.atom, path, ns.disambiguator)
//   }
// }




// func (t *Treedoc) Content() []string {
//   var content []string{}
//   return t.infixNode(content)
// }

// // Print the Treedoc infix
// func (t *Treedoc) String() string {
//   if t == nil {
//     return ""
//   }
//
//   var buf bytes.Buffer
//   return t.infix(&buf).String()
// }





// Return the number of visible nodes
func (t *Treedoc) size() int {
  count := 0
  miniNodes := t.node
  for _,m := range miniNodes {
    if !m.tombstone {
      count++
    }
  }

  if t.left != nil {
    count += t.left.size()
  }
  if t.right != nil {
    count += t.right.size()
  }
  return count
}


// Return the posId of the tree element based on the display location (infix order)
func (t *Treedoc) posId(disPos int) (*posId, error) {
  return &posId{}, nil
}
func (t *Treedoc) parent(s *Treedoc) bool {
  return false
}
func (t *Treedoc) ancestor(s *Treedoc) bool {
  return false
}

// func (t *Treedoc) Insert(pos int, n *node) *Treedoc {
//   // error if p < 1 || p > size
//   if t == nil {
//     return &Treedoc{nil, []*node{n}, nil}
//   }
//
//
//   var path = [][]bool{}
//   tree.infixPath([]bool{}, &path)
//
//   posBefore := path[pos]
//   posAfter := path[pos+1]
//
//   newPos := newId(posBefore, posAfter)
//
//   for _,branch := newPos[:len(newPos-2)] {
//     if branch {
//       t = t.right
//     } else {
//       t = t.left
//     }
//   }
//
//   // if node already exists, append as mini-node
//   // else
//   if newPos[len(newPos-1)] {
//     t.right = &Treedoc{nil, atom, nil}
//   } else {
//     t.left = &Treedoc{nil, atom, nil}
//   }
//
//   return nil
// }

func (t *Treedoc) insertTree(path, node) {

}


// // Use bytes.Buffer as a mutable string (StringBuilder in Java)
// func (t *Treedoc) infix(str *bytes.Buffer) *bytes.Buffer {
//   // showTombstone := true
//
//   if t.left != nil {
//     str = t.left.infix(str)
//   }
//
//
//   //nstr := fmt.Sprintf("%v", *t.node[0])
//   str.WriteString(fmt.Sprintf("%v", t.node))
//
//   if t.right != nil {
//     str = t.right.infix(str)
//   }
//   return str
// }






/*



//https://www.reddit.com/r/golang/comments/25aeof/building_a_stack_in_go_slices_vs_linked_list/?st=iviuwwi5&sh=d9a088ef
type stackPath []bool
func (s stackPath) Empty() bool {
  return len(s) == 0
}

func (s stackPath) Peek() bool {
  return s[len(s)-1]
}

func (s *stackPath) Put(b bool) {
  *s = append(*s, b)
}

func (s *stackPath) Pop() bool {
  d := (*s)[len(*s)-1]
  *s = (*s)[:len(*s)-1]
  return d
}

type stack []Treedoc
func (s stack) Empty() bool {
  return len(s) == 0
}

func (s stack) Peek() Treedoc {
  return s[len(s)-1]
}

func (s *stack) Put(t Treedoc) {
  *s = append(*s, t)
}

func (s *stack) Pop() *Treedoc {
  d := (*s)[len(*s)-1]
  *s = (*s)[:len(*s)-1]
  return &d
}






// Implement the Treedoc CRDT
type node struct {
  disambiguator string
  value string
}
type path []bool
func (p *path) prefix(q *path) bool {
  // error checking
  if len(p) > len(q) {
    return false
  }
  for i,e := range p {
    if e != q[i] {
      return false
    }
  }
  return true
}

type Treedoc struct {
  left *Treedoc
  //node []node
  value string
  right *Treedoc
}

func (t *Treedoc) height() int {
  h := 0
  if t.left != nil {
    t.left.height()
    h++
  }

  if t.right != nil {
    t.right.height()
    h++
  }
  h++
  return h
}






func (t *Treedoc) print(prefix string, isTail bool) {
  if isTail {
    fmt.Println(prefix + "└── " + t.value)
  } else {
    fmt.Println(prefix + "├── " + t.value)
  }
  h := t.height()
  for i := 0; i < ; i++ {

  }
}





/*
// // Create a new treedoc and return a pointer to it
// func NewTreedoc() *Treedoc {
//
// }


func createPosPath(posStart path, posEnd path) path {
  if posStart.prefix(posEnd) {
    return append(posEnd, false)
  } else if posEnd.prefix(posStart) {
    return append(posStart, true)
  } else {
    return append(posStart, true)
  }

}

// Add an atom to the tree after the specified atom (where the index is number in
// infix order).
// atom: the element to insert into the tree
// after: the position (in infix order) that the
// site: the identifier for the site or process making the insert
func (t *Treedoc) Insert(atom string, after int, site string) error {
  // TODO Error checking
  ltree := t.left.infixPath(make([]bool))
  rtree := t.right.infixPath(make([]bool))

  lheight := len(ltree)
  rheight := len(rtree)

  // TODO infixPath should return two arrays
  if after > lheight {
    posBefore = rtree[:(after-lheight)]
    posAfter = rtree[:(after-lheight+1)]
  } else {
    posBefore = ltree[:after]
    posAfter = rtree[:(after+1)]
  }

  newPos := newId(posBefore, posAfter)

  for _,branch := newPos[:len(newPos-2)] {
    if branch {
      t = t.right
    } else {
      t = t.left
    }
  }

  // if node already exists, append as mini-node
  // else
  if newPos[len(newPos-1)] {
    t.right = &Treedoc{nil, atom, nil}
  } else {
    t.left = &Treedoc{nil, atom, nil}
  }

  return nil
}

// func (t *Treedoc) insertAtom(atom string, position []bool) {
//   for _,branch := position[:len(position-2)] {
//     if branch {
//       t = t.right
//     } else {
//       t = t.left
//     }
//   }
//
//   // if node already exists, append as mini-node
//   // else
//   if positions[len(position-1)] {
//     t.right = &Treedoc{nil, atom, nil}
//   } else {
//     t.left = &Treedoc{nil, atom, nil}
//   }
// }


// Find path to value within before looking at at most before when examining the
// tree infix order.
// Return the path to the node closest to before in infix order
func (t *Treedoc) Find(value string, before int) ([]bool, error) {
  s := stack{}
  p := stackPath{}
  for !s.Empty() || t != nil {
    fmt.Println("HHH")

    if t.left != nil {
      fmt.Println("SDA")
      s.Put(*t)
      p.Put(false)
      t = t.left
    } else {
      fmt.Println("Backtracking one level.")
      t = s.Pop()
      for t.right == nil {
        fmt.Println("Backtracking one level.")
        t = s.Pop()
      }

      before--
      if before < 0 {
        return []bool{}, fmt.Errorf("Value not found 1.")
      }
      if t.value == value {
        return p, nil
      }
      p.Pop()
      fmt.Println("Going down right branch.")
      t = t.right

    }

  }
  fmt.Println("ASASAS")

  return []bool{false}, fmt.Errorf("Value not found.")
}

func (t *Treedoc) findInfix(value string, limit int, path []bool) {

}

func (t *Treedoc) Path(value string, path *[]bool) bool  {
  if t == nil {
    return false
  }
  switch {
  case t.value == value:
    return true
  case t.left.Path(value, path):
    *path = append(*path, false)
    return true
  case t.right.Path(value, path):
    *path = append(*path, true)
    return true
  }
  return false
}




func (t *Treedoc) Insert(v string) *Treedoc {
  if t == nil || *t == (Treedoc{}) {
    return &Treedoc{nil, v, nil}
  }
  if v < t.value {
    t.left = t.left.Insert(v)
    return t
  }
  t.right = t.right.Insert(v)
  return t
}
*/
// Print the Treedoc infix
func (t *Treedoc) String() string {
  if t == nil {
    return ""
  }

  var buf bytes.Buffer
  return t.infix(&buf).String()
}

// Use bytes.Buffer as a mutable string (StringBuilder in Java)
func (t *Treedoc) infix(str *bytes.Buffer) *bytes.Buffer {
  if t.left != nil {
    str = t.left.infix(str)
  }
  str.WriteString(t.node[0].atom)
  if t.right != nil {
    str = t.right.infix(str)
  }
  return str
}

// func (t *Treedoc) infixPath(path *[]bool) { // *[]bool {
//   if t.left != nil {
//     *path = append(*path, false)
//     t.left.infixPath(path)
//   }
//   if t.right != nil {
//     *path = append(*path, true)
//     t.right.infixPath(path)
//   }
//   //return str
// }









// Implement vector clock CRDT
//http://www.vs.inf.ethz.ch/edu/VS/exercises/A3/DVC_Landes.pdf
type DynamicVectorClock map[string]int

// Increment the value of the specified key
func (v DynamicVectorClock) Increment(key string) error {
  // check that the key exists
  _,ok := v[key]
  if !ok {
    return fmt.Errorf("<W> DynamicVectorClock::Increment(%v) - Key does not exist.", key)
  }

  // Increment the counter associated with the specified key
  v[key]++

  return nil
}

// Add a new key to the dynamic vector clock
func (v DynamicVectorClock) Append(key string) error {
  // check that key doesn't already exist
  _,ok := v[key]
  if ok {
    return fmt.Errorf("<W> DynamicVectorClock::Append(%v) - Key already exists.", key)
  }

  // append the key setting its value to 0
  v[key] = 0

  return nil
}

func (v DynamicVectorClock) Equals(w DynamicVectorClock) bool {
  if v == nil && w == nil {
    return true
  }

  if len(v) != len(w) {
    return false
  }

  for k,clock := range v {
    if w[k] != clock {
      return false
    }
  }

  return true
}


// Compare two vector clocks to determine if v happened before w.
// The happened-before relation is defined as:
// v < w iff v[k] <= w[k] for all k in v, and there is at least one k where v[k] < w[k]
func (v DynamicVectorClock) Before(w DynamicVectorClock) bool {
  if v.Equals(w) {
    return false
  }

  // v must have witnessed at least one more event then w, so v could not be before w.
  // this checks that there is at least one k where v[k] < w[k]
  if len(v) > len(w) {
    return false
  }

  // this checks that v[k] <= w[k] for all k in v
  for k,clock := range v {
    // Note that if w does not have key k then w[k] == 0, so clock > w[k] is true
    if clock > w[k] {
      return false
    }
  }

  return true
}

/*
func (v *DynamicVectorClock) Caused(w *DynamicVectorClock) bool {

}

func (v *DynamicVectorClock) Concurrent(w *DynamicVectorClock) bool {

}
*/
