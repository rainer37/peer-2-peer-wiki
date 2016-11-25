/* TODO Packaage description
 *
 */

package crdt


import "fmt"

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

func (v DynamicVectorClock) Value(key string) int {
  // // check that the key exists
  // val,ok := v[key]
  // if !ok {
  //   return 0
  // }


  return v[key]
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




// This is a standard binary tree except that each node can contain many sibling
// nodes.
type Treedoc struct {
  miniNodes []*node
  left *Treedoc
  right *Treedoc
}


// Walk the tree rooted at t in infix order.
// Return the atoms of the non-tombstone nodes.
func (t *Treedoc) Contents() []string {
  var contents []string

  for _, node := range t.traverse() {
    if !node.tombstone {
      contents = append(contents, node.value)
    }
  }

  return contents
}


// Prevent a value of a node from being shown to a user.
func (t *Treedoc) Delete(pos int, site string) error {
  nodes := t.traverse()
  var siteNodes []*node
  for _,node := range nodes {
    if node.id.disambiguator == site {
      siteNodes = append(siteNodes, node)
    }
  }

  if len(siteNodes) < pos {
    return fmt.Errorf("Treedoc::Delete(...) - Position is invalid.")  // too big
  }
  return t.deleteNode(siteNodes[pos-1], site)
  //return t.deleteNode(nodes[pos-1], nodes[pos-1].id.disambiguator)
}


// Insert text at the specified position and identify the mini node that is created
// with the site parameter.
func (t *Treedoc) Insert(atom string, pos int, site string) error {
  majorNodes := t.getMajorNodes()
  newNode := node{atom,posId{path{}, site},false}

  // NOTE in calls to newUid, only send the 1st mini-node. This works because newUid,
  // only checks the path to the major node (mini-nodes appear only when newUid
  // concurrently generates the same path on DIFFERENT clients (newUid will never
  // generate the same path on the same client)).

  // NOTE since we do not perform GC on the tree currently, no need to worry about
  // creating missing ancestor nodes.

  switch {
  case pos <= 1:
    if len(majorNodes) == 1 || len(majorNodes[0][0].id.path) == 0 { // only root or infix leftmost node is root
      newNode.id.path = path{false}
    } else {
      p,err := t.newUid(&node{}, majorNodes[0][0])  // (root, left-most path (1st mini-node))
      if err != nil {
        return fmt.Errorf("Treedoc::Insert(...) - Failed to generate new node id.")
      }
      newNode.id.path = p
    }
  case pos > len(majorNodes):
    if len(majorNodes) == 1 || len(majorNodes[len(majorNodes)-1][0].id.path) == 0 { // only root or infix rightmost node is root
      newNode.id.path = path{true}
    } else {
      p,err := t.newUid(majorNodes[len(majorNodes)-1][0], &node{})  // (right-most path (1st mini-node), root)
      if err != nil {
        return fmt.Errorf("Treedoc::Insert(...) - Failed to generate new node id.")
      }
      newNode.id.path = p
    }
  default:
    p,err := t.newUid(majorNodes[pos-2][0], majorNodes[pos-1][0])
    if err != nil {
      return fmt.Errorf("Treedoc::Insert(...) - Failed to generate new node id.")
    }
    newNode.id.path = p
  }

  return t.insertNode(&newNode)
}


// Implements the treedoc delete function
func (t *Treedoc) deleteNode(n *node, site string) error {
  // NOTE since we do not perform GC on the tree currently, no need to worry about
  // creating missing ancestor nodes.

  n.tombstone = true
  return nil
}


// Returns major nodes and their mini nodes as a 2D array
func (t *Treedoc) getMajorNodes() [][]*node {
  var majorNodes [][]*node
  t.infix(&majorNodes)
  return majorNodes
}


// Build a list of nodes in infix order
func (t *Treedoc) infix(n *[][]*node) {
  if t.left != nil {
    t.left.infix(n)
  }

  *n = append(*n, t.miniNodes)

  if t.right != nil {
    t.right.infix(n)
  }
}


// Helper function to insert a node into the treedoc.
// Require: 1) non-empty path for insertion; 2) path must be unique
func (t *Treedoc) insertNode(n *node) error {
  path := n.id.path
  sid := n.id.disambiguator

  // error checking
  if len(path) < 1 {
    fmt.Errorf("Treedoc::insertNode(%v) - Empty path.", *n)
  }

  if len(path) > 1 {
    // Iterate over the path to set the t pointer to the correct node
    for i := range path[:len(path)-1] {
      if path[i] {
        t = t.right
      } else {
        t = t.left
      }
    }
  }

  // Insert the node by setting the left or right pointer to the node's address
  // if a node doesn't already exist. Otherwise, append new node as a mini-node.
  next := &t.left
  if path[len(path)-1] {
    next = &t.right
  }

  if *next == nil {
    newtd := Treedoc{}
    newtd.miniNodes = append(newtd.miniNodes, n)
    *next = &newtd
  } else {
    pos := -1
    for i,miniNode := range (*next).miniNodes {
      if miniNode.id.disambiguator > sid {
        pos = i
      } else if miniNode.id.disambiguator == sid {
        fmt.Errorf("Treedoc::insertNode(%v) - Attempted to insert node with duplicate disambiguator.", *n)
      }
    }
    (*next).miniNodes = append((*next).miniNodes, n)
    if pos > -1 { // insert sibling node in order of disabmiguator
      copy((*next).miniNodes[pos+1:], (*next).miniNodes[pos:])
      (*next).miniNodes[pos] = n
    }
  }

  return nil
}


// Generate a unique path for a new node to be inserted between nodes p and f
// Require: p < f (where < is the posId.before operation)
func (t *Treedoc) newUid(p *node, f *node) (path, error) {
  uidp := p.id
  uidf := f.id
  if !uidp.before(&uidf) {
    return path{}, fmt.Errorf("Treedoc::newUid(p:%v, f:%v) - p.podId !< f.posId", *p, *f)
  }

  majorNodes := t.getMajorNodes()
  var m *node

  for _,majorNode := range majorNodes {
    uidm := majorNode[0].id  // first mini-node
    if uidp.before(&uidm) && uidm.before(&uidf) {
      m = majorNode[0]
      break
    }
  }

  switch {
  case m != nil: // elements between p and f
    return t.newUid(p, m)
  case p.ancestor(f):
    return append(uidf.path, false), nil
  case f.ancestor(p):
    return append(uidp.path, true), nil
  default:
    return append(uidp.path, true), nil
  }
}


// Return an array of all nodes in the tree in infix order
func (t *Treedoc) traverse() []*node {
  var miniNodes []*node
  majorNodes := t.getMajorNodes()

  // flatten the tree (only use mini nodes)
  for _,majorNode := range majorNodes {
    for _,miniNode := range majorNode {
      miniNodes = append(miniNodes, miniNode)
    }
  }

  return miniNodes
}




// A node in the treedoc. Nodes have a value, path, disambiguator (siteId) and
// an indicator of whether the node is visible (deleted).
type node struct {
  value string
  id posId  // once the node is inserted, the path will never change so store it
  tombstone bool  // true if node has been deleted
}

// u is a parent of v if they have the same common path and v's path is one longer
// than u's path.
func (u *node) parent(v *node) bool {
  upath := u.id.path
  vpath := v.id.path
  prefixLen := upath.commonPrefix(&vpath)
  if len(upath) == prefixLen && len(vpath) == len(upath) + 1 {
    return true
  }
  return false
}

// u is an ancestor of v if they share any common path and v's path is strictly
// longer than u's path. The root node (empty path) is an ancestor of all nodes.
func (u *node) ancestor(v *node) bool {
  upath := u.id.path
  vpath := v.id.path
  if upath.prefix(&vpath) && len(upath) < len(vpath) {
    return true
  }
  return false
}

// u is a mini-sibling (or side-node) of v if they both have the same path and
// different disambiguators.
func (u *node) miniSibling(v *node) bool {
  uid := u.id
  vid := v.id
  return uid.path.equals(&vid.path) && uid.disambiguator != vid.disambiguator
}


// Nodes are identified in a treedoc by their path and their disambiguator (siteId)
type posId struct {
  path path
  disambiguator string
}
func (p *posId) before(q *posId) bool {
  u := p.path
  v := q.path
  i := u.commonPrefix(&v)

  // check that u != v
  if len(u) == i && len(v) == i && p.disambiguator == q.disambiguator {
    return false
  }

  switch {
  case len(u) == 0 && len(v) > 0 && v[0] == true:  //(e, 1...)
    return true
  case len(u) > 0 && len(v) == 0 && u[0] == false: //(0..., e)
    return true
  case i > 0 && len(u) > i && len(v) == i && u[i] == false: //(c1...cn0..., c1...cn) NOT (0, e)
    return true
  case i > 0 && len(u) == i && len(v) > i && v[i] == true: //(c1...cn, c1...cn1...) NOT (e, 1)
    return true
  case len(u) > i && len(v) > i:
    if u[i] == false && v[i] == true {
      return true
    }
    if u[i] == v[i] && p.disambiguator < q.disambiguator {
      return true
    }
  case len(u) == i && len(v) == i && p.disambiguator < q.disambiguator:
    return true
  }

  return false
}

// A treedoc is a binary tree so a path is a bitstring (represented as an array)
// starting from the root where a 0 indicates a left branch and a 1 indicates a
// right branch.
type path []bool

// Two paths are equal if they have the same length and agree in every position.
func (p *path) equals(q *path) bool {
  if len(*p) == len(*q) {
    for i := range *p {
      if (*p)[i] != (*q)[i] {
        return false
      }
    }
    return true
  }
  return false
}

// p is a prefix of q if p is the root node OR p and q agree in every position up
// to the length of p. Note that we must have p <= q in terms of length.
func (p *path) prefix(q *path) bool {
  if len(*p) == 0 { // root
    return true
  }
  if len(*p) > len(*q) {
    return false
  }

  for i := range *p {
    if (*p)[i] != (*q)[i] {
      return false
    }
  }
  return true
}
// return the length of the longest common prefix of p and q where the common
// prefix is the first position in the bitstrings where the values disagree.
func (p *path) commonPrefix(q *path) int {
  if len(*p) <= len(*q) {
    for i := range *p {
      if (*p)[i] != (*q)[i] {
        return i
      }
    }
    return len(*p)
  } else {
    for i := range *q {
      if (*p)[i] != (*q)[i] {
        return i
      }
    }
    return len(*q)
  }
}
