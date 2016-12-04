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
  vpaths,_ := t.traverseVisible()
  paths,nodes := t.traverse()


  switch {
  case len(nodes) == 0:  // empty tree
    pid = append(pid, PosId{Empty, site})
    fmt.Println("Insert - inserting (in empty tree) at ", pid.String())
  case pos <= 1:  // inserting to the left of the tree
    lid := paths[0]  // get path to leftmost node
    lid[len(lid)-1].Site = ""  // remove the last disambiguator (so we insert to left of major node)
    pid = append(lid, PosId{Left, site})
    fmt.Println("Insert - inserting (leftmost) at ", pid.String())
  case pos > len(nodes):  // inserting to the right of the tree
    rid := paths[len(paths)-1]  // get path to rightmost node
    rid[len(rid)-1].Site = ""
    pid = append(rid, PosId{Right, site})
    fmt.Println("Insert - inserting (rightmost) at ", pid.String())
  default:  // inserting between two existing nodes
    p := vpaths[pos-2]
    f := vpaths[pos-1]

    if p.before(&f) {
      pid,_ = t.newUid(&p, &f, site)
      fmt.Printf("1) New path between %s and %s is %s.\n", p.String(), f.String(), pid.String())
    } else {
      pid,_ = t.newUid(&f, &p, site)
      fmt.Printf("2) New path between %s and %s is %s.\n", f.String(), p.String(), pid.String())
    }
  }

  return pid, t.insertNode(pid, &Node{atom, site, false, nil, nil})

/*


// special cases:
// Empty tree

// 1) insert between left node of major node and left mini-node
// 2) insert between right mini-node and right node of major node


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





    // p is left of a major node containing f and f is (the leftmost) mini-node

    // f is right of a major node containing p and p is (the rightmost) mini-node

    p := vpaths[pos-2]//t.path(pos-1)
    f := vpaths[pos-1]//t.path(pos)

    ppos,pIsLast := t.miniNodePos(p)
    fpos,fIsLast := t.miniNodePos(f)

    fmt.Println(ppos, pIsLast, fpos, fIsLast)
    switch {
    case ppos == 0 && pIsLast && fpos == 0 && !fIsLast:
      // remove f's Disambiguator
      f[len(f)-1].Site = ""
    case ppos > 0 && pIsLast && fpos == 0 && !fIsLast:
      // remove p's Disambiguator
      fmt.Println("&&&&&&&&&&&&&&&&&&&&&&&&&")
      p[len(p)-1].Site = ""
    }

    // if err != nil {
    //   return pid, fmt.Errorf("Treedoc::Insert(...) - Could not find path for given position.")
    // }
    // for _,path := range paths {
    //   fmt.Println(path.String())
    // }
    // for _,path := range vpaths {
    //   fmt.Println(path.String())
    // }
fmt.Printf("New path between %s and %s\n", p.String(), f.String())
    var q Path
    if p.before(&f) {
      q,_ = t.newUid(&p, &f, site)
      fmt.Printf("New path between %s and %s is %s\n", p.String(), f.String(), q.String())
    } else {
      fmt.Println("**********")
      q,_ = t.newUid(&f, &p, site)
      fmt.Printf("New path between %s and %s is %s\n", f.String(), p.String(), q.String())
    }

    // if err != nil {
    //   fmt.Println("Here si affsgfsjlkfdiofdikdh")
    //   return pid, fmt.Errorf("Treedoc::Insert(...) - Failed to find a path to insert.")
    // }
    pid =q[:]

  }
  fmt.Println(pid.String())
  return pid, t.insertNode(pid, &Node{atom, site, false, nil, nil})

  // // TODO
  return Path{}, nil
*/
}


func (t *Treedoc) miniNodePos(path Path) (int, bool) {
  fmt.Println("Got path", path)

   // special case when path is just empty node

  next := &t
  if len(path) > 1 {
    if len(path) == 2 && path[0].Dir == Empty {
      // through mininodes
      for _,m := range t.MiniNodes {
        if m.Site == path[0].Site {
         switch path[1].Dir {
         case Left:
           t = m.Left
         case Right:
           t = m.Right
         }
        }
      }
    } else {
      next = t.walk(path[:len(path)-1])
    }
  } else {
    switch path[0].Dir {
    case Left:
      t = t.Left
    case Right:
      t = t.Right
    }
    next = &t
  }

  // switch path[len(path)-1].Dir {
  // case Left:
  //   t = t.Left
  // case Right:
  //   t = t.Right
  // }
  // next = &t


fmt.Println((*next).MiniNodes)

  for i,p := range (*next).MiniNodes {
    if p.Site == path[len(path)-1].Site {
      if i == len((*next).MiniNodes) - 1 {
        return i, true
      } else {
        return i, false
      }
    }
  }
  return 0, true
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

  site := path[len(path)-1].Site
  for _,m := range (*n).MiniNodes {
    if m.Site == site {
      m.Tombstone = true
    }
  }

  return nil
}

func (t *Treedoc) walk(path Path) **Treedoc {
  pLen := len(path)

  if pLen == 0 || pLen == 1 && path[0].Dir == Empty {
    return &t
  }

  var midPath Path
  if path[0].Dir == Empty {  // walk through the root mini-nodes
    for _,m := range t.MiniNodes {
      if m.Site == path[0].Site {
        switch path[1].Dir {
        case Left:
          t = m.Left
        case Right:
          t = m.Right
        }
        break
      }
    }
    midPath = path[2:pLen-1]
  } else {
    midPath = path[:pLen-1]
  }

  // Walk the path
  if len(midPath) > 0 {
    for i,p := range midPath {
      switch {
      case i > 0 && midPath[i-1].Site != "":
        for _,m := range t.MiniNodes {
          if m.Site == midPath[i-1].Site {
            switch p.Dir {
            case Left:
              t = m.Left
            case Right:

              t = m.Right
            }
            break
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
  }

  // Set the return pointer
  switch {
  case pLen > 1 && path[pLen-2].Site != "":
    for _,m := range t.MiniNodes {
      if m.Site == path[pLen-2].Site {
        switch path[pLen-1].Dir {
        case Left:
          return &m.Left
        case Right:
          return &m.Right
        }
        break
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
  b := make(Path, len(p))
  copy(b, p)

  if t.Left != nil {
    b = append(p, PosId{Left, ""})
    t.Left.infix(b, paths, n)
  }

  for _,m := range t.MiniNodes {
    c := make(Path, len(p))
    copy(c, p)

    if len(c) == 0 {
      if len(t.MiniNodes) > 1 {
        c = append(p, PosId{Empty, m.Site})
      } else {
        c = append(p, PosId{Empty, ""})
      }
    }
    //  var site Disambiguator
    // // // Check if this is a disambiguated mini-node
    // if len(t.MiniNodes) > 1 {
    //   if len(c) > 1 {
    //     c[len(c)-1].Site = m.Site
    //   } else {
    //     c[0].Site = m.Site
    //   }
    //
    //   //site = m.Site
    //   //fmt.Println("Should append disambiguator ", m.Site, p)
    // }
    // //
    // // // Check if this is a leaf node
    // if t.Left == nil && t.Right == nil || len(t.MiniNodes) > 1 {
    //   fmt.Println("Should APPEND disambiguator ", m.Site, p)
    //   site = m.Site
    // }


    if m.Left != nil {
      if len(m.Left.MiniNodes) > 1 {
        b = append(c, PosId{Left, m.Site})
      } else {
        b = append(c, PosId{Left, ""})
      }

      if len(t.MiniNodes) > 1 {
        if len(b) > 1 {
          b[len(b)-2].Site = m.Site
        }
      }
      m.Left.infix(b, paths, n)
    }


    *n = append(*n, m)


    // if need disambiguate root:
    // if len(p) == 0 {
    //    c = append(p, PosId{Empty, m.Site})
    //  }
     if len(c) >= 1 {
       c[len(c)-1].Site = m.Site
     }
    // if len(p) >= 1 {
    //    //p[len(p)-1].Site = site
    //    //fmt.Println("Path is", b)
    //
    //
    //    c[len(c)-1].Site = m.Site
    //   *paths = append(*paths, c)
    // } else {
    //   *paths = append(*paths, append(c, PosId{Empty, m.Site}))
    // }
*paths = append(*paths, c)


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



  // if len(p) > 0 {
  //   p[len(p)-1].Site = ""
  // }

  if t.Right != nil {

    // if len(b) > 0 {
    //    b[len(b)-1].Site = ""
    // }

    //fmt.Println("Path is ", b, t.Right.MiniNodes[0].Value)
    b = append(p, PosId{Right, ""})
    t.Right.infix(b, paths, n)
  }
}


// func (t *Treedoc) infix(p Path, paths *[]Path, n *[]*Node) {
//   b := make(Path, len(p))
//   copy(b, p)
//   fmt.Printf("Infix(%v, %v, %v)\n", p, paths, n)
//   if t.Left != nil {
//     t.Left.infix(append(b, PosId{Left, ""}), paths, n)
//   }
//
//   for _,m := range t.MiniNodes {
//
//     var site Disambiguator
//     // Check if this is a disambiguated mini-node
//     if len(t.MiniNodes) > 1 {
//       site = m.Site
//     }
//
//     // Check if this is a leaf node
//     if t.Left == nil && t.Right == nil {
//       site = m.Site
//     }
//
//
//     if m.Left != nil {
//       m.Left.infix(append(b, PosId{Left, m.Site}), paths, n)
//     }
//
//
//     *n = append(*n, m)
//
//     if len(b) >= 1 {
//        b[len(b)-1].Site = site
//        //fmt.Println("Path is", b)
//       *paths = append(*paths, b)
//     } else {
//       *paths = append(*paths, append(b, PosId{Empty, m.Site}))
//     }
//
//
//
//     if m.Right != nil {
//       m.Right.infix(append(b, PosId{Right, m.Site}), paths, n)
//     }
//   }
//
//
//
//   // if len(p) > 0 {
//   //   p[len(p)-1].Site = ""
//   // }
//
//   if t.Right != nil {
//
//     // if len(b) > 0 {
//     //    b[len(b)-1].Site = ""
//     // }
//
//     //fmt.Println("Path is ", b, t.Right.MiniNodes[0].Value)
//     t.Right.infix(append(b, PosId{Right, ""}), paths, n)
//   }
// }





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

  // // Require uidp < uidf
  // if !uidp.before(uidf) {
  //   fmt.Println("***Here***")
  //   *uidp, *uidf = *uidf, *uidp
  //   //return p, fmt.Errorf("Treedoc::newUid() - uidp not before uidf.")
  // }

  // // Check if there is a node between uidp and uidf. If so, call newUid on the
  // // leftmost node such that uidp < uidm < uidf.
  // paths,_ := t.traverse()
  // for i,path := range paths {
  //   if path.equal(uidp) {
  //     if i < len(paths)-1 {
  //       if !paths[i+1].equal(uidf) {
  //         return t.newUid(uidp, &paths[i+1], site)
  //       }
  //     } else {
  //       // can't possibly be correct since uidp is at the end of the tree
  //     }
  //   }
  // }



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
    p = append(*uidf, PosId{Left, site})
  case uidf.ancestor(uidp):
    p = append(*uidp, PosId{Right, site})
  default:
    p = append(*uidp, PosId{Right, site})
  }


  // ppos, pIsLast := t.miniNodePos(p[:len(p)-1])
  // fmt.Println("Sending path", p[:len(p)-1],ppos, pIsLast)
  // if ppos == 0 && pIsLast {
  //   p[len(p)-2].Site = ""
  // }

  // TODO @Nick Check if we need this disabmiguator
p[len(p)-2].Site = ""
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
  case p.prefix(q) > 0 && len(*p) < len(*q):  // the two nodes have some part of the path in common
    return true
  default:
    return false
  }
}

// p < q
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


// returns the longest common prefix
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
