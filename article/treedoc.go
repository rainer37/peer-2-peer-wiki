package article

import (
  //"fmt"
  "strings"
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

// Represents the smallest unit that can be modified atomically
type Atom string

// A globally-unique identifier for the process making the action
type Disambiguator string


// A treedoc is a binary tree so a path is a bitstring (represented as an array)
// starting from the root where a 0 indicates a left branch and a 1 indicates a
// right branch.
type Path []PosId

// Nodes are identified in a treedoc by their path and their disambiguator (siteId)
type PosId struct {
  Dir Direction
  Site Disambiguator
}

// The branch direction in the tree
type Direction int8
const (
  Empty Direction = -1  // root
  Left Direction = 0  // left branch
  Right Direction = 1  // right branch
)



// Walk the tree rooted at t in infix order.
// Return the atoms of the non-tombstone (visible) nodes.
func (t *Treedoc) Contents() []Atom {
  var contents []Atom
  _,nodes := t.traverseVisible()
  for _, node := range nodes {
    contents = append(contents, node.Value)
  }

  return contents
}


// Prevent a value of a node from being shown to a user.
// Does not remove node from tree.
func (t *Treedoc) Delete(pos int, site Disambiguator) (Path, error) {
  paths,nodes := t.traverseVisible()

  switch {
  case pos < 1:
    pos = 1
  case pos > len(nodes):
    pos = len(nodes)
  }
  nodes[pos-1].Tombstone = true

  return paths[pos-1], nil
}


// Insert a node at the specified position (as ordered by the Contents func).
func (t *Treedoc) Insert(pos int, atom Atom, site Disambiguator) (Path, error) {
  pid := Path{}
  vpaths,_ := t.traverseVisible()
  paths,nodes := t.traverse()


  switch {
  case len(nodes) == 0:  // empty tree
    pid = append(pid, PosId{Empty, site})
    //fmt.Println("Insert - inserting (in empty tree) at ", pid.String())
  case pos <= 1:  // inserting to the left of the tree
    lid := paths[0]  // get path to leftmost node
    lid[len(lid)-1].Site = ""  // remove the last disambiguator (so we insert to left of major node)
    pid = append(lid, PosId{Left, site})
    //fmt.Println("Insert - inserting (leftmost) at ", pid.String())
  case pos > len(nodes):  // inserting to the right of the tree
    rid := paths[len(paths)-1]  // get path to rightmost node
    rid[len(rid)-1].Site = ""
    pid = append(rid, PosId{Right, site})
    //fmt.Println("Insert - inserting (rightmost) at ", pid.String())
  default:  // inserting between two existing nodes
    p := vpaths[pos-2]
    f := vpaths[pos-1]

    // Check if there are hidden nodes between p and f. If so, set f to be the
    // leftmost (possibly hidden) node after p.
    if !f.equal(&paths[pos-1]) {
      for i,q := range paths[pos-2:] {
        if p.equal(&q) {
          f = paths[i+pos-1]
        }
      }
    }

    if p.before(&f) {
      pid,_ = t.newUid(&p, &f, site)
      //fmt.Printf("1) New path between %s and %s is %s.\n", p.String(), f.String(), pid.String())
    } else {
      pid,_ = t.newUid(&f, &p, site)
      //fmt.Printf("2) New path between %s and %s is %s.\n", f.String(), p.String(), pid.String())
    }
  }

  return pid, t.insertNode(pid, &Node{atom, site, false, nil, nil})
}


// Format path to make it look nice
func (p *Path) String() string {
  var str []string
  for _,pid := range *p {
    dir := ""
    site := string(pid.Site)
    switch pid.Dir {
    case Empty:
      dir = "e"
    case Left:
      dir = "0"
    case Right:
      dir = "1"
    }
    if site != "" {
      str = append(str, "("+dir+":"+site+")")
    } else {
      str = append(str, dir)
    }
  }
  return strings.Join(str, "")
}


// Insert the node at the specified path
func (t *Treedoc) insertNode(path Path, n *Node) error {
  next := t.walk(path)
  inserted := false

  if *next == nil {
    *next = &Treedoc{}
  } else {
    for i,m := range (*next).MiniNodes {
      if m.Site > n.Site {
        (*next).MiniNodes = append((*next).MiniNodes, nil)
        copy((*next).MiniNodes[i+1:], (*next).MiniNodes[i:])
        (*next).MiniNodes[i] = n
        inserted = true
        break
      }
    }
  }

  // Append to the end of the list of mini-nodes because new disambiguator > existing ones
  // Or, this is a new node and there are no other mini-nodes
  if !inserted {
    (*next).MiniNodes = append((*next).MiniNodes, n)
  }

  return nil
}


// Set the tombstone field on the nodes specified by path to true (make it invisible)
func (t *Treedoc) deleteNode(path Path) error {
  n := t.walk(path)

  for _,m := range (*n).MiniNodes {
    if m.Site == path[len(path)-1].Site {
      m.Tombstone = true
    }
  }

  return nil
}

// Walk the tree following the path and return a pointer to the last node's Left
// or right field.
func (t *Treedoc) walk(path Path) **Treedoc {
  var s **Treedoc
  pLen := len(path)

  if pLen == 0 || pLen == 1 && path[0].Dir == Empty {
    return &t
  }

  for i,p := range path {
    // disambiguate mini-nodes (using the site of the previous path element)
    if i > 0 && path[i-1].Site != "" {
      for _,m := range t.MiniNodes {
        if m.Site == path[i-1].Site {
          switch p.Dir {
          case Left:
            s = &m.Left
            t = m.Left
          case Right:
            s = &m.Right
            t = m.Right
          }
          break
        }
      }
    } else {
      switch p.Dir {
      case Left:
        s = &t.Left
        t = t.Left
      case Right:
        s = &t.Right
        t = t.Right
      }
    }
  }

  return s
}


// Build a list of nodes (and their paths) in infix order
func (t *Treedoc) infix(p Path, paths *[]Path, n *[]*Node) {
  b := make(Path, len(p))
  copy(b, p)

  // Go left from major node
  if t.Left != nil {
    b = append(p, PosId{Left, ""})
    t.Left.infix(b, paths, n)
  }

  for _,m := range t.MiniNodes {
    c := make(Path, len(p))
    copy(c, p)

    // handle the root node
    if len(c) == 0 {
      if len(t.MiniNodes) > 1 {
        c = append(p, PosId{Empty, m.Site})
      } else {
        c = append(p, PosId{Empty, ""})
      }
    }

    // Go left from mini-node
    if m.Left != nil {
      // Include disambiguator if next node has mini-nodes
      if len(m.Left.MiniNodes) > 1 {
        b = append(c, PosId{Left, m.Site})
      } else {
        b = append(c, PosId{Left, ""})
      }
      // Set the disambiguator of the previous position if I'm a mini-node
      if len(t.MiniNodes) > 1 {
        if len(b) > 1 {
          b[len(b)-2].Site = m.Site
        }
      }
      m.Left.infix(b, paths, n)
    }


    *n = append(*n, m)



     if len(c) >= 1 {
       c[len(c)-1].Site = m.Site
     }
     *paths = append(*paths, c)


     // Go right from mini-node
    if m.Right != nil {
      // Include disambiguator if next node has mini-nodes
      if len(m.Right.MiniNodes) > 1 {
        b = append(c, PosId{Right, m.Site})
      } else {
          b = append(c, PosId{Right, ""})
      }
      // Set the disambiguator of the previous position if I'm a mini-node
      if len(t.MiniNodes) > 1 {
        if len(b) > 1 {
          b[len(b)-2].Site = m.Site
        }
      }
      m.Right.infix(b, paths, n)
    }
  }

  // Go right from major node
  if t.Right != nil {
    b = append(p, PosId{Right, ""})
    t.Right.infix(b, paths, n)
  }
}


// Return an array of all nodes (and their paths) in the tree in infix order
func (t *Treedoc) traverse() ([]Path, []*Node) {
  var paths []Path
  var nodes []*Node
  t.infix(Path{}, &paths, &nodes)
  return paths, nodes
}


// Return an array of only visible nodes (and their paths) in the tree in infix order
func (t *Treedoc) traverseVisible() ([]Path, []*Node) {
  var paths []Path
  var nodes []*Node
  t.infix(Path{}, &paths, &nodes)

  var visNodes []*Node
  var visPaths []Path
  for i,n := range nodes {
    if !n.Tombstone {
      visNodes = append(visNodes, n)
      visPaths = append(visPaths, paths[i])
    }
  }

  return visPaths, visNodes
}


// Generate a new path Id that will be between uidp and uidf
func (t *Treedoc) newUid(uidp *Path, uidf *Path, site Disambiguator) (Path, error) {
  var s **Treedoc
  p := Path{}

  switch {
  case uidp.ancestor(uidf):
    p = append(*uidf, PosId{Left, site})
    s = t.walk(*uidf)
  case uidf.ancestor(uidp):
    p = append(*uidp, PosId{Right, site})
    s = t.walk(*uidp)
  default:
    p = append(*uidp, PosId{Right, site})
    s = t.walk(*uidp)
  }

  // Check if we can remove the disambiguator from the end of the previous path
  // i.e. it is not a mini-node
  if len((*s).MiniNodes) == 1 {
    p[len(p)-2].Site = ""
  }

  return p, nil
}





// Path functions

// Is the node at p an ancestor of the node at q
func (p *Path) ancestor(q *Path) bool {
  switch {
  case (*p)[0].Dir == Empty && (*q)[0].Dir != Empty:  // p is root and q is not
    return true
  case p.prefix(q) > 0 && len(*p) < len(*q):  // the two nodes have some part of the path in common
    return true
  default:
    return false
  }
}


// Does p represent a node that is shown before the node at q
func (p *Path) before(q *Path) bool {
  //fmt.Println("Calling before with ", p, q)
  switch {
  case p.prefix(q) == len(*p) && (*q)[len(*p)].Dir == Right:
    return true
  case q.prefix(p) == len(*q) && (*p)[len(*q)].Dir == Left:
    return true
  default:
    preLen := p.prefix(q)
    i1 := (*p)[preLen]
    j1 := (*q)[preLen]
    //fmt.Printf("$$$$ len:%d, i1:%s, j1:%s\n.", preLen,i1.String(),j1.String())
    return i1.before(j1)
  }
}


// Does path p == path q
func (p *Path) equal(q *Path) bool {
  if len(*p) != len(*q) {
    return false
  }

  for i := range *p {
    if (*p)[i].Dir != (*q)[i].Dir || (*p)[i].Site != (*q)[i].Site {
      return false
    }
  }
  return true
}


// returns the longest common prefix of two paths
func (p *Path) prefix(q *Path) int {
  // if (*p)[0].Dir == (*q)[0].Dir && (*p)[0].Site == Empty {
  //   return 0
  // }
  if len(*p) >= len(*q) {
    //fmt.Println("Ranging over q")
    for i := range *q {
      //fmt.Printf("Checking: %v != %v || %s < %s\n.", (*p)[i].Dir,(*q)[i].Dir,(*p)[i].Site,(*q)[i].Site)
      if (*p)[i].Dir != (*q)[i].Dir || ((*p)[i].Site != "" && (*q)[i].Site != "" && (*p)[i].Site != (*q)[i].Site) {
        return i
      }
    }
    return len(*q)
  } else {
    for i := range *p {
      if (*p)[i].Dir != (*q)[i].Dir || ((*p)[i].Site != "" && (*q)[i].Site != "" && (*p)[i].Site != (*q)[i].Site) {
        return i
      }
    }
    return len(*p)
  }
}


// Does the path element r come before the path element s
func (r *PosId) before(s PosId) bool {
  switch {
  case r.Site == "" && s.Site == "" && r.Dir == Left && s.Dir == Right:
    return true
    //r.Site != "" && s.Site != "" &&
  case ((r.Dir == Left && s.Dir == Right) || (r.Dir == s.Dir && r.Site < s.Site)):
    return true
  case r.Dir == Left && s.Site != "":
    return true
  case r.Site != "" && s.Dir == Right:
    return true
  }
  return false
}
