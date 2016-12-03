package article

import (
  "fmt"
  "strings"
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







// Walk the tree rooted at t in infix order.
// Return the atoms of the non-tombstone nodes.
func (t *Treedoc) Contents() []Atom {
  var contents []Atom
  _,nodes := t.traverseVisible()
  for _, node := range nodes {
    contents = append(contents, node.Value)
  }

  return contents
}


// Prevent a value of a node from being shown to a user.
// Does not remove node from tree,
func (t *Treedoc) Delete(pos int, site Disambiguator) (Path, error) {
  path := Path{}
  paths,nodes := t.traverseVisible()
  if pos > len(nodes) {
    return path, fmt.Errorf("Treedoc::Delete(...) - Position is invalid.")
  }
  nodes[pos-1].Tombstone = true

  return paths[pos-1], nil
}


func (t *Treedoc) Insert(pos int, atom Atom, site Disambiguator) (Path, error) {
  pid := Path{}
  vpaths,vnodes := t.traverseVisible()
  _,nodes := t.traverse()

  if pos > len(vnodes)+1 || pos < 0 {
    return pid, fmt.Errorf("Treedoc::Insert(...) - Position is invalid.")
  }


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
    p := vpaths[pos-2]//t.path(pos-1)
    f := vpaths[pos-1]//t.path(pos)
    // if err != nil {
    //   return pid, fmt.Errorf("Treedoc::Insert(...) - Could not find path for given position.")
    // }

    pid,err := t.newUid(&p, &f, site)
    if err != nil {
      return pid, fmt.Errorf("Treedoc::Insert(...) - Failed to find a path to insert.")
    }
  }

  return pid, t.insertNode(pid, &Node{atom, site, false, nil, nil})

  // // TODO
  // return Path{}, nil
}


func (t *Treedoc) insertNode(path Path, n *Node) error {
  next := t.walk(path)
  inserted := false

  if *next == nil {
    *next = &Treedoc{}
  } else {
    for i,m := range (*next).MiniNodes {
      if m.Site > n.Site {
        inserted = true
        (*next).MiniNodes = append((*next).MiniNodes, nil)
        copy((*next).MiniNodes[i+1:], (*next).MiniNodes[i:])
        (*next).MiniNodes[i] = n
        break
      }
    }
  }

  if !inserted {
    (*next).MiniNodes = append((*next).MiniNodes, n)
  }

  return nil
}

// Last position in path must have disabmiguator
func (t *Treedoc) deleteNode(path Path) error {
  n := t.walk(path)
  // if err != nil {
  //   return fmt.Errorf("Treedoc::deleteNode(...) - Invalid path.")
  // }

  site := path[len(path)-1].Site
  for _,m := range (*n).MiniNodes {
    if m.Site == site {
      m.Tombstone = true
    }
  }

  return nil
}

// // TODO @Nick check for nil pointers while walking through tree
// // Will return a pointer to the last major node
// func (t *Treedoc) walk(path Path) (*Treedoc, error) {
//   if len(p) == 0 || (len(p) == 1 && p[0].Dir == Empty) {
//     return t
//   }
//
//   for _,p := range path {
//     switch {
//     case p.Dir == Empty:
//       for _,m := range t.MiniNodes {
//         if m.Site == p.Site {
//           // FIXME @Nick
//           t = m  // change pointer from root major node to root mini-node
//         }
//       }
//     case p.Dir == Right && p.Site == "":  // right on major node
//       t = t.Right
//     case p.Dir == Left && p.Site == "":  // left on major node
//       t = t.Left
//     case p.Dir == Right:  // right on sibling node
//       for _,m := range t.MiniNodes {
//         if m.Site == p.Site {
//           t = m.Right
//         }
//       }
//     case p.Dir == Left:  // left on sibling node
//       for _,m := range t.MiniNodes {
//         if m.Site == p.Site {
//           t = m.Left
//         }
//       }
//     }
//   }
//
//   return t, nil
// }


// func (t *Treedoc) walk(path Path) (**Treedoc, error) {
//   //s := t
//   if len(path) == 0 || (len(path) == 1 && path[0].Dir == Empty) {
//     return &t, nil
//   }
//
//   // Ignore disambiguator for the last element   path[len(path)-1].Site = "kkk"
//
//   for i,p := range path[:len(path)-1] {
//     if i > 0 && path[i-1].Site != "" {
//       for _,m := range t.MiniNodes {
//         if m.Site == path[i-1].Site {
//           switch {
//           case p.Dir == Left:
//             t = m.Left
//           case p.Dir == Right:
//             t = m.Right
//           }
//         }
//       }
//     } else {
//       switch {
//       case p.Dir == Left:
//         fmt.Println("LEFT")
//         t = t.Left
//       case p.Dir == Right:
//         fmt.Println("RIGHT")
//         t = t.Right
//       }
//     }
//   }
//
//
//     if path[len(path)-2].Site != "" {
//       for _,m := range t.MiniNodes {
//         if m.Site == path[len(path)-2].Site {
//           if path[len(path)-1].Dir == Left {
//             fmt.Println("Returning mininode LEFT", m)
//             return &m.Left, nil
//           } else {
//             return &m.Right, nil
//           }
//         }
//       }
//     } else {
//
//     switch path[len(path)-1].Dir {
//     case Left:
//       fmt.Println("Returning LEFT")
//       return &t.Left, nil
//     case Right:
//       return &t.Right, nil
//     }
//   }
//   return &t, nil
// }

func (t *Treedoc) walk(path Path) **Treedoc {
  pLen := len(path)

  if pLen == 0 || pLen == 1 && path[0].Dir == Empty {
    return &t
  }

  // Walk the path
  for i,p := range path[:pLen-1] {
    switch {
    case i > 0 && path[i-1].Site != "":
      for _,m := range t.MiniNodes {
        if m.Site == path[i-1].Site {
          switch p.Dir {
          case Left:
            t = m.Left
          case Right:
            t = m.Right
          }
        }
      }
    default:
      switch p.Dir {
      case Left:
        t = t.Left
      case Right:
        t = t.Right
      }
    }
  }

  // Set the return pointer
  switch {
  case path[pLen-2].Site != "":
    for _,m := range t.MiniNodes {
      if m.Site == path[pLen-2].Site {
        switch path[pLen-1].Dir {
        case Left:
          return &m.Left
        case Right:
          return &m.Right
        }
      }
    }
  default:
    switch path[pLen-1].Dir {
    case Left:
      return &t.Left
    case Right:
      return &t.Right
    }
  }
  return nil // not called
}



// TODO @Nick this needs to be updated to match the interface
// Build a list of nodes in infix order
func (t *Treedoc) infix(p Path, paths *[]Path, n *[]*Node) {
  //fmt.Println(paths)
  if t.Left != nil {
    var site Disambiguator
    t.Left.infix(append(p, PosId{Left, site}), paths, n)

  }

  for _,m := range t.MiniNodes {
    if m.Left != nil {
      m.Left.infix(append(p, PosId{Left, m.Site}), paths, n)
    }


    *n = append(*n, m)

    //fmt.Println("Site is ", m.Site, m.Value)
    var s Disambiguator
    if len(t.MiniNodes) > 1 {
      s = m.Site
    }

    if len(p) >= 1 {
      p[len(p)-1].Site = s
      // fmt.Println("Appending site", p)
      q := Path{}
      q = p
      //copy(q, p)
      *paths = append(*paths, q)
    } else {
      //fmt.Println("Appending path", p)
      //p = append(p,PosId{Empty, site})
      *paths = append(*paths, append(p, PosId{Empty, s}))
    }

    // var site Disambiguator
    // if len(t.MiniNodes) > 1 {
    //   site = m.Site
    // }
    //
    // if len(p) > 1 {
    //   p[len(p)-1].Site = site
    //   *paths = append(*paths, p)
    // } else {
    //   *paths = append(*paths, append(p, PosId{Empty, site}))
    // }




    // if len(t.MiniNodes) > 1 {  // Annotate position with disambiguator
    //   if len(p) == 0 {
    //     //p = append(p, PosId{Empty, m.Site})
    //     *paths = append(*paths, append(p, PosId{Empty, m.Site}))
    //   } else {
    //   //p[len(p)-1].Site = m.Site
    // }
    // } else {
    //   *paths = append(*paths, p)
    // }




    if m.Right != nil {
      m.Right.infix(append(p, PosId{Right, m.Site}), paths, n)
    }
  }

  if t.Right != nil {
    var site Disambiguator
    t.Right.infix(append(p, PosId{Right, site}), paths, n)
  }
}


func (t *Treedoc) isEmpty() bool {
  return len(t.MiniNodes) == 0
}


// Return an array of all nodes in the tree in infix order
func (t *Treedoc) traverse() ([]Path, []*Node) {
  var paths []Path
  var nodes []*Node
  t.infix(Path{}, &paths, &nodes)
  return paths, nodes
}

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





func (t *Treedoc) newUid(uidp *Path, uidf *Path, site Disambiguator) (Path, error) {
  p := Path{}

  // Require uidp < uidf
  if !uidp.before(uidf) {
    return p, fmt.Errorf("Treedoc::newUid() - uidp not before uidf.")
  }

  // Check if there is a node between uidp and uidf. If so, call newUid on the
  // leftmost node such that uidp < uidm < uidf.
  paths,_ := t.traverse()
  for i,path := range paths {
    if path.equal(uidp) {
      if i < len(paths)-1 {
        if !paths[i+1].equal(uidf) {
          return t.newUid(uidp, &paths[i+1], site)
        }
      } else {
        // can't possibly be correct since uidp is at the end of the tree
      }
    }
  }

  // // TODO add the disabmiguator
  // // Check if there is a node between uidp and uidf. If so, call newUid on the
  // // leftmost node such that uidp < uidm < uidf.
  // if len(nodes) > 2 {
  //   for i := 0; i < len(nodes)-2; i++ {
  //     if nodes[i].Id.Path.equals(uidp.Path) {
  //       if !nodes[i+1].id.Path.equals(uidf) {
  //         return t.newUid(uidp, nodes[i+1].Id, site)
  //       }
  //     }
  //   }
  // }


  switch {
  case uidp.ancestor(uidf):
    p = append(*uidp, PosId{Left, site})
  case uidf.ancestor(uidp):
    p = append(*uidf, PosId{Right, site})
  default:
    p = append(*uidp, PosId{Right, site})
  }

  return p, nil
}





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
      //strt = "("+string(pid.Dir)+":"+string(pid.Site)+")"
    } else {
      str = append(str, dir)
    }
  }
  return strings.Join(str, "")
}

// Path functions
func (p *Path) ancestor(q *Path) bool {
  switch {
  case (*p)[0].Dir == Empty && (*q)[0].Dir != Empty:  // p is root and q is not
    return true
  case p.prefix(q) > 0:  // the two nodes have some part of the path in common
    return true
  default:
    return false
  }
}

// p < q
func (p *Path) before(q *Path) bool {
  switch {
  case p.prefix(q) == len(*p) && (*q)[len(*p)].Dir == Right:
    return true
  case q.prefix(p) == len(*q) && (*p)[len(*q)].Dir == Left:
    return true
  default:
    preLen := p.prefix(q)
    i1 := (*p)[preLen]
    j1 := (*q)[preLen]
    return i1.before(j1)
  }
}

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

func (r *PosId) before(s PosId) bool {
  switch {
  case r.Site == "" && s.Site == "" && r.Dir == Left && s.Dir == Right:
    return true
    case r.Site != "" && s.Site != "" && (
      (r.Dir == Left && s.Dir == Right) || (r.Dir == s.Dir && r.Site < s.Site)):
      return true
    case r.Dir == Left && s.Site != "":
      return true
    case r.Site != "" && s.Dir == Right:
      return true
  }
  return false
}


// returns the longest common prefix
func (p *Path) prefix(q *Path) int {
  if len(*p) >= len(*q) {
    for i := range *q {
      if (*p)[i].Dir != (*q)[i].Dir || (*p)[i].Site != (*q)[i].Site {
        return i
      }
    }
    return len(*q)
  } else {
    for i := range *p {
      if (*p)[i].Dir != (*q)[i].Dir || (*p)[i].Site != (*q)[i].Site {
        return i
      }
    }
    return len(*p)
  }
}


/*






// // u is a parent of v if they have the same common path and v's path is one longer
// // than u's path.
// func (u *node) parent(v *node) bool {
//   upath := u.id.path
//   vpath := v.id.path
//   prefixLen := upath.commonPrefix(&vpath)
//   if len(upath) == prefixLen && len(vpath) == len(upath) + 1 {
//     return true
//   }
//   return false
// }

// u is an ancestor of v if they share any common path and v's path is strictly
// longer than u's path. The root node (empty path) is an ancestor of all nodes.
func (u *Path) ancestor(v *Path) bool {

  if upath.prefix(&vpath) && len(upath) < len(vpath) {
    return true
  }
  return false
}

// // u is a mini-sibling (or side-node) of v if they both have the same path and
// // different disambiguators.
// func (u *node) miniSibling(v *node) bool {
//   uid := u.id
//   vid := v.id
//   return uid.path.equals(&vid.path) && uid.site != vid.site
// }






func (p *Path) before(q *Path) bool {
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

// // p is a prefix of q if p is the root node OR p and q agree in every position up
// // to the length of p. Note that we must have p <= q in terms of length.
// func (p *path) prefix(q *path) bool {
//   if len(*p) == 0 { // root
//     return true
//   }
//   if len(*p) > len(*q) {
//     return false
//   }
//
//   for i := range *p {
//     if (*p)[i] != (*q)[i] {
//       return false
//     }
//   }
//   return true
// }
// // return the length of the longest common prefix of p and q where the common
// // prefix is the first position in the bitstrings where the values disagree.
// func (p *path) commonPrefix(q *path) int {
//   if len(*p) <= len(*q) {
//     for i := range *p {
//       if (*p)[i] != (*q)[i] {
//         return i
//       }
//     }
//     return len(*p)
//   } else {
//     for i := range *q {
//       if (*p)[i] != (*q)[i] {
//         return i
//       }
//     }
//     return len(*q)
//   }
// }*/
