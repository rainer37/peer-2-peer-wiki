package article

import (
  "fmt"
)

type Direction int8
const (
  Empty Direction = -1  // root
  Left Direction = 0  // left branch
  Right Direction = 1  // right branch
)

// This is a standard binary tree except that each node can contain many sibling
// nodes.
type Treedoc struct {
  MiniNodes []*Node
  Left *Treedoc
  Right *Treedoc
}

// A node in the treedoc. Nodes have a value, path, disambiguator (siteId) and
// an indicator of whether the node is visible (deleted).
type Node struct {
  Value Atom
  Site Disambiguator
  Tombstone bool  // true if node has been deleted
  Left *Treedoc
  Right *Treedoc
}

// Nodes are identified in a treedoc by their path and their disambiguator (siteId)
type PosId struct {
  Dir Direction
  Site Disambiguator
}

// Represents the smallest unit that can be modified atomically
type Atom string

// A treedoc is a binary tree so a path is a bitstring (represented as an array)
// starting from the root where a 0 indicates a left branch and a 1 indicates a
// right branch.
type Path []PosId

// A globally-unique identifier for the process making the action
type Disambiguator string

// A list of insert/delete commands executed locally
type OpLog []Operation

type Operation struct {
  Command string  // "insert" or "delete"
  Path Path
  Site Disambiguator
}




// Walk the tree rooted at t in infix order.
// Return the atoms of the non-tombstone nodes.
func (t *Treedoc) Contents() []Atom {
  var contents []Atom

  for _, node := range t.traverseVisible() {
    contents = append(contents, node.Value)
  }

  return contents
}


// Prevent a value of a node from being shown to a user.
// Does not remove node from tree,
func (t *Treedoc) Delete(pos int, site Disambiguator) (*Node, error) {
  nodes := t.traverseVisible()
  if pos > len(nodes) {
    return fmt.Errorf("Treedoc::Delete(...) - Position is invalid.")
  }
  nodes[pos-1].Tombstone = true

  return nil
}


func (t *Treedoc) Insert(pos int, atom Atom, site Disambiguator) (*Node, error) {
  nodes := t.traverse()

  if pos > len(t.traverseVisible())+1 || pos < 0 {
    return fmt.Errorf("Treedoc::Insert(...) - Position is invalid.")
  }

  pid := Path{}
  switch len(nodes) {
  case 0:  // empty tree
    pid = append(pid, PosId{Empty, site})
  case 1:  // tree has only one node
    if pos == 1 {
      pid = append(pid, PosId{Left, site})
    } else {
      pid = append(pid, PosId{Right, site})
    }
  default:
    p,err := t.path(pos-1)
    f,err := t.path(pos)
    if err != nil {
      return fmt.Errorf("Treedoc::Insert(...) - Could not find path for given position.")
    }

    pid,err := newUid(p, f)
    if err != nil {
      return fmt.Errorf("Treedoc::Insert(...) - Failed to find a path to insert.")
    }
  }

  return t.insertNode(pid, &Node{atom, site, false, nil, nil})
}


func (t *Treedoc) insertNode(path Path, n *Node) error {
  next,err := t.walk(path)
  if err != nil {
    return fmt.Errorf("Treedoc::insertNode(...) - Invalid path.")
  }

  newTd := Treedoc{}
  newTd.MiniNodes = append(newTd.MiniNodes, n)
  next = &newTd

  return nil
}


func (t *Treedoc) deleteNode(path Path) {
  n,err := t.walk(path)
  if err != nil {
    return fmt.Errorf("Treedoc::deleteNode(...) - Invalid path.")
  }

  n.Tombstone = true

  return nil
}


func (t *Treedoc) walk(path Path) (*Treedoc, error) {
  lenP := len(path)

  for i,p := range path {
    switch {
    case p.Dir == Empty:
      for _,m := range t.MiniNodes {
        if m.Site == p.Site {
          t = m  // change pointer from root major node to root mini-node
        }
      }
    case p.Dir == Right && p.Site == "":  // right on major node
      t = t.Right
    case p.Dir == Left && p.Site == "":  // left on major node
      t = t.Left
    case p.Dir == Right:  // right on sibling node
      for _,m := range t.MiniNodes {
        if m.Site == path.Site {
          t = m.Right
        }
      }
    case path[i].Dir == Left:  // left on sibling node
      for _,m := range t.MiniNodes {
        if m.Site == path.Side {
          t = m.Left
        }
      }
    }
  }

  return t
}


func (t *Treedoc) path(pos int) (Path, error) {
  // translate pos (which is for visible) to infix position of all nodes
  // call infix to the limit
}


// Build a list of nodes in infix order
func (t *Treedoc) infix(p *Path, n *[]*node, depth int) {
  if t.left != nil {
    t.left.infix(n)
  }

  for _,m := range t.miniNodes {
    if m.left != nil {
      m.left.infix(n)
    }

    *n = append(*n, m)

    if m.right != nil {
      m.right.infix(n)
    }
  }

  if t.right != nil {
    t.right.infix(n)
  }
}


func (t *Treedoc) isEmpty() bool {
  return len(t.miniNodes) == 0
}


// Return an array of all nodes in the tree in infix order
func (t *Treedoc) traverse() []*node {
  var nodes []*node
  t.infix(&nodes, -1)
  return nodes
}

func (t *Treedoc) traverseVisible() []*node {
  var visNodes []*node
  var nodes []*node
  t.infix(&nodes)

  for _,n := range nodes {
    if !n.tombstone {
      visNodes = append(visNodes, n)
    }
  }

  return visNodes
}





func (t *Treedoc) newUid(uidp *Path, uidf *Path) (Path, error) {
  newPosId := PosId{}
  nodes := t.traverse()

  // Require uidp < uidf
  if !uidp.before(uidf) {
    return newPosId, fmt.errorf("Treedoc::newUid() - uidp not before uidf.")
  }


  // TODO add the disabmiguator
  // Check if there is a node between uidp and uidf. If so, call newUid on the
  // leftmost node such that uidp < uidm < uidf.
  if len(nodes) > 2 {
    for i := 0; i < len(nodes)-2; i++ {
      if nodes[i].Id.Path.equals(uidp.Path) {
        if !nodes[i+1].id.Path.equals(uidf) {
          return newUid(uidp, nodes[i+1].Id)
        }
      }
    }
  }

  // TODO decide if a disabmiguator needs to be included in the PosId: if uidp is
  // a mini-node then yes
  switch {
  case uidp.ancestor(uidf):
    newPosId = PosId{append(uidp.Path, false), uidp.Site}
  case uidf.ancestor(uidp):
    newPosId = PosId{append(uidf.Path, true), udif.Site}
  default:
    newPosId = PosId{append(uidp.Path, true), uidp.Site}
  }

  return newPosId, nil
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
  return uid.path.equals(&vid.path) && uid.site != vid.site
}






func (p *posId) before(q *posId) bool {
  u := p.path
  v := q.path
  i := u.commonPrefix(&v)

  // check that u != v
  if len(u) == i && len(v) == i && p.site == q.site {
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
    if u[i] == v[i] && p.site < q.site {
      return true
    }
  case len(u) == i && len(v) == i && p.site < q.site:
    return true
  }

  return false
}






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
