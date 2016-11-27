package article

import (
  "fmt"
  "encoding/json"
  "strings"
  "io/ioutil"
  //"github.com/nickbradley/p2pwiki/util"
)

type LocalBuffer struct {
  Title string
  Paras []string  // paragraphs
  Log OpLog  // log of operations
}

func NewLocalBuffer(title string, paragraphs []string, site string) *LocalBuffer {
  opLog := OpLog{}
  opLog.Title = title
  opLog.Site = site
  a := LocalBuffer{title, paragraphs, opLog}
  return &a
}

func OpenLocalBuffer(path string, title string) (*LocalBuffer, error) {
  var lb LocalBuffer
  dat,err := ioutil.ReadFile(path + title + ".json")
  err = json.Unmarshal(dat, &lb)

  return &lb, err
}

// Insert a paragraph at the location specified
func (a *LocalBuffer) Insert(pos int, text string) error {
  switch {
  case pos < 1:
    fmt.Errorf("LocalBuffer::Insert(...) - Position must be greater than 1.")
  case pos > len(a.Paras):
    a.Paras = append(a.Paras, text)
  default:
    i := pos-1
    a.Paras = append(a.Paras, "")
    copy(a.Paras[i+1:], a.Paras[i:])
    a.Paras[i] = text
  }
  a.Log.Append("insert", OpArg{pos, text})

  return nil
}

// Remove the paragraph at the specified position
func (a *LocalBuffer) Delete(pos int) error {
  switch {
  case pos < 1:
    fmt.Errorf("LocalBuffer::Insert(...) - Position must be greater than 1.")
  case pos > len(a.Paras):
    _,a.Paras = a.Paras[len(a.Paras)-1], a.Paras[:len(a.Paras)-1]
  default:
    i := pos-1
    copy(a.Paras[i:], a.Paras[i+1:])
    a.Paras[len(a.Paras)-1] = ""
    a.Paras = a.Paras[:len(a.Paras)-1]
  }
  a.Log.Append("delete", OpArg{pos,""})

  return nil
}

func (a *LocalBuffer) Save(path string) error {
  b,err := json.Marshal(*a)
  if err != nil {
    return err
  }

  err = ioutil.WriteFile(path + a.Title + ".json", b, 0644)
  return err
}


func (a *LocalBuffer) String() string {
  return strings.Join(a.Paras, "\n")
}


type OpLog struct {
  Title string
  Site string
  Operations []string
  OpArgs []OpArg
}
func (r *OpLog) Append(operation string, args OpArg) {
  r.Operations = append(r.Operations, operation)
  r.OpArgs = append(r.OpArgs, args)
}
func (r *OpLog) Remove(upto int) error {
  i := upto - 1

  switch {
  case upto < 1:
    fmt.Errorf("OpLog::Remove() - Invalid position.")
  case upto > len(r.Operations):
    i = len(r.Operations) - 1
  }

  // delete from Operations array
  copy(r.Operations[i:], r.Operations[i+1:])
  r.Operations[len(r.Operations)-1] = ""
  r.Operations = r.Operations[:len(r.Operations)-1]

  // delete from OpArgs array
  copy(r.OpArgs[i:], r.OpArgs[i+1:])
  r.OpArgs[len(r.OpArgs)-1] = OpArg{}
  r.OpArgs = r.OpArgs[:len(r.OpArgs)-1]

  return nil
}

type OpArg struct {
  Pos int
  Text string
}
