package treedoc

import "fmt"

// This is a standard binary tree except that each node can contain many sibling
// nodes.
type Treedoc struct {
  node *[]node
  left *Treedoc
  right *Treedoc
}

// Walk the tree rooted at t in infix order.
// Return the atoms of the non-tombstone nodes.
func (t *Treedoc) Contents() []string {
  var contents []string

  majorNodes := t.traverse()
  for _,majorNode := range majorNodes {
    for _,miniNode := range *majorNode {
      if !miniNode.tombstone {
        contents = append(contents, miniNode.value)
      }
    }
  }

  return contents
}

func (t *Treedoc) Delete(pos int, site string) error {
  nodes := t.traverse()
  (*nodes[pos-1])[0].tombstone = true
  return nil
}

/*
func (t *Treedoc) Insert(atom string, pos int, site string) error {
  fmt.Printf("Treedoc::Insert(%s, %d, %s)\n", atom, pos, site)
  if atom == "" {
    return fmt.Errorf("")
  }
  nodes := t.traverse()
  size := len(nodes) // where not tombstone and only first side-node
  newNode := node{atom,posId{path{}, site},false}
  var newId path

  // bootstrap
  // only root
  if len(nodes) == 1 {
    if pos <= 1 {
      fmt.Println("Inserting left.")
      newNode.id.path = path{false}
      //t.insertNode(&node{atom,posId{path{false}, site},false})
      //t.left = &Treedoc{&[]node{node{atom,posId{path{false}, site},false}}, nil, nil}
    } else {
      fmt.Println("Inserting right.")
      newNode.id.path = path{true}
      //t.insertNode(&node{atom,posId{path{true}, site},false})
      //t.right = &Treedoc{&[]node{node{atom,posId{path{true}, site},false}}, nil, nil}
    }
    t.insertNode(&newNode)
    return nil
  } else {
    switch {
    // generate newId that will put node in leftmost position of tree (first item dispalyed)
    case pos <= 1:
      fmt.Println("Inserting front")
      if len((*nodes[0])[0].id.path) == 0 { // leftmost node is root
        fmt.Println("Bootstrap left")
        t.left = &Treedoc{&[]node{node{atom,posId{path{false}, site},false}}, nil, nil}
        return nil
      }
      newId,_ = t.newUid(&node{}, &(*nodes[0])[0])
    // generate newId that will put node in rightmost position of tree (last item displayed)
    case pos > size:
      if len((*nodes[len(nodes)-1])[0].id.path) == 0 { // rightmost node is root
        t.right = &Treedoc{&[]node{node{atom,posId{path{true}, site},false}}, nil, nil}
        return nil
      }
      fmt.Println("Inserting end")
      newId,_ = t.newUid(&(*nodes[len(nodes)-1])[0], &node{})
    default:
      fmt.Println("Inserting middle")
      newId,_ = t.newUid(&(*nodes[pos-2])[0], &(*nodes[pos-1])[0])
    }
  }



  // //newId,_ := t.newUid(&(*nodes[pos-1])[0], &(*nodes[pos])[0])
  // fmt.Printf("Inserting '%s' at %v.\n", atom, newId)
  // for i := range newId[:len(newId)-1] {
  //   if newId[i] {
  //     fmt.Println("Going right")
  //     t = t.right
  //   } else {
  //     fmt.Println("Going left")
  //     t = t.left
  //   }
  // }
  // last := newId[len(newId)-1]
  // if last {
  //   fmt.Println("Set right")
  //   t.right = &Treedoc{&[]node{node{atom,posId{newId, site},false}}, nil, nil}
  // } else {
  //   fmt.Println("set left")
  //   t.left = &Treedoc{&[]node{node{atom,posId{newId, site},false}}, nil, nil}
  // }
  newNode = node{atom,posId{newId, site},false}
  t.insertNode(&newNode)

  return nil
}
*/


func (t *Treedoc) Insert(atom string, pos int, site string) error {
  nodes := t.traverse()
  //size := len(nodes) // where not tombstone and only first side-node
  newNode := node{atom,posId{path{}, site},false}

  switch {
  case pos <= 1:
    if len(nodes) == 1 || len((*nodes[0])[0].id.path) == 0 { // only root or leftmost node is root
      newNode.id.path = path{false}
    } else {
      p,_ := t.newUid(&node{}, &(*nodes[0])[0])
      newNode.id.path = p
    }
  case pos > len(nodes):
    if len(nodes) == 1 || len((*nodes[len(nodes)-1])[0].id.path) == 0 { // only root or rightmost node is root
      newNode.id.path = path{true}
    } else {
      p,_ := t.newUid(&(*nodes[len(nodes)-1])[0], &node{})
      newNode.id.path = p
    }
  default:
    p,_ := t.newUid(&(*nodes[pos-2])[0], &(*nodes[pos-1])[0])
    newNode.id.path = p
  }

  t.insertNode(&newNode)

  return nil
}






func (t *Treedoc) infix(n *[]*[]node) {
  if t.left != nil {
    t.left.infix(n)
  }
  *n = append(*n, t.node)
  if t.right != nil {
    t.right.infix(n)
  }
}


func (t *Treedoc) insertNode(n *node) {
  fmt.Printf("Treedoc::insertNode(%v)\n", n)
  path := n.id.path
  sid := n.id.disambiguator

  // TODO error checking

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
    *next = &Treedoc{&[]node{*n}, nil, nil}
  } else {
    // for _,miniNode := range *(*next).node {
    //   if miniNode.id.disambiguator == sid {
    //     fmt.Errorf("Mini-node with disambiguator already exists.")
    //   }
    // }
    // *(*next).node = append(*(*next).node, *n)
    pos := -1
    for i,miniNode := range *(*next).node {
      if miniNode.id.disambiguator > sid {
        pos = i
      } else if miniNode.id.disambiguator == sid {
        //error
      }
    }
    *(*next).node = append(*(*next).node, *n)
    if pos > -1 {
      copy((*(*next).node)[pos+1:], (*(*next).node)[pos:])
      (*(*next).node)[pos] = *n
    }
  }


  // if path[len(path)-1] {
  //   if t.right == nil {
  //     t.right = &Treedoc{&[]node{*n}, nil, nil}
  //   }
  // } else {
  //   t.left = &Treedoc{&[]node{*n}, nil, nil}
  // }
}


func (t *Treedoc) newUid(p *node, f *node) (path, error) {

  fmt.Printf("Treedoc::newUid(%v, %v)\n", p, f)

  uidp := p.id
  uidf := f.id
  if !uidp.before(&uidf) {
    fmt.Errorf("Position of p must come before position of f.")
  }

  var found bool
  var m *node
  nodes := t.traverse()
  for _,n := range nodes {
    // TODO range over all mininodes (not just first one)
    if uidp.before(&(*n)[0].id) && (*n)[0].id.before(&uidf) {
      found = true
      m = &(*n)[0]
      break
    }
  }
  fmt.Println("Found node:", m)

  switch {
  case found://no elements between p and f
    fmt.Println("newUid Case 1")
    return t.newUid(p, m)
  case p.ancestor(f):
    fmt.Println("newUid Case 2:",uidf.path)

    return append(uidf.path, false), nil
  case f.ancestor(p):
    fmt.Println("newUid Case 3")

    return append(uidp.path, true), nil
  default:
    fmt.Println("newUid Case default")

    return append(uidp.path, true), nil
  }
}
func (t *Treedoc) traverse() []*[]node {
  var node []*[]node
  t.infix(&node)
  return node
}





type node struct {
  //mininode map[string]string  // (disambiguator, atom)
  value string
  id posId  // once the node is inserted, the path will never change so store it
  tombstone bool
}
func (u *node) parent(v *node) bool {
  uid := u.id.path
  vid := v.id.path
  prefixLen := uid.commonPrefix(&vid)
  if len(vid) == prefixLen && len(uid) == len(vid) + 1 {
    return true
  }
  return false
}
func (u *node) ancestor(v *node) bool {
  uid := u.id.path
  vid := v.id.path
  if uid.prefix(&vid) {
    return true
  }
  return false
}
func (u *node) miniSibling(v *node) bool {
  uid := u.id.path
  vid := v.id.path
  return uid.equals(&vid)
}



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
    // fmt.Println("1")
    return true
  case len(u) > 0 && len(v) == 0 && u[0] == false: //(0..., e)
  // fmt.Println("2")
    return true
  case i > 0 && len(u) > i && len(v) == i && u[i] == true: //(c1...cn1..., c1...cn) NOT (1, e)
  // fmt.Println("3")
    return true
  case i > 0 && len(u) == i && len(v) > i && v[i] == false: //(c1...cn, c1...cn0...) NOT (e, 0)
  // fmt.Println("4")
    return true
  case len(u) > i && len(v) > i:
    if u[i] == false && v[i] == true {
      // fmt.Println("5")
      return true
    }
    if u[i] == v[i] && p.disambiguator < q.disambiguator {
      // fmt.Println("6")
      return true
    }
  case len(u) == i && len(v) == i && p.disambiguator < q.disambiguator:
    // fmt.Println("7")
    return true
  }

  return false
}

// func (p *posId) before(q *posId) bool {
//   if len(p.path) == 0 && len(q.path) == 0 {
//     return false
//   }
//   cm := p.path.commonPrefix(&q.path)
//   switch {
//   case len(p.path) == cm && q.path[cm] == false:
//     return false
//   case len(p.path) == cm && q.path[cm] == true:
//     return true
//   case len(q.path) == cm && p.path[cm] == true:
//     return false
//   case len(q.path) == cm && p.path[cm] == false:
//     return true
//   case p.path[cm] == false:
//     //return true
//     if cm > 0 { // have a common prefix
//       if len(p.path) == cm && len(q.path) > cm {
//         return true  // i1 does not exist so i1 < j1 trivially
//       }
//
//       i1 := p.path[cm]  // first element after common prefix
//       j1 := q.path[cm]
//       di := p.disambiguator
//       dj := q.disambiguator
//       i1LessThanj1 := i1 == false && j1 == true
//
//       switch {
//       case di == "" && dj == "" && i1LessThanj1://pi < pj:
//         return true
//       case di != "" && dj != "" && (i1LessThanj1 || (i1 == j1 && di < dj))://(pi < pj || (pi == pj && di < dj)):
//         return true
//       case di == "" && dj != "" && i1 == false: //0
//         return true
//       case di != "" && dj == "" && j1 == true: //1
//         return true
//       default:
//         return false
//       }
//     } else {
//       return false
//     }
//   default:
//     return false
//   }
// }


// func (p *posId) before(q *posId) bool {
//   switch {
//   case p.path.prefix(&q.path) && q.path[len(p.path)] == true: //1
//     return true
//   case p.path.prefix(&q.path) && q.path[len(p.path)] == false: //0
//     return false
//   default:
//     cplen := p.path.commonPrefix(&q.path)
//     if cplen > 0 { // have a common prefix
//       if len(p.path) == cplen && len(q.path) > cplen {
//         return true  // i1 does not exist so i1 < j1 trivially
//       }
//
//       i1 := p.path[cplen]  // first element after common prefix
//       j1 := q.path[cplen]
//       di := p.disambiguator
//       dj := q.disambiguator
//       i1LessThanj1 := i1 == false && j1 == true
//
//       switch {
//       case di == "" && dj == "" && i1LessThanj1://pi < pj:
//         return true
//       case di != "" && dj != "" && (i1LessThanj1 || (i1 == j1 && di < dj))://(pi < pj || (pi == pj && di < dj)):
//         return true
//       case di == "" && dj != "" && i1 == false: //0
//         return true
//       case di != "" && dj == "" && j1 == true: //1
//         return true
//       default:
//         return false
//       }
//     } else {
//       return false
//     }
//   }
// }



type path []bool
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
// return the length of the longest common prefix of p and q
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
